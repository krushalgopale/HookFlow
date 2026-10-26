import os

from dotenv import load_dotenv
from openai import AsyncOpenAI

load_dotenv()

_client = None


def get_groq_client() -> AsyncOpenAI:
    global _client

    if _client is None:
        api_key = os.getenv("GROQ_API_KEY")

        if not api_key:
            raise RuntimeError("GROQ_API_KEY is not configured")

        _client = AsyncOpenAI(
            api_key=api_key,
            base_url="https://api.groq.com/openai/v1",
        )

    return _client


def get_model() -> str:
    return os.getenv("GROQ_MODEL", "openai/gpt-oss-20b")
