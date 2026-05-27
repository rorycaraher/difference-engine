FROM python:3.12-bookworm

COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

RUN apt-get update && apt-get install -y ffmpeg \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY pyproject.toml .
RUN uv sync --no-dev

COPY . .

ENV PYTHONUNBUFFERED=1

EXPOSE 8000

CMD ["uv", "run", "fastapi", "run", "main.py", "--port", "8000"]
