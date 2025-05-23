FROM python:3.12-slim

ENV PYTHONUNBUFFERED True

ENV APP_HOME /app

ENV PORT 8080

WORKDIR $APP_HOME

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY src/ ./src/
COPY .env.* ./

CMD exec gunicorn --bind :$PORT --workers 1 --threads 8 --timeout 0 src.main:app