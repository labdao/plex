import typer
import os
import pandas as pd
import numpy as np
from transformers import AutoTokenizer, AutoModelForMaskedLM, AutoModelForCausalLM, AutoModel, pipeline
from Bio import SeqIO
from Bio.Seq import Seq
from Bio.SeqRecord import SeqRecord
import torch
import re
import json
import csv

app = typer.Typer()

def softmax(arr):
    return np.exp(arr) / np.sum(np.exp(arr), axis=0)

def generate_embedding(record, tokenizer, model, output_path):
    sequence = re.sub(r"[UZOB]", "X", str(record.seq))
    spaced_sequence = " ".join(sequence)
    print("embedding sequence: ", spaced_sequence)
    encoded_input = tokenizer(spaced_sequence, return_tensors='pt')
    model_output = model(**encoded_input)
    embeddings = model_output.last_hidden_state.detach().numpy()

    output_file = os.path.join(output_path, f"{record.id}_embedding.csv")
    np.savetxt(output_file, embeddings[0], delimiter=",")

    print(f"embedding saved to {output_file}")
    return(output_file)

def fill_mask(record, unmasker, output_path):
    sequence = str(record.seq)
    spaced_sequence = " ".join(sequence)
    spaced_masked_sequence = spaced_sequence.replace("X", "[MASK]")
    print("filling masked sequence: ", spaced_masked_sequence)
    predictions = unmasker(spaced_masked_sequence)
    
    output_file = os.path.join(output_path, f"{record.id}_filled_mask.json")
    with open(output_file, "w") as f:
        json.dump(predictions, f)

    print(f"filled mask saved to {output_file}")
    return(output_file)

def conditional_probability_matrix(record, tokenizer, masked_model, output_path):
    sequence = str(record.seq)
    print("generating scoring matrix for sequence: ", sequence)

    all_token_scores = []
    for idx in range(len(sequence)):
        x_sequence = sequence[:idx] + "X" + sequence[idx+1:]
        spaced_sequence = " ".join(x_sequence)
        spaced_masked_sequence = spaced_sequence.replace("X", "[MASK]")
        print("tokenising masked sequence: ", spaced_masked_sequence)
        encoded_input = tokenizer(spaced_masked_sequence, return_tensors='pt')

        with torch.no_grad():
            model_output = masked_model(**encoded_input)

        scores = model_output.logits
        print("scored sequence: ", scores)
        mask_position = torch.tensor([idx], dtype=torch.long)
        mask_scores = scores[0, mask_position, :]

        token_scores = torch.softmax(mask_scores, dim=-1).squeeze()

        token_score_dict = {}
        for token, token_score in zip(tokenizer.vocab.keys(), token_scores):
            token_score_dict[token] = token_score.item()

        all_token_scores.append(token_score_dict)

    # Write the scoring matrix to a CSV file
    output_file = os.path.join(output_path, f"{record.id}_conditional_probability_matrix.csv")
    with open(output_file, 'w', newline='') as csvfile:
        print("writing scoring matrix to " + output_file)
        fieldnames = ['position', 'identity'] + list(tokenizer.vocab.keys())
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)

        writer.writeheader()
        for i, token_scores in enumerate(all_token_scores):
            row = {'position': i+1, 'identity': sequence[i]}
            row.update(token_scores)
            writer.writerow(row)

    print(f"conditional probability matrix saved to {output_file}")
    return(output_file)

def joint_probability_score():
    #TODO
    print("coming soon")

def generate(record, tokenizer, generator_model, output_path, max_length=50, top_k=5):
    # prompt
    sequence = str(record.seq)
    spaced_sequence = " ".join(sequence)
    spaced_masked_sequence = spaced_sequence.replace("X", "[MASK]")
    print("prompt sequence: ", spaced_masked_sequence)

    # generate
    generated_sequences = generator_model(spaced_masked_sequence, max_length=max_length, top_k=top_k)
    seq_records = []
    for i, seq in enumerate(generated_sequences):
        masked_sequence = seq['generated_text'].replace(' ', '')
        generated_sequence = masked_sequence.replace("[MASK]", 'X')
        seq_record = SeqRecord(Seq(generated_sequence), id=f"sequence_{i+1}", description="")
        seq_records.append(seq_record)
    output_file = os.path.join(output_path, "generated_sequences.fasta")
    with open(output_file, "w") as output_handle:
        SeqIO.write(seq_records, output_handle, "fasta")
    print(f"Generated sequences saved to {output_file}")

@app.command()
def main(
        input: str = typer.Argument(..., help="Path to the input fasta file."),
        output_path: str = typer.Argument(..., help="Path to the output directory."),
        huggingface_model_name: str = typer.Argument("Rostlab/prot_bert_bfd", help="The model name to use. Supply a Hugginface identifier. Default is 'Rostlab/prot_bert'."),
        mode: str = typer.Option(
            ...,
            help="Mode of operation. Choose from 'embedding', 'fill-mask', 'conditional-probability', 'joint-probability', or 'generate'.",
            prompt="Select mode of operation",
            case_sensitive=False,
            show_choices=True,
            show_default=True,
            metavar="MODE",
            callback=lambda value: value.lower(),
            autocompletion=lambda: ["embedding", "fill-mask", "generate", "conditional-probability", "joint-probability"],
        ),
    ):
    # Create output directory if it doesn't exist
    if not os.path.exists(output_path):
        os.makedirs(output_path)

    # Provided model
    print(f"Using model {huggingface_model_name}")
    
    # Load models and tokenizer
    tokenizer = AutoTokenizer.from_pretrained(huggingface_model_name, do_lower_case=False)
    masked_model = AutoModelForMaskedLM.from_pretrained(huggingface_model_name)
    generator_model = AutoModelForCausalLM.from_pretrained(huggingface_model_name)
    plain_model = AutoModel.from_pretrained(huggingface_model_name)
    unmasker = pipeline('fill-mask', model=masked_model, tokenizer=tokenizer)
    generator = pipeline('text-generation', model=generator_model, tokenizer=tokenizer)

    # Load input sequence
    record = SeqIO.read(input, "fasta")

    if mode == "embedding":
        generate_embedding(record, tokenizer, plain_model, output_path)
    elif mode == "fill-mask":
        fill_mask(record, unmasker, output_path)
    elif mode == "conditional-probability":
        conditional_probability_matrix(record, tokenizer, masked_model, output_path)
    elif mode == "joint-probability":
        joint_probability_score()
    elif mode == "generate":
        generate(record, tokenizer, generator, output_path)
    else:
        typer.echo("Invalid mode. Please choose from 'embedding', 'fill-mask', 'scoring-matrix', 'top-k', or 'sample-n'.")

if __name__ == "__main__":
    typer.run(main)
