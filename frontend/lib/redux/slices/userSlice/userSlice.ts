import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface UserState {
  walletAddress: string;
  isLoading: boolean;
  error: string | null;
  isLoggedIn: boolean;
}

const initialState: UserState = {
  walletAddress: '',
  isLoading: false,
  error: null,
  isLoggedIn: false,
};

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setWalletAddress: (state, action: PayloadAction<string>) => {
      state.walletAddress = action.payload;
    },
    startLoading: (state) => {
      state.isLoading = true;
    },
    endLoading: (state) => {
      state.isLoading = false;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setIsLoggedIn: (state, action: PayloadAction<boolean>) => {
      state.isLoggedIn = action.payload;
    },
  },
});

export const {
  setWalletAddress,
  setError,
  startLoading,
  endLoading,
  setIsLoggedIn,
} = userSlice.actions;

export default userSlice.reducer;