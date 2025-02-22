from typing import Any, cast

from aiogram import types
from aiogram_dialog import Dialog, DialogManager, Window
from aiogram_dialog.widgets.input import TextInput
from aiogram_dialog.widgets.text import Const
from aiogram_dialog.widgets.widget_event import ensure_event_processor

from admin_bot.misc import CANCEL
from admin_bot.states import ChangeDateSG
from ads_api import AdvertiserApiClient


@ensure_event_processor
async def date_handler(
    message: types.Message,
    _: Any,
    dialog_manager: DialogManager,
    date: int,
) -> None:
    api = cast(AdvertiserApiClient, dialog_manager.middleware_data["api"])
    await api.update_date(date)
    await message.answer("Обновлено.")
    await dialog_manager.done()


change_date_dialog = Dialog(
    Window(
        Const("Введите новую дату:"),
        CANCEL,
        TextInput("enter_date", type_factory=int, on_success=date_handler),
        state=ChangeDateSG.date,
    )
)
