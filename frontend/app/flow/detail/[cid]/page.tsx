'use client'

import React, { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux';
import {
  AppDispatch,
  flowDetailThunk,
  flowPatchDetailThunk,
  selectFlowDetailLoading,
  selectFlowDetailError,
  selectFlowDetail,
} from '@/lib/redux'
import Link from 'next/link'
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


export default function FlowDetail() {
  const dispatch = useDispatch<AppDispatch>()

  const flow = useSelector(selectFlowDetail)
  const loading = useSelector(selectFlowDetailLoading)
  const error = useSelector(selectFlowDetailError)

  useEffect(() => {
    const flowCid = window.location.href.split('/').pop()
    if (flowCid) {
      dispatch(flowDetailThunk(flowCid))
    }
  }, [dispatch])

  return (
    <Box maxWidth={800} margin="auto">
      <Typography gutterBottom>
        <a href={`${process.env.NEXT_PUBLIC_GATEWAY_ENDPOINT}${flow.CID}/`}>
          Name: {flow.Name}
        </a>
      </Typography>
      <Typography gutterBottom>
        <a href={`${process.env.NEXT_PUBLIC_GATEWAY_ENDPOINT}${flow.CID}/`}>
          CID: {flow.CID}
        </a>
      </Typography>
      <Typography gutterBottom>
        Wallet Address: {flow.WalletAddress}
      </Typography>
      <Button
        variant="contained"
        color="primary"
        onClick={() => dispatch(flowPatchDetailThunk(flow.CID))}
        disabled={loading}
      >
        {loading ? "Submitting..." : "Update"}
      </Button>
      <Typography gutterBottom>
        Jobs:
      </Typography>
      {error && (
        <Box my={2}>
          <Alert severity="error" variant="filled">
            <Typography align="center">{error}</Typography>
          </Alert>
        </Box>
      )}
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Bacalhau Id</TableCell>
              <TableCell>Tool Name</TableCell>
              <TableCell>Tool CID</TableCell>
              <TableCell>State</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {flow.Jobs.map((job, index) => (
              <TableRow key={index}>
                <TableCell>
                  <Link href={`/job/detail/${job.BacalhauJobID}`}>
                    {job.BacalhauJobID}
                  </Link>
                </TableCell>
                <TableCell>{job.Tool.Name}</TableCell>
                <TableCell>
                  <a href={`${process.env.NEXT_PUBLIC_GATEWAY_ENDPOINT}${job.Tool.CID}/`}>
                    {job.Tool.CID}
                  </a>
                </TableCell>
                <TableCell>{job.State}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  )
}
