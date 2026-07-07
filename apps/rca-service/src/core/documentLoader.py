import json
from langchain_community.document_loaders import JSONLoader
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_core.documents.base import Document

def load_json_data(file_path:str)-> str :
  loader = JSONLoader(file_path, jq_schema=".[] | {timestamp: .timestamp, message: .message, level: .level, service: .service, request_id: .request_id, aws_request_id: .aws_request_id, trace_id: .trace_id, environment: .environment, context: .context}", text_content=False)
  documents = loader.load()
  print("Loaded documents:", documents);
  # Recursive Text Splitting
  splitter = RecursiveCharacterTextSplitter(chunk_size=400, chunk_overlap=60);
  print("Splitting documents into chunks...");
  join_chunks = splitter.split_documents(documents);
  
  return join_chunks;
