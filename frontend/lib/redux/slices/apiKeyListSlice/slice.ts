import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface ApiKey {
  id: string;
  key: string;
  name: string;
}

interface ApiKeyListSliceState {
  apiKeys: ApiKey[];
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: ApiKeyListSliceState = {
  apiKeys: [],
  loading: false,
  error: null,
  success: false,
};

export const apiKeyListSlice = createSlice({
  name: 'ApiKeyList',
  initialState,
  reducers: {
    setApiKeyList: (state, action: PayloadAction<ApiKey[]>) => {
      state.apiKeys = action.payload;
    },
    setApiKeyListLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setApiKeyListError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setApiKeyListSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
  },
});

export const {
  setApiKeyList,
  setApiKeyListLoading,
  setApiKeyListError,
  setApiKeyListSuccess,
} = apiKeyListSlice.actions;

export default apiKeyListSlice.reducer;