# Conditional Colabdesign
Based on Sergey Ovchinnikov Colab notebooks

## Run an Example


```
# initiate a flow
./plex init -t tools/colabdesign-conditional/_colabdesign_conditional-dev.json -i '{"binder_protein_template": ["tools/colabdesign-conditional/input_target.pdb"], "config": ["tools/colabdesign-conditional/params.yaml"], "protein": ["tools/colabdesign-conditional/input_binder.pdb"]}' --scatteringMethod=dotProduct --autoRun=false
```

```
# run an existing flow
./plex run -i QmQ5vzPD6AoXT7MUMFmwFFUaLK3eu3x2j4odziyKdgTDVN
```
