from hashids import Hashids
from dotenv import load_dotenv
import os

load_dotenv()

alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
salt = os.getenv("saltseckey")
hashids = Hashids(salt=salt, alphabet=alphabet, min_length=4)

def encode_id(num: int) -> str:
    return hashids.encode(num)

def decode_id(hash_str: str) -> int:
    decoded = hashids.decode(hash_str)
    return decoded[0] if decoded else None