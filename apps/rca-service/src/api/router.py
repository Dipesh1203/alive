from fastapi import APIRouter 
from api.routes.uploadLog import router as upload_log_router
from api.routes.chat import router as chat_vector_store

router = APIRouter()

router.include_router(upload_log_router, prefix="/upload-log", tags=["Upload Log"],)
router.include_router(chat_vector_store, prefix="/chat", tags=["Chat"],)
