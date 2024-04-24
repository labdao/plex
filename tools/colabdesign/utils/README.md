READ ME

Hello! visulaisers.py is a script that will render a .pdb file in Blender 
using the MolecularNodes addon and output a .png

Necessary packages:
To run visulaisers.py you need to be using

* a python version between 3.10 and 3.11 *

You'll also need to have Blender and the molecularnodes addon installed.
Notably this script will only work for molecular nodes v.4.0.12

pip install bpy
pip install molecularnodes==4.0.12

Notable bits:

(line 229) is where you will find the input:
pdb=FilePath("/Users/lily/Downloads/Condition 52 Design 0.pdb")

(Line 35) has the output directory.
output_directory = os.path.join(os.getcwd(), "outputs/blender")

(Line 186) changes the render engine
bpy.context.scene.render.engine = 'CYCLES'
There are different blender render engines that vary in the amount of time they take to run.
Read more here: https://renderguide.com/blender-eevee-vs-cycles-tutorial/
'CYCLES' is a more powerful rendering engine, it takes longer and I've commented
it out in the script for now.
As a point of reference it takes up to 10 minutes to run the visualiser on cycles using CPU.
(it takes about 15 seconds with eevee)
I would love if someone could use a few differnt .pdbs and, running on GPU,test how long it takes  using CYCLES.