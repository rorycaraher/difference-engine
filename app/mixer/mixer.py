import os
import random
from datetime import datetime

import boto3
import ffmpeg
from botocore.config import Config

from mixer.config import (
    OUTPUT_DIR,
    R2_ACCOUNT_ID,
    R2_ACCESS_KEY_ID,
    R2_SECRET_ACCESS_KEY,
    R2_STEMS_BUCKET,
)


def _r2_client():
    return boto3.client(
        "s3",
        endpoint_url=f"https://{R2_ACCOUNT_ID}.r2.cloudflarestorage.com",
        aws_access_key_id=R2_ACCESS_KEY_ID,
        aws_secret_access_key=R2_SECRET_ACCESS_KEY,
        region_name="auto",
        config=Config(signature_version="s3v4"),
    )


class Mixer:
    def get_stems(self, directory, selected_stems):
        if R2_STEMS_BUCKET:
            client = _r2_client()
            return [
                client.generate_presigned_url(
                    "get_object",
                    Params={"Bucket": R2_STEMS_BUCKET, "Key": f"{stem}.mp3"},
                    ExpiresIn=300,
                )
                for stem in selected_stems
            ]
        return [os.path.join(directory, f"{stem}.mp3") for stem in selected_stems]

    def create_mixdown(self, input_files, random_volumes=None):
        if random_volumes is None:
            random_volumes = [round(random.uniform(0.5, 1), 3) for _ in range(len(input_files))]

        os.makedirs(OUTPUT_DIR, exist_ok=True)
        timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
        output_path = os.path.join(OUTPUT_DIR, f"output-{timestamp}.mp3")

        streams = [
            ffmpeg.input(f).filter("volume", vol)
            for f, vol in zip(input_files, random_volumes)
        ]
        combined = ffmpeg.filter(streams, "amix", inputs=len(streams), duration="longest").output(
            output_path, acodec="libmp3lame", audio_bitrate="192k"
        )
        ffmpeg.run(combined)
        return output_path


if __name__ == "__main__":
    from mixer.config import STEMS_DIR

    m = Mixer()
    files = m.get_stems(STEMS_DIR, [8, 3, 7])
    m.create_mixdown(files, [0.644, 0.517, 0.522])
