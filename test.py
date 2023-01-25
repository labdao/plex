import unittest

from process import InputError, format_args, validate_instructions


class TestCalculations(unittest.TestCase):
    def test_format_args(self):
        input_args = {"gpus": "all", "inference_steps": 15}

        # it works with single hyphen prefix
        self.assertEqual(" -gpus all -inference_steps 15", format_args(input_args, "-"))

        # it works with double hyphen prefix
        self.assertEqual(
            " --gpus all --inference_steps 15", format_args(input_args, "--")
        )

    def test_validate_instructions(self):
        invalid_instructions = {"de": "sci"}
        self.assertRaises(InputError, validate_instructions, invalid_instructions)


if __name__ == "__main__":
    unittest.main()
