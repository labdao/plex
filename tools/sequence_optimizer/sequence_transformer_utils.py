from transformers import AutoTokenizer, EsmForMaskedLM
import torch
import matplotlib.pyplot as plt
import numpy as np
import os
import argparse

def generate_log_likelihoods_matrix(protein_sequence, start_pos, end_pos, model, tokenizer, amino_acids):
    # Initialize matrix for log likelihoods
    log_likelihoods = np.zeros((20, end_pos - start_pos + 1))

    # Tokenize the input sequence, including special tokens
    input_ids = tokenizer.encode(protein_sequence, return_tensors="pt")
    sequence_length = input_ids.shape[1] - 2  # Adjust for special tokens

    # Calculate log likelihoods for each position and amino acid
    for position in range(start_pos, end_pos + 1):
        # Adjust position for special tokens
        position_index = position + 1  # Assuming one special token at the start

        # Mask the target position
        masked_input_ids = input_ids.clone()
        masked_input_ids[0, position_index] = tokenizer.mask_token_id

        # Get logits for the masked token
        with torch.no_grad():
            logits = model(masked_input_ids).logits
        
        # Calculate log probabilities
        probabilities = torch.nn.functional.softmax(logits[0, position_index], dim=0)
        log_probabilities = torch.log(probabilities)
        
        # Store log probabilities for each amino acid
        for i, amino_acid in enumerate(amino_acids):
            aa_token_id = tokenizer.convert_tokens_to_ids(amino_acid)
            log_likelihoods[i, position - start_pos] = log_probabilities[aa_token_id].item()
    
    return log_likelihoods

def generate_log_likelihood_ratios_matrix(log_likelihoods, protein_sequence, tokenizer, start_pos = 1):
    # Initialize matrix for log likelihood ratios
    log_likelihood_ratios = np.zeros(log_likelihoods.shape)
    
    # Calculate LLR for each variant
    print(range(start_pos, start_pos + log_likelihoods.shape[1]))

    for aa in amino_acids:
        encoded_aa = tokenizer.encode(aa, add_special_tokens=False)[0] -4
        print(f"Amino acid: {aa}, Token ID: {encoded_aa}")

    print("Starting likelihood ratio determination")
    for position in range(start_pos, start_pos + log_likelihoods.shape[1]):
        # Get the log probability of the wild-type residue
        if position != 0:
            print("Moving to next position")
        print(f"WT identity: {protein_sequence[position - 1]}, WT index: {tokenizer.encode(protein_sequence[position - 1], add_special_tokens=False)}")
        wt_residue = tokenizer.encode(protein_sequence[position - 1], add_special_tokens=False)[0] -4
        print(wt_residue)
        log_prob_wt = log_likelihoods[wt_residue, position - start_pos]
        print(f"Position: {position}, wt_residue token ID: {wt_residue}, log_likelihoods shape: {log_likelihoods.shape}, log_likelihood_wt: {log_prob_wt}")
        
        
        # Calculate LLR for each variant
        for i in range(log_likelihoods.shape[0]):
            print(f"Residue: {i}, Position: {position - start_pos}")
            log_likelihood_ratio = log_likelihoods[i, position - start_pos] - log_prob_wt
            print(f"Log Likelihood Ratio: {log_likelihood_ratio}")
            log_likelihood_ratios[i, position - start_pos] = log_likelihood_ratio
    
    return log_likelihood_ratios

def plot_heatmap(matrix, protein_sequence, start_pos, end_pos, amino_acids, title, label, filename = "heatmap.png"):
    plt.figure(figsize=(15, 5))
    plt.imshow(matrix, cmap="viridis_r", aspect="auto")
    plt.xticks(range(end_pos - start_pos + 1), list(protein_sequence[start_pos-1:end_pos]))
    plt.yticks(range(20), amino_acids)
    plt.xlabel("Position in Protein Sequence")
    plt.ylabel("Amino Acid Mutations")
    plt.title(title)
    plt.colorbar(label=label)
    plt.show()
    # Save the plot to a temporary file and return the file path
    plt.savefig(filename)
    plt.close()
    return filename

# Load the model and tokenizer outside the functions to avoid reloading them each time
model_name = "facebook/esm2_t6_8M_UR50D"
tokenizer = AutoTokenizer.from_pretrained(model_name)
model = EsmForMaskedLM.from_pretrained(model_name)

# Example usage
protein_sequence = "MSKGEELFTGVVPILVELDGDVNGHKFSVSGEGEGDATYGKLTLKFICTTGKLPVPWPTLVTTFSYGVQCFSRYPDHMKQHDFFKSAMPEGYVQERTIFFKDDGNYKTRAEVKFEGDTLVNRIELKGIDFKEDGNILGHKLEYNYNSHNVYIMADKQKNGIKVNFKIRHNIEDGSVQLADHYQQNTPIGDGPVLLPDNHYLSTQSALSKDPNEKRDHMVLLEFVTAAGITHGMDELYK"
start_pos = 1
end_pos = len(protein_sequence)  # or any other end position you want to analyze
amino_acids = list("LAGVSERTIDPKQNFYMHWC")

log_likelihoods = generate_log_likelihoods_matrix(protein_sequence, start_pos, end_pos, model, tokenizer, amino_acids)
log_likelihood_ratios = generate_log_likelihood_ratios_matrix(log_likelihoods, protein_sequence, tokenizer)
print(log_likelihood_ratios)

# Plot the log likelihoods heatmap
plot_heatmap(log_likelihoods, protein_sequence, start_pos, end_pos, amino_acids, "Log Likelihoods Heatmap", "Log Likelihood")
plot_heatmap(log_likelihood_ratios, protein_sequence, start_pos, end_pos, amino_acids, "Log Likelihood Ratios Heatmap", "Log Likelihood Ratio (LLR)")
