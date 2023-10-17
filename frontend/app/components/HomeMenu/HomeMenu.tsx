'use client'

import React from 'react'
import List from '@mui/material/List'
import ListItem from '@mui/material/ListItem'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'
import { useRouter } from 'next/navigation'
import { usePrivy } from '@privy-io/react-auth';

export const HomeMenu = () => {
  const router = useRouter()

  const handleNavigation = (path: string) => {
    router.push(path)
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
      <List>
        <ListItem disablePadding>
          <ListItemButton onClick={() => handleNavigation('/datafile/add')}>
            <ListItemText primary="Add Data" />
          </ListItemButton>
        </ListItem>
        <ListItem disablePadding>
          <ListItemButton onClick={() => handleNavigation('/datafile/list')}>
            <ListItemText primary="View Data" />
          </ListItemButton>
        </ListItem>
        <ListItem disablePadding>
          <ListItemButton onClick={() => handleNavigation('/tool/add')}>
            <ListItemText primary="Add Tool" />
          </ListItemButton>
        </ListItem>
        <ListItem disablePadding>
          <ListItemButton onClick={() => handleNavigation('/tool/list')}>
            <ListItemText primary="View Tools" />
          </ListItemButton>
        </ListItem>
        <ListItem disablePadding>
          <ListItemButton onClick={() => handleNavigation('/flow/add')}>
            <ListItemText primary="Add a Flow" />
          </ListItemButton>
        </ListItem>
        <ListItem disablePadding>
          <ListItemButton onClick={() => handleNavigation('/flow/list')}>
            <ListItemText primary="View Flows" />
          </ListItemButton>
        </ListItem>
        <ListItem disablePadding>
          <ListItemButton onClick={() => handleNavigation('/philosophy')}>
            <ListItemText primary="Philosophy" />
          </ListItemButton>
        </ListItem>
      </List>
    </div>
  )
}
