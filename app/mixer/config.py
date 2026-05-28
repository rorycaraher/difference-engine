import os
from dotenv import load_dotenv

load_dotenv()

STEMS_DIR = os.getenv("STEMS_DIR", "/Users/rca/nltl/lvg-bucket/mp3/first-principles")
OUTPUT_DIR = os.getenv("OUTPUT_DIR", "./output")

# R2 — only set in production. Absence of R2_STEMS_BUCKET means local filesystem is used.
R2_ACCOUNT_ID      = os.getenv("R2_ACCOUNT_ID")
R2_ACCESS_KEY_ID   = os.getenv("R2_ACCESS_KEY_ID")
R2_SECRET_ACCESS_KEY = os.getenv("R2_SECRET_ACCESS_KEY")
R2_STEMS_BUCKET    = os.getenv("R2_STEMS_BUCKET")
