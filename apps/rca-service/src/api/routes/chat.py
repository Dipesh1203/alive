from fastapi import APIRouter, HTTPException
from db.chroma_connection import get_vector_store
from core.llm import analyze_logs
import traceback

router = APIRouter();

@router.get("/", summary="Chat Vector Store Endpoint")
def chat_vector_store(search: str):
    try:
      if not search:
          raise HTTPException(status_code=400, detail="Search query is required.")
      vectorStore = get_vector_store();
      retriever = vectorStore.as_retriever(search_kwargs={"k":2});
      result = retriever.invoke(search);
      log_analyzer = analyze_logs()
      llm_output = log_analyzer.invoke(result)
      return {"status": "success", "search_query": search,  "llm_output": llm_output}
    except Exception as e:
      traceback.print_exc()
      raise HTTPException(status_code=500, detail=f"An error occurred while processing the chat request: {str(e)}")