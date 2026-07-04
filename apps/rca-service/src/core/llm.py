from langchain_google_genai import ChatGoogleGenerativeAI
from dotenv import load_dotenv
from langchaing_core.prompts import PromptTemplate
from core.documentLoader import load_json_data

load_dotenv()

template= "You are a Senior developer and a DevOps Engineer. You are an expert in analyzing logs and identifying root causes of issues. You have been given a set of logs from a system that is experiencing issues. Your task is to analyze the logs and identify the root cause of the issue. Please provide a detailed explanation of your analysis and the steps you took to identify the root cause. Here are the logs:{logs}. Please provide your analysis in a clear and concise manner, and include any relevant information that may help in understanding the issue.
"

model = ChatGoogleGenerativeAI(model="gemini-2.5-flash");

promptTemplate = PromptTemplate(
    input_variables=["logs"],
    template=template
)

data = load_json_data("data.json");
com

chain = data | promptTemplate | model

