'use client'

import { useEffect } from 'react'
import LoginComponent from '../components/Web3AuthLogin/LoginComponent'

import {
   useSelector,
   useDispatch,
   setUsername,
   setWalletAddress,
   setError,
   startLoading,
   endLoading,
   selectUsername,
   selectIsLoggedIn,
   selectWalletAddress,
   selectUserFormError,
   selectUserFormIsLoading,
   saveUserAsync,
} from '@/lib/redux'

import { useRouter } from 'next/navigation'

import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Button from '@mui/material/Button'
import TextField from '@mui/material/TextField'
import Alert from '@mui/material/Alert'
import Typography from '@mui/material/Typography'

export default function LoginPage() {
  const isUserLoggedIn = useSelector(selectIsLoggedIn)
  const router = useRouter()

  useEffect(() => {
    if (isUserLoggedIn) {
      router.push('/')
    }
  }, [router, isUserLoggedIn])

  return (
    <div>
      <LoginComponent />
    </div>
  )
}
