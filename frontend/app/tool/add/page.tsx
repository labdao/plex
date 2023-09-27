'use client'

import React, { useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { addToolAsync, selectWalletAddress } from '@/lib/redux'
import Link from '@mui/material/Link';
import {
  selectToolError,
  selectToolIsLoading,
  selectToolIsUploaded
} from '@/lib/redux/slices/toolAddSlice/selectors'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Checkbox from '@mui/material/Checkbox'
import FormControlLabel from '@mui/material/FormControlLabel'


export default function AddTool() {
  const dispatch = useDispatch()

  const isLoading = useSelector(selectToolIsLoading);
  const error = useSelector(selectToolError);
  const isUploaded = useSelector(selectToolIsUploaded);
  const walletAddress = useSelector(selectWalletAddress)

  const [name, setName] = useState("")
  const [description, setDescription] = useState("")
  const [author, setAuthor] = useState("")
  const [colabNotebook, setColabNotebook] = useState("")
  const [gpuBool, setGpuBool] = useState(false)
  const [inputs, setInputs] = useState([""])
  const [outputs, setOutputs] = useState([""])

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  };

  const handleDescriptionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setDescription(e.target.value);
  };

  const handleColabNotebookChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setColabNotebook(e.target.value);
  };

  const handleGpuBoolChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setGpuBool(e.target.checked)
  };

  const handleInputsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    try {
      const value = JSON.parse(e.target.value);
      setInputs(value);
    } catch (error) {
      console.error(error);
    }
  };

  const handleOutputsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    try {
      const value = JSON.parse(e.target.value);
      setOutputs(value);
    } catch (error) {
      console.error(error);
    }
  };
  
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    console.log("Submitting tool config")
    console.log("Wallet address: ", walletAddress)

    const toolConfig = {
      "name": name,
      "description": description,
      "author": author,
      "gpuBool": gpuBool,
      "inputs": inputs,
      "outputs": outputs
    };

    // @ts-ignore
    dispatch(addToolAsync({ toolData: toolConfig, walletAddress }));
  };

  return (
    <form onSubmit={handleSubmit}>
      <Box maxWidth={500} margin="auto">
        <Grid container direction="column" spacing={2}>
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
    label="Colab Notebook"
    variant="outlined"
    value={colabNotebook}
    onChange={handleColabNotebookChange}
    helperText={
      <span>
        Link to a Google Colab Notebook that follows the plex schema. See the{' '}
        <Link href="https://your-link.com" target="_blank" rel="noopener">
          Template Tool Notebook
        </Link> for more details.
      </span>
    }
  />
          </Grid>
          <Grid item>
            <FormControlLabel
              control={
                <Checkbox
                  checked={gpuBool}
                  onChange={handleGpuBoolChange}
                />
              }
              label="Require GPU"
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
