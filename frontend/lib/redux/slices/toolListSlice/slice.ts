import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { ToolDetail } from "../toolDetailSlice";

export interface Tool {
  CID: string;
  WalletAddress: string;
  Name: string;
}

interface ToolListSliceState {
  tools: ToolDetail[];
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: ToolListSliceState = {
  tools: [],
  loading: false,
  error: null,
  success: false,
};

export const toolListSlice = createSlice({
  name: "toolList",
  initialState,
  reducers: {
    setToolList: (state, action: PayloadAction<ToolDetail[]>) => {
      state.tools = action.payload;
    },
    setToolListLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setToolListError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setToolListSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
  },
});

export const { setToolList, setToolListLoading, setToolListError, setToolListSuccess } = toolListSlice.actions;

export default toolListSlice.reducer;
