'use client'

import React, { useContext, useEffect } from 'react'

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
  setWalletAddress,
  setEmailAddress,
  setIsLoggedIn,
  selectEmailAddress,
} from '@/lib/redux'
import { usePrivy } from '@privy-io/react-auth';
import { PrivyAuthContext } from '../../../lib/PrivyContext';

export const TopNav = () => {
  const dispatch = useDispatch()
  const router = useRouter()
  const isLoggedIn = useSelector(selectIsLoggedIn)
  const walletAddress = useSelector(selectWalletAddress)

  const { logout } = usePrivy();

  // State and handlers for the dropdown menu
  const [anchorEl, setAnchorEl] = React.useState<null | SVGSVGElement>(null)

  const handleClick = (event: React.MouseEvent<SVGSVGElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

  const handleNavigation = (path: string) => {
    router.push(path)
  }

  const handleLogout = async () => {
    logout();
    localStorage.removeItem('walletAddress');
    dispatch(setWalletAddress(''));
    dispatch(setIsLoggedIn(false));
    handleClose();
    router.push('/login');
  }

  return (
    <nav className={styles.navbar}>
      <span className={styles.link} onClick={() => handleNavigation('/')}>
        plex
      </span>
      {isLoggedIn && (
        <div className={styles.userContainer}>
          <MenuIcon style={{ color: 'white', marginLeft: '10px' }} onClick={(e: any) => handleClick(e)} />
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
    </nav>
  )
}