'use client'

import React from 'react'

import MenuIcon from '@mui/icons-material/Menu'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'

import { useRouter } from 'next/navigation'

import styles from './topnav.module.css'
import {
  useDispatch,
  useSelector,
  selectWalletAddress,
  selectIsLoggedIn,
  selectUsername,
  setUsername,
  setWalletAddress,
  setIsLoggedIn,
} from '@/lib/redux'

export const TopNav = () => {
  const dispatch = useDispatch()
  const router = useRouter()
  const isLoggedIn = useSelector(selectIsLoggedIn)
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
    dispatch(setIsLoggedIn(false))
    handleClose()
    router.push('/login')
  }


  return (
    <nav className={styles.navbar}>
      <span className={styles.link}>Plex</span>
      {isLoggedIn && (
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
