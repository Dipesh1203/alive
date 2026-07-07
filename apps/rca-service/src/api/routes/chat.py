from fastapi import APIRouter, HTTPException
from db.chroma_connection import get_vector_store
from core.llm import analyze_logs

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
      return {"status": "success", "search_query": search,  "llm_output": llm_output.content}
    except Exception as e:
      raise HTTPException(status_code=500, detail=f"An error occurred while accessing the vector store: {str(e)}")