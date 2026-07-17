from fastapi import FastAPI
from api.router import router
from fastapi.middleware.cors import CORSMiddleware


app = FastAPI(title="Alive Telemetry Ingestion Hub",
    version="1.0.0")

origins = [
    "http://localhost:3000", 
    "http://localhost:8000",
    "http://localhost:8001",
]
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


app.include_router(router,prefix="/api")


@app.get("/health",tags=["Health Check"], summary="Health Check Endpoint")
async def health_check():
    return {"Status": "Alive Telemetry Ingestion Hub is running!"}

