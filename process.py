import argparse
import json
import subprocess


REQUIRED_INSTRUCTION_FIELDS = {"container_id"}


class InputError(Exception):
    def __init__(self, msg: str) -> None:
        self.msg = msg


def validate_instructions(instructions: dict) -> None:
    for key in REQUIRED_INSTRUCTION_FIELDS:
        if key not in instructions:
            raise InputError(f"Missing required input field {key}")


def build_docker_cmd(instructions: dict) -> str:
    return 'docker run hello_world'


def main(instructions: dict) -> None:
    validate_instructions(instructions)
    docker_cmd = build_docker_cmd(instructions)
    result = subprocess.run(docker_cmd, capture_output=True, text=True)
    import pdb; pdb.set_trace()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("instructions_json", help="JSON string with job input params")
    args = parser.parse_args()
    instructions_json = json.loads(args.instructions_json)
    main(instructions_json)
