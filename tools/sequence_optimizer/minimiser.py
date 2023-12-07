import os
import pandas as pd
import hydra
from dataclasses import dataclass, field
from hydra.core.config_store import ConfigStore
from omegaconf import DictConfig
from pdbfixer import PDBFixer
from openmm import app, Platform, LocalEnergyMinimizer
import openmm as mm
from openmm.app import PDBFile
from openmm.unit import *
import datetime

#TODO test with different forcefields and implicit solvent models - http://docs.openmm.org/latest/userguide/application/02_running_sims.html#implicit-solvent 
@dataclass
class Config:
    pdb_files: list = field(default_factory=list)
    output_dir: str = "/outputs"
    forcefield_model: field(default_factory=list) 
    implicit_solvent_model: field(default_factory=list) 

def prepare_pdbs_for_minimization(pdb_file_paths, output_dir=None):
    fixed_pdb_paths = []

    for pdb_path in pdb_file_paths:
        # Determine output file path
        pdb_directory, pdb_filename = os.path.split(pdb_path)
        pdb_filename_without_ext, _ = os.path.splitext(pdb_filename)
        fixed_filename = pdb_filename_without_ext + '_fixed.pdb'

        if output_dir:
            fixed_pdb_path = os.path.join(output_dir, fixed_filename)
        else:
            fixed_pdb_path = os.path.join(pdb_directory, fixed_filename)

        # Fix PDB
        fixer = PDBFixer(filename=pdb_path)
        fixer.findMissingResidues()
        fixer.findNonstandardResidues()
        fixer.replaceNonstandardResidues()
        fixer.removeHeterogens(keepWater=False)
        fixer.findMissingAtoms()
        fixer.addMissingAtoms()
        fixer.addMissingHydrogens(7.0)

        # Write fixed PDB
        with open(fixed_pdb_path, 'w') as f:
            PDBFile.writeFile(fixer.topology, fixer.positions, f)

        fixed_pdb_paths.append(fixed_pdb_path)

    return fixed_pdb_paths

def minimize_pdbs(pdb_file_paths, forcefield_list, implicit_solvent_list, output_dir=None, platform_name='CUDA', write_csv=True):
    if len(forcefield_list) != len(implicit_solvent_list):
        raise ValueError("Forcefield and implicit solvent lists must be the same length")

    if output_dir is None:
        output_dir = os.getcwd()  # default to current working directory

    minimized_pdb_paths = []
    energy_data = []

    for pdb_path in pdb_file_paths:
        try:
            pdb = app.PDBFile(pdb_path)
            pdb_directory, pdb_filename = os.path.split(pdb_path)
            pdb_filename_without_ext, _ = os.path.splitext(pdb_filename)

            for i in range(len(forcefield_list)):
                minimized_filename = f"{pdb_filename_without_ext}_minimized_{i}_{datetime.datetime.now().strftime('%Y%m%d%H%M%S')}.pdb"
                minimized_pdb_path = os.path.join(output_dir, minimized_filename)

                run_local_minimization(pdb, forcefield_list[i], implicit_solvent_list[i], platform_name, minimized_pdb_path, energy_data, pdb_filename, minimized_filename)

            minimized_pdb_paths.append(minimized_pdb_path)

        except Exception as e:
            print(f"Error processing {pdb_path}: {e}")

    # DataFrame and CSV handling
    energy_df = pd.DataFrame(energy_data)
    if write_csv:
        csv_path = os.path.join(output_dir, "minimization_energy_data.csv")
        energy_df.to_csv(csv_path, index=False)
        print(f"Energy data written to CSV at: {csv_path}")

    return minimized_pdb_paths, energy_df

def run_local_minimization(pdb, forcefield_model, implicit_solvent_model, platform_name, minimized_pdb_path, energy_data, pdb_filename, minimized_filename):
    forcefield = app.ForceField(forcefield_model, implicit_solvent_model)
    system = forcefield.createSystem(pdb.topology, nonbondedMethod=app.NoCutoff, constraints=app.HBonds)
    
    platform = mm.Platform.getPlatformByName(platform_name)

    context = mm.Context(system, mm.VerletIntegrator(1.0 * mm.unit.femtoseconds), platform)
    context.setPositions(pdb.positions)

    initial_energy = get_energy(context)
    mm.LocalEnergyMinimizer.minimize(context, tolerance=0.001)
    final_energy = get_energy(context)

    state = context.getState(getPositions=True)
    with open(minimized_pdb_path, 'w') as outfile:
        app.PDBFile.writeFile(pdb.topology, state.getPositions(), outfile)

    energy_data.append({
        "Original Protein": pdb_filename,
        "Minimized Protein": minimized_filename,
        "Initial Energy (kJ/mol)": initial_energy,
        "Final Energy (kJ/mol)": final_energy, 
        "Energy Difference (kJ/mol)": final_energy - initial_energy,
        "Forcefield": forcefield_model,
        "Implicit Solvent": implicit_solvent_model
    })

def get_energy(context):
    state = context.getState(getEnergy=True)
    return state.getPotentialEnergy().value_in_unit(kilojoules_per_mole)



@hydra.main(config_path="conf", config_name="config")
def main(cfg: Config) -> None:
    print("Configuration:", cfg)  # Print the configuration
    output_directory = cfg.output_dir if cfg.output_dir else None
    print("Original PDB files:", cfg.pdb_files)
    fixed_pdbs = prepare_pdbs_for_minimization(cfg.pdb_files, cfg.forcefield_model, cfg.implicit_solvent_model, output_directory)
    print("Fixed PDB files:", fixed_pdbs)
    minimized_pdbs, energy_summary = minimize_pdbs(fixed_pdbs, output_directory)
    print("Minimized PDB files:", minimized_pdbs)
    print("Energy Overview:", energy_summary)

if __name__ == "__main__":
    main()