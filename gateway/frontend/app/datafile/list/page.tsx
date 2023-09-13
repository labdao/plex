'use client'

import React, { useEffect, useState } from 'react'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'

import backendUrl from 'lib/backendUrl'

export default function ListDataFiles() {
  interface DataFile {
    CID: string;
    WalletAddress: string;
    Filename: string;
    IsPublic: boolean;
    IsVisible: boolean;
  }

  const [datafiles, setDataFiles] = useState<DataFile[]>([]);

  useEffect(() => {
    fetch(`${backendUrl()}/get-datafiles`)
      .then(response => {
        if (!response.ok) {
          throw new Error(`HTTP error ${response.status}`);
        }
        return response.json();
      })
      .then(data => {
        console.log('Fetched datafiles:', data);
        setDataFiles(data);
      })
  }, [])

  return (
    <TableContainer>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>CID</TableCell>
            <TableCell>Uploader Wallet</TableCell>
            <TableCell>Filename</TableCell>
            <TableCell>Is Public</TableCell>
            <TableCell>Is Visible</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {datafiles.map((datafile, index) => (
            <TableRow key={index}>
              <TableCell>
                <a href={`http://bacalhau.labdao.xyz:8080/ipfs/${datafile.CID}/`}>
                  {datafile.CID}
                </a>
              </TableCell>
              <TableCell>{datafile.WalletAddress}</TableCell>
              <TableCell>{datafile.Filename}</TableCell>
              <TableCell>{datafile.IsPublic ? 'Yes' : 'No'}</TableCell>
              <TableCell>{datafile.IsVisible ? 'Yes' : 'No'}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}