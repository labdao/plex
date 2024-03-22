import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface ToolDetail {
  CID: string;
  WalletAddress: string;
  Name: string;
  DefaultTool: boolean;
  ToolJson: {
    inputs: {};
    outputs: {
      [key: string]: {
        glob: string[];
        item: string;
        type: string;
      };
    } | null;
    name: string;
    author: string;
    description: string;
    github: string;
    paper: string;
    guide: string;
    checkpointCompatible: boolean;
  };
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
      guide: "",
      checkpointCompatible: false,
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
    resetToolDetail: () => {
      return initialState;
    },
  },
});

export const { setToolDetail, setToolDetailLoading, setToolDetailError, setToolDetailSuccess, resetToolDetail } = toolDetailSlice.actions;

export default toolDetailSlice.reducer;
