import asyncio
import operator
from typing import Any, cast

from aiogram import Bot, F, types
from aiogram_dialog import Dialog, DialogManager, ShowMode, Window
from aiogram_dialog.widgets.input import TextInput
from aiogram_dialog.widgets.kbd import Button, Column, Next, Select, SwitchTo
from aiogram_dialog.widgets.text import Const, Format, Jinja
from aiogram_dialog.widgets.text.base import Or

from admin_bot.misc import BACK, CANCEL
from admin_bot.states import CampaignsSG, CreateCampaignSG
from ads_api import AdvertiserApiClient, models


async def campaigns_create_getter(
    api: AdvertiserApiClient, dialog_manager: DialogManager, **_: Any
) -> dict[str, Any]:
    comment = dialog_manager.find("suggest_comment").get_value()
    suggestions = dialog_manager.dialog_data.get("suggestions")

    return {
        "advertiser": await api.get_advertiser(),
        "date": await api.get_date(),
        "ad_title": dialog_manager.find("ad_title").get_value(),
        "comment": comment,
        "comment_render": comment or "нет",
        "suggestions": suggestions,
        "suggestion_indexes": list(range(len(suggestions))) if suggestions else [],
    }


async def action_generator(bot: Bot, chat_id: int, action: str) -> None:
    # 30 x 4 => 120s
    for _ in range(30):
        await bot.send_chat_action(chat_id, action)
        await asyncio.sleep(4)


async def start_generating_handler(event: Any, _: Any, manager: DialogManager) -> None:
    if isinstance(event, types.Message):
        await event.answer("Генерируем текст...")
    elif isinstance(event, types.CallbackQuery):
        await event.answer()
        await event.message.answer("Генерируем текст...")
    manager.show_mode = ShowMode.SEND

    future = asyncio.ensure_future(
        action_generator(manager.middleware_data["bot"], event.from_user.id, "typing")
    )

    try:
        api = cast(AdvertiserApiClient, manager.middleware_data["api"])
        result = await api.ai_suggest_text(
            manager.find("ad_title").get_value(),
            manager.find("suggest_comment").get_value(),
        )
    finally:
        if not (future.done() or future.cancelled()):
            future.cancel()

    manager.dialog_data["suggestions"] = result


async def on_suggestion_selected(
    _: Any, __: Any, manager: DialogManager, pos: str
) -> None:
    text = manager.dialog_data["suggestions"][int(pos)]
    manager.find("ad_text").widget.set_widget_data(manager, text)
    await manager.next()


async def gte_zero_filter(message: types.Message, **_: Any) -> bool:
    try:
        value = float(message.text)
    except ValueError:
        await message.answer("Введите число.")
        return False

    if value < 0:
        await message.answer("Значение не может быть меньше 0.")
        return False

    return True


async def start_date_filter(
    message: types.Message, api: AdvertiserApiClient, **_: Any
) -> bool:
    if not message.text.isdigit():
        await message.answer("Введите число.")
        return False

    date = await api.get_date()
    if int(message.text) < date:
        await message.answer(f"Дата не может быть в прошлом. Текущий день: {date}")
        return False

    return True


async def end_date_filter(
    message: types.Message, dialog_manager: DialogManager, **_: Any
) -> bool:
    if not message.text.isdigit():
        await message.answer("Введите число.")
        return False

    start_date = dialog_manager.find("start_date").get_value()
    if int(message.text) < start_date:
        await message.answer("Дата окончания не может быть раньше даты начала.")
        return False

    return True


async def gender_callback(_: Any, __: Any, manager: DialogManager, gender: str) -> None:
    manager.dialog_data["targeting_gender"] = gender
    await manager.next()


async def age_to_filter(
    message: types.Message, dialog_manager: DialogManager, **_: Any
) -> bool:
    if not message.text.isdigit():
        await message.answer("Введите число.")
        return False

    age_from = dialog_manager.find("targeting_age_from").get_value() or 0
    if int(message.text) < age_from:
        await message.answer("Максимальный возраст не может быть меньше минимального.")
        return False

    return True


async def confirm_handler(_: Any, __: Any, manager: DialogManager) -> None:
    api = cast(AdvertiserApiClient, manager.middleware_data["api"])

    campaign_input = models.CampaignEditable(
        ad_title=manager.find("ad_title").get_value(),
        ad_text=manager.find("ad_text").get_value(),
        clicks_limit=manager.find("clicks_limit").get_value(),
        impressions_limit=manager.find("impressions_limit").get_value(),
        cost_per_click=manager.find("cost_per_click").get_value(),
        cost_per_impression=manager.find("cost_per_impression").get_value(),
        start_date=manager.find("start_date").get_value(),
        end_date=manager.find("end_date").get_value(),
        targeting=models.CampaignTargeting(
            gender=manager.dialog_data["targeting_gender"],
            age_from=manager.find("targeting_age_from").get_value(),
            age_to=manager.find("targeting_age_to").get_value(),
            location=manager.find("targeting_location").get_value(),
        ),
    )
    campaign = await api.create_campaign(campaign_input)

    await manager.done(show_mode=ShowMode.NO_UPDATE)
    await manager.start(
        CampaignsSG.show, {"campaign_id": campaign.campaign_id}, show_mode=ShowMode.AUTO
    )


