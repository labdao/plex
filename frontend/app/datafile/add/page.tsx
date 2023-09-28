'use client'

import React, { useState } from 'react'
import {
   useSelector,
   useDispatch,
   setError,
   startLoading,
   endLoading,
   selectCID,
   selectDataFileError,
   selectDataFileIsLoading,
   saveDataFileAsync,
   selectWalletAddress
} from '@/lib/redux'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Button from '@mui/material/Button'
import Alert from '@mui/material/Alert'
import Typography from '@mui/material/Typography'
import { useRouter } from 'next/navigation'

export default function DataFileForm() {
  const dispatch = useDispatch()

  const router = useRouter()
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

  const handleSuccess = () => {
    router.push('/datafile/list')
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (file === null) {
      dispatch(setError("Please select a file"))
      return
    }

    dispatch(startLoading())
    dispatch(setError(null))
    const metadata = { walletAddress };

    try {
      await dispatch(saveDataFileAsync({ file, metadata, handleSuccess }))
      dispatch(endLoading())
    } catch (error) {
      dispatch(setError("Error uploading file"))
      dispatch(endLoading())
    }
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
        </Grid>
      </Box>
    </form>
  )
}
