FROM python:3.12-bookworm

# Install FFmpeg
RUN apt-get update && apt-get install -y ffmpeg \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Set the working# # Install Python dependencies
WORKDIR /app

COPY requirements.txt /app/
RUN pip install -r /app/requirements.txt

COPY . /app

ENV PYTHONUNBUFFERED=1

EXPOSE 8000

CMD ["fastapi", "run", "main.py", "--port", "8000"]
