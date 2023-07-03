# ProtBERT: Predicting Protein Residues with Language Models

ProtBERT is a powerful tool designed for predicting unknown protein residues using language models. It takes as input a protein sequence and can operate in several modes, including fill-mask, embedding, and conditional-probability.

This tool is specifically designed to be run on PLEX, a client for distributed computation. It is built to utilize the distributed compute and storage capabilities of the PLEX network to enable fast, scalable, and highly reproducible protein residue predictions.

## Running ProtBERT with PLEX

To run ProtBERT as a PLEX job, you can use the provided Docker container and the JSON configuration file, `protbert.json`. The following is an example command to create and auto-run a ProtBERT job on PLEX:

```
./plex create -t tools/protbert.json -i your_input_file.fasta --autoRun=True
```

Replace `your_input_file.fasta` with the path to your input file.

## Inputs and Outputs

ProtBERT takes as input a file containing a protein sequence. The expected file format is FASTA, and it should have one of the following extensions: `.fasta`, `.fa`, `.faa`, `.fna`, `.ffn`, `.frn`, `.fsa`, `.fas`.

The outputs from ProtBERT include the completed protein sequence (`*_mask.json`), the embedded protein sequence (`*_encoded.csv`), and the protein sequence conditional probability (`*_scoring_matrix.csv`). These output files will be generated in the `/outputs` directory.

The inputs and outputs are defined in the JSON configuration file (`protbert.json`) and are handled by PLEX.

## Distributed Computation with PLEX

ProtBERT is optimized to take advantage of the distributed computation capabilities of PLEX. This means you can expect fast, scalable, and highly reproducible results when running ProtBERT jobs, regardless of the size of your protein sequence. By leveraging the PLEX network, ProtBERT can deliver predictions quickly and at any scale.

---

Please review this draft and let me know if there are any other changes you'd like to see.