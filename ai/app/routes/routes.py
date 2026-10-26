from fastapi import APIRouter

from app.handlers.analyze import handle_analyze
from app.schemas.analyze import FailureResponse, FailureRequest

router = APIRouter()


@router.post("/analyze", response_model=FailureResponse)
async def analyze(data: FailureRequest):
    return await handle_analyze(data)

