import numpy as np
import logging
from utils import determine_acceptance

class SequenceSelector:
    def __init__(self, cfg):
        self.cfg = cfg

    def _sample_isotropic_dirichlet(self, d, alpha=1.0):

        alpha_vector = np.full(d, alpha)
        sample = np.random.dirichlet(alpha_vector)
        return sample

    def run(self, t, df, accept_flag):
        if self.cfg.params.basic_settings.bouncer_flag == 'open-door':
            acceptance_probability = 1.0
            logging.info(f"acceptance probability: {acceptance_probability}")
            accept_flag = True
            logging.info(f"accept_flag: {accept_flag}")
            return True

        elif self.cfg.params.basic_settings.bouncer_flag == 'boltzmann':
            T = self.cfg.params.basic_settings.temperature
            E = 0.
            scoring_metrics = self.cfg.params.basic_settings.scoring_metrics
            d = len(scoring_metrics)
            weights = self._sample_isotropic_dirichlet(d, alpha=10.0)
            for i, metric in enumerate(scoring_metrics.split(',')):
                ref_row = df[df['acceptance_flag'] == True].iloc[-1:]
                previousStep_sequence_metric = ref_row[metric].values[0]
                proposed_sequence_metric = df.iloc[-1][metric]
                DeltaE = proposed_sequence_metric - previousStep_sequence_metric
                if metric in ['pseudolikelihood']:
                    DeltaE *= -weights[i]
                E += DeltaE
            p_mod = np.exp(-E / T)
            acceptance_probability = np.minimum(1, p_mod)
            logging.info(f"acceptance probability: {acceptance_probability}")
            accept_flag = determine_acceptance(acceptance_probability)
            logging.info(f"accept_flag: {accept_flag}")
            return accept_flag