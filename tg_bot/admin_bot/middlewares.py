from collections.abc import Awaitable, Callable
from typing import Any

from aiogram import BaseMiddleware
from aiogram.enums import ChatType
from aiogram.types import TelegramObject

from admin_bot import config
from ads_api import AdvertiserApiClient, uuid_from_id


class ChatTypeMiddleware(BaseMiddleware):
    async def __call__(
        self,
        handler: Callable[[TelegramObject, dict[str, Any]], Awaitable[Any]],
        event: TelegramObject,
        data: dict[str, Any],
    ) -> Any:
        event_chat = data.get("event_chat")
        if event_chat and event_chat.type != ChatType.PRIVATE and event.message:
            return None

        return await handler(event, data)


class ApiInjectorMiddleware(BaseMiddleware):
    async def __call__(
        self,
        handler: Callable[[TelegramObject, dict[str, Any]], Awaitable[Any]],
        event: TelegramObject,
        data: dict[str, Any],
    ) -> Any:
        event_from_user = data.get("event_from_user")
        if not event_from_user:
            return await handler(event, data)

        advertiser_id = uuid_from_id(event_from_user.id)
        client = AdvertiserApiClient(config.API_BASE_URL, advertiser_id)
        async with client:
            return await handler(event, data | {"api": client})
