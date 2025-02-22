FROM python:3.13.1-alpine3.19

WORKDIR /opt

ENV POETRY_VERSION=1.8
ENV PYTHONUNBUFFERED=1
ENV POETRY_VIRTUALENVS_CREATE=0

STOPSIGNAL SIGKILL

RUN pip install "poetry==$POETRY_VERSION"

COPY pyproject.toml poetry.lock ./
RUN poetry install --only=main --no-interaction --no-ansi

COPY . .

ENTRYPOINT ["python3", "main.py"]