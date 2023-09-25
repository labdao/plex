
# Run any Google Colab notebook with Plex

You can run any colab notebook with plex in headless mode. The following command will run the notebook in headless mode with a config file and generic input file. 


````
# Run the notebook in headless mode with a config file
./plex init -t tools/colab/colab-config.json -i '{"notebook": ["tools/colab/colab_plex.ipynb"], "config": ["tools/colab/config-1.yaml"], "generic_input": ["tools/colab/input.txt"]}' --scatteringMethod=dotProduct --autoRun=true
````

You can also scale the execution of the notebook across multiple config files. For example, if you have 3 config files, you can run the notebook in parallel with the following command:


````
./plex init -t tools/colab/colab-config.json -i '{"notebook": ["tools/colab/colab_plex.ipynb"], "config": ["tools/colab/config-1.yaml", "tools/colab/config-2.yaml", "tools/colab/config-3.yaml"], "generic_input": ["tools/colab/input.txt"]}' --scatteringMethod=crossProduct --autoRun=true --concurrency=3
````


