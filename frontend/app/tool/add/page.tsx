'use client'

import React, { useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { addToolAsync, selectWalletAddress } from '@/lib/redux'
import {
  selectToolIsLoading,
} from '@/lib/redux/slices/toolAddSlice/selectors'
import Button from '@mui/material/Button'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import { JsonInput } from '@mantine/core'
import { MantineProvider } from '@mantine/core';


export default function AddTool() {
  const dispatch = useDispatch()

  const isLoading = useSelector(selectToolIsLoading);
  const walletAddress = useSelector(selectWalletAddress)

  const [toolJson, setToolJson] = useState("")

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log("Submitting tool.json: ", toolJson);

    try {
      const toolConfig = JSON.parse(toolJson)

      // @ts-ignore
      dispatch(addToolAsync({ toolData: toolConfig, walletAddress }))
    } catch (error) {
      console.error("Invalid JSON format", error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <Box maxWidth={500} margin="auto">
        <Grid container direction="column" spacing={2} justifyContent="center" alignItems="center">
          <Grid style={{ width: '100%' }} item>
            <MantineProvider>
              <JsonInput
                label="Tool Definition"
                placeholder="Paster your tool's JSON definition here."
                validationError="Invalid JSON"
                formatOnBlur
                autosize
                minRows={10}
                value={toolJson}
                onChange={setToolJson}
                styles={{
                  input: { 'width': '100%' },
                }}
              />
            </MantineProvider>
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
