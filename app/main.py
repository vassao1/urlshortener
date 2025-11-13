import os
from typing import Optional
from contextlib import asynccontextmanager
import fastapi
from fastapi.responses import JSONResponse, RedirectResponse
from pydantic import BaseModel, HttpUrl
from redis.asyncio import Redis as AsyncRedis
from psycopg_pool import AsyncConnectionPool
from dotenv import load_dotenv
from hashing import encode_id

load_dotenv()

@asynccontextmanager
async def lifespan(app: fastapi.FastAPI):
    global db_pool, redis_client
    
    db_pool = AsyncConnectionPool(
        conninfo=get_conninfo(),
        min_size=5,
        max_size=100,
    )

    redis_client = AsyncRedis(
        host=os.getenv("REDIS_HOST", "counter"), 
        port=int(os.getenv("REDIS_PORT", 6379)),
        decode_responses=True
    )
    
    async with db_pool.connection() as conn:
        async with conn.cursor() as cursor:
            await cursor.execute("""
                CREATE TABLE IF NOT EXISTS urls (
                    hash TEXT PRIMARY KEY,
                    original_url TEXT
                );
            """)
    await redis_client.setnx('Counter', 0)
    
    yield
    
    if db_pool:
        await db_pool.close()
    if redis_client:
        await redis_client.close()


app = fastapi.FastAPI(lifespan=lifespan)

db_pool: Optional[AsyncConnectionPool] = None
redis_client: Optional[AsyncRedis] = None


def get_conninfo() -> str:
    user = os.getenv("POSTGRES_USER")
    password = os.getenv("POSTGRES_PASSWORD")
    host = os.getenv("DB_HOST", "db")
    port = os.getenv("DB_PORT", 5432)
    database = os.getenv("POSTGRES_DB")
    return f"postgresql://{user}:{password}@{host}:{port}/{database}"


class UrlRequest(BaseModel):
    url: HttpUrl


@app.post("/short")
async def encode(data: UrlRequest):
    if not db_pool or not redis_client:
        raise fastapi.HTTPException(status_code=503, detail="Serviço indisponível")

    num = await redis_client.incr('Counter')
    hash_str = encode_id(num)
    
    async with db_pool.connection() as conn:
        async with conn.cursor() as cursor:
            await cursor.execute(
                "INSERT INTO urls (hash, original_url) VALUES (%s, %s)",
                (hash_str, str(data.url))
            )

    base_url = os.getenv("BASE_URL")
    return JSONResponse(
        status_code=201,
        content={"shortened_url": f"{base_url}/{hash_str}"}
    )

@app.get("/short/{hash_str}")
async def decode(hash_str: str):
    if not db_pool:
        raise fastapi.HTTPException(status_code=503, detail="Serviço indisponível")

    async with db_pool.connection() as conn:
        async with conn.cursor() as cursor:
            await cursor.execute(
                "SELECT original_url FROM urls WHERE hash = %s",
                (hash_str,)
            )
            result = await cursor.fetchone()

    if result:
        original_url = result[0]
        return RedirectResponse(url=original_url, status_code=301)
    else:
        return JSONResponse(
            status_code=404,
            content={"detail": "URL not found"}
        )