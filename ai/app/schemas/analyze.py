from enum import Enum

from pydantic import BaseModel, ConfigDict, Field


class FailureRequest(BaseModel):
    status_code: int = Field(ge=100, le=599)
    response_body: str = Field(max_length=5000)
    attempt: int = Field(ge=1)


class Confidence(str, Enum):
    low = "low"
    medium = "medium"
    high = "high"


class FailureResponse(BaseModel):
    model_config = ConfigDict(extra="forbid")

    likely_cause: str
    evidence: list[str]
    recommended_action: str
    confidence: Confidence
