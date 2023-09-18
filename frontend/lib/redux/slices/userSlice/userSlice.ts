import { createSlice, PayloadAction } from '@reduxjs/toolkit'

const initialState: UserSliceState = {
  username: '',
  walletAddress: '',
  isLoading: false,
  error: null,
  isLoggedIn: false,
}

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setUsername: (state, action: PayloadAction<string>) => {
      state.username = action.payload
    },
    setWalletAddress: (state, action: PayloadAction<string>) => {
      state.walletAddress = action.payload
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    startLoading: (state) => {
      state.isLoading = true
    },
    endLoading: (state) => {
      state.isLoading = false
    },
    setIsLoggedIn: (state, action: PayloadAction<boolean>) => {
      state.isLoggedIn = action.payload
    },
  }
})


export const {
  setUsername,
  setWalletAddress,
  setError,
  startLoading,
  endLoading,
  setIsLoggedIn,
} = userSlice.actions;


export interface UserSliceState {
  username: string
  walletAddress: string
  isLoading: boolean
  error: string | null
  isLoggedIn: boolean
}
