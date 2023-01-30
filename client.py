import json


def generate_diffdock_instructions(
    debug_logs=True,
    protein_path="inputs/test.pdb",
    fasta_out_file="outputs/prepared_for_esm.fasta",
    ligand="inputs/test.sdf",
    repr_layers=33,
    out_dir="outputs",
    inference_steps=20,
    samples_per_complex=40,
    batch_size=10,
    actual_steps=18,
) -> dict:
    instructions = {
        "container_id": "ghcr.io/labdao/diffdock:main",
        "debug_logs": debug_logs,
        "short_args": {},
        "long_args": {"gpus": "all"},
        "cmd": (
            '/bin/bash -c \"'
            f'python datasets/esm_embedding_preparation.py --protein_path {protein_path} --out_file {fasta_out_file} '
            f'&& HOME=esm/model_weights python esm/scripts/extract.py esm2_t33_650M_UR50D {fasta_out_file} '
            f'outputs/esm2_output --repr_layers {repr_layers} --include per_tok && '
            f'python -m inference --protein_path {protein_path} --ligand {ligand} --out_dir {out_dir} '
            f'--inference_steps {inference_steps} --samples_per_complex {samples_per_complex} --batch_size {batch_size} '
            f'--actual_steps {actual_steps} --esm_embeddings_path outputs/esm2_output --no_final_step_noise\"'
        ),
    }
    return json.dumps(instructions)

def generate_vina_instructions(
    debug_logs=True,
    protein="1a30/1a30_protein.pdb",
    ligand="1a30/1a30_ligand.sdf",
    output="1a30/1a30_scored_vina.sdf.gz",
    cnn_scoring="none",
    modifier="score_only",
) -> dict:
    instructions = {
        "container_id": "gnina/gnina:latest",
        "debug_logs": debug_logs,
        "short_args": {"v": "/home/ubuntu/casf-2016:/inputs"},
        "long_args": {"gpus": 0},
        "cmd": (
            "gnina -r"
            f" /inputs/{protein} -l /inputs/{ligand} -o"
            f" /inputs/{output}"
            f" --autobox_ligand /inputs/{protein} --cnn_scoring {cnn_scoring} --exhaustiveness 64"
            f" --{modifier}"
        ),
    }
    return json.dumps(instructions)

def generate_gnina_instructions(
    debug_logs=True,
    protein="1a30/1a30_protein.pdb",
    ligand="1a30/1a30_ligand.sdf",
    output="1a30/1a30_scored_gnina.sdf.gz",
    cnn_scoring="rescore",
    modifier="", # --score_only for blind docking
) -> dict:
    instructions = {
        "container_id": "gnina/gnina:latest",
        "debug_logs": debug_logs,
        "short_args": {"v": "/home/ubuntu/casf-2016:/inputs"},
        "long_args": {"gpus": 0},
        "cmd": (
            "gnina -r"
            f" /inputs/{protein} -l /inputs/{ligand} -o"
            f" /inputs/{output}"
            f" --autobox_ligand /inputs/{protein} --cnn_scoring {cnn_scoring} --exhaustiveness 64"
            f" {modifier}"
        ),
    }
    return json.dumps(instructions)

if __name__ == "__main__":
    print(generate_vina_instructions())
