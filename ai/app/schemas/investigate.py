from datetime import datetime

from pydantic import BaseModel, ConfigDict, Field


class DeliveryAttempt(BaseModel):
    model_config = ConfigDict(extra="ignore")

    attempt_number: int = Field(ge=1)
    status: str
    response_status: int | None = None
    response_body: str | None = None
    error: str | None = None
    created_at: datetime


class Delivery(BaseModel):
    model_config = ConfigDict(extra="ignore")

    id: str
    event_type: str
    destination_name: str
    status: str
    response_status: int | None = None
    response_body: str | None = None
    error: str | None = None
    attempt_count: int = Field(ge=0)


class InvestigationRequest(BaseModel):
    model_config = ConfigDict(extra="forbid")

    delivery: Delivery
    attempts: list[DeliveryAttempt] = Field(max_length=100)


class InvestigationResponse(BaseModel):
    model_config = ConfigDict(extra="forbid")

    likely_cause: str
    evidence: list[str]
    recommended_action: str
    confidence: str
