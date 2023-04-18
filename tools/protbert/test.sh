#!/bin/bash
python3 app.py test.fasta /test --mode embedding
python3 app.py test.fasta /test --mode fill-mask
python3 app.py test.fasta /test --mode conditional-probability
python3 app.py test.fasta /test --mode joint-probability
python3 app.py test.fasta /test --mode generate