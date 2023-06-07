import bpy
import os
import sys
from math import radians
import random
from mathutils import Vector, Matrix

def get_random_color(pastel_factor=0.4):
    return [(x + pastel_factor) / (1.0 + pastel_factor) for x in [random.uniform(0, 1.0) for i in [1, 2, 3]]]

def color_distance(c1, c2):
    return sum([abs(x[0] - x[1]) for x in zip(c1, c2)])

def generate_new_color(existing_colors, pastel_factor=0.5):
    max_distance = None
    best_color = None
    for i in range(0, 100):
        color = get_random_color(pastel_factor=pastel_factor)
        if not existing_colors:
            return color
        best_distance = min([color_distance(color, c) for c in existing_colors])
        if not max_distance or best_distance > max_distance:
            max_distance = best_distance
            best_color = color
    return best_color

def main(input_file_path, output_file_path):
    bpy.ops.preferences.addon_enable(module='io_mesh_atomic')
    
    pdb_file_path = input_file_path
    
    if os.path.exists(pdb_file_path):
        bpy.ops.object.select_all(action='SELECT')
        bpy.ops.object.delete()
        bpy.ops.import_mesh.pdb(filepath=pdb_file_path)
    else:
        print("PDB file not found:", pdb_file_path)
        return

    camera_data = bpy.data.cameras.new(name="New_Camera")
    camera_object = bpy.data.objects.new("New_Camera", camera_data)
    scene = bpy.context.scene
    scene.collection.objects.link(camera_object)
    scene.camera = camera_object
    camera_object.location = (58.991, -11.855, 30.663)
    camera_object.rotation_mode = 'ZYX'
    camera_object.rotation_euler = (radians(24.4), radians(57.9), radians(-291))
    camera_object.scale = (1.0, 1.0, 1.0)
    bpy.ops.object.select_all(action='DESELECT')
    camera_object.select_set(True)
    bpy.context.view_layer.objects.active = camera_object

    colors = []
    for material in bpy.data.materials:
        new_color = generate_new_color(colors, pastel_factor=0.9)
        if material.use_nodes:
            for node in material.node_tree.nodes:
                if node.type == 'BSDF_PRINCIPLED':
                    node.inputs['Base Color'].default_value = (new_color[0], new_color[1], new_color[2], 1)
                    node.inputs['Roughness'].default_value = 0
                    node.inputs['Sheen Tint'].default_value = 0.7
                    node.inputs['Metallic'].default_value = 0.8

    bpy.context.view_layer.update()

    camera_location = Vector((58.991, -11.855, 30.663))
    camera_rotation = (radians(24.4), radians(57.9), radians(-291))
    light_location = camera_location + Vector((0, 0, 3)) - 3 * (Matrix.Rotation(camera_rotation[0], 3, 'X') @ Matrix.Rotation(camera_rotation[1], 3, 'Y') @ Matrix.Rotation(camera_rotation[2], 3, 'Z') @ Vector((0, 1, 0))) + Vector((2, 0, 0))

    light_data = bpy.data.lights.new(name="New_Sun_Light", type='SUN')
    light_object = bpy.data.objects.new("New_Sun_Light", light_data)
    scene.collection.objects.link(light_object)
    light_object.location = light_location
    light_object.rotation_mode = 'XYZ'
    light_rotation = Vector(camera_rotation)
    random_y_rotation_offset = radians(random.uniform(-12, -8))
    light_rotation.y += random_y_rotation_offset
    light_object.rotation_euler = light_rotation
    light_data.energy = 10

    bpy.context.scene.render.image_settings.file_format = 'PNG'
    bpy.context.scene.render.filepath = output_file_path
    bpy.ops.render.render(write_still=True)

if __name__ == "__main__":
    if "--" not in sys.argv:
        argv = []  # as if no args are passed
    else:
        argv = sys.argv[sys.argv.index("--") + 1:]  # get all args after "--"
    input_file_path = argv[0]
    output_file_path = argv[1]
    main(input_file_path, output_file_path)


# blender --background --python /inputs/app.py -- /inputs/5RGA.pdb /outputs/file.png
# sudo docker run -it -v /home/jupyter-niklas/plex/tools/blender:/inputs -v /home/jupyter-niklas/plex/tools/blender:/outputs nytimes/blender:3.3.1-cpu-ubuntu18.04 bash
# sudo docker run -it -v /Users/rindtorff/github/labdao/plex/tools/blender:/inputs -v /Users/rindtorff/github/labdao/plex/tools/blender:/outputs nytimes/blender:3.3.1-cpu-ubuntu18.04 bash
# blender --background --python tools/blender/app.py -- tools/blender/5RGA.pdb tools/blender/test.png