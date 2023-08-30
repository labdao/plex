'use client'

import React, { useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { addToolAsync, setError, startFileUpload, endFileUpload  } from '@/lib/redux'
import {
  selectToolError,
  selectToolIsLoading,
  selectToolIsUploaded
} from '@/lib/redux/slices/toolAddSlice/selectors'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
// import { useRouter } from 'next/router'

export default function AddTool() {
  const dispatch = useDispatch()
  // const router = useRouter()

  const isLoading = useSelector(selectToolIsLoading);
  const error = useSelector(selectToolError);
  const isUploaded = useSelector(selectToolIsUploaded);

  const sampleToolConfig = {
    "class": "CommandLineTool",
    "name": "equibind",
    "description": "Docking of small molecules to a protein",
    "author": "@misc{stärk2022equibind,\n      title={EquiBind: Geometric Deep Learning for Drug Binding Structure Prediction}, \n      author={Hannes Stärk and Octavian-Eugen Ganea and Lagnajit Pattanaik and Regina Barzilay and Tommi Jaakkola},\n      year={2022},\n      eprint={2202.05146},\n      archivePrefix={arXiv},\n      primaryClass={q-bio.BM}\n}",
    "baseCommand": ["/bin/bash", "-c"],
    "arguments": [
      "mkdir -p /tmp-inputs/tmp;",
      "mkdir -p /tmp-outputs/tmp;",
      "cp /inputs/* /tmp-inputs/tmp/;",
      "ls /tmp-inputs/tmp;",
      "cd /src && python /src/inference.py --config=/src/configs_clean/bacalhau.yml;",
      "mv /tmp-outputs/tmp/* /outputs/;",
      "mv /outputs/lig_equibind_corrected.sdf /outputs/$(inputs.protein.basename)_$(inputs.small_molecule.basename)_docked.$(inputs.small_molecule.ext);",
      "mv /tmp-inputs/tmp/*.pdb /outputs/;"],
    "dockerPull": "ghcr.io/labdao/equibind:main@sha256:21a381d9ab1ff047565685044569c8536a55e489c9531326498b28d6b3cc244f",
    "gpuBool": false,
    "networkBool": false,
    "inputs": {
      "protein": {
        "type": "File",
        "item": "",
        "glob": ["*.pdb"]
      },
      "small_molecule": {
        "type": "File",
        "item": "",
        "glob": ["*.sdf", "*.mol2"]
      }
    },
    "outputs": {
      "best_docked_small_molecule": {
        "type": "File",
        "item": "",
        "glob": ["*_docked.sdf", "*_docked.mol2"]
      },
      "protein": {
        "type": "File", 
        "item": "",
        "glob": ["*.pdb"]
      }
    }
  }

  const [toolClass, setToolClass] = useState(sampleToolConfig.class)
  const [name, setName] = useState(sampleToolConfig.name);
  const [description, setDescription] = useState(sampleToolConfig.description);
  const [author, setAuthor] = useState(sampleToolConfig.author);
  const [baseCommand, setBaseCommand] = useState(sampleToolConfig.baseCommand);
  const [toolArguments, setToolArguments] = useState(sampleToolConfig.arguments);
  const [dockerPull, setDockerPull] = useState(sampleToolConfig.dockerPull);
  const [gpuBool, setGpuBool] = useState(sampleToolConfig.gpuBool);
  const [networkBool, setNetworkBool] = useState(sampleToolConfig.networkBool);
  const [inputs, setInputs] = useState(sampleToolConfig.inputs);
  const [outputs, setOutputs] = useState(sampleToolConfig.outputs);
    
  const handleToolClassChange = (e: any) => {
    setToolClass(e.target.value);
  };

  const handleNameChange = (e: any) => {
    setName(e.target.value);
  };

  const handleDescriptionChange = (e: any) => {
    setDescription(e.target.value);
  };

  const handleAuthorChange = (e: any) => {
    setAuthor(e.target.value);
  };

  const handleBaseCommandChange = (e: any) => {
    setBaseCommand(e.target.value);
  };

  const handleToolArgumentsChange = (e: any) => {
    setToolArguments(e.target.value);
  };

  const handleDockerPullChange = (e: any) => {
    setDockerPull(e.target.value);
  };

  const handleGpuBoolChange = (e: any) => {
    setGpuBool(e.target.value);
  };

  const handleNetworkBoolChange = (e: any) => {
    setNetworkBool(e.target.value);
  };

  const handleInputsChange = (e: any) => {
    setInputs(e.target.value);
  };

  const handleOutputsChange = (e: any) => {
    setOutputs(e.target.value);
  };
  
  const handleSubmit = async (e: any) => {
    e.preventDefault();
    
    const toolConfig = {
      "name": name,
      "description": description,
      "author": author,
      "baseCommand": baseCommand,
      "arguments": toolArguments,
      "dockerPull": dockerPull,
      "gpuBool": gpuBool,
      "networkBool": networkBool,
      "inputs": inputs,
      "outputs": outputs
    };

    dispatch(addToolAsync({ toolData: toolConfig }));
    // router.push('/tool/list');
  };

  return (
    <form onSubmit={handleSubmit}>
      <Box maxWidth={500} margin="auto">
        <Grid container direction="column" spacing={2}>
          <Grid item>
            <TextField
              fullWidth
              label="Tool Class"
              variant="outlined"
              value={toolClass}
              onChange={handleToolClassChange}
            />
          </Grid>
          <Grid item>
            <TextField
              fullWidth
              label="Tool Name"
              variant="outlined"
              value={name}
              onChange={handleNameChange}
            />
          </Grid>
          <Grid item>
            <TextField
              fullWidth
              label="Description"
              variant="outlined"
              value={description}
              onChange={handleDescriptionChange}
            />
          </Grid>
          <Grid item>
            <TextField
              fullWidth
              label="Author"
              variant="outlined"
              value={author}
              onChange={handleAuthorChange}
            />
          </Grid>
          <Grid item>
            <TextField
              fullWidth
              label="Base Command"
              variant="outlined"
              value={baseCommand}
              onChange={handleBaseCommandChange}
            />
          </Grid>
          <Grid item>
            <TextField
              fullWidth
              label="Arguments"
              variant="outlined"
              value={toolArguments}
              onChange={handleToolArgumentsChange}
            />
          </Grid>
          <Grid item>
            <TextField
              fullWidth
              label="Docker Pull"
              variant="outlined"
              value={dockerPull}
              onChange={handleDockerPullChange}
            />
          </Grid>
          <Grid item>
            <TextField
              fullWidth
              label="GPU Bool"
              variant="outlined"
              value={gpuBool}
              onChange={handleGpuBoolChange}
            />
          </Grid>
          <Grid item>
            <TextField
              fullWidth
              label="Network Bool"
              variant="outlined"
              value={networkBool}
              onChange={handleNetworkBoolChange}
            />
          </Grid>
          <Grid item>
            <TextField
              fullWidth
              label="Inputs"
              variant="outlined"
              value={inputs}
              onChange={handleInputsChange}
            />
          </Grid>
          <Grid item>
            <TextField
              fullWidth
              label="Outputs"
              variant="outlined"
              value={outputs}
              onChange={handleOutputsChange}
            />
          </Grid>
          <Grid item container justifyContent="center">
            <Button variant="contained" color="primary" type="submit">
              {isLoading ? "Submitting..." : "Submit"}
            </Button>
          </Grid>
        </Grid>
      </Box>
    </form>
  )
}