campaigns_create_dialog = Dialog(
    Window(
        Format(
            "Создаём кампанию в организации {advertiser.name}.\n\n"
            "Введите название кампании:"
        ),
        CANCEL,
        TextInput("ad_title", on_success=Next()),
        state=CreateCampaignSG.ad_title,
        getter=campaigns_create_getter,
    ),
    Window(
        Const("Введите текст кампании:"),
        SwitchTo(
            Const("✨ Сгенерировать текст"),
            "suggest_text",
            CreateCampaignSG.suggest_ad_text,
        ),
        BACK,
        TextInput(
            "ad_text",
            on_success=SwitchTo(Const(""), "sw", CreateCampaignSG.cost_per_impression),
        ),
        state=CreateCampaignSG.ad_text,
        getter=campaigns_create_getter,
    ),
    Window(
        Format(
            "Вы можете сгенерировать рекламный текст при помощи ИИ.\n\n"
            "Параметры:\n"
            "Название организации: {advertiser.name}\n"
            "Название кампании: {ad_title}\n"
            "Комментарий: {comment_render}. "
            "Отправьте новый сообщением, чтобы изменить его",
            when=~F["suggestions"],
        ),
        Jinja(
            """
Варианты:

{% for suggestion in suggestions %}
{{ loop.index }}. {{ suggestion }}
{% endfor %}

Комментарий: {{ comment_render }} (чтобы изменить, отправьте новый сообщением)
""",
            when="suggestions",
        ),
        Select(
            Format("{pos}"),
            id="suggestions_selector",
            item_id_getter=lambda x: x,
            items="suggestion_indexes",
            on_click=on_suggestion_selected,
        ),
        Button(
            Or(
                Const("✨ Сгенерировать", when=~F["suggestions"]),
                Const("✨ Сгенерировать новые"),
            ),
            "start_generating",
            start_generating_handler,
        ),
        BACK,
        TextInput("suggest_comment"),
        state=CreateCampaignSG.suggest_ad_text,
        getter=campaigns_create_getter,
        preview_add_transitions=[Next()],
    ),
    Window(
        Const("Введите стоимость за показ:"),
        TextInput(
            "cost_per_impression",
            type_factory=float,
            on_success=Next(),
            filter=gte_zero_filter,
        ),
        SwitchTo(Const("⬅️ Назад"), "__back__", CreateCampaignSG.ad_text),
        state=CreateCampaignSG.cost_per_impression,
    ),
    Window(
        Const("Введите лимит показов:"),
        TextInput(
            "impressions_limit",
            type_factory=int,
            on_success=Next(),
            filter=gte_zero_filter,
        ),
        BACK,
        state=CreateCampaignSG.impressions_limit,
    ),
    Window(
        Const("Введите стоимость за клик:"),
        TextInput(
            "cost_per_click",
            type_factory=float,
            on_success=Next(),
            filter=gte_zero_filter,
        ),
        BACK,
        state=CreateCampaignSG.cost_per_click,
    ),
    Window(
        Const("Введите лимит кликов:"),
        TextInput(
            "clicks_limit",
            type_factory=int,
            on_success=Next(),
            filter=gte_zero_filter,
        ),
        BACK,
        state=CreateCampaignSG.clicks_limit,
    ),
    Window(
        Const("Введите дату начала (включительно):"),
        TextInput(
            "start_date", type_factory=int, on_success=Next(), filter=start_date_filter
        ),
        BACK,
        state=CreateCampaignSG.start_date,
    ),
    Window(
        Const("Введите дату окончания (включительно):"),
        TextInput(
            "end_date", type_factory=int, on_success=Next(), filter=end_date_filter
        ),
        BACK,
        state=CreateCampaignSG.end_date,
    ),
    Window(
        Const("Выберите пол целевой аудитории:"),
        Column(
            Select(
                Format("{item[0]}"),
                id="gender_selector",
                item_id_getter=operator.itemgetter(1),
                items=[
                    ("Мужской", "MALE"),
                    ("Женский", "FEMALE"),
                    ("Не таргетировать по полу", "ALL"),
                ],
                on_click=gender_callback,
            ),
        ),
        BACK,
        state=CreateCampaignSG.targeting_gender,
        preview_add_transitions=[Next()],
    ),
    Window(
        Const("Введите минимальный возраст целевой аудитории (включительно):"),
        TextInput(
            "targeting_age_from",
            type_factory=int,
            on_success=Next(),
            filter=gte_zero_filter,
        ),
        Next(Const("Пропустить")),
        BACK,
        state=CreateCampaignSG.targeting_age_from,
    ),
    Window(
        Const("Введите максимальный возраст целевой аудитории (включительно):"),
        TextInput(
            "targeting_age_to",
            type_factory=int,
            on_success=Next(),
            filter=age_to_filter,
        ),
        Next(Const("Пропустить")),
        BACK,
        state=CreateCampaignSG.targeting_age_to,
    ),
    Window(
        Const("Введите локацию целевой аудитории:"),
        TextInput("targeting_location", on_success=Next()),
        Next(Const("Пропустить")),
        BACK,
        state=CreateCampaignSG.targeting_location,
    ),
    Window(
        Const("Почти готово! Нажмите на кнопку, чтобы создать кампанию."),
        Button(Const("🔥 Создать"), "confirm", confirm_handler),
        BACK,
        state=CreateCampaignSG.confirm,
    ),
)
