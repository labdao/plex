from .base_generator import BaseGenerator

class trivial_Generator(BaseGenerator):

    def generate(self, args):

        return ''.join(args.sequence), ''.join(args.permissibility_vector)