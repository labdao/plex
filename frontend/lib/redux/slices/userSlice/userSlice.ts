import { createSlice, PayloadAction } from '@reduxjs/toolkit';
// import { initWeb3Auth } from './thunks';

interface UserState {
  web3Auth: any;
  username: string;
  walletAddress: string;
  isLoading: boolean;
  error: string | null;
  isLoggedIn: boolean;
}

const initialState: UserState = {
  web3Auth: null,
  username: '',
  walletAddress: '',
  isLoading: false,
  error: null,
  isLoggedIn: false,
};

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setUsername: (state, action: PayloadAction<string>) => {
      state.username = action.payload;
    },
    setWalletAddress: (state, action: PayloadAction<string>) => {
      state.walletAddress = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    startLoading: (state) => {
      state.isLoading = true;
    },
    endLoading: (state) => {
      state.isLoading = false;
    },
    setIsLoggedIn: (state, action: PayloadAction<boolean>) => {
      state.isLoggedIn = action.payload;
    },
  },
  // extraReducers: (builder) => {
  //   builder.addCase(initWeb3Auth.fulfilled, (state, action: PayloadAction<any>) => {
  //     state.web3Auth = action.payload;
  //   });
  // },
});

export const {
  setUsername,
  setWalletAddress,
  setError,
  startLoading,
  endLoading,
  setIsLoggedIn,
} = userSlice.actions;

export default userSlice.reducer;