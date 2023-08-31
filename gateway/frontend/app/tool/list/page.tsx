'use client'

import React, { useEffect, useState } from 'react'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'

export default function ListToolFiles() {
  interface Tool {
    CID: string;
    ToolJSON: string;
  }

  const [tools, setTools] = useState<Tool[]>([]);

  useEffect(() => {
    fetch('http://localhost:8080/get-tools')
      .then(response => {
        if (!response.ok) {
          throw new Error(`HTTP error ${response.status}`);
        }
        return response.json();
      })
      .then(data => {
        console.log('Fetched tools:', data);
        setTools(data);
      })
      .catch(error => {
        console.error('Error fetching tools:', error);
      });
  }, []);

  return (
    <TableContainer>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>CID</TableCell>
            <TableCell>Serialized Tool Config</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {tools.map((tool, index) => ( 
            <TableRow key={index}>
              <TableCell>{tool.CID}</TableCell>
              <TableCell>{JSON.stringify(tool.ToolJSON)}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}
