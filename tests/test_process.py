import unittest

from process import InputError, build_docker_cmd, format_args, validate_instructions


class TestProcess(unittest.TestCase):
    def test_format_args(self):
        input_args = {"gpus": "all", "inference_steps": 15}

        # it works with single hyphen prefix
        self.assertEqual(" -gpus all -inference_steps 15", format_args(input_args, "-"))

        # it works with double hyphen prefix
        self.assertEqual(
            " --gpus all --inference_steps 15", format_args(input_args, "--")
        )

    def test_build_docker_cmd(self):
        self.maxDiff = None
        instructions = {
            "container_id": "compbio:latest",
            "short_args": {'p': 5000},
            "long_args": {"gpus": "all"},
            "cmd": (
                'python design_drug.py'
            ),
        }
        expected_output = (
            "docker run -v /home/ubuntu/inputs:/inputs -v /home/ubuntu/outputs:/outputs "
            "--gpus all -p 5000 compbio:latest python design_drug.py"
        )
        self.assertEqual(expected_output, build_docker_cmd(instructions))

    def test_validate_instructions(self):
        missing_field_instructions = {"cmd": "echo DeSci"}

        # it raises error when required field is missing
        self.assertRaises(InputError, validate_instructions, missing_field_instructions)

        # it raises error when extra field is required
        extra_field_instructions = {"container_id": "hello-world", "de": "sci"}
        self.assertRaises(InputError, validate_instructions, extra_field_instructions)

        # it raises no error for valid input
        valid_instructions = {"container_id": "hello-world", "cmd": "echo DeSci"}
        self.assertIsNone(validate_instructions(valid_instructions))


if __name__ == "__main__":
    unittest.main()
