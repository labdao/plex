# script.py
import ray
import subprocess
# import pandas

@ray.remote
def hello_world():
    return "hello world"
    # return subprocess.run([""]) 

# Automatically connect to the running Ray cluster.
ray.init()
print(ray.get(hello_world.remote()))
