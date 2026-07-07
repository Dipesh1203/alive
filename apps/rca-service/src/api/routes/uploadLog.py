from fastapi import FastAPI,APIRouter, UploadFile, File, HTTPException
import os
import shutil
from db.chroma_connection import get_vector_store
from core.documentLoader import load_json_data
from uuid import uuid4

router = APIRouter();

TEMP_DIR = "./temp_logs"
os.makedirs(TEMP_DIR, exist_ok=True)

@router.post("/upload-log", summary="Upload Log Endpoint")
def upload_log(file:UploadFile = File(...)):
  print(f"Received file: {file.filename}")
  if not file.filename.endswith((".json", ".jsonl")):
    raise HTTPException(status_code=400, detail="Invalid file format. Only JSON files are allowed.")
  
  temp_file_path = os.path.join(TEMP_DIR, file.filename)
  try:
    with open(temp_file_path, "wb") as temp_file:
      shutil.copyfileobj(file.file, temp_file)
    loader = load_json_data(temp_file_path)
    vectorStore = get_vector_store();
    uuids = [str(uuid4()) for _ in range(len(loader))]
    result=vectorStore.add_documents(documents=loader,uuids=uuids);
    return {"status": "success","filename": file.filename, "message": "File uploaded and processed successfully.", "data": result}
  except Exception as e:
    raise HTTPException(status_code=500, detail=f"An error occurred while processing the file: {str(e)}")
  
