def squeeze_seq(new_sequence):
    return ''.join(filter(lambda x: x != '-', new_sequence))

def generate_contig(action_mask, target, starting_target_residue=None, end_target_residue=None):
    if starting_target_residue is None:
        starting_target_residue = 1
    if end_target_residue is None:
        end_target_residue = len(target)
    
    # Initialize variables
    action_mask_contig = ''
    current_group = ''
    alphabet = 'LAGVSERTIDPKQNFYMHWC'
    position = 0  # Position within the action_mask
    
    # Iterate over the squeezed_action_mask to form groups
    for char in action_mask:
        if char in alphabet:
            if current_group == '' or current_group[-1] in alphabet:
                current_group += char  # Continue the current alphabet group
            elif current_group[-1]=='X':
                action_mask_contig += f'{len(current_group)}/'
                current_group = char
        elif char=='X':  # char is 'X'
            if current_group == '' or current_group[-1] == 'X':
                current_group += char  # Continue the current X group
            elif current_group[-1] in alphabet:
                action_mask_contig += f'A{position-len(current_group)+1}-{position}/'
                current_group = char
        elif char=='-':
            if current_group!='' and current_group[-1] in alphabet:
                action_mask_contig += f'A{position-len(current_group)+1}-{position}/'
                current_group = ''
            elif current_group!='' and current_group[-1]=='X':
                action_mask_contig += f'{len(current_group)}/'
                current_group = ''

        position += 1
    
    # Add the last group to the contig
    if current_group:
        if current_group[-1] == 'X':
            action_mask_contig += f'{len(current_group)}/'  # X group
        else:
            action_mask_contig += f'B{position-len(current_group)+1}-{position}/'  # Alphabet group
    
    # Remove the trailing '/' if it exists
    if action_mask_contig.endswith('/'):
        action_mask_contig = action_mask_contig[:-1]
    
    # Construct the final contig string
    contig = f'A{starting_target_residue}-{end_target_residue}/0 {action_mask_contig}'
    
    return contig
