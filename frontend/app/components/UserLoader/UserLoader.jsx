'use client'

import { useState, useEffect } from 'react';

import {
  useDispatch,
  useSelector,
  selectWalletAddress,
  selectEmailAddress,
  setWalletAddress,
  setIsLoggedIn,
  setEmailAddress,
} from '@/lib/redux'

import { useRouter } from 'next/navigation'

export const UserLoader = ({ children }) => {
  const dispatch = useDispatch()
  const router = useRouter();
  const [isLoaded, setIsLoaded] = useState(false);

  const walletAddressFromRedux = useSelector(selectWalletAddress)
  // const emailAddressFromRedux = useSelector(selectEmailAddress)

  useEffect(() => {
    const walletAddressFromLocalStorage = localStorage.getItem('walletAddress')

    if (!walletAddressFromRedux && walletAddressFromLocalStorage) {
      dispatch(setWalletAddress(walletAddressFromLocalStorage))
    }

    if (!walletAddressFromLocalStorage) {
      router.push('/login')
    } else {
      dispatch(setIsLoggedIn(true))
    }

    setIsLoaded(true)
  }, [dispatch])

  if (!isLoaded) return null

  return children
}
