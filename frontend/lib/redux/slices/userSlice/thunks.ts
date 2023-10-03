import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { setError, setIsLoggedIn } from '@/lib/redux'
import { saveUserDataToServer } from './actions'

interface UserPayload {
  username: string
  walletAddress: string
}

export const saveUserAsync = createAppAsyncThunk(
  'user/saveUserData',
  async ({username, walletAddress}: UserPayload, { dispatch }) => {
    try {
      const response = await saveUserDataToServer(username, walletAddress)

      if (response.username && response.walletAddress) {
        localStorage.setItem('username', username)
        localStorage.setItem('walletAddress', walletAddress)
        dispatch(setIsLoggedIn(true))
      } else {
        dispatch(setError('Failed to save user data.'))
      }
      return response
    } catch (error: unknown) {
      const errorMessage = typeof error === 'object' && error !== null && 'message' in error 
        ? (error as { message?: string }).message 
        : undefined;

      dispatch(setError(errorMessage || 'An error occurred.'));
      return false;
    }
  }
)
