import asyncio
import json
import subprocess

from subprocess import PIPE, STDOUT, Popen
from websockets.server import WebSocketServerProtocol, serve

from ipfs import DEFAULT_DATA_CIDS, load_cids_to_inputs
from process import validate_instructions, build_docker_cmd


async def isomorphic_print(ws: WebSocketServerProtocol, msg: str) -> None:
    print(msg, end="")  # `end=''` removes extra newline
    await ws.send(msg)


async def run_docker_example(ws: WebSocketServerProtocol) -> None:
    await isomorphic_print(ws, "A new request was received\n")
    await isomorphic_print(ws, "About to run some containerized code\n")
    with Popen(
        docker_cmd,
        shell=True,
        stdout=PIPE,
        stderr=STDOUT,
        bufsize=1,
        universal_newlines=True,
    ) as proc:
        for line in proc.stdout:
            await isomorphic_print(ws, line)
    await isomorphic_print(ws, "Finished running containerized code\n")


async def process_instructions(ws: WebSocketServerProtocol) -> None:
    await isomorphic_print(ws, "A new socket was opened\n")
    await isomorphic_print(ws, "Waiting for instructions...\n")
    instructions = await ws.recv()
    await isomorphic_print(ws, f"Validating instructions: {instructions}")
    instructions = json.loads(instructions)
    validate_instructions(instructions)  # raises error if invalid
    await isomorphic_print(ws, f"Instructions are valid")
    docker_cmd = build_docker_cmd(instructions)
    await isomorphic_print(ws, f"About to run {docker_cmd}\n")
    with Popen(
        docker_cmd,
        shell=True,
        stdout=PIPE,
        stderr=STDOUT,
        bufsize=1,
        universal_newlines=True,
    ) as proc:
        for line in proc.stdout:
            await isomorphic_print(ws, line)
    await isomorphic_print(ws, "Finished running containerized code\n")


async def run_server() -> None:
    async with serve(process_instructions, "localhost", 8765):
        print('Server is now accepting requests')
        await asyncio.Future()


if __name__ == "__main__":
    print('Loading default data')
    load_cids_to_inputs(DEFAULT_DATA_CIDS)
    asyncio.run(run_server())
