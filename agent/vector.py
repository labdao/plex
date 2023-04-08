from langchain.embeddings.openai import OpenAIEmbeddings
from langchain.vectorstores import Chroma
from langchain.text_splitter import CharacterTextSplitter
from langchain.llms import OpenAI
from langchain.chains import RetrievalQA
from langchain.document_loaders import TextLoader, WebBaseLoader
from langchain.document_loaders import ObsidianLoader
from pathlib import Path

llm = OpenAI(temperature=0)

# get the path to the text file
relevant_parts = []
for p in Path(".").absolute().parts:
    relevant_parts.append(p)
    if relevant_parts[-3:] == ["langchain", "docs", "modules"]:
        break
doc_path = str(Path(*relevant_parts) / "state_of_the_union.txt")

# load the text and obsidian file
loader = TextLoader(doc_path)
documents = loader.load()
loader_obsidian = ObsidianLoader("")
documents += loader_obsidian.load()

# split the documents into smaller chunks
text_splitter = CharacterTextSplitter(chunk_size=1000, chunk_overlap=0)
texts = text_splitter.split_documents(documents)

# create an embedding model
embeddings = OpenAIEmbeddings()

# create a document vector store for the state-of-the-union documents
docsearch = Chroma.from_documents(texts, embeddings, collection_name="state-of-union")

# create a retrieval-based question-answering model for the state-of-the-union documents
state_of_union = RetrievalQA.from_chain_type(llm=llm, chain_type="stuff", retriever=docsearch.as_retriever())

# load the web page
loader = WebBaseLoader("https://beta.ruff.rs/docs/faq/")
docs = loader.load()

# split the web page into smaller chunks
ruff_texts = text_splitter.split_documents(docs)

# create a document vector store for the Ruff web page
ruff_db = Chroma.from_documents(ruff_texts, embeddings, collection_name="ruff")

# create a retrieval-based question-answering model for the Ruff web page
ruff = RetrievalQA.from_chain_type(llm=llm, chain_type="stuff", retriever=ruff_db.as_retriever())
