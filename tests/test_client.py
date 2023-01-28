import json
import unittest

from client import generate_diffdock_instructions


class TestClient(unittest.TestCase):
    def test_format_args(self):
        self.maxDiff = None
        expected_instructions = json.dumps(
            {
                "container_id": "ghcr.io/labdao/diffdock:main",
                "debug_logs": True,
                "short_args": {},
                "long_args": {"gpus": "all"},
                "cmd": (
                    '/bin/bash -c "python datasets/esm_embedding_preparation.py '
                    '--protein_path inputs/test.pdb --out_file outputs/prepared_for_esm.fasta '
                    '&& HOME=esm/model_weights python esm/scripts/extract.py esm2_t33_650M_UR50D '
                    'outputs/prepared_for_esm.fasta outputs/esm2_output --repr_layers 33 --include per_tok && '
                    'python -m inference --protein_path inputs/test.pdb --ligand inputs/test.sdf --out_dir outputs '
                    '--inference_steps 20 --samples_per_complex 40 --batch_size 10 --actual_steps 18 '
                    '--esm_embeddings_path outputs/esm2_output --no_final_step_noise"'
                ),
            }
        )
        self.assertEqual(generate_diffdock_instructions(), expected_instructions)
