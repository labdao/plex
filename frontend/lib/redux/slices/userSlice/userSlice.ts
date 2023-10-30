import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { fetchUserMemberStatusAsync } from './thunks';

interface UserState {
  walletAddress: string;
  isLoading: boolean;
  error: string | null;
  isLoggedIn: boolean;
  isMember: boolean;
};

const initialState: UserState = {
  walletAddress: '',
  isLoading: false,
  error: null,
  isLoggedIn: false,
  isMember: false,
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
    setIsMember: (state, action: PayloadAction<boolean>) => {
      state.isMember = action.payload;
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchUserMemberStatusAsync.fulfilled, (state, action: PayloadAction<boolean>) => {
        state.isMember = action.payload;
      })
      .addCase(fetchUserMemberStatusAsync.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchUserMemberStatusAsync.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to fetch member status';
      });
  },
});

export const {
  setWalletAddress,
  setError,
  startLoading,
  endLoading,
  setIsLoggedIn,
  setIsMember
} = userSlice.actions;

export default userSlice.reducer;