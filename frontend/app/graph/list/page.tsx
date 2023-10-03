'use client'

import React, { useEffect, useState } from 'react'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'

import backendUrl from 'lib/backendUrl'

export default function ListToolFiles() {
  interface Tool {
    CID: string;
    Name: string;
    WalletAddress: string;
  }

  const [tools, setTools] = useState<Tool[]>([]);

  useEffect(() => {
    fetch(`${backendUrl()}/graphs`)
      .then(response => {
        if (!response.ok) {
          throw new Error(`HTTP error ${response.status}`);
        }
        return response.json();
      })
      .then(data => {
        console.log('Fetched graphs:', data);
        setTools(data);
      })
      .catch(error => {
        console.error('Error fetching graphs:', error);
      });
  }, []);

  return (
    <TableContainer>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>CID</TableCell>
            <TableCell>Uploader Wallet Address</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {tools.map((graph, index) => (
            <TableRow key={index}>
              <TableCell>
                <a href={`http://bacalhau.labdao.xyz:8080/ipfs/${graph.CID}/`}>
                  {graph.CID}
                </a>
              </TableCell>
              <TableCell>{graph.WalletAddress}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}
