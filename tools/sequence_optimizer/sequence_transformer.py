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

        # print('to gpu or not to gpu, that is the question', torch.cuda.is_available())

        if end_pos is None:
            end_pos = len(protein_sequence)

        log_likelihoods = np.zeros((20, end_pos - start_pos + 1))

        input_ids = self.tokenizer.encode(protein_sequence, return_tensors="pt").to(self.device)
        sequence_length = input_ids.shape[1] - 2

        for position in range(start_pos, end_pos + 1):
            position_index = position + 1

            masked_input_ids = input_ids.clone()
            masked_input_ids[0, position_index] = self.tokenizer.mask_token_id

            with torch.no_grad():
                logits = self.model(masked_input_ids).logits

            probabilities = torch.nn.functional.softmax(logits[0, position], dim=0)
            log_probabilities = torch.log(probabilities)

            for i, amino_acid in enumerate(self.amino_acids):
                aa_token_id = self.tokenizer.convert_tokens_to_ids(amino_acid)
                log_likelihoods[i, position - start_pos] = log_probabilities[aa_token_id].item()

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
            reference_residue = self.tokenizer.encode(protein_sequence[position - 1], add_special_tokens=False)[0] -4
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


### old version of ESM2 Runner class without GPU support ###

# from transformers import AutoTokenizer, EsmForMaskedLM
# import torch
# import matplotlib.pyplot as plt
# import numpy as np
# import os
# import argparse
# import logging

# # Set up logging
# logging.basicConfig(level=logging.INFO)
# logger = logging.getLogger(__name__)

# class ESM2Runner:
#     def __init__(self, model_name="facebook/esm2_t6_8M_UR50D"):
#         verbosity = os.getenv('VERBOSITY', 'WARNING').upper()
#         logger.setLevel(verbosity)

#         self.tokenizer = AutoTokenizer.from_pretrained(model_name)
#         self.model = EsmForMaskedLM.from_pretrained(model_name)
#         self.amino_acids = list("LAGVSERTIDPKQNFYMHWC")
        
#         logger.debug(f"ESM2Runner initialized with model {model_name} and verbosity {verbosity}")

#     def token_masked_marginal_log_likelihood_matrix(self, protein_sequence, start_pos = 1, end_pos = None):
#         """
#         # return matrix with n x m, where n is the number of tokens and m is the length of the sequence
#         # the entries are the log likelihood for every token conditioned on all others
#         # retrieve the amino_acids reference with self.amino_acids
#         """

#         print('to gpu or not to gpu, that is the question', torch.cuda.is_available())

#         # set end position
#         if end_pos is None:
#             end_pos = len(protein_sequence)
#         # Initialize matrix for log likelihoods
#         log_likelihoods = np.zeros((20, end_pos - start_pos + 1))

#         # Tokenize the input sequence, including special tokens
#         input_ids = self.tokenizer.encode(protein_sequence, return_tensors="pt")
#         sequence_length = input_ids.shape[1] - 2  # Adjust for special tokens

#         # Calculate log likelihoods for each position and amino acid
#         for position in range(start_pos, end_pos + 1):
#             # Adjust position for special tokens
#             position_index = position + 1  # Assuming one special token at the start

#             # Mask the target position
#             masked_input_ids = input_ids.clone()
#             masked_input_ids[0, position_index] = self.tokenizer.mask_token_id

#             # Get logits for the masked token
#             with torch.no_grad():
#                 logits = self.model(masked_input_ids).logits
            
#             # Calculate log probabilities
#             probabilities = torch.nn.functional.softmax(logits[0, position], dim=0)
#             log_probabilities = torch.log(probabilities)
            
#             # Store log probabilities for each amino acid
#             for i, amino_acid in enumerate(self.amino_acids):
#                 aa_token_id = self.tokenizer.convert_tokens_to_ids(amino_acid)
#                 log_likelihoods[i, position - start_pos] = log_probabilities[aa_token_id].item()
        
#         return log_likelihoods
    
#     def _compute_likelihood_ratio_and_pseudolikelihood_matrix(self, protein_sequence, start_pos = 1, end_pos = None):
#         """
#         the entries are the log likelihood for every token conditioned on all others subtracted from the wt identity;
#         positive values indicate that a substitution is increasing the likelihood, negative values indicate that a substitution is decreasing the likelihood
#         Initialize matrix for log likelihood ratios
#         """
#         # calculate log likelihoods
#         log_likelihoods = self.token_masked_marginal_log_likelihood_matrix(protein_sequence)
#         # initiate the ratio matrix
#         log_likelihood_ratios = np.zeros(log_likelihoods.shape)
#         pseudolikelihood = []

#         for aa in self.amino_acids:
#             encoded_aa = self.tokenizer.encode(aa, add_special_tokens=False)[0] -4
#             logger.debug(f"Amino acid: {aa}, Token ID: {encoded_aa}")

