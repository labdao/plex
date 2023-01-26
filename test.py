import unittest

from process import InputError, build_docker_cmd, format_args, validate_instructions


class TestCalculations(unittest.TestCase):
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
            "container_id": "ghcr.io/labdao/diffdock:main",
            "short_args": {"v": "/home/ubuntu/diffdock:/diffdock"},
            "long_args": {"gpus": "all"},
            "cmd": (
                "python datasets/esm_embedding_preparation.py --protein_path"
                " test/test.pdb --out_file data/prepared_for_esm.fasta &&"
                " HOME=esm/model_weights python esm/scripts/extract.py"
                " esm2_t33_650M_UR50D data/prepared_for_esm.fasta data/esm2_output"
                " --repr_layers 33 --include per_tok && python -m inference"
                " --protein_path test/test.pdb --ligand test/test.sdf --out_dir"
                " /outputs --inference_steps 20 --samples_per_complex 40 --batch_size"
                " 10 --actual_steps 18 --no_final_step_noise"
            ),
        }
        expected_output = (
            "docker run --gpus all -v /home/ubuntu/diffdock:/diffdock"
            ' ghcr.io/labdao/diffdock:main /bin/bash -c "python'
            " datasets/esm_embedding_preparation.py --protein_path test/test.pdb"
            " --out_file data/prepared_for_esm.fasta && HOME=esm/model_weights python"
            " esm/scripts/extract.py esm2_t33_650M_UR50D data/prepared_for_esm.fasta"
            " data/esm2_output --repr_layers 33 --include per_tok && python -m"
            " inference --protein_path test/test.pdb --ligand test/test.sdf --out_dir"
            " /outputs --inference_steps 20 --samples_per_complex 40 --batch_size 10"
            ' --actual_steps 18 --no_final_step_noise"'
        )
        self.assertEqual(expected_output, build_docker_cmd(instructions))

    def test_validate_instructions(self):
        invalid_instructions = {"de": "sci"}
        self.assertRaises(InputError, validate_instructions, invalid_instructions)


if __name__ == "__main__":
    unittest.main()
