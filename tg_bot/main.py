import asyncio
import logging

from admin_bot.handlers import routers as admin_routers
from admin_bot.misc import bot, dispatcher

logging.basicConfig(
    level=logging.INFO,
    format="[%(levelname)s] %(asctime)s - %(name)s: %(message)s",
)


async def main() -> None:
    dispatcher.include_routers(*admin_routers)

    await bot.delete_webhook(drop_pending_updates=True)
    await dispatcher.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())
