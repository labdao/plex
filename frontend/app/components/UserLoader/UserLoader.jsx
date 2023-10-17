'use client'

import { useState, useEffect, useContext } from 'react';

import {
  useDispatch,
  useSelector,
  selectWalletAddress,
  setWalletAddress,
  setIsLoggedIn,
} from '@/lib/redux'
import { usePrivy } from '@privy-io/react-auth'
import { useRouter } from 'next/navigation'

export const UserLoader = ({ children }) => {
  const dispatch = useDispatch()
  const router = useRouter();
  const [isLoaded, setIsLoaded] = useState(false);
  const { ready, authenticated } = usePrivy();


  const walletAddressFromRedux = useSelector(selectWalletAddress)

  useEffect(() => {
    const walletAddressFromLocalStorage = localStorage.getItem('walletAddress')

    if (!walletAddressFromRedux && walletAddressFromLocalStorage) {
      dispatch(setWalletAddress(walletAddressFromLocalStorage))
    }

    if (ready) {
      if (!authenticated) {
        router.push('/login')
      } else {
        dispatch(setIsLoggedIn(true))
      }
    }
    setIsLoaded(true)
  }, [dispatch, ready, authenticated])

  if (!isLoaded) return null

  return children
}
