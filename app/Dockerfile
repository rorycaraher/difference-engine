# Stage 1: Use jrottenberg/ffmpeg as the first base image
FROM jrottenberg/ffmpeg:4.4-ubuntu2004 AS ffmpeg

# Stage 2: Use python:3.12-slim as the second base image
FROM python:3.12-slim

# Install necessary tools and dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    apt-utils \
    ffmpeg \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy ffmpeg binaries from the first stage to the second
COPY --from=ffmpeg /usr/local /usr/local

WORKDIR /app
COPY . /app
RUN pip install --no-cache-dir -r requirements.txt

ENV PYTHONUNBUFFERED=1

EXPOSE 5000
CMD ["python", "app.py"]
