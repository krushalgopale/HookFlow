from fastapi import HTTPException
from openai import OpenAIError
from pydantic import ValidationError

from app.schemas.investigate import InvestigationRequest, InvestigationResponse
from app.services.investigate import investigate_delivery


async def handle_investigate(data: InvestigationRequest) -> InvestigationResponse:
    try:
        return await investigate_delivery(data)
    except RuntimeError as exc:
        raise HTTPException(
            status_code=500,
            detail="AI service configuration is missing",
        ) from exc
    except OpenAIError as exc:
        raise HTTPException(
            status_code=502,
            detail="The AI provider request failed",
        ) from exc
    except (ValueError, ValidationError) as exc:
        raise HTTPException(
            status_code=502,
            detail="The AI provider returned an invalid response",
        ) from exc
