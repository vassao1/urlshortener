import fastapi
from fastapi.responses import JSONResponse, RedirectResponse
from pydantic import BaseModel, HttpUrl, ValidationError, field_validator
from hashing import encode_id, decode_id
import redis
from psycopg2 import pool
from dotenv import load_dotenv
import os

load_dotenv()

db_pool = pool.SimpleConnectionPool(
    minconn=5,
    maxconn=100,
    user=os.getenv("POSTGRES_USER"),
    password=os.getenv("POSTGRES_PASSWORD"),
    host=os.getenv("DB_HOST"),
    port=os.getenv("DB_PORT"),
    database=os.getenv("POSTGRES_DB")
)

conn = db_pool.getconn()
try:
    with conn.cursor() as cursor:
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS urls (
                hash TEXT PRIMARY KEY,
                original_url TEXT
            );
        """)
    conn.commit()
finally:
    db_pool.putconn(conn)

r = redis.Redis(host='counter', port=6379)
r.setnx('Counter', 0)

app = fastapi.FastAPI()

class UrlRequest(BaseModel):
    url: HttpUrl
    
@app.post("/short")
def encode(data: UrlRequest):
    num = r.incr('Counter')
    hash_str = encode_id(num)
    conn = db_pool.getconn()
    try:
        with conn.cursor() as cursor:
            cursor.execute(
                "INSERT INTO urls (hash, original_url) VALUES (%s, %s)",
                (hash_str, str(data.url))
            )
        conn.commit()
    finally:
        db_pool.putconn(conn)
    return JSONResponse(
        status_code=201,
        content={"shortened_url": f"{os.getenv('BASE_URL')}/{hash_str}"}
    )

@app.get("/short/{hash_str}")
def decode(hash_str: str):
    conn = db_pool.getconn()
    try:
        with conn.cursor() as cursor:
            cursor.execute(
                "SELECT original_url FROM urls WHERE hash = %s",
                (hash_str,)
            )
            result = cursor.fetchone()
            if result:
                original_url = result[0]
                return RedirectResponse(url=original_url, status_code=301)
            else:
                return JSONResponse(
                    status_code=404,
                    content={"detail": "URL not found"}
                )
    finally:
        db_pool.putconn(conn)