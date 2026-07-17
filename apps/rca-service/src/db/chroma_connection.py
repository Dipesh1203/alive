from functools import lru_cache 
from langchain_chroma import Chroma
from langchain_google_genai import GoogleGenerativeAIEmbeddings
from fastapi import Depends
from langchain_huggingface import HuggingFaceEmbeddings
from utils.config import CHROMA_API_KEY, CHROMA_TENANT, CHROMA_DATABASE

@lru_cache(maxsize=1)
def get_embeddings_google() -> GoogleGenerativeAIEmbeddings:
    return GoogleGenerativeAIEmbeddings(model="models/gemini-embedding-001")

@lru_cache(maxsize=1)
def get_embeddings_huggingface() -> HuggingFaceEmbeddings:
    return HuggingFaceEmbeddings(model_name="all-MiniLM-L6-v2")

@lru_cache(maxsize=1)
def get_vector_store() -> Chroma:
	return Chroma(
        collection_name="alive_rca_collection_huggingface",
        embedding_function=get_embeddings_huggingface(),
        chroma_cloud_api_key=CHROMA_API_KEY,
        tenant=CHROMA_TENANT,
        database=CHROMA_DATABASE,
    )