import json
from langchain_community.document_loaders import JSONLoader
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_core.documents.base import Document

def load_json_data(file_path:str)-> str :
  loader = JSONLoader(file_path, jq_schema=".[] | {timestamp: .timestamp, message: .message, level: .level, service: .service, request_id: .request_id}", text_content=False)
  documents = loader.load()
  # Recursive Text Splitting
  splitter = RecursiveCharacterTextSplitter(chunk_size=150, chunk_overlap=0);
  
  join_chunks = splitter.split_documents(documents);
  
  return join_chunks;


load_json_data("data.json");