import os
import json
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns


def get_plex_job_inputs():
    # Retrieve the environment variable
    json_str = os.getenv("PLEX_JOB_INPUTS")

    # Check if the environment variable is set
    if json_str is None:
        raise ValueError("PLEX_JOB_INPUTS environment variable is missing.")

    # Convert the JSON string to a Python dictionary
    try:
        data = json.loads(json_str)
        return data
    except json.JSONDecodeError:
        # Handle the case where the string is not valid JSON
        raise ValueError("PLEX_JOB_INPUTS is not a valid JSON string.")


def select_best_sample(df):
    def process_group(group):
        # Filter for designs with plddt >= 0.8 and rmsd <= 2
        filtered = group[(group['plddt'] >= 0.8) & (group['rmsd'] <= 1)]

        # If some designs meet both criteria, select the one with the lowest i_pae
        if not filtered.empty:
            return filtered.nsmallest(1, 'i_pae')

        # If no designs meet both criteria, then filter for plddt >= 0.8
        # and select the one with the lowest i_pae
        filtered = group[group['plddt'] >= 0.8]
        if not filtered.empty:
            return filtered.nsmallest(1, 'i_pae')

        # If no designs meet plddt >= 0.8 criterion, select the one with the highest plddt
        return group.nlargest(1, 'plddt')

    # Apply the selection logic to each group
    return df.groupby('global_design').apply(process_group).reset_index(drop=True)


def create_unique_design_id(df, file_idx):
    df['global_design'] = str(file_idx) + "_" + df['design'].astype(str)
    return df


def main():
    # Get the job inputs from the environment variable
    try:
        job_inputs = get_plex_job_inputs()
        print("Job Inputs:", job_inputs)
    except ValueError as e:
        print(e)
        sys.exit(1)

    # Initialize an empty list to hold data from each CSV file
    all_data = []

    for i, file_path in enumerate(job_inputs["csv_result_files"]):
        print(f'Reading file: {file_path}, {i + 1} out of {len(job_inputs["csv_result_files"])}')
        df = pd.read_csv(file_path)
        all_data.append(df)

        df = create_unique_design_id(df, i)

        # Log the number of unique designs in the current CSV file
        num_designs = df['design'].nunique()
        print(f'File: {file_path}, Number of Designs: {num_designs}')

    # Concatenate all data into a single DataFrame
    aggregated_df = pd.concat(all_data, ignore_index=True)

    # Apply the selection logic to the aggregated data
    best_samples_df = select_best_sample(aggregated_df)

    # Create /outputs directory if it doesn't exist
    os.makedirs("/outputs", exist_ok=True)

    # Save the selected best samples to a new CSV file
    best_samples_df.to_csv("/outputs/aggregated.csv", index=False)

    # Apply the final filter criteria
    filtered_df = best_samples_df[(best_samples_df['i_pae'] <= 10) & (best_samples_df['plddt'] >= 0.8) & (best_samples_df['rmsd'] <= 1) & (best_samples_df['affinity'] <= -8.5)]

    # Print the details of designs that pass the final filter
    for index, row in filtered_df.iterrows():
        print(f"Passing Design: {row['global_design']}, i_pae: {row['i_pae']}, plddt: {row['plddt']}, rmsd: {row['rmsd']}, affinity: {row['affinity']}, sequence: {row['seq']}")

    # Calculate the ratio of designs meeting the criteria using best_samples_df
    total_designs = len(best_samples_df)
    print(f'Total Unique Designs After Selection: {total_designs}')
    print(f'Passing Designs: {len(filtered_df)}')
    ratio_designs_passed = len(filtered_df) / total_designs

    # Setting up the new subplot layout with 2 rows and 2 columns
    fig, axs = plt.subplots(2, 2, figsize=(12, 10))  # Adjusted figure size for better layout

    # Define a consistent and visually appealing color palette
    palette = sns.color_palette("husl", 3)

    # pLDDT Plot
    sns.histplot(best_samples_df['plddt'], bins=20, kde=True, color=palette[2], ax=axs[0, 0])
    axs[0, 0].fill_betweenx([0, axs[0, 0].get_ylim()[1]], 0.8, axs[0, 0].get_xlim()[1], color='lightgreen', alpha=0.3)
    axs[0, 0].set_title('pLDDT (passing is greater than 0.8)')
    axs[0, 0].set_xlabel('plddt')
    axs[0, 0].set_ylabel('Frequency')

    # Interaction PAE Plot
    sns.histplot(best_samples_df['i_pae'], bins=20, kde=True, color=palette[2], ax=axs[0, 1])
    axs[0, 1].fill_betweenx([0, axs[0, 1].get_ylim()[1]], 5, 10, color='lightgreen', alpha=0.3)
    axs[0, 1].set_title('Interaction PAE (passing is less than 10)')
    axs[0, 1].set_xlabel('i_pae')
    axs[0, 1].set_ylabel('Frequency')

    # RMSD Plot
    sns.histplot(best_samples_df['rmsd'], bins=1000, kde=True, color=palette[2], ax=axs[1, 0])
    axs[1, 0].fill_betweenx([0, axs[1, 0].get_ylim()[1]], 0, 1.0, color='lightgreen', alpha=0.3)
    axs[1, 0].set_title('Zoomed RMSD of ProteinMPNN and AF2 (Passing is less than 1.0)')
    axs[1, 0].set_xlabel('rmsd')
    axs[1, 0].set_ylabel('Frequency')
    # axs[1, 0].set_xlim(0, 3)

    # Affinity Plot
    sns.histplot(best_samples_df['affinity'], bins=20, kde=True, color=palette[2], ax=axs[1, 1])
    axs[1, 1].fill_betweenx([0, axs[1, 1].get_ylim()[1]], -11, -8.5, color='lightgreen', alpha=0.3)
    axs[1, 1].set_title('Affinity (passing is less than -8.5)')
    axs[1, 1].set_xlabel('affinity')
    axs[1, 1].set_ylabel('Frequency')

    # Save the plots
    plt.tight_layout()
    plt.savefig('/outputs/distribution.png')


if __name__ == "__main__":
    main()

