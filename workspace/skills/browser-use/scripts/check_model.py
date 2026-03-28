# pip install langchain-nvidia-ai-endpoints
api_key = "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"

from langchain_nvidia_ai_endpoints import ChatNVIDIA

llm = ChatNVIDIA(
    model="meta/llama-3.1-70b-instruct",
    api_key=api_key,
    temperature=0.6
)
#print(llm.available_models)
print("可用模型列表：")
for model in llm.available_models:
    print(f"- {model.id}")
