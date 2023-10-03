'use client'

import { useState, useEffect } from 'react';

import {
  useDispatch,
  useSelector,
  selectUsername,
  selectWalletAddress,
  setUsername,
  setWalletAddress,
  setIsLoggedIn,
} from '@/lib/redux'

import { useRouter } from 'next/navigation'

export const UserLoader = ({ children }) => {
  const dispatch = useDispatch()
  const router = useRouter();
  const [isLoaded, setIsLoaded] = useState(false);
  const userNameFromRedux = useSelector(selectUsername)
  const walletAddressFromRedux = useSelector(selectWalletAddress)

  useEffect(() => {
    const usernameFromLocalStorage = localStorage.getItem('username')
    const walletAddressFromLocalStorage = localStorage.getItem('walletAddress')

    if (!userNameFromRedux && usernameFromLocalStorage) {
      dispatch(setUsername(usernameFromLocalStorage));
    }

    if (!walletAddressFromRedux && walletAddressFromLocalStorage) {
      dispatch(setWalletAddress(walletAddressFromLocalStorage))
    }

    if (!usernameFromLocalStorage || !walletAddressFromLocalStorage) {
      router.push('/login')
    } else {
      dispatch(setIsLoggedIn(true))
    }

    setIsLoaded(true)
  }, [dispatch])

  if (!isLoaded) return null

  return children
}
