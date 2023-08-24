'use client'

import { useEffect } from 'react'

import {
   useSelector,
   useDispatch,
   setUsername,
   setWalletAddress,
   setError,
   startLoading,
   endLoading,
   selectUsername,
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
  const isUserLoggedIn = useSelector((state) => state.user.username)
  const router = useRouter()

  useEffect(() => {
    const usernameFromLocalStorage = localStorage.getItem('username')
    const walletAddressFromLocalStorage = localStorage.getItem('walletAddress')

    if (usernameFromLocalStorage && walletAddressFromLocalStorage) {
      router.push('/')
    }
  }, [router])

  const dispatch = useDispatch()

  const username = useSelector(selectUsername)
  const walletAddress = useSelector(selectWalletAddress)
  const errorMessage = useSelector(selectUserFormError)
  const isLoading = useSelector(selectUserFormIsLoading)

  const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    dispatch(setUsername(e.target.value))
  }

  const handleWalletAddressChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    dispatch(setWalletAddress(e.target.value))
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    dispatch(startLoading())
    dispatch(setError(null))
    await dispatch(saveUserAsync({ username, walletAddress }))
    dispatch(endLoading())
    router.push('/')
  }

  return (
    <form onSubmit={handleSubmit}>
      <Box maxWidth={500} margin="auto"> {/* Constrain max width and center */}
        <Grid container direction="column" spacing={2}>
          <Grid item>
            <TextField
              fullWidth
              label="Username"
              variant="outlined"
              value={username}
              onChange={handleUsernameChange}
            />
          </Grid>
          <Grid item>
          <TextField
            fullWidth
            label="Eth Wallet Address"
            type="text"
            variant="outlined"
            value={walletAddress}
            onChange={handleWalletAddressChange}
          />
          </Grid>
          {errorMessage && (
            <Box my={2}>
              <Alert severity="error" variant="filled">
                <Typography align="center">{errorMessage}</Typography>
              </Alert>
            </Box>
          )}
          <Grid item container justifyContent="center">
            <Button variant="contained" color="primary" type="submit">
              {isLoading ? "Submitting..." : "Submit"}
            </Button>
          </Grid>
        </Grid>
      </Box>
    </form>
  )
}
