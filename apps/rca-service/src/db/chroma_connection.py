from functools import lru_cache 
from langchain_chroma import Chroma
from langchain_google_genai import GoogleGenerativeAIEmbeddings
from fastapi import Depends
from utils.config import CHROMA_API_KEY, CHROMA_TENANT, CHROMA_DATABASE

@lru_cache(maxsize=1)
def get_embeddings() -> GoogleGenerativeAIEmbeddings:
    return GoogleGenerativeAIEmbeddings(model="models/gemini-embedding-001")

@lru_cache(maxsize=1)
def get_vector_store() -> Chroma:
	return Chroma(
        collection_name="alive_rca_collection",
        embedding_function=get_embeddings(),
        chroma_cloud_api_key=CHROMA_API_KEY,
        tenant=CHROMA_TENANT,
        database=CHROMA_DATABASE,
    )