import uvicorn
from fastapi import FastAPI
from app.routes.routes import router

app = FastAPI(title="HookFlow AI Service")
app.include_router(router)


if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
