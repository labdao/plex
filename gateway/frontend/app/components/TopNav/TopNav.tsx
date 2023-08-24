'use client'
import React from 'react'

import Link from 'next/link'
import MenuIcon from '@mui/icons-material/Menu'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'

import styles from './topnav.module.css'
import {
  useSelector,
  useDispatch,
  selectUsername,
  selectWalletAddress,
  setUsername,
  setWalletAddress,
} from '@/lib/redux'

export const TopNav = () => {
  const dispatch = useDispatch()
  const username = useSelector(selectUsername)
  const walletAddress = useSelector(selectWalletAddress)

  // State and handlers for the dropdown menu
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

  const handleLogout = () => {
    // Clear data from localStorage
    localStorage.removeItem('username')
    localStorage.removeItem('walletAddress')
    dispatch(setUsername(''))
    dispatch(setWalletAddress(''))
    handleClose()
  };


  return (
    <nav className={styles.navbar}>
      <span className={styles.link}>Plex</span>
      {username && (
        <div className={styles.userContainer}>
          <span className={styles.username}>{username}</span>
          <MenuIcon style={{ color: 'white', marginLeft: '10px' }} onClick={handleClick} />
          <Menu
            anchorEl={anchorEl}
            keepMounted
            open={Boolean(anchorEl)}
            onClose={handleClose}
          >
            <MenuItem onClick={handleClose}>Wallet: { walletAddress }</MenuItem>
            <MenuItem onClick={handleLogout}>Logout</MenuItem>
          </Menu>
        </div>
      )}
      {/* Other links or elements can be added here if required */}
    </nav>
  )
}
