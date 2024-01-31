import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface ToolDetail {
  CID: string;
  WalletAddress: string;
  Name: string;
  DefaultTool: boolean;
  ToolJson: { inputs: {}; outputs: {}; name: string; author: string; description: string; github: string; paper: string };
}

export interface ToolDetailSliceState {
  tool: ToolDetail;
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: ToolDetailSliceState = {
  tool: {
    CID: "",
    WalletAddress: "",
    Name: "",
    DefaultTool: false, 
    ToolJson: {
      inputs: {},
      outputs: {},
      name: "",
      author: "",
      description: "",
      github: "",
      paper: "",
    },
  },
  loading: true,
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
