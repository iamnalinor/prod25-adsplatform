from http import HTTPStatus
from typing import Any, cast

from aiogram import Router, types
from aiogram.filters import CommandStart
from aiogram_dialog import Dialog, DialogManager, ShowMode, StartMode, Window
from aiogram_dialog.widgets.kbd import Button, Start
from aiogram_dialog.widgets.text import Case, Const, Format

from admin_bot.states import (
    CampaignsSG,
    ChangeDateSG,
    CreateCampaignSG,
    MainSG,
    RegisterSG,
    StatsSG,
)
from ads_api.advertiser import AdvertiserApiClient, AdvertiserApiError


async def advertiser_getter(api: AdvertiserApiClient, **_: Any) -> dict[str, Any]:
    return {
        "advertiser": await api.get_advertiser(),
        "date": await api.get_date(),
        "is_moderation_enabled": await api.is_moderation_enabled(),
    }


async def toggle_moderation_handler(_: Any, __: Any, manager: DialogManager) -> None:
    api = cast(AdvertiserApiClient, manager.middleware_data["api"])
    enabled = await api.is_moderation_enabled()
    await api.set_moderation_enabled(not enabled)


start_dialog = Dialog(
    Window(
        Format(
            "Главное меню\n\n"
            "ID рекламодателя: {advertiser.advertiser_id}\n"
            "Имя: {advertiser.name}\n"
            "Текущая дата: {date}",
        ),
        Start(
            Const("➕ Создать кампанию"), "create_campaign", CreateCampaignSG.ad_title
        ),
        Start(Const("📣 Кампании"), "campaigns", CampaignsSG.list),
        Start(Const("🪪 Пройти регистрацию заново"), "register", RegisterSG.name),
        Start(Const("✏️ Изменить дату"), "change_date", ChangeDateSG.date),
        Button(
            Case(
                {
                    True: Const("🟢 Модерация включена"),
                    False: Const("🔴 Модерация выключена"),
                },
                "is_moderation_enabled",
            ),
            "toggle_moderation",
            toggle_moderation_handler,
        ),
        Start(Const("📊 Статистика"), "stats", StatsSG.show),
        getter=advertiser_getter,
        state=MainSG.start,
    )
)

start_router = Router()


@start_router.message(CommandStart())
async def start_cmd(
    message: types.Message,
    dialog_manager: DialogManager,
    api: AdvertiserApiClient,
) -> None:
    try:
        await api.get_advertiser()
    except AdvertiserApiError as e:
        if e.status_code == HTTPStatus.NOT_FOUND:
            await dialog_manager.start(RegisterSG.name, mode=StartMode.RESET_STACK)
            return
        raise e

    await dialog_manager.start(
        MainSG.start,
        mode=StartMode.RESET_STACK,
        show_mode=ShowMode.SEND,
    )
