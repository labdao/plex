import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface ApiKeyAddSliceState {
  key: string;
  name: string;
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: ApiKeyAddSliceState = {
  key: "",
  name: "",
  loading: false,
  error: null,
  success: false,
};

export const apiKeyAddSlice = createSlice({
  name: "apiKeyAdd",
  initialState,
  reducers: {
    setApiKey: (state, action: PayloadAction<string>) => {
      state.key = action.payload;
    },
    setApiKeyName: (state, action: PayloadAction<string>) => {
      state.name = action.payload;
    },
    setApiKeyLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setApiKeyError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setApiKeySuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
    // You might need additional reducers depending on your form and API key creation logic
  },
});

export const { setApiKey, setApiKeyName, setApiKeyLoading, setApiKeyError, setApiKeySuccess } =
  apiKeyAddSlice.actions;

export default apiKeyAddSlice.reducer;