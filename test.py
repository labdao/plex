import unittest

from process import InputError, validate_instructions

class TestCalculations(unittest.TestCase):

    def test_validate_instructions(self):
        invalid_instructions = {"de": "sci"}
        self.assertRaises(InputError, validate_instructions, invalid_instructions)

if __name__ == '__main__':
    unittest.main()
