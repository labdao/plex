# run from basedir
# pull container images for basic CASF 2016 evaluation
docker pull ghcr.io/labdao/diffdock:main
docker pull gnina/gnina
docker pull ghcr.io/labdao/casf-2016-eval:main

# pull data for basic CASF 2016 evaluation
mkdir data
cd data
ipfs get bafybeiaqyjf65cs2slhilsrqvo3mo6ckdqnr5spplcts7svq7256hiiguy
cd ..
# avoided by adding some examples into git
# aws s3 sync s3://labdao-benchmark/CASF-2016/coreset/1a30/ data


## interactive for debug
# docker run -it -v /home/ubuntu/casf-2016-evaluator/data/coreset:/inputs -v /home/ubuntu/outputs:/outputs gnina/gnina bash
# NOTRUN dock gnina
# NOTRUN gnina -r /inputs/1a30/1a30_protein.pdb -l /inputs/1a30/1a30_ligand.sdf --autobox_ligand /inputs/1a30/1a30_protein.pdb --cnn_scoring rescore -o /outputs/whole_docked.sdf.gz --exhaustiveness 64
# NOTRUN score gnina
# NOTRUN gnina -r /inputs/1a30/1a30_protein.pdb -l /inputs/1a30/1a30_ligand.sdf --autobox_ligand /inputs/1a30/1a30_protein.pdb --cnn_scoring rescore -o /outputs/scored_gnina.sdf.gz --exhaustiveness 64 --score_only 
# score vina
# gnina -r /inputs/1a30/1a30_protein.pdb -l /inputs/1a30/1a30_ligand.sdf --autobox_ligand /inputs/1a30/1a30_protein.pdb --cnn_scoring none -o /outputs/scored_vina.sdf.gz --exhaustiveness 64 --score_only 

# run scoring
docker run \
    -v /home/ubuntu/bafybeiaqyjf65cs2slhilsrqvo3mo6ckdqnr5spplcts7svq7256hiiguy:/inputs \
    -v /home/ubuntu/outputs:/outputs \
    gnina/gnina \
    gnina   --autobox_ligand /inputs/1a30/1a30_protein.pdb \
            --cnn_scoring none \
            --exhaustiveness 64 \
            --score_only  \
            -r /inputs/1a30/1a30_protein.pdb \
            -l /inputs/1a30/1a30_ligand.sdf \
            -o /outputs/1a30/1a30_scored_vina.sdf.gz
