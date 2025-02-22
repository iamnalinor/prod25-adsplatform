import asyncio
import json
import typing
import uuid
from http.client import BAD_REQUEST, NETWORK_AUTHENTICATION_REQUIRED
from typing import Any

import aiohttp

from ads_api import models


class AdvertiserApiError(Exception):
    def __init__(self, status_code: int, message: str) -> None:
        self.status_code = status_code
        self.message = message


class AdvertiserApiClient:
    def __init__(self, base_url: str, advertiser_id: uuid.UUID) -> None:
        self.advertiser_id = advertiser_id
        self.session = aiohttp.ClientSession(base_url=base_url + "/")

    async def __aenter__(self) -> typing.Self:
        await self.session.__aenter__()
        return self

    async def __aexit__(self, exc_type: Any, exc_val: Any, exc_tb: Any) -> None:
        await self.session.__aexit__(exc_type, exc_val, exc_tb)

    async def request(
        self, method: str, url: str, **kwargs: Any
    ) -> dict[str, Any] | str:
        url = url.lstrip("/")
        async with self.session.request(method, url, **kwargs) as response:
            if BAD_REQUEST <= response.status <= NETWORK_AUTHENTICATION_REQUIRED:
                try:
                    data = await response.json()
                    message = data["message"]
                except (aiohttp.ContentTypeError, json.JSONDecodeError):
                    message = await response.text()
                raise AdvertiserApiError(response.status, message)

            if response.content_type == "application/json":
                return await response.json()

            return await response.text()

    async def get(self, url: str, **params: Any) -> Any:
        return await self.request("GET", url, params=params)

    async def post(self, url: str, data: Any) -> dict[str, Any]:
        return await self.request("POST", url, json=data)

    async def put(self, url: str, data: Any) -> dict[str, Any]:
        return await self.request("PUT", url, json=data)

    async def delete(self, url: str) -> None:
        await self.request("DELETE", url)

    async def get_advertiser(self) -> models.Advertiser:
        result = await self.get(f"/advertisers/{self.advertiser_id}")
        return models.Advertiser(**result)

    async def upsert_advertiser(self, data: models.Advertiser) -> None:
        await self.post("/advertisers/bulk", [data.model_dump(mode="json")])

    async def get_date(self) -> int:
        result = await self.get("/time")
        return result["current_date"]

    async def update_date(self, date: int) -> None:
        await self.post("/time/advance", {"current_date": date})

    async def get_campaigns(self, size: int, page: int) -> list[models.Campaign]:
        result = await self.get(
            f"/advertisers/{self.advertiser_id}/campaigns",
            size=size,
            page=page,
        )
        return [models.Campaign(**item) for item in result]

    async def create_campaign(self, data: models.CampaignEditable) -> models.Campaign:
        result = await self.post(
            f"/advertisers/{self.advertiser_id}/campaigns",
            data.model_dump(mode="json"),
        )
        return models.Campaign(**result)

    async def get_campaign_by_id(self, campaign_id: uuid.UUID) -> models.Campaign:
        result = await self.get(
            f"/advertisers/{self.advertiser_id}/campaigns/{campaign_id}"
        )
        return models.Campaign(**result)

    async def update_campaign(
        self, campaign_id: uuid.UUID, data: models.CampaignEditable
    ) -> models.Campaign:
        result = await self.put(
            f"/advertisers/{self.advertiser_id}/campaigns/{campaign_id}",
            data.model_dump(mode="json"),
        )
        return models.Campaign(**result)

    async def delete_campaign(self, campaign_id: str) -> None:
        await self.delete(f"/advertisers/{self.advertiser_id}/campaigns/{campaign_id}")

    async def ai_suggest_text(
        self, ad_title: str, comment: str | None = None
    ) -> list[str]:
        task = await self.post(
            f"/ai/advertisers/{self.advertiser_id}/suggestText",
            {"ad_title": ad_title, "comment": comment},
        )
        task_id = task["task_id"]

        # Short-polling with 5 seconds interval and 120 seconds timeout
        for _ in range(120 // 5):
            result = await self.get(f"/ai/tasks/{task_id}")
            if result["completed"]:
                return result["suggestions"]
            await asyncio.sleep(5)

        raise AdvertiserApiError(504, f"task {task_id} hasn't completed withing 120s")

    async def is_moderation_enabled(self) -> bool:
        res = await self.get("/ai/moderation/enabled")
        return res["enabled"]

    async def set_moderation_enabled(self, enabled: bool) -> None:
        await self.post("/ai/moderation/enabled", {"enabled": enabled})

    async def get_full_stats(self) -> models.Stats:
        res = await self.get(f"/stats/advertisers/{self.advertiser_id}/campaigns")
        return models.Stats(**res)

    async def get_stats_for_campaign(self, campaign_id: str) -> models.Stats:
        res = await self.get(f"/stats/campaigns/{campaign_id}")
        return models.Stats(**res)
