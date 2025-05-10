# from FastAPI docs
from typing import Union
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel


import os
import json
from datetime import datetime
import uuid
import psycopg2
from psycopg2 import pool

from mixer import mixer


class Item(BaseModel):
    name: str
    description: str | None = None

app = FastAPI()

origins = [
    "http://localhost",
    "http://localhost:8000",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.mount("/site", StaticFiles(directory="site"), name="site")

class Values(BaseModel):
    de_values: dict

mixer = mixer.Mixer()
stems_dir = "/Users/rca/nltl/lvg-bucket/mp3/first-principles"
output_dir = "./output"

@app.get("/")
def read_root():
    return {"Hello": "World"}

@app.post('/mixdown')
async def test_mixdown(values: Values):
    data = values.dict()
    if 'de_values' in data:
        job_id = str(uuid.uuid4())
        timestamp = datetime.now()

        stems = data['de_values']['stems']
        volumes = data['de_values']['volumes'][0:len(stems)]

        rounded_volumes = [round(i, 3) for i in volumes]
        input_files = mixer.get_stems(stems_dir, stems)
        
        # do db write here

        file_response = mixer.create_mixdown(input_files, rounded_volumes)
        return json.dumps({"message": f"Success! {data}"}), 200
    else:
        return json.dumps({"error": "Invalid data"}), 400
