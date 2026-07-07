from langchain_google_genai import ChatGoogleGenerativeAI
from dotenv import load_dotenv
from langchain_core.prompts import PromptTemplate
from langchain_core.runnables import RunnableLambda
from core.documentLoader import load_json_data
from langchain_core.documents.base import Document

load_dotenv()

def analyze_logs():
    template= "You are a Senior developer and a DevOps Engineer. You are an expert in analyzing logs and identifying root causes of issues. You have been given a set of logs from a system that is experiencing issues. Your task is to analyze the logs and identify the root cause of the issue. Please provide a detailed explanation of your analysis and the steps you took to identify the root cause. Here are the logs:{logs}. Please provide your analysis in a clear and concise manner, and include any relevant information that may help in understanding the issue."

    model = ChatGoogleGenerativeAI(model="gemini-2.5-flash");

    promptTemplate = PromptTemplate(
        input_variables=["logs"],
        template=template
    )
    def extract_logs_text(data: list[Document]) -> str:
        return "\n".join([doc.page_content for doc in data])

    chain = RunnableLambda(extract_logs_text)|(lambda text: {"logs": text}) | promptTemplate | model
    print("Chain type:",type(chain)," chain : ", chain);
    return chain
    



