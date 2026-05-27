import asyncio
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel

from mixer import mixer
from mixer.config import STEMS_DIR

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost", "http://localhost:8000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

_mixer = mixer.Mixer()


class Values(BaseModel):
    de_values: dict


@app.post("/mixdown")
async def mixdown(values: Values):
    data = values.model_dump()
    stems = data["de_values"]["stems"]
    volumes = data["de_values"]["volumes"][: len(stems)]
    rounded_volumes = [round(v, 3) for v in volumes]
    input_files = _mixer.get_stems(STEMS_DIR, stems)

    loop = asyncio.get_event_loop()
    file_path = await loop.run_in_executor(
        None, _mixer.create_mixdown, input_files, rounded_volumes
    )

    return FileResponse(file_path, media_type="audio/mpeg", filename="mixdown.mp3")


app.mount("/", StaticFiles(directory="site", html=True), name="site")
