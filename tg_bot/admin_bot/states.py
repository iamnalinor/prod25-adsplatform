from aiogram.fsm.state import State, StatesGroup


class MainSG(StatesGroup):
    start = State()


class RegisterSG(StatesGroup):
    name = State()
    done = State()


class ChangeDateSG(StatesGroup):
    date = State()


class CreateCampaignSG(StatesGroup):
    ad_title = State()
    ad_text = State()
    suggest_ad_text = State()
    cost_per_impression = State()
    impressions_limit = State()
    cost_per_click = State()
    clicks_limit = State()
    start_date = State()
    end_date = State()
    targeting_gender = State()
    targeting_age_from = State()
    targeting_age_to = State()
    targeting_location = State()
    confirm = State()


class CampaignsSG(StatesGroup):
    list = State()
    show = State()
    delete_confirm = State()
    edit_title = State()
    edit_text = State()
    edit_view_cost = State()
    edit_click_cost = State()
    edit_view_limit = State()
    edit_click_limit = State()


class StatsSG(StatesGroup):
    show = State()
