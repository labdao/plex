'use client'

import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useRouter } from 'next/navigation'
import {
  AppDispatch,
  addFlowThunk,
  selectWalletAddress,
  selectFlowAddLoading,
  selectFlowAddKwargs,
  selectFlowAddSuccess,
  selectFlowAddError,
  selectFlowAddDataFiles,
  setFlowAddError,
  setFlowAddKwargs,
  setFlowAddLoading,
  setFlowAddTool,
  setFlowAddDataFiles,
  setFlowAddSuccess,
  dataFileListThunk,
  selectDataFileListError,
  selectDataFileList,
  selectFlowAddTool,
  toolListThunk,
  selectToolListError,
  selectToolList,
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
  const dispatch = useDispatch<AppDispatch>()
  const router = useRouter()

  const walletAddress = useSelector(selectWalletAddress)
  const loading = useSelector(selectFlowAddLoading)
  const error = useSelector(selectFlowAddError)
  const kwargs = useSelector(selectFlowAddKwargs)
  const success = useSelector(selectFlowAddSuccess)
  const selectedTool = useSelector(selectFlowAddTool)
  const toolListError = useSelector(selectToolListError)
  const dataFileListError = useSelector(selectDataFileListError)
  const dataFiles = useSelector(selectDataFileList)
  const tools = useSelector(selectToolList)
  console.log(selectedTool)

  useEffect(() => {
    if (success) {
      dispatch(setFlowAddSuccess(false))
      dispatch(setFlowAddKwargs({}))
      dispatch(setFlowAddTool({ CID: '', WalletAddress: '', Name: '', ToolJson: { inputs: {} }}))
      dispatch(setFlowAddError(null))
      router.push('/flows/list')
      return
    }
    dispatch(toolListThunk())
    dispatch(dataFileListThunk())
  }, [success, dispatch])

  const handleToolChange = (event: any) => {
    dispatch(setFlowAddTool(tools[parseInt(event.target.value)]))
  }

  const handleKwargsChange = (event: any) => {
    console.log(event)
    dispatch(setFlowAddSelectedDataFiles({
      ...selectedDataFiles,
      [event.target.name]: event.target.value
    }))
  }

  const isValidForm = (): boolean => {
    if (selectedTool.CID === '') return false;
    for (const key in selectedTool.ToolJson.inputs) {
      if (!selectedDataFiles[key]) return false;
    }
    return true;
  };

  const handleSubmit = async (event: any) => {
    event.preventDefault()
    console.log('Submitting flow')
    dispatch(setFlowAddLoading(true))
    dispatch(setFlowAddError(null))
    await dispatch(addFlowThunk({
      walletAddress,
      toolCid: selectedTool.CID,
      scatteringMethod: "dotproduct",
      kwargs: selectedDataFiles,
    }))
    dispatch(setFlowAddLoading(false))
  }

  return (
    <form onSubmit={handleSubmit}>
      <Box maxWidth={500} margin="auto"> {/* Constrain max width and center */}
        <Grid container direction="column" spacing={2}>
          <Grid item>
            <InputLabel>Tool</InputLabel>
            <Select
              value={selectedTool}
              label="Tool"
              onChange={e => handleToolChange(e)}
            >
              {tools.map((tool, index) => {
                return (
                  <MenuItem key={index} value={index}>{tool.Name}</MenuItem>
                )
              })}
            </Select>
          </Grid>
          {Object.keys(selectedTool.ToolJson.inputs).map(key => {
            return (
              <Grid item key={key}>
                <InputLabel>{key}</InputLabel>
                <Select
                  value={key}
                  onChange={e => handleKwargsChange(e)}
                >
                  {dataFiles.map(dataFile => (
                    <MenuItem key={dataFile.cid} value={dataFile.cid + "/" + dataFile.filename}>{dataFile.filename}</MenuItem>
                  ))}
                </Select>
              </Grid>
            )
          })}
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
