import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface UserState {
  error: string | null;
}

const initialState: UserState = {
  error: null,
};

export const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
  },
});

export const { setError } = userSlice.actions;

export default userSlice.reducer;
