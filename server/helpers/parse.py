import argparse
import json
import os
import subprocess

from helpers.settings import PROJECT_ROOT


REQUIRED_INSTRUCTION_FIELDS = ["container_id"]
ALLOWED_INSTRUCTION_FIELDS = REQUIRED_INSTRUCTION_FIELDS + [
    "short_args",
    "long_args",
    "cmd",
    "debug_logs",
]


class InputError(Exception):
    def __init__(self, msg: str) -> None:
        self.msg = msg


def validate_instructions(instructions: dict) -> None:
    for key in REQUIRED_INSTRUCTION_FIELDS:
        if key not in instructions:
            raise InputError(f"Missing required input field {key}")

    invalid_fields = [
        key for key in instructions.keys() if key not in ALLOWED_INSTRUCTION_FIELDS
    ]
    if invalid_fields:
        raise InputError(
            f"The following fields are not allowed {','.join(invalid_fields)}"
        )


def format_args(instruction_args: dict, prefix: str) -> str:
    arg_flags = ""
    for key, val in instruction_args.items():
        arg_flags += f" {prefix}{key} {val}"
    return arg_flags


def build_docker_cmd(instructions: dict) -> str:
    input_vol = f'{os.path.join(PROJECT_ROOT, "inputs")}:/inputs'
    output_vol = f'{os.path.join(PROJECT_ROOT, "outputs")}:/outputs'

    return (
        f'docker run -v {input_vol} -v {output_vol}'
        f'{format_args(instructions.get("long_args", {}), "--")}'
        f'{format_args(instructions.get("short_args", {}), "-")}'
        f' {instructions["container_id"]} {instructions["cmd"]}'
    )
