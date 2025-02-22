from typing import Any, cast

from aiogram import types
from aiogram_dialog import Dialog, DialogManager, StartMode, Window
from aiogram_dialog.widgets.input import MessageInput
from aiogram_dialog.widgets.kbd import Start
from aiogram_dialog.widgets.text import Const

from admin_bot.states import MainSG, RegisterSG
from ads_api import AdvertiserApiClient, models


async def name_handler(
    message: types.Message, _: Any, dialog_manager: DialogManager
) -> None:
    api = cast(AdvertiserApiClient, dialog_manager.middleware_data["api"])
    await api.upsert_advertiser(
        models.Advertiser(advertiser_id=api.advertiser_id, name=message.text)
    )
    await dialog_manager.next()


register_dialog = Dialog(
    Window(
        Const(
            "Добро пожаловать в админ-панель рекламной платформы.\n\n"
            "Введите имя рекламодателя:"
        ),
        MessageInput(name_handler, content_types=types.ContentType.TEXT),
        state=RegisterSG.name,
    ),
    Window(
        Const("Регистрация пройдена!"),
        Start(Const("🏠 Домой"), "home", MainSG.start, mode=StartMode.RESET_STACK),
        state=RegisterSG.done,
    ),
)
