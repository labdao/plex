import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { fetchUserDataAsync } from "./thunks";

interface UserState {
  error: string | null;
  walletAddress: string | null;
  did: string | null;
  isAdmin: boolean | null;
}

const initialState: UserState = {
  error: null,
  walletAddress: null,
  did: null,
  isAdmin: null,
};

export const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(fetchUserDataAsync.fulfilled, (state, action) => {
      state.walletAddress = action.payload.walletAddress;
      state.did = action.payload.did;
      state.isAdmin = action.payload.isAdmin;
    });
  }
});

export const { setError } = userSlice.actions;

export default userSlice.reducer;
