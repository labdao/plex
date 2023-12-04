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

#TODO test with different forcefields and implicit solvent models - http://docs.openmm.org/latest/userguide/application/02_running_sims.html#implicit-solvent 
@dataclass
class Config:
    pdb_files: list = field(default_factory=list)
    output_dir: str = "/outputs"
    forcefield_model: str = "amber99sbildn.xml"
    implicit_solvent_model: str = "amber99_obc.xml"

def prepare_pdb_for_minimization(pdb_file_paths, output_dir=None):
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

def minimize_protein(pdb_file_paths, forcefield_model, implicit_solvent_model, output_dir=None, write_csv=True):
    minimized_pdb_paths = []
    energy_data = []

    for pdb_path in pdb_file_paths:
        pdb_directory, pdb_filename = os.path.split(pdb_path)
        pdb_filename_without_ext, _ = os.path.splitext(pdb_filename)
        minimized_filename = pdb_filename_without_ext + '_minimized.pdb'

        output_dir = output_dir if output_dir else pdb_directory
        minimized_pdb_path = os.path.join(output_dir, minimized_filename)

        # Load the PDB structure
        pdb = app.PDBFile(pdb_path)
        print("Minimizing: " + pdb_path)

        # Specify the force field - AMBER99SB-ILDN with OBC implicit solvent
        forcefield = app.ForceField(forcefield_model, implicit_solvent_model)

        # Create the system
        system = forcefield.createSystem(pdb.topology, nonbondedMethod=app.NoCutoff, constraints=app.HBonds)

        # Specify the platform; use CUDA for NVIDIA GPUs or OpenCL for other GPUs
        platform = Platform.getPlatformByName('CUDA')  # Change 'CUDA' to 'OpenCL' if needed

        # Create a Context for the system
        context = mm.Context(system, mm.VerletIntegrator(1.0*mm.unit.femtoseconds), platform)
        context.setPositions(pdb.positions)

        # Get initial energy
        state = context.getState(getEnergy=True)
        initial_energy = state.getPotentialEnergy().value_in_unit(kilojoules_per_mole)
        print(f"Initial energy: {initial_energy}")

        # Perform local energy minimization using LocalEnergyMinimizer
        mm.LocalEnergyMinimizer.minimize(context, tolerance=0.001)

        # Get final energy
        state = context.getState(getEnergy=True)
        final_energy = state.getPotentialEnergy().value_in_unit(kilojoules_per_mole)
        print(f"Final energy: {final_energy}")

        # Save the minimized structure
        state = context.getState(getPositions=True)
        with open(minimized_pdb_path, 'w') as outfile:
            app.PDBFile.writeFile(pdb.topology, state.getPositions(), outfile)

        # Store energy data
        energy_data.append({
            "Original Protein": pdb_filename,
            "Minimized Protein": minimized_filename,
            "Initial Energy (kJ/mol)": initial_energy,
            "Final Energy (kJ/mol)": final_energy
        })

        minimized_pdb_paths.append(minimized_pdb_path)

    # Create DataFrame
    energy_df = pd.DataFrame(energy_data)

    # Optionally write to CSV
    if write_csv and output_dir:
        csv_path = os.path.join(output_dir, "minimization_energy_data.csv")
        energy_df.to_csv(csv_path, index=False)
        print(f"Energy data written to CSV at: {csv_path}")

    return minimized_pdb_paths, energy_df


@hydra.main(config_path="conf", config_name="config")
def main(cfg: Config) -> None:
    print("Configuration:", cfg)  # Print the configuration
    output_directory = cfg.output_dir if cfg.output_dir else None
    print("Original PDB files:", cfg.pdb_files)
    fixed_pdbs = prepare_pdb_for_minimization(cfg.pdb_files, cfg.forcefield_model, cfg.implicit_solvent_model, output_directory)
    print("Fixed PDB files:", fixed_pdbs)
    minimized_pdbs, energy_summary = minimize_protein(fixed_pdbs, output_directory)
    print("Minimized PDB files:", minimized_pdbs)
    print("Energy Overview:", energy_summary)

if __name__ == "__main__":
    main()