#         logger.debug("Starting likelihood ratio determination")
#         for position in range(start_pos, start_pos + log_likelihoods.shape[1]):
#             # Get the log probability of the wild-type residue
#             if position != 0:
#                 logger.debug("Moving to next position")
#             logger.debug(f"Reference residue identity: {protein_sequence[position - 1]}, WT index: {self.tokenizer.encode(protein_sequence[position - 1], add_special_tokens=False)}")
#             reference_residue = self.tokenizer.encode(protein_sequence[position - 1], add_special_tokens=False)[0] -4
#             logger.debug(reference_residue)
#             log_prob_wt = log_likelihoods[reference_residue, position - start_pos]
#             pseudolikelihood.append(log_prob_wt)
#             logger.debug(f"Position: {position}, reference_residue token ID: {reference_residue}, log_likelihoods shape: {log_likelihoods.shape}, log_likelihood_wt: {log_prob_wt}")
            
            
#             # Calculate LLR for each variant
#             for i in range(log_likelihoods.shape[0]):
#                 logger.debug(f"Residue: {i}, Position: {position - start_pos}")
#                 log_likelihood_ratio = log_likelihoods[i, position - start_pos] - log_prob_wt
#                 logger.debug(f"Log Likelihood Ratio: {log_likelihood_ratio}")
#                 log_likelihood_ratios[i, position - start_pos] = log_likelihood_ratio
        
#         logger.debug(f"Pseudolikelihood: {pseudolikelihood}")
#         return log_likelihood_ratios, pseudolikelihood

#     def token_masked_marginal_log_likelihood_ratio_matrix(self, protein_sequence):
#         log_likelihood_ratios, pseudolikelihood = self._compute_likelihood_ratio_and_pseudolikelihood_matrix(protein_sequence)
#         return log_likelihood_ratios

#     def sequence_pseudo_log_likelihoods_scalar(self, protein_sequence):
#         log_likelihood_ratios, pseudolikelihood = self._compute_likelihood_ratio_and_pseudolikelihood_matrix(protein_sequence)
#         return sum(pseudolikelihood)

#     def sequence_average_log_likelihood_scalar(self, protein_sequence):
#         """
#         returns a scalar
#         loss * number of residues
#         """
#         # Tokenize the input sequence, including special tokens
#         input_ids = self.tokenizer.encode(protein_sequence, return_tensors="pt")
#         labels = input_ids.clone()  # The labels are the input ids themselves

#         # Mask the input ids, the model will try to predict the original ids
#         masked_input_ids = input_ids.clone()
#         masked_input_ids[1:-1] = self.tokenizer.mask_token_id  # Mask all tokens except special tokens

#         # Get the logits from the model
#         with torch.no_grad():
#             outputs = self.model(masked_input_ids, labels=labels)
#             loss = outputs.loss  # The model outputs the loss directly

#         # Convert the loss to the average log likelihood (negative of loss)
#         average_log_likelihood = -loss.item()

#         return average_log_likelihood

#     def sequence_scaled_average_log_likelihood_scalar(self, protein_sequence):
#         """
#         returns a scalar
#         loss * number of residues
#         """
#         # estimate the average log likelihood of the sequence
#         average_log_likelihood = self.sequence_average_log_likelihood(protein_sequence)
#         logger.debug(f"Sequence Length: {len(protein_sequence)}")
#         scaled_average_log_likelihood = average_log_likelihood * len(protein_sequence)

#         return scaled_average_log_likelihood

# if __name__ == "__main__":
#     # Example protein sequence
#     example_sequence = "MKTLLVLLPLVSSQCVNLTTRTQLPPAYTNSFTRGVYYPDKVFRSSVLHSTQDLFLPFFSNVTWFHAIHVSGTNGTKRFDNPVLPFNDGVYFASTEKSNIIRGWIFGTTLDSKTQSLLIVNNATNVVIKVCEFQFCNDPFLGVYYHKNNKSWMESEFRVYSSANNCTFEYVSQPFLMDLEGKQGNFKNLREFVFKNIDGYFKIYSKHTPINLVRDLPQGFSALEPLVDLPIGINITRFQTLLALHRSYLTPGDSSSGWTAGAAAYYVGYLQPRTFLLKYNENGTITDAVDCSQNPLAELKCSVKSFEIDKGIYQTSNQVVDC"

#     # Initialize the ESM2Runner with the default model
#     runner = ESM2Runner()

#     # Calculate the average log likelihood for the example sequence
#     #matrix = runner.token_masked_marginal_log_likelihood_ratio_matrix(example_sequence)
#     scalar = runner.sequence_pseudo_log_likelihoods_scalar(example_sequence)
#     # Print the result
#     #print(f"Log Likelihood Matrix for the example sequence: {matrix}")
#     print(f"Pseudolikelihood Scalar for the example sequence: {scalar}")