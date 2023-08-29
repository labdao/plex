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

  const [name, setName] = useState('');

  const handleNameChange = (e: any) => {
    setName(e.target.value);
  };

  const handleSubmit = async (e: any) => {
    e.preventDefault();
    
    const toolConfig = {
      "name": name,
      "description": "Sample description",
      "author": "Sample author",
      "baseCommand": ["/bin/sample", "-c"],
      "arguments": ["arg1", "arg2"],
      "dockerPull": "sample/docker:image",
      "gpuBool": false,
      "networkBool": false,
      "inputs": {
        "protein": {
          "type": "File",
          "glob": ["*.sample"]
        },
        "small_molecule": {
          "type": "File",
          "glob": ["*.sample2"]
        }
      },
      "outputs": {
        "best_docked_small_molecule": {
          "type": "File",
          "glob": ["*_output.sample"]
        },
        "protein": {
          "type": "File",
          "glob": ["*.sample"]
        }
      }
    };

    dispatch(addToolAsync({ toolData: toolConfig }));
    // router.push('/tool/list');
  };

  return (
    <form onSubmit={handleSubmit}>
      <Box maxWidth={500} margin="auto">
        <Grid container direction-="column" spacing={2}>
          <Grid item>
            <TextField
              fullWidth
              label="Tool Name"
              variant="outlined"
              value={name}
              onChange={handleNameChange}
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
