# FROM jrottenberg/ffmpeg:4.4
FROM jrottenberg/ffmpeg:arm64

# Install PostgreSQL development libraries
RUN apt-get update && apt-get install -y libpq-dev python3-pip

# Install Python dependencies
COPY requirements.txt /app/
RUN pip install -r /app/requirements.txt

WORKDIR /app
COPY . /app

ENV PYTHONUNBUFFERED=1

EXPOSE 5000
CMD ["python", "app.py"]
