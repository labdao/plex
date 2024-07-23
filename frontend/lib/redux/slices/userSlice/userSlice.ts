import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { fetchUserDataAsync, saveUserAsync } from "./thunks";

interface UserState {
  error: string | null;
  walletAddress: string | null;
  did: string | null;
  tier: 'Free' | 'Paid' | null;
  isAdmin: boolean | null;
  subscriptionStatus: 'active' | 'trialing' | 'inactive' | null;
}

const initialState: UserState = {
  error: null,
  walletAddress: null,
  did: null,
  tier: null,
  isAdmin: null,
  subscriptionStatus: null,
};

export const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setUserTier: (state, action: PayloadAction<'Free' | 'Paid' | null>) => {
      state.tier = action.payload;
    },
    setSubscriptionStatus: (state, action: PayloadAction<'active' | 'trialing' | 'inactive' | null>) => {
      state.subscriptionStatus = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchUserDataAsync.fulfilled, (state, action) => {
        state.walletAddress = action.payload.walletAddress;
        state.did = action.payload.did;
        state.tier = action.payload.tier;
        state.isAdmin = action.payload.isAdmin;
        state.subscriptionStatus = action.payload.subscriptionStatus;
      })
      .addCase(saveUserAsync.fulfilled, (state, action) => {
        state.walletAddress = action.payload.walletAddress;
        state.did = action.payload.did;
        state.tier = action.payload.tier;
        state.isAdmin = action.payload.isAdmin;
        state.subscriptionStatus = action.payload.subscriptionStatus;
      });
  }
});

export const { setError, setSubscriptionStatus, setUserTier } = userSlice.actions;

export default userSlice.reducer;