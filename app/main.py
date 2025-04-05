# from FastAPI docs
from typing import Union
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel

# existing lvg stuff
import os
import json
from datetime import datetime
import uuid
import psycopg2
from psycopg2 import pool
# mixdown stuff
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
    lvg_values: dict

mixer = mixer.Mixer()
stems_dir = "/Users/rca/nltl/lvg-bucket/mp3/first-principles"
output_dir = "./output"

@app.get("/")
def read_root():
    return {"Hello": "World"}

@app.post('/mixdown')
async def test_mixdown(values: Values):
    # return await request.json()
    # print(values)
    data = values.dict()
    # print(data)
    if 'lvg_values' in data:
        job_id = str(uuid.uuid4())
        timestamp = datetime.now()
        # cursor = conn.cursor()

        stems = data['lvg_values']['stems']
        volumes = data['lvg_values']['volumes'][0:len(stems)] # only need volume for each stem
        rounded_volumes = [round(i, 3) for i in volumes]
        input_files = mixer.get_stems(stems_dir, stems)
        print(input_files)
        file_response = mixer.create_mixdown(input_files, rounded_volumes)
        # return send_file(
        #     file_response,
        #     as_attachment=True,
        #     download_name=file_response,
        #     mimetype='audio/mpeg'
        # ), 200
        # return jsonify({"message": f"Success! {data}"}), 200
        return "Success"
    else:
        # return jsonify({"error": "Invalid data"}), 400    
        return "Invalid data"