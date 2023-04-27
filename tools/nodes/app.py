import bpy

# Install and enable nodes
# Register the addon and enable it
bpy.ops.preferences.addon_install(filepath='./MolecularNodes_2.6.0.zip')
bpy.ops.preferences.addon_enable(module='MolecularNodes_2.6.0')

# Use the addon directly from Python
# ...