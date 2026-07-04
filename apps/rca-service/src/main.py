from fastapi import FastAPI
from api.router import router


app = FastAPI(title="Alive Telemetry Ingestion Hub",
    version="1.0.0")

app.include_router(router,prefix="/api")


@app.get("/health",tags=["Health Check"], summary="Health Check Endpoint")
async def health_check():
    return {"Status": "Alive Telemetry Ingestion Hub is running!"}

