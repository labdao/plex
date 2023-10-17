'use client'

import React, { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux';
import {
  AppDispatch,
  flowListThunk,
  selectFlowList,
} from '@/lib/redux'
import Link from 'next/link'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'


export default function ListToolFiles() {
  const dispatch = useDispatch<AppDispatch>()

  const flows = useSelector(selectFlowList)

  useEffect(() => {
    dispatch(flowListThunk())
  }, [dispatch])

  return (
    <TableContainer>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Name</TableCell>
            <TableCell>CID</TableCell>
            <TableCell>Uploader Wallet Address</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {flows.map((flow, index) => (
            <TableRow key={index}>
              <TableCell>
                <Link href={`/flow/detail/${flow.CID}`}>
                  {flow.Name}
                </Link>
              </TableCell>
              <TableCell>
                <a href={`http://bacalhau.labdao.xyz:8080/ipfs/${flow.CID}/`}>
                  {flow.CID}
                </a>
              </TableCell>
              <TableCell>{flow.WalletAddress}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}
