import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface UserState {
  error: string | null;
  // isLoggedIn: boolean;
  // authToken: string;
}

const initialState: UserState = {
  error: null,
  // isLoggedIn: false,
  // authToken: '',
};

export const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    // setIsLoggedIn: (state, action: PayloadAction<boolean>) => {
    //   state.isLoggedIn = action.payload;
    // },
    // setAuthToken: (state, action: PayloadAction<string>) => {
    //   state.authToken = action.payload;
    // },
  },
});

// export const {
//   setWalletAddress,
//   setError,
//   startLoading,
//   endLoading,
//   setIsLoggedIn,
//   setAuthToken,
// } = userSlice.actions;

export const { setError } = userSlice.actions;

export default userSlice.reducer;
