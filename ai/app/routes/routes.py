from fastapi import APIRouter

from app.handlers.analyze import handle_analyze
from app.schemas.analyze import FailureResponse, FailureRequest

from app.handlers.investigate import handle_investigate
from app.schemas.investigate import InvestigationRequest, InvestigationResponse

router = APIRouter()


@router.post("/analyze", response_model=FailureResponse)
async def analyze(data: FailureRequest):
    return await handle_analyze(data)


@router.post("/investigate", response_model=InvestigationResponse)
async def investigate(data: InvestigationRequest):
    return await handle_investigate(data)
