import json


def generate_diffdock_instructions(
    debug_logs=True,
    protein_path="test/test.pdb",
    out_file="data/prepared_for_esm.fasta",
    repr_layers=33,
    ligand="test/test.sdf",
    inference_steps=20,
    samples_per_complex=40,
    batch_size=10,
    actual_steps=18,
) -> dict:
    instructions = {
        "container_id": "ghcr.io/labdao/diffdock:main",
        "debug_logs": debug_logs,
        "short_args": {"v": "/home/ubuntu/diffdock:/diffdock"},
        "long_args": {"gpus": "all"},
        "cmd": (
            '/bin/bash -c "python datasets/esm_embedding_preparation.py --protein_path'
            f" {protein_path} --out_file {out_file} && HOME=esm/model_weights python"
            " esm/scripts/extract.py esm2_t33_650M_UR50D data/prepared_for_esm.fasta"
            f" data/esm2_output --repr_layers {repr_layers} --include per_tok && python"
            f" -m inference --protein_path {protein_path} --ligand {ligand} --out_dir"
            f" /outputs --inference_steps {inference_steps} --samples_per_complex"
            f" {samples_per_complex} --batch_size {batch_size} --actual_steps"
            f' {actual_steps} --no_final_step_noise"'
        ),
    }
    return json.dumps(instructions)

#TODO: #26 introduce guidance on volume mounting, especially mounting multiple volumes
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

if __name__ == "__main__":
    print(generate_vina_instructions())
