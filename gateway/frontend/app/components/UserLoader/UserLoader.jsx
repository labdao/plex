'use client'

import { useState, useEffect } from 'react';

import {
  useDispatch,
  useSelector,
  selectUsername,
  selectWalletAddress,
  setUsername,
  setWalletAddress,
 } from '@/lib/redux'


export const UserLoader = ({ children }) => {
  const dispatch = useDispatch()
  const [isLoaded, setIsLoaded] = useState(false);
  const userNameFromRedux = useSelector(selectUsername)
  const walletAddressFromRedux = useSelector(selectWalletAddress)

  useEffect(() => {
    const usernameFromLocalStorage = localStorage.getItem('username')
    if (!userNameFromRedux && usernameFromLocalStorage) {
      dispatch(setUsername(usernameFromLocalStorage));
    }

    const walletAddressFromLocalStorage = localStorage.getItem('walletAddress')
    if (!walletAddressFromRedux && walletAddressFromLocalStorage) {
      dispatch(setWalletAddress(walletAddressFromLocalStorage))
    }

    setIsLoaded(true)
  }, [dispatch])

  if (!isLoaded) return null

  return children
}
