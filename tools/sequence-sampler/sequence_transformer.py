from transformers import AutoTokenizer, EsmForMaskedLM
import torch
import matplotlib.pyplot as plt
import numpy as np
import os
import argparse
import logging

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class ESM2Runner:
    def __init__(self, model_name="facebook/esm2_t6_8M_UR50D"):
        verbosity = os.getenv('VERBOSITY', 'WARNING').upper()
        logger.setLevel(verbosity)

        self.tokenizer = AutoTokenizer.from_pretrained(model_name)
        self.model = EsmForMaskedLM.from_pretrained(model_name)

        # Check if GPU is available and move the model to GPU
        self.device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        self.model.to(self.device)

        self.amino_acids = list("LAGVSERTIDPKQNFYMHWC")
        
        logger.debug(f"ESM2Runner initialized with model {model_name} and verbosity {verbosity}")

    def token_masked_marginal_log_likelihood_matrix(self, protein_sequence, start_pos=1, end_pos=None):

        # print('gpu', torch.cuda.is_available())

        if end_pos is None:
            end_pos = len(protein_sequence)

        log_likelihoods = np.zeros((20, end_pos - start_pos + 1))
        # print('log_likelihoods shape', log_likelihoods.shape)

        # print('protein seq length', len(protein_sequence))
        # print('protein seq', protein_sequence)
        input_ids = self.tokenizer.encode(protein_sequence, return_tensors="pt").to(self.device)
        # print('input_ids', input_ids)
        # print('input_ids shape', input_ids.shape)

        for position in range(start_pos, end_pos + 1):
            # # position_index = position + 1 # part of original version
            # print('position', position)

            masked_input_ids = input_ids.clone()
            # # masked_input_ids[0, position_index] = self.tokenizer.mask_token_id # part of original version
            # print('before masking masked_input_ids[0, position]', masked_input_ids[0, position].item())
            masked_input_ids[0, position] = self.tokenizer.mask_token_id
            # print('after masking masked_input_ids[0, position]', masked_input_ids[0, position].item())
            # print('masked_input_ids shape', masked_input_ids.shape)
            # print('after masking masked_input_ids', masked_input_ids[0,:])
            # print('after masking masked_input_ids full', masked_input_ids)

            with torch.no_grad():
                logits = self.model(masked_input_ids).logits

            # # probabilities = torch.nn.functional.softmax(logits[0, position_index], dim=0) # part of original version
            # print('logits', logits)
            # print('logits[0, position]', logits[0, position])
            probabilities = torch.nn.functional.softmax(logits[0, position], dim=0)
            # print('probabilities', probabilities)
            # print('probabilities checksum', probabilities.sum())
            log_probabilities = torch.log(probabilities)

            for i, amino_acid in enumerate(self.amino_acids):
                # print('i', i)
                # print('position - start_pos', position - start_pos)
                aa_token_id = self.tokenizer.convert_tokens_to_ids(amino_acid)
                log_likelihoods[i, position - start_pos] = log_probabilities[aa_token_id].item()
        # print('')

        return log_likelihoods

    def _compute_likelihood_ratio_and_pseudolikelihood_matrix(self, protein_sequence, start_pos = 1, end_pos = None):
        """
        the entries are the log likelihood for every token conditioned on all others subtracted from the wt identity;
        positive values indicate that a substitution is increasing the likelihood, negative values indicate that a substitution is decreasing the likelihood
        Initialize matrix for log likelihood ratios
        """
        # calculate log likelihoods
        log_likelihoods = self.token_masked_marginal_log_likelihood_matrix(protein_sequence)
        # initiate the ratio matrix
        log_likelihood_ratios = np.zeros(log_likelihoods.shape)
        pseudolikelihood = []

        for aa in self.amino_acids:
            encoded_aa = self.tokenizer.encode(aa, add_special_tokens=False)[0] -4
            logger.debug(f"Amino acid: {aa}, Token ID: {encoded_aa}")

        logger.debug("Starting likelihood ratio determination")
        for position in range(start_pos, start_pos + log_likelihoods.shape[1]):
            # Get the log probability of the wild-type residue
            if position != 0:
                logger.debug("Moving to next position")
            logger.debug(f"Reference residue identity: {protein_sequence[position - 1]}, WT index: {self.tokenizer.encode(protein_sequence[position - 1], add_special_tokens=False)}")
            reference_residue = self.tokenizer.encode(protein_sequence[position - 1], add_special_tokens=False)[0] - 4
            logger.debug(reference_residue)
            log_prob_wt = log_likelihoods[reference_residue, position - start_pos]
            pseudolikelihood.append(log_prob_wt)
            logger.debug(f"Position: {position}, reference_residue token ID: {reference_residue}, log_likelihoods shape: {log_likelihoods.shape}, log_likelihood_wt: {log_prob_wt}")
            
            
            # Calculate LLR for each variant
            for i in range(log_likelihoods.shape[0]):
                logger.debug(f"Residue: {i}, Position: {position - start_pos}")
                log_likelihood_ratio = log_likelihoods[i, position - start_pos] - log_prob_wt
                logger.debug(f"Log Likelihood Ratio: {log_likelihood_ratio}")
                log_likelihood_ratios[i, position - start_pos] = log_likelihood_ratio
        
        logger.debug(f"Pseudolikelihood: {pseudolikelihood}")
        return log_likelihood_ratios, pseudolikelihood

    def token_masked_marginal_log_likelihood_ratio_matrix(self, protein_sequence):
        log_likelihood_ratios, pseudolikelihood = self._compute_likelihood_ratio_and_pseudolikelihood_matrix(protein_sequence)
        return log_likelihood_ratios

    def sequence_pseudo_log_likelihoods_scalar(self, protein_sequence):
        log_likelihood_ratios, pseudolikelihood = self._compute_likelihood_ratio_and_pseudolikelihood_matrix(protein_sequence)
        return sum(pseudolikelihood)

    def sequence_average_log_likelihood_scalar(self, protein_sequence):
        input_ids = self.tokenizer.encode(protein_sequence, return_tensors="pt").to(self.device)
        labels = input_ids.clone()  # The labels are the input ids themselves

        masked_input_ids = input_ids.clone()
        masked_input_ids[1:-1] = self.tokenizer.mask_token_id  # Mask all tokens except special tokens

        with torch.no_grad():
            outputs = self.model(masked_input_ids, labels=labels)
            loss = outputs.loss

        average_log_likelihood = -loss.item()
        return average_log_likelihood

    def sequence_scaled_average_log_likelihood_scalar(self, protein_sequence):
        average_log_likelihood = self.sequence_average_log_likelihood_scalar(protein_sequence)
        logger.debug(f"Sequence Length: {len(protein_sequence)}")
        scaled_average_log_likelihood = average_log_likelihood * len(protein_sequence)

        return scaled_average_log_likelihood