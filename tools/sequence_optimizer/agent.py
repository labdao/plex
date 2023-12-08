import pandas as pd
import random
import sys
import os

# Assuming your script is in a directory that contains the 'proximal-exploration' folder
script_directory = os.path.dirname(os.path.abspath(__file__))
proximal_exploration_path = os.path.join(script_directory, 'proximal-exploration')

# Add 'proximal-exploration' to the Python path
sys.path.append(proximal_exploration_path)

import numpy as np
from landscape import get_landscape, task_collection, landscape_collection
from algorithm import get_algorithm, algorithm_collection
from model import get_model, model_collection
from model.ensemble import ensemble_rules
from utils.os_utils import get_arg_parser
from utils.eval_utils import Runner


def get_args():
    parser = get_arg_parser()
                    
    parser.add_argument('--device', help='device', type=str, default='cuda')
                    
    # landscape arguments
    parser.add_argument('--task', help='fitness landscape', type=str, default='avGFP', choices=task_collection.keys())
    parser.add_argument('--oracle_model', help='oracle model of fitness landscape', type=str, default='tape', choices=landscape_collection.keys())

    # algorithm arguments
    parser.add_argument('--alg', help='exploration algorithm', type=str, default='pex', choices=algorithm_collection.keys())
    parser.add_argument('--num_rounds', help='number of query rounds', type=np.int32, default=10)
    parser.add_argument('--num_queries_per_round', help='number of black-box queries per round', type=np.int32, default=100)
    parser.add_argument('--num_model_queries_per_round', help='number of model predictions per round', type=np.int32, default=2000)
                    
    # model arguments
    parser.add_argument('--net', help='surrogate model architecture', type=str, default='mufacnet', choices=model_collection.keys())
    parser.add_argument('--lr', help='learning rate', type=np.float32, default=1e-3)
    parser.add_argument('--batch_size', help='batch size', type=np.int32, default=256)
    parser.add_argument('--patience', help='number of epochs without improvement to wait before terminating training', type=np.int32, default=10)
    parser.add_argument('--ensemble_size', help='number of model instances in ensemble', type=np.int32, default=3)
    parser.add_argument('--ensemble_rule', help='rule to aggregate the ensemble predictions', type=str, default='mean', choices=ensemble_rules.keys())

    args, _ = parser.parse_known_args()
                    
    # PEX arguments
    if args.alg == 'pex':
        parser.add_argument('--num_random_mutations', help='number of amino acids to mutate per sequence', type=np.int32, default=2)
        parser.add_argument('--frontier_neighbor_size', help='size of the frontier neighbor', type=np.int32, default=5)
                    
    # MuFacNet arguments
    if args.net == 'mufacnet':
        parser.add_argument('--latent_dim', help='dimension of latent mutation embedding', type=np.int32, default=32)
        parser.add_argument('--context_radius', help='the radius of context window', type=np.int32, default=10)

    args = parser.parse_args()
    return args

class Agent:

    def __init__(self, t, df, reward, policy_flag):
        self.t = t
        self.df = df
        self.reward = reward
        self.policy_flag = policy_flag

    def apply_policy(self):

        # select subset of sequences (those of the current step) on which to which to apply the policy
        if self.t == 0:
            df_action = self.df[self.df['t'] == 0][['t', 'sequence_number', 'seq']].copy()
        else:
            df_action = self.df[self.df['t'] == self.t - 1][['t', 'sequence_number', 'seq']].copy()        

        if self.policy_flag == 'proximal_exploration':

            # load the wild-type sequence and its score f(wt) from the first row (zeroth iteration entry) of df.
            # nb: may want to move the code for df_action above into the policies where it is relevant.

            for index, row in df_action.iterrows():

                args = get_args()

                # # Some pseudo code / this is def run() function from the evals_utils script adapted and put in the context of our code.
                # D^{t-1} is given by the entire (!) df table

                # train the model f^hat on D^{t-1}

                # run pex with f^hat and D^{t-1}, this gives a new set of M sequences, s_i^(t). Probably store them in df_action.

            # # Concatenate self.df with df_action
            # result_df = pd.concat([self.df, df_action])
        
            # next step is to return results_df and the oracle is computing the truth f(s_i^(t)) for us and supplement the full df table with it.
            return result_df, df_action     

        elif self.policy_flag == 'random_mutation':

            amino_acids = 'ACDEFGHIKLMNPQRSTVWY'  # List of amino acids in one-letter code   

            # Update the 't' column in df_action to the current value of self.t
            df_action['t'] = self.t

            for index, row in df_action.iterrows():
                seq = list(row['seq'])  # Convert string to list for mutation

                if len(seq) > 0:
                    mutation_pos = random.randint(0, len(seq) - 1)  # Select random position
                    original_residue = seq[mutation_pos]

                    # Select a new amino acid different from the original
                    new_residue = random.choice([aa for aa in amino_acids if aa != original_residue])

                    # Perform the mutation
                    seq[mutation_pos] = new_residue

                    # Update the sequence in the DataFrame without brackets
                    df_action.at[index, 'seq'] = ''.join(seq)

                    print('df_action', df_action)

                    # Optional: Print the original and mutated sequences with highlighted mutations
                    original_seq_for_printing = ''.join(seq[:mutation_pos]) + '[' + original_residue + ']' + ''.join(seq[mutation_pos + 1:])
                    mutated_seq_for_printing = ''.join(seq[:mutation_pos]) + '[' + new_residue + ']' + ''.join(seq[mutation_pos + 1:])
                    print(f"Original: {original_seq_for_printing} -> Mutated: {mutated_seq_for_printing}")

            # Concatenate self.df with df_action
            result_df = pd.concat([self.df, df_action])

            return result_df, df_action
        
        elif self.policy_flag == 'random_deletion':

            for index, row in df_action.iterrows():
                seq = list(row['seq'])  # Convert string to list for deletion

                if len(seq) > 1:  # Ensure there is at least one residue to delete
                    deletion_pos = random.randint(0, len(seq) - 1)  # Select random position

                    original_residue = seq[deletion_pos]
                    original_seq_for_printing = ''.join(seq[:deletion_pos]) + '[' + original_residue + ']' + ''.join(seq[deletion_pos + 1:])
                    seq_after_delete_for_printing = ''.join(seq[:deletion_pos]) + ''.join(seq[deletion_pos + 1:])

                    # Perform the deletion
                    del seq[deletion_pos]

                    # Update the sequence in the DataFrame after deletion
                    df_action.at[index, 'seq'] = ''.join(seq)

                    # Optional: Print the sequence before and after deletion
                    print(f"Sequence before deletion: {original_seq_for_printing} -> Sequence after deletion: {seq_after_delete_for_printing}")

            # Update the 't' column in df_action to the current value of self.t
            df_action['t'] = self.t

            # Concatenate self.df with df_action
            result_df = pd.concat([self.df, df_action])

            return result_df, df_action
        

### some old code snippets
            # ## some pseudo code:
            # # replacing landscape, alphabet, starting_sequence = get_landscape(args):
            # alphabet = 'ACDEFGHIKLMNPQRSTVWY'
            # starting_sequence = row['seq']
            # # and landscape is now given by the plddt from 

            # print('landscape', landscape)
            # print('alphabet', alphabet)
            # print('starting_seq', starting_sequence)
            # model = get_model(args, alphabet=alphabet, starting_sequence=starting_sequence)
            # explorer = get_algorithm(args, model=model, alphabet=alphabet, starting_sequence=starting_sequence)

            # runner = Runner(args)
            # runner.run(landscape, starting_sequence, model, explorer)   