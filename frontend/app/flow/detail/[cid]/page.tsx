'use client'

import React, { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux';
import Router from 'next/router'
import {
  AppDispatch,
  flowDetailThunk,
  selectFlowDetail,
} from '@/lib/redux'
import Box from '@mui/material/Box'
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

  useEffect(() => {
    const flowCid = window.location.href.split('/').pop()
    if (flowCid) {
      dispatch(flowDetailThunk(flowCid))
    }
  }, [dispatch])

  return (
    <Box maxWidth={800} margin="auto">
      <Typography gutterBottom>
        <a href={`http://bacalhau.labdao.xyz:8080/ipfs/${flow.CID}/`}>
          Name: {flow.Name}
        </a>
      </Typography>
      <Typography gutterBottom>
        <a href={`http://bacalhau.labdao.xyz:8080/ipfs/${flow.CID}/`}>
          CID: {flow.CID}
        </a>
      </Typography>
      <Typography gutterBottom>
        Wallet Address: {flow.WalletAddress}
      </Typography>
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Bacalhau Id</TableCell>
              <TableCell>Tool Name</TableCell>
              <TableCell>Tool CID</TableCell>
              <TableCell>State</TableCell>
              <TableCell>Error</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {flow.Jobs.map((job, index) => (
              <TableRow key={index}>
                <TableCell>{job.BacalhauJobID}</TableCell>
                <TableCell>{job.Tool.Name}</TableCell>
                <TableCell>
                  <a href={`http://bacalhau.labdao.xyz:8080/ipfs/${job.Tool.CID}/`}>
                    {job.Tool.CID}
                  </a>
                </TableCell>
                <TableCell>{job.State}</TableCell>
                <TableCell>{job.Error}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  )
}
