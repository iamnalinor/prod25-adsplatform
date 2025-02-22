import os
from pathlib import Path

ADMIN_BOT_TOKEN = (
    os.environ.get("ADMIN_BOT_TOKEN")
    or Path("/run/secrets/admin_bot_token").read_text().strip()
)

API_BASE_URL = os.environ["API_BASE_URL"]
REDIS_URL = os.environ["REDIS_URL"]
