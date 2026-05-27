import os
from dotenv import load_dotenv

load_dotenv()

STEMS_DIR = os.getenv("STEMS_DIR", "/Users/rca/nltl/lvg-bucket/mp3/first-principles")
OUTPUT_DIR = os.getenv("OUTPUT_DIR", "./output")
