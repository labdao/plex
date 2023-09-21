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
// import { useRouter } from 'next/router'

export default function DataFileForm() {
  const dispatch = useDispatch()
  // const router = useRouter()

  const cid = useSelector(selectCID)
  const errorMessage = useSelector(selectDataFileError)
  const isLoading = useSelector(selectDataFileIsLoading)
  const walletAddress = useSelector(selectWalletAddress)

  const [file, setFile] = useState<File | null>(null)
  const [isPublic, setIsPublic] = useState<boolean>(true)
  const [isVisible, setIsVisible] = useState<boolean>(true)

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const uploadedFile = e.target.files && e.target.files[0]
    if (uploadedFile) {
      setFile(uploadedFile)
    }
  }

  const handlePublicChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setIsPublic(!e.target.checked)
  }

  const handleVisibleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setIsVisible(!e.target.checked)
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (file === null) {
      dispatch(setError("Please select a file"))
      return
    }
    
    dispatch(startLoading())
    dispatch(setError(null))
    const metadata = { walletAddress, isPublic, isVisible };

    try {
      await dispatch(saveDataFileAsync({ file, metadata }))
      dispatch(endLoading())

      // router.push('/data/list')
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
          <Grid item container justifyContent="center">
            <label>
              <input
                type="checkbox"
                checked={!isPublic}
                onChange={handlePublicChange}
              />
              File should be private
            </label>
            <label>
              <input
                type="checkbox"
                checked={!isVisible}
                onChange={handleVisibleChange}
              />
              File should be hidden
            </label>
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
