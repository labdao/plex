def squeeze_seq(new_sequence):
    return ''.join(filter(lambda x: x != '-', new_sequence))