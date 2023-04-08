from langchain.agents import load_tools, initialize_agent, Tool
from langchain.agents import AgentType
from langchain.tools import BaseTool
from langchain.chat_models import ChatOpenAI
from langchain.llms import OpenAI
from langchain.utilities import BashProcess
from langchain import LLMMathChain, SerpAPIWrapper

import typer
from vector import state_of_union, ruff

app = typer.Typer()

manual_tools = [
    Tool(
        name = "State of Union QA System",
        func=state_of_union.run,
        description="useful for when you need to answer questions about the most recent state of the union address. Input should be a fully formed question."
    ),
    Tool(
        name = "Ruff QA System",
        func=ruff.run,
        description="useful for when you need to answer questions about ruff (a python linter). Input should be a fully formed question."
    ),
]

def main(question: str):
    chat = ChatOpenAI(temperature=0)
    llm = OpenAI(temperature=0)
    tools = load_tools(["serpapi", "llm-math", "human"], llm=llm)
    all_tools = manual_tools + tools
    agent = initialize_agent(all_tools, chat, agent=AgentType.CHAT_ZERO_SHOT_REACT_DESCRIPTION, verbose=True)
    response = agent.run(question)
    typer.echo(response)

if __name__ == "__main__":
    typer.run(main)