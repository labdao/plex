'use client'

import React, { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useRouter } from 'next/navigation'
import {
  AppDispatch,
  setAddToolError,
  setAddToolLoading,
  setAddToolJson,
  setAddToolSuccess,
  selectWalletAddress,
  selectAddToolLoading,
  selectAddToolError,
  selectAddToolJson,
  selectAddToolSuccess,
  createToolThunk,
  toolListThunk,
  selectToolListError,
  dataFileListThunk,
  selectDataFileListError,
} from '@/lib/redux'
import Button from '@mui/material/Button'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Alert from '@mui/material/Alert'
import Typography from '@mui/material/Typography'
import { JsonInput } from '@mantine/core'
import { MantineProvider } from '@mantine/core';



export default function AddTool() {
  const dispatch = useDispatch<AppDispatch>()
  const router = useRouter()

  const walletAddress = useSelector(selectWalletAddress)
  const loading = useSelector(selectAddToolLoading)
  const error = useSelector(selectAddToolError)
  const toolJson = useSelector(selectAddToolJson)
  const toolSuccess = useSelector(selectAddToolSuccess)

  const toolListError = useSelector(selectToolListError)
  const dataFileListError = useSelector(selectDataFileListError)

  useEffect(() => {
    if (toolSuccess) {
      dispatch(setAddToolSuccess(false))
      dispatch(setAddToolJson(""))
      router.push('/tool/list')
      return
    }
    dispatch(toolListThunk())
    dispatch(dataFileListThunk())
  }, [toolSuccess, dispatch])

  const handleToolJsonChange = (toolJsonInput: string) => {
    dispatch(setAddToolJson(toolJsonInput))
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log("Submitting tool.json: ", toolJson);
    dispatch(setAddToolLoading(true))
    dispatch(setAddToolError(""))
    try {
      const toolJsonParsed = JSON.parse(toolJson)
      await dispatch(createToolThunk({ walletAddress, toolJson: toolJsonParsed }))
    } catch (error) {
      console.error("Error creating tool", error)
      dispatch(setAddToolError("Error creating tool"))
    }
    dispatch(setAddToolLoading(false))
  }

  return (
    <form onSubmit={handleSubmit}>
      <Box maxWidth={500} margin="auto">
        <Grid container direction="column" spacing={2} justifyContent="center" alignItems="center">
          <Grid style={{ width: '100%' }} item>
            <MantineProvider>
              <JsonInput
                label="Tool Definition"
                placeholder="Paste your tool's JSON definition here."
                validationError="Invalid JSON"
                autosize
                minRows={10}
                value={toolJson}
                onChange={handleToolJsonChange}
                styles={{
                  input: { 'width': '100%' },
                }}
              />
            </MantineProvider>
          </Grid>
          {error && (
            <Box my={2}>
              <Alert severity="error" variant="filled">
                <Typography align="center">{error}</Typography>
              </Alert>
            </Box>
          )}
          <Grid item container justifyContent="center">
            <Button variant="contained" color="primary" type="submit">
              {loading ? "Submitting..." : "Submit"}
            </Button>
          </Grid>
        </Grid>
      </Box>
    </form>
  )
}
