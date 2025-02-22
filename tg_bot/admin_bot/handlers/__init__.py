from .campaigns import campaigns_dialog
from .campaigns_create import campaigns_create_dialog
from .change_date import change_date_dialog
from .register import register_dialog
from .start import start_dialog, start_router
from .stats import stats_dialog

routers = [
    start_router,
    campaigns_dialog,
    campaigns_create_dialog,
    change_date_dialog,
    register_dialog,
    start_dialog,
    stats_dialog,
]

__all__ = ["routers"]
