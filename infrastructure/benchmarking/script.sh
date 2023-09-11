echo "Running equibind..."

./plex init -t tools/equibind.json -i '{"protein": ["testdata/binding/6d08_protein_processed.pdb"], "small_molecule": ["testdata/binding/6d08_ligand.sdf"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking

echo "Running diffdock..."

./plex init -t tools/diffdock.json -i '{"protein": ["testdata/binding/6d08_protein_processed.pdb"], "small_molecule": ["testdata/binding/6d08_ligand.sdf"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking

echo "Running colabfold..."

declare -A proteins_to_fold

# proteins with their amino acid (AA) lengths
# longer lengths will fold more slowly

proteins_to_fold["testdata/folding/histatin-3.fasta"]="51"
proteins_to_fold["testdata/folding/tax1-binding_protein_3.fasta"]="124"
proteins_to_fold["testdata/folding/c-reactive_protein.fasta"]="224"
proteins_to_fold["testdata/folding/gap_junction_protein.fasta"]="396"
proteins_to_fold["testdata/folding/phospholipid_glycerol_acyltransferase_domain-containing_protein.fasta"]="795"
proteins_to_fold["testdata/folding/rims-binding_protein_3b.fasta"]="1639"
proteins_to_fold["testdata/folding/dna_helicase.fasta"]="3011"
proteins_to_fold["testdata/folding/hect-type_e3_ubiquitin_transferase.fasta"]="4374"

for protein in "${!proteins_to_fold[@]}"; do
    echo "Folding $protein with increasing AA length ${proteins_to_fold[$protein]}..."
    ./plex init -t tools/colabfold-mini.json -i "{\"sequence\": [\"$protein\"]}" --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking
done

echo "Running RFDiffusion-colabdesign..."
echo "Run: small"
./plex init -t tools/colabdesign/_colabdesign-dev.json -i '{"protein": ["tools/colabdesign/6vja_stripped.pdb"], "config": ["tools/colabdesign/small-config.yaml"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking
echo "Run: medium"
./plex init -t tools/colabdesign/_colabdesign-dev.json -i '{"protein": ["tools/colabdesign/6vja_stripped.pdb"], "config": ["tools/colabdesign/config.yaml"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking
echo "Run: large"
./plex init -t tools/colabdesign/_colabdesign-dev.json -i '{"protein": ["tools/colabdesign/6vja_stripped.pdb"], "config": ["tools/colabdesign/large-config.yaml"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking
