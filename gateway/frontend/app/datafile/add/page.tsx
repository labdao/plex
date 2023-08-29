'use client'

import { useEffect, useState } from 'react'

import {
   useSelector,
   useDispatch,
   setCid,  // assuming you have an action to set cid
   setError,
   startLoading,
   endLoading,
   selectCID,
   selectDataFileError,
   selectDataFileIsLoading,
   saveDataFileAsync,
   selectWalletAddress
} from '@/lib/redux'

import { useRouter } from 'next/router'  // changed 'next/navigation' to 'next/router'

import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Button from '@mui/material/Button'
import Alert from '@mui/material/Alert'
import Typography from '@mui/material/Typography'

export default function DataFileForm() {
  const dispatch = useDispatch()

  const cid = useSelector(selectCID)
  const errorMessage = useSelector(selectDataFileError)
  const isLoading = useSelector(selectDataFileIsLoading)
  const walletAddress = useSelector(selectWalletAddress)

  const [file, setFile] = useState<File | null>(null)

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const uploadedFile = e.target.files && e.target.files[0]
    if (uploadedFile) {
      setFile(uploadedFile)
    }
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    dispatch(startLoading())
    dispatch(setError(null))
    await dispatch(saveDataFileAsync({ file, metadata: { walletAddress }}))
    dispatch(endLoading())
  }

  return (
    <form onSubmit={handleSubmit}>
      <Box maxWidth={500} margin="auto"> {/* Constrain max width and center */}
        <Grid container direction="column" spacing={2}>
          <Grid item container justifyContent="center">
            {file && <Typography variant="subtitle1">{`Selected File: ${file.name}`}</Typography>}
          </Grid>
          <Grid item container justifyContent="center">
            <Button variant="outlined" component="label">
              Upload File
              <input
                type="file"
                hidden
                onChange={handleFileChange}
              />
            </Button>
          </Grid>
          {errorMessage && (
            <Box my={2}>
              <Alert severity="error" variant="filled">
                <Typography align="center">{errorMessage}</Typography>
              </Alert>
            </Box>
          )}
          <Grid item container justifyContent="center">
            <Button variant="contained" color="primary" type="submit">
              {isLoading ? "Submitting..." : "Submit"}
            </Button>
          </Grid>
          <Grid item container justifyContent="center">
            ToDo: Write Go App Gateway Endpoint
          </Grid>
        </Grid>
      </Box>
    </form>
  )
}
