import enum
import uuid

from pydantic import BaseModel, Field


class GenderEnum(enum.Enum):
    MALE = "MALE"
    FEMALE = "FEMALE"


class GenderAllEnum(enum.Enum):
    MALE = "MALE"
    FEMALE = "FEMALE"
    ALL = "ALL"


class Client(BaseModel):
    id: uuid.UUID = Field(alias="client_id")
    login: str
    age: int
    location: str
    gender: GenderEnum


class Advertiser(BaseModel):
    advertiser_id: uuid.UUID
    name: str


class CampaignTargeting(BaseModel):
    age_from: int | None
    age_to: int | None
    gender: GenderAllEnum | None
    location: str | None


class CampaignEditable(BaseModel):
    ad_text: str
    ad_title: str
    clicks_limit: int
    impressions_limit: int
    cost_per_click: float
    cost_per_impression: float
    start_date: int
    end_date: int
    targeting: CampaignTargeting | None


class CampaignModerationResult(BaseModel):
    acceptable: bool
    reason: str


class Campaign(CampaignEditable):
    campaign_id: str
    advertiser_id: str
    image_path: str
    moderation_result: CampaignModerationResult | None

    @property
    def ad_title_with_exclamation(self) -> str:
        title = self.ad_title
        if self.moderation_result is not None and not self.moderation_result.acceptable:
            title += " ⚠️"
        return title


class Stats(BaseModel):
    impressions_count: int
    clicks_count: int
    spent_impressions: float
    spent_clicks: float
    spent_total: float
    conversion: float
    date: int | None = None
