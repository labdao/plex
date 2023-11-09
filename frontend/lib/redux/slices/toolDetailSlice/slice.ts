import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface ToolDetail {
  CID: string;
  WalletAddress: string;
  Name: string;
  ToolJson: { inputs: {}; name: string; author: string; description: string };
}

export interface ToolDetailSliceState {
  tool: ToolDetail;
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: ToolDetailSliceState = {
  tool: { CID: "", WalletAddress: "", Name: "", ToolJson: { inputs: {}, name: "", author: "", description: "" } },
  loading: false,
  error: null,
  success: false,
};

export const toolDetailSlice = createSlice({
  name: "toolDetail",
  initialState,
  reducers: {
    setToolDetail: (state, action: PayloadAction<ToolDetail>) => {
      state.tool = action.payload;
    },
    setToolDetailLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setToolDetailError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setToolDetailSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
  },
});

export const { setToolDetail, setToolDetailLoading, setToolDetailError, setToolDetailSuccess } = toolDetailSlice.actions;

export default toolDetailSlice.reducer;
