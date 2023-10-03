'use client'

import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
    selectWalletAddress,
  } from '@/lib/redux'

import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Button from '@mui/material/Button'
import TextField from '@mui/material/TextField'
import Alert from '@mui/material/Alert'
import InputLabel from '@mui/material/InputLabel'
import Select from '@mui/material/Select'
import MenuItem from '@mui/material/MenuItem'

import backendUrl from 'lib/backendUrl';


export default function AddGraph() {
  const dispatch = useDispatch();

  const walletAddress = useSelector(selectWalletAddress)

  interface Tool {
    Name: string;
    CID: string;
    ToolJson: {
      inputs: {
        [key: string]: {
          type: string;
          item: string;
          glob: string[];
        };
      };
    };
  }

  interface DataFile {
      ID: number;
      CID: string;
      WalletAddress: string;
      Filename: string;
      Timestamp: Date;
  }

  const [tools, setTools] = useState<Tool[]>([]);
  const [selectedToolIndex, setSelectedToolIndex] = useState('')
  const [dataFiles, setDataFiles] = useState<DataFile[]>([]);
  const [selectedDataFiles, setSelectedDataFiles] = useState<string[]>([]);
  const [loading, setLoading] = useState(false)

  useEffect(() => {
      fetch(`${backendUrl()}/get-tools`)
          .then(response => response.json())
          .then(data => setTools(data))
          .catch(error => console.error('Error fetching tools:', error));

      fetch(`${backendUrl()}/get-datafiles`)
          .then(response => response.json())
          .then(data => setDataFiles(data))
          .catch(error => console.error('Error fetching data files:', error));
  }, []);

  const isValidForm = (): boolean => {
    if (selectedToolIndex === '') return false;
  
    const selectedToolInputs = tools[parseInt(selectedToolIndex)].ToolJson.inputs;
    for (const key in selectedToolInputs) {
      if (!selectedDataFiles[key] || selectedDataFiles[key].length === 0) return false;
    }
  
    return true;
  };

  const handleSubmit = (event: any) => {
      event.preventDefault()
      setLoading(true)

      const data = {
          walletAddress: walletAddress,
          toolCid: tools[parseInt(selectedToolIndex)].CID,
          scatteringMethod: "dotProduct",
          kwargs: selectedDataFiles,
      };
      console.log('data:', data)

      fetch(`${backendUrl()}/graphs`, {
          method: 'POST',
          headers: {
              'Content-Type': 'application/json',
          },
          body: JSON.stringify(data),
      })
      .then(response => response.json())
      .then(data => console.log('Job initialized:', data))
      .catch((error) => console.error('Error initializing job:', error))
      setLoading(false)
  }

  return (
    <form onSubmit={handleSubmit}>
      <Box maxWidth={500} margin="auto"> {/* Constrain max width and center */}
        <Grid container direction="column" spacing={2}>
          <Grid item>
            <InputLabel>Tool</InputLabel>
            <Select
              value={selectedToolIndex}
              label="Tool"
              onChange={e => setSelectedToolIndex(e.target.value)}
            >
              {tools.map((tool, index) => {
                return (
                  <MenuItem key={index} value={index}>{tool.Name}</MenuItem>
                )
              })}
            </Select>
          </Grid>
          {selectedToolIndex !== '' && Object.entries(tools[parseInt(selectedToolIndex)].ToolJson.inputs).map(([key, input]) => (
            <Grid item key={key}>
              <InputLabel>{key}</InputLabel>
              <Select
                value={selectedDataFiles[key] || ''}
                onChange={e => setSelectedDataFiles({
                    ...selectedDataFiles,
                    [key]: [...(selectedDataFiles[key] || []), e.target.value]
                  })}
              >
                {dataFiles.map(file => (
                  <MenuItem key={file.CID} value={file.CID + "/" + file.Filename}>{file.Filename}</MenuItem>
                ))}
              </Select>
            </Grid>
          ))}
          <Grid item container justifyContent="center">
            <Button 
              variant="contained" 
              color="primary" 
              type="submit"
              disabled={loading || !isValidForm()}
            >
              {loading ? "Submitting..." : "Submit"}
            </Button>
          </Grid>
         </Grid>
      </Box>
    </form>
  )
}
