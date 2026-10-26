from fastapi import HTTPException
from openai import OpenAIError
from pydantic import ValidationError

from app.schemas.analyze import FailureResponse, FailureRequest
from app.services.analyze import analyze_failure


async def handle_analyze(data: FailureRequest) -> FailureResponse:
    try:
        return await analyze_failure(data)
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
