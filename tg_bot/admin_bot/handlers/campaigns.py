import operator
from typing import Any, cast

from aiogram import F, types
from aiogram.enums import ParseMode
from aiogram_dialog import Dialog, DialogManager, Window
from aiogram_dialog.widgets.input import ManagedTextInput, TextInput
from aiogram_dialog.widgets.kbd import Button, Column, Row, Select, SwitchTo
from aiogram_dialog.widgets.text import Const, Format, Jinja

from admin_bot.misc import BACK, BACK_CANCEL, CANCEL
from admin_bot.states import CampaignsSG, StatsSG
from admin_bot.widgets.start_same_data import StartWithSameData
from ads_api import AdvertiserApiClient

AD_TEXT_SHOW_LIMIT = 64


async def campaigns_getter(api: AdvertiserApiClient, **_: Any) -> dict[str, Any]:
    return {
        "campaigns": await api.get_campaigns(size=100, page=1),
    }


async def campaign_getter(
    api: AdvertiserApiClient, dialog_manager: DialogManager, **_: Any
) -> dict[str, Any]:
    campaign_id = dialog_manager.start_data["campaign_id"]
    return {
        "campaign": await api.get_campaign_by_id(campaign_id),
        "current_date": await api.get_date(),
    }


async def campaign_select_handler(
    _: Any, __: Any, manager: DialogManager, campaign_id: Any
) -> None:
    await manager.start(CampaignsSG.show, {"campaign_id": campaign_id})


async def confirm_delete_callback(
    call: types.CallbackQuery, __: Any, manager: DialogManager
) -> None:
    campaign_id = manager.start_data["campaign_id"]
    api = cast(AdvertiserApiClient, manager.middleware_data["api"])
    await api.delete_campaign(campaign_id)
    await call.answer("Удалено.")
    await manager.done()


async def edit_field_handler(
    message: types.Message,
    widget: ManagedTextInput,
    manager: DialogManager,
    value: str | int | float,
) -> None:
    if isinstance(value, int | float) and value < 0:
        await message.answer("Число не может быть меньше 0.")
        return

    api = cast(AdvertiserApiClient, manager.middleware_data["api"])
    campaign = await api.get_campaign_by_id(manager.start_data["campaign_id"])
    setattr(campaign, widget.widget_id, value)
    await api.update_campaign(manager.start_data["campaign_id"], campaign)
    await manager.done()


campaigns_dialog = Dialog(
    Window(
        Const("📣 Кампании:"),
        Column(
            Select(
                Format("{item.ad_title_with_exclamation}"),
                "campaigns_selector",
                item_id_getter=operator.attrgetter("campaign_id"),
                items="campaigns",
                on_click=campaign_select_handler,
            ),
        ),
        BACK_CANCEL,
        state=CampaignsSG.list,
        getter=campaigns_getter,
    ),
    Window(
        Jinja("""
<b>Кампания {{ campaign.ad_title }}</b>
ID: <code>{{ campaign.campaign_id }}</code>

<blockquote expandable>{{ campaign.ad_text }}</blockquote>

<b>Просмотры</b>: {{ campaign.impressions_limit }} x {{ campaign.cost_per_impression }}$ каждый
<b>Клики</b>: {{ campaign.clicks_limit }} x {{ campaign.cost_per_click }}$ каждый
<b>Даты</b>: с {{ campaign.start_date }} по {{ campaign.end_date }}

{% if campaign.targeting.gender.value == 'MALE' %}
Пол ЦА: 👨 мужской
{% elif campaign.targeting.gender == 'FEMALE' %}
Пол ЦА: 👩 женский
{% else %}
Пол ЦА: любой
{% endif %}
{% if campaign.targeting.age_from is not none and campaign.targeting.age_to is not none %}
Возраст ЦА: с {{ campaign.targeting.age_from }} по {{ campaign.targeting.age_to }}
{% elif campaign.targeting.age_from is not none %}
Возраст ЦА: с {{ campaign.targeting.age_from }}
{% elif campaign.targeting.age_to is not none %}
Возраст ЦА: по {{ campaign.targeting.age_to }}
{% endif %}
{% if campaign.targeting.location is not none %}
Локация ЦА: {{ campaign.targeting.location }}
{% endif %}

{% if campaign.moderation_result is not none and not campaign.moderation_result.acceptable %}
⚠️ Не прошло модерацию: {{ campaign.moderation_result.reason }}
{% endif %}
"""),
        Row(
            StartWithSameData(
                Const("✏️ Название"), "edit_title", CampaignsSG.edit_title
            ),
            StartWithSameData(Const("✏️ Текст"), "edit_text", CampaignsSG.edit_text),
        ),
        Row(
            StartWithSameData(
                Const("✏️ Цена просмотра"),
                "edit_view_cost",
                CampaignsSG.edit_view_cost,
            ),
            StartWithSameData(
                Const("✏️ Цена клика"),
                "edit_click_cost",
                CampaignsSG.edit_click_cost,
            ),
        ),
        Row(
            StartWithSameData(
                Const("✏️ Лимит просмотров"),
                "edit_view_limit",
                CampaignsSG.edit_view_limit,
            ),
            StartWithSameData(
                Const("✏️ Лимит кликов"),
                "edit_click_limit",
                CampaignsSG.edit_click_limit,
            ),
            when=F["current_date"] < F["campaign"].start_date,
        ),
        StartWithSameData(Const("📊 Статистика"), "stats", StatsSG.show),
        SwitchTo(Const("🗑 Удалить"), "delete_campaign", CampaignsSG.delete_confirm),
        BACK_CANCEL,
        state=CampaignsSG.show,
        getter=campaign_getter,
        parse_mode=ParseMode.HTML,
    ),
    Window(
        Format(
            "Вы действительно хотите удалить кампанию {campaign.ad_title}? "
            "Это действие нельзя отменить."
        ),
        Button(Const("🗑 Да, удалить"), "delete_campaign", confirm_delete_callback),
        BACK,
        state=CampaignsSG.delete_confirm,
        getter=campaign_getter,
    ),
    Window(
        Const("Введите новое название кампании:"),
        CANCEL,
        TextInput("ad_title", on_success=edit_field_handler),
        state=CampaignsSG.edit_title,
    ),
    Window(
        Const("Введите новый текст кампании:"),
        CANCEL,
        TextInput("ad_text", on_success=edit_field_handler),
        state=CampaignsSG.edit_text,
    ),
    Window(
        Const("Введите новую цену за просмотр:"),
        CANCEL,
        TextInput(
            "cost_per_impression", type_factory=float, on_success=edit_field_handler
        ),
        state=CampaignsSG.edit_view_cost,
    ),
    Window(
        Const("Введите новую цену за клик:"),
        CANCEL,
        TextInput("cost_per_click", type_factory=float, on_success=edit_field_handler),
        state=CampaignsSG.edit_click_cost,
    ),
    Window(
        Const("Введите новый лимит просмотров:"),
        CANCEL,
        TextInput("impressions_limit", type_factory=int, on_success=edit_field_handler),
        state=CampaignsSG.edit_view_limit,
    ),
    Window(
        Const("Введите новый лимит кликов:"),
        CANCEL,
        TextInput("clicks_limit", type_factory=int, on_success=edit_field_handler),
        state=CampaignsSG.edit_click_limit,
    ),
)
