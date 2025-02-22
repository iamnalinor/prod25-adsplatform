from aiogram import Bot, Dispatcher
from aiogram.fsm.storage.base import DefaultKeyBuilder
from aiogram.fsm.storage.redis import RedisStorage
from aiogram_dialog import setup_dialogs
from aiogram_dialog.widgets.kbd import Back, Cancel
from aiogram_dialog.widgets.text import Const

from admin_bot.config import ADMIN_BOT_TOKEN, REDIS_URL
from admin_bot.middlewares import ApiInjectorMiddleware, ChatTypeMiddleware

bot = Bot(token=ADMIN_BOT_TOKEN)
dispatcher = Dispatcher(
    storage=RedisStorage.from_url(
        REDIS_URL,
        key_builder=DefaultKeyBuilder(with_bot_id=True, with_destiny=True),
    )
)
dispatcher.update.middleware(ChatTypeMiddleware())
dispatcher.update.middleware(ApiInjectorMiddleware())
setup_dialogs(dispatcher)

CANCEL = Cancel(Const("❌ Отмена"))
BACK = Back(Const("⬅️ Назад"))
BACK_CANCEL = Cancel(Const("⬅️ Назад"))

__all__ = ["bot", "dispatcher", "CANCEL", "BACK", "BACK_CANCEL"]
