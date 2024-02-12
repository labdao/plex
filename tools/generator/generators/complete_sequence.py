import os
import numpy as np
import logging
from .base_generator import BaseGenerator
import sequence_transformer

class complete_sequence_Generator(BaseGenerator):

    def generate(self, args):

        sequence = args.sequence
        permissibility_vector = args.permissibility_vector
        masked_sequence = sequence
        outputs_directory = args.outputs_directory
        generator_name = args.generator_name

        generator_directory = os.path.join(outputs_directory, generator_name)
        if not os.path.exists(generator_directory):
            os.makedirs(generator_directory, exist_ok=True)

        logging.info(f"Running {generator_name}")

        runner = sequence_transformer.ESM2Runner()

        predicted_sequence = runner.predict_masked_sequence(masked_sequence)

        logging.info(f"Original sequence: {masked_sequence}")
        logging.info(f"Predicted sequence: {predicted_sequence}")

        return ''.join(predicted_sequence), ''.join(permissibility_vector)