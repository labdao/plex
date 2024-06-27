import os
import uuid
import random
import bpy
from mathutils import Vector
import molecularnodes as mn
import math
from base_classes import ProteinBinderTargetStructure
from pydantic import FilePath
mn.register()

class MolecularNodesRender(ProteinBinderTargetStructure):
    png: FilePath

    def to_dict(self):
        return {
            "png": str(self.png),
            "pdb": str(self.pdb) # this one is inherited from the ProteinBinderTargetStructure class
        }

class ProteinComplexRender():
    def __init__(self, output_directory=None):
        if output_directory is None:
            output_directory = os.getcwd()  # Use current working directory if none provided
        self.output_directory = output_directory 
        self.style_cartoon = "cartoon"
        self.style_surface = "surface"
        self.new_style = "MN Default Copy"
        self.metallic_value = 0.8
        self.alpha_value = 0.07
        self.roughness_value = 0
        self.resolution_x = 2000
        self.resolution_y = 2000
        self.file_format = "PNG"    

    def _set_background_color(self, color=(0.0, 0.0, 0.0, 1.0)):
        """
        Change the background colour of the outputted .png image    
        """
        bpy.data.worlds["World"].node_tree.nodes["Background"].inputs[0].default_value = color

    def _clear_scene(self):
        """
        Clear the scene of any objects (always necessary for new bleder world)
        """
        bpy.ops.object.select_all(action="DESELECT")
        bpy.ops.object.select_by_type(type="MESH")
        bpy.ops.object.delete()

    def _import_pdb_file(self, style, file_path, object_name):
        """
        Import the .pdb file. Future directions is to allow other types of files. 
        It would also be very easy to use Molecular Nodes built in import to pull from a database he uses if we ever wanted to do that
        """
        bpy.context.scene.MN_import_local_path = file_path
        bpy.context.scene.MN_import_local_name = object_name
        bpy.context.scene.MN_import_style = style
        bpy.ops.mn.import_protein_local()

    def _duplicate_and_modify_material(self, original_material_name, new_material_name, metallic_value, alpha_value, roughness_value):
        """
        Here we duplicate and modify the material of that is used for NewMolecule and NewMoleculeSurface
        This new transparent, glossy and metalic material will be applied to NewMoleculeSurface later
        Get the original material by its name
        """
        original_material = bpy.data.materials.get(original_material_name)
        duplicated_material = original_material.copy()# Duplicate the material
        duplicated_material.name = new_material_name # Rename the duplicated material
        duplicated_material.use_nodes = True # Ensure the duplicated material uses nodes
        nodes = duplicated_material.node_tree.nodes # Find the Principled BSDF node
        principled_bsdf = next((node for node in nodes if node.type == 'BSDF_PRINCIPLED'), None)
        # Set Principled BSDF node Metallic, roughness and Alpha(transparency) values
        principled_bsdf.inputs['Metallic'].default_value = metallic_value
        principled_bsdf.inputs['Alpha'].default_value = alpha_value
        principled_bsdf.inputs['Roughness'].default_value = roughness_value
        duplicated_material.blend_method = 'BLEND' # Set the blend mode to 'Alpha blend'
        duplicated_material.shadow_method = 'HASHED' # Or 'NONE' for no shadows from transparent areas

    def _set_seed_for_color_attribute_random_nodes(self, obj_name, modifier_name, seed_value):
        """
        Set the seed value for all 'MN_color_attribute_random' nodes in the specified object's Geometry Nodes modifier.
        Random colours are picked from a list of colours that can be found in material.py in the molecular nodes files.
        https://github.com/BradyAJohnston/MolecularNodes/releases/ 
        A future direction could be modifying these colours or adding to them.
        """
        obj = bpy.data.objects.get(obj_name)
        nodes_modifier = obj.modifiers.get(modifier_name)
        node_tree = nodes_modifier.node_group
        for node in node_tree.nodes:
            if 'MN_color_attribute_random' in node.name and 'Seed' in node.inputs:
                node.inputs['Seed'].default_value = seed_value
                print(f"Updated Seed value for {node.name} in '{obj_name}' to {seed_value}")

    def _update_seed_values(self,seed_value):
        """
        Update seed values for MN_color_attribute_random nodes for both objects.
        """
        self._set_seed_for_color_attribute_random_nodes("NewMolecule", "MolecularNodes", seed_value)
        self._set_seed_for_color_attribute_random_nodes("NewMoleculeSurface", "MolecularNodes", seed_value)

    def _apply_material_to_node(self, node_tree_name, node_name, material_name):
        """
        Apply a specified material to a specific node in a given node tree.
        """
        new_material = bpy.data.materials.get(material_name)# Get the new material
        node_tree = bpy.data.node_groups.get(node_tree_name) # Find the node tree
        style_surface_node = node_tree.nodes.get(node_name)# Find the specific node
        material_input = style_surface_node.inputs.get('Material') # Find the input socket for the material
        material_input.default_value = new_material    # Assign the material object directly to the socket's default_value

    def _calculate_bounding_box(self, obj_name):
        """
        Calculate the bounding box center and size for a given object.
        """
        obj = bpy.data.objects.get(obj_name)
        bbox_corners = [obj.matrix_world @ Vector(corner) for corner in obj.bound_box]
        bbox_center = sum(bbox_corners, Vector()) / 8
        max_dimension = max((max(bbox_corners, key=lambda c: c[i])[i] - min(bbox_corners, key=lambda c: c[i])[i] for i in range(3)))
        return bbox_center, max_dimension, bbox_corners

    def _set_camera_position(self, camera_name, bbox_center, max_dimension):
        """
        Set the camera position based on the bounding box center and size.
        """
        camera = bpy.data.objects[camera_name]
        camera_distance = max_dimension * 2.0# Calculate the distance needed to fit the bounding box within the camera's view
        # (This is still being tweaked a bit)
        camera.location = bbox_center + Vector((0, -camera_distance, camera_distance / 2))
        camera.data.clip_start = 0.1 # Can also adjust the camera's clipping distances
        camera.data.clip_end = camera_distance * 3
        direction = bbox_center - camera.location# Point the camera towards the bounding box center
        rot_quat = direction.to_track_quat('-Z', 'Y')# Convert the direction to a rotation (quaternion) that points the camera's '-Z' axis towards the molecule
        camera.rotation_euler = rot_quat.to_euler()# Assign the quaternion to the camera's rotation_euler

    def _set_light_position(self, light_name, bbox_center, max_dimension):
        """
        Set the light position and orientation based on the bounding box center and size.
        """
        light = bpy.data.objects[light_name]
        light.data.energy = 100 # Increase the light's energy (brightness)
        light.data.type = 'AREA' # Change the light type to 'AREA'
        phi = math.pi / 4 # Set a fixed polar angle, e.g. 45 degrees elevation
        theta = random.uniform(0, 2 * math.pi) # Generate a random azimuthal angle
        # Calculate the light's new position using spherical to Cartesian coordinates conversion
        light_distance = max_dimension * 2.0 # distance can be adjusted as needed
        x = light_distance * math.sin(phi) * math.cos(theta)
        y = light_distance * math.sin(phi) * math.sin(theta)
        z = light_distance * math.cos(phi)
        new_light_location = bbox_center + Vector((x, y, z))
        light.location = new_light_location
        direction = bbox_center - light.location # Calculate the direction from the light to the molecule's center
        rot_quat = direction.to_track_quat('-Z', 'Y') # Convert the direction to a rotation (quaternion)
        light.rotation_euler = rot_quat.to_euler() # Convert the quaternion to Euler angles since object rotation in Blender is typically represented in Euler angles

    def _create_plane_below_object(self, obj_name, bbox_corners):
        """
        Create a plane below the molecule to add something underneath it (nice shadow)
        Can adjust the offset, 0.009 currently
        """
        bpy.ops.mesh.primitive_plane_add(size=1, enter_editmode=False, align='WORLD', location=(0, 0, 0))
        plane = bpy.context.active_object  # Reference plane
        plane.scale = (10, 10, 1) # Scale the plane by 10 in the X and Y axes
        molecule = bpy.data.objects.get(obj_name)
        min_z = min(corner.z for corner in bbox_corners) # Find the lowest Z value among the bounding box corners
        plane.location = (molecule.location.x, molecule.location.y, min_z - 0.009) # Set the location of the plane to be just below the lowest point of the bounding box
        black_material = bpy.data.materials.new(name="BlackMaterial") # Create a new dark material 
        black_material.diffuse_color = (0.01, 0.01, 0.01, 1) # RGB and Alpha
        plane.data.materials.append(black_material) # Assign the material to the plane
    
    def check_gpu_availability(self):

        bpy.context.scene.render.engine = "CYCLES"
        bpy.context.preferences.addons["cycles"].preferences.compute_device_type = "CUDA"
        # Force refresh of available devices
        bpy.context.preferences.addons["cycles"].preferences.get_devices()

        # Explicitly enable GPU devices
        has_gpu = False
        for device in bpy.context.preferences.addons["cycles"].preferences.devices:
            if device.type == 'CUDA':  # Check if the device is a GPU
                device.use = True
                has_gpu = True
                print(f"Enabled CUDA Device: {device.name}")

        bpy.context.scene.cycles.device = 'GPU'
        
        if not has_gpu:
            print("No CUDA devices were found or enabled! Fallback to BLENDER_EEVEE")
            bpy.context.scene.render.engine = "BLENDER_EEVEE"
            bpy.context.preferences.addons["cycles"].preferences.compute_device_type = 'NONE'  # Disable GPU compute type
            bpy.context.scene.cycles.device = 'CPU'

    def _render_image(self, path, resolution_x, resolution_y, file_format):
        """
        These two lines set the resolution of the image, 
        they can be adjusted to change the output width and height too.
        """
        # bpy.context.scene.render.engine = 'CYCLES' #This is a more powerful rendering engine, it takes much longer
        self.check_gpu_availability()
        bpy.context.scene.render.resolution_x = resolution_x
        bpy.context.scene.render.resolution_y = resolution_y
        bpy.context.scene.render.image_settings.file_format = file_format
        bpy.context.scene.render.filepath = path
        bpy.ops.render.render(write_still=True)
        bpy.data.images["Render Result"].save_render(filepath=bpy.context.scene.render.filepath)

    def visualise(self, protein_binder_target_structure: ProteinBinderTargetStructure) -> MolecularNodesRender:
        """
        Render a molecule from a .pdb file and save the result as an image.
        """
        pdb_file_path = str(protein_binder_target_structure.pdb)
        pdb_file_name, extension = os.path.splitext(os.path.basename(pdb_file_path))
        output_image_filename = f"{pdb_file_name}.png"
        output_image_path = os.path.join(self.output_directory, output_image_filename)

        self._set_background_color()
        self._clear_scene()
        self._import_pdb_file(self.style_cartoon, pdb_file_path, "NewMolecule")
        self._import_pdb_file(self.style_surface, pdb_file_path, "NewMoleculeSurface")
        self._duplicate_and_modify_material("MN Default", self.new_style, self.metallic_value, self.alpha_value, self.roughness_value)
        seed_value = random.randint(1, 20)
        self._update_seed_values(seed_value)
        self._apply_material_to_node("MN_NewMoleculeSurface", "MN_style_surface", self.new_style)
        bbox_center, max_dimension, bbox_corners = self._calculate_bounding_box("NewMolecule")
        self._set_camera_position("Camera", bbox_center, max_dimension)
        self._set_light_position("Light", bbox_center, max_dimension)
        self._create_plane_below_object("NewMolecule", bbox_corners)
        self._render_image(output_image_path, self.resolution_x, self.resolution_y, self.file_format)
         # Create MolecularNodesRender instance with additional data
        visualisation = MolecularNodesRender(
            png=FilePath(output_image_path),
            pdb=protein_binder_target_structure.pdb,
            binder_sequence=protein_binder_target_structure.binder_sequence,
            target_sequence=protein_binder_target_structure.target_sequence
        )
        return visualisation


def visualise_protein_complex(pdb_file_path: str, output_directory: str):
    protein_structure = ProteinBinderTargetStructure(
        binder_sequence="ACDEFGHIKLMNPQRSTVWY",
        target_sequence="ACDEFGHIKLMNPQRSTVWY",
        pdb=FilePath(pdb_file_path)
    )
    output_directory = os.path.join(os.path.dirname(output_directory),"visualizations")
    renderer = ProteinComplexRender(output_directory)
    result = renderer.visualise(protein_structure)
    print(result.to_dict())
    return result.to_dict()