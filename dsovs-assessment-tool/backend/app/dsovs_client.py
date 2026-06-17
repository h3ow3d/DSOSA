from __future__ import annotations

import hashlib
import json
from typing import Any

import httpx

from app.config import settings


async def fetch_dsovs() -> tuple[dict[str, Any], str]:
    """Fetch the DSOVS JSON catalogue and return (data, sha256_hex)."""
    async with httpx.AsyncClient(timeout=30) as client:
        resp = await client.get(settings.dsovs_api_url)
        resp.raise_for_status()
        raw = resp.content

    sha256 = hashlib.sha256(raw).hexdigest()
    data = json.loads(raw)
    return data, sha256
