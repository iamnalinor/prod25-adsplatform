from typing import Any

from aiogram.enums import ParseMode
from aiogram_dialog import Dialog, DialogManager, Window
from aiogram_dialog.widgets.text import Format

from admin_bot.misc import BACK_CANCEL
from admin_bot.states import StatsSG
from ads_api import AdvertiserApiClient


async def stats_getter(
    api: AdvertiserApiClient, dialog_manager: DialogManager, **_: Any
) -> dict[str, Any]:
    if dialog_manager.start_data:
        campaign = await api.get_campaign_by_id(
            dialog_manager.start_data["campaign_id"]
        )
        stats = await api.get_stats_for_campaign(campaign.campaign_id)
        title = f"Статистика по кампании {campaign.ad_title}"
    else:
        stats = await api.get_full_stats()
        title = "Статистика по всем кампаниям текущего рекламодателя"

    return {
        "title": title,
        "s": stats,
    }


stats_dialog = Dialog(
    Window(
        Format("""
📊 <b>{title}</b>

Количество просмотров: {s.impressions_count}
Доход с просмотров: {s.spent_impressions:.2f}$

Количество кликов: {s.clicks_count}
Доход с кликов: {s.spent_clicks:.2f}$

Общий доход: {s.spent_total:.2f}$
Конверсия: {s.conversion:.2f}%
"""),
        BACK_CANCEL,
        getter=stats_getter,
        state=StatsSG.show,
        parse_mode=ParseMode.HTML,
    )
)
