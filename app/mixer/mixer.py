
import os
import random
from datetime import datetime
import ffmpeg

stems_dir = "/Users/rca/nltl/lvg-bucket/mp3/first-principles"
output_dir = "./output"

class Mixer:
    def __init__(self):
        stems_dir = "/Users/rca/nltl/lvg-bucket/mp3/first-principles"
        output_dir = "./output"

    def get_stems(self, directory, selected_stems):
        selected_file_paths = [os.path.join(directory, f"{file}.mp3") for file in selected_stems]
        return selected_file_paths

    def create_mixdown(self, input_files, random_volumes=[round(random.uniform(0.5, 1), 3) for _ in range(3)]):
        timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
        # Apply volume filters to each input file
        streams = []
        for i, input_file in enumerate(input_files):
            stream = (
                ffmpeg.input(input_file)
                .filter('volume', random_volumes[i])
            )
            streams.append(stream)

        # Combine the inputs using amix
        combined = (
            ffmpeg.filter(streams, 'amix', inputs=len(streams), duration='longest')
            .output(f'output/output-{timestamp}.mp3', acodec='libmp3lame', audio_bitrate='192k')
        )
        # Run the ffmpeg command
        ffmpeg.run(combined)
        return f'output/output-{timestamp}.mp3'

if __name__ == '__main__':
    mixer = Mixer()
    input_files = mixer.get_stems(stems_dir, [8, 3, 7])
    mixer.create_mixdown(input_files, [0.644, 0.517, 0.522])

