import argparse
import json
import subprocess


REQUIRED_INSTRUCTION_FIELDS = {"container_id", "short_args", "long_args", "cmd"}


class InputError(Exception):
    def __init__(self, msg: str) -> None:
        self.msg = msg


def validate_instructions(instructions: dict) -> None:
    for key in REQUIRED_INSTRUCTION_FIELDS:
        if key not in instructions:
            raise InputError(f"Missing required input field {key}")


def format_args(instruction_args: dict, prefix: str) -> str:
    arg_flags = ""
    for key, val in instruction_args.items():
        arg_flags += f" {prefix}{key} {val}"
    return arg_flags


def build_docker_cmd(instructions: dict) -> str:
    return (
        "docker"
        f' run{format_args(instructions["long_args"], "--")}'
        f'{format_args(instructions["short_args"], "-")}'
        f' {instructions["container_id"]} /bin/bash'
        f' -c "{instructions["cmd"]}"'
    )


def main(instructions: dict) -> None:
    validate_instructions(instructions)
    docker_cmd = build_docker_cmd(instructions)
    if instructions.get('debug_logs'):
        result = subprocess.run(docker_cmd, shell=True)
    else:
        result = subprocess.run(docker_cmd, capture_output=True, shell=True, text=True)
    print(result)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("instructions_json", help="JSON string with job input params")
    args = parser.parse_args()
    instructions_json = json.loads(args.instructions_json)
    main(instructions_json)
