'use client'

import React, { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux';
import {
  AppDispatch,
  jobDetailThunk,
  jobPatchDetailThunk,
  selectJobDetailLoading,
  selectJobDetailError,
  selectJobDetail,
} from '@/lib/redux'
import Alert from '@mui/material/Alert'
import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import Typography from '@mui/material/Typography'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'


export default function JobDetail() {
  const dispatch = useDispatch<AppDispatch>()

  const job = useSelector(selectJobDetail)
  const loading = useSelector(selectJobDetailLoading)
  const error = useSelector(selectJobDetailError)

  useEffect(() => {
    const jobBacalhauID = window.location.href.split('/').pop()
    if (jobBacalhauID) {
      dispatch(jobDetailThunk(jobBacalhauID))
    }
  }, [dispatch])

  return (
    <Box maxWidth={800} margin="auto">
      <Typography gutterBottom>
        Bacalhau ID: {job.BacalhauJobID}
      </Typography>
      <Typography gutterBottom>
        State: {job.State}
      </Typography>
      <Typography gutterBottom>
        Error: {job.Error}
      </Typography>
      <Typography gutterBottom>
        Tool CID: {job.ToolID}
      </Typography>
      <Typography gutterBottom>
        Flow Initial CID: {job.FlowID}
      </Typography>
      <Button
        variant="contained"
        color="primary"
        onClick={() => dispatch(jobPatchDetailThunk(job.BacalhauJobID))}
        disabled={loading}
      >
        {loading ? "Submitting..." : "Update"}
      </Button>
      {error && (
        <Box my={2}>
          <Alert severity="error" variant="filled">
            <Typography align="center">{error}</Typography>
          </Alert>
        </Box>
      )}
    </Box>
  )
}
