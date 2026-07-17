from langchain_google_genai import ChatGoogleGenerativeAI
from dotenv import load_dotenv
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.runnables import RunnableLambda
from core.documentLoader import load_json_data
from langchain_core.documents.base import Document
from langchain_ollama import ChatOllama
from langchain_core.output_parsers import StrOutputParser
from utils.config import NGROK_TUNNEL_URL

load_dotenv()

def analyze_logs(query:str = None):
    # system_instruction = "You are a Senior developer and a DevOps Engineer. You are an expert in analyzing logs and identifying root causes of issues."
    
    # user_instruction = (
    #     "You have been given a set of logs from a system that is experiencing issues. "
    #     "Your task is to analyze the logs and identify the root cause of the issue. "
    #     "Please provide a detailed explanation of your analysis and the steps you took to identify the root cause.\n\n"
    #     "Here are the logs:\n{logs}\n\n"
    #     "Please provide your analysis in a clear and concise manner, and include any relevant information that may help in understanding the issue."
    # )

    system_instruction = """
    You are an expert production incident analyst.
    Find the most probable root cause from logs.
    Be concise.
    """
    
    user_instruction = """
    Logs:
    {logs}
    Query: {query}
    
    Provide only:
    
    - Root Cause
    - Evidence
    - Fix
    - Query answers if asked in query.(Optional)
    
    Maximum 100 words.
    """

    model = ChatGoogleGenerativeAI(model="gemini-2.5-flash");
    llm = ChatOllama(
        base_url=NGROK_TUNNEL_URL,
        model="gemma4:12b", # or "batiai/gemma4-12b"
        temperature=0.3
    )
    promptTemplate = ChatPromptTemplate.from_messages([
        ("system", system_instruction),
        ("user", user_instruction)
    ])
    def extract_logs_text(data: list[Document]) -> str:
        return "\n".join([doc.page_content for doc in data])

    chain = RunnableLambda(extract_logs_text)|(lambda text: {"logs": text, "query": query}) | promptTemplate | llm |StrOutputParser()
    print("Chain type:",type(chain)," chain : ", chain);
    return chain
    



