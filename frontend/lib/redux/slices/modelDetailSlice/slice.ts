import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface ModelDetail {
  CID: string;
  WalletAddress: string;
  Name: string;
  DefaultModel: boolean;
  ModelJson: {
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
    taskCategory?: string;
    maxRunningTime?: number;
  };
}

export interface ModelDetailSliceState {
  model: ModelDetail;
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: ModelDetailSliceState = {
  model: {
    CID: "",
    WalletAddress: "",
    Name: "",
    DefaultModel: false,
    ModelJson: {
      inputs: {},
      outputs: {},
      name: "",
      author: "",
      description: "",
      github: "",
      paper: "",
      guide: "",
      checkpointCompatible: false,
      maxRunningTime: 2700,
    },
  },
  loading: true,
  error: null,
  success: false,
};

export const modelDetailSlice = createSlice({
  name: "modelDetail",
  initialState,
  reducers: {
    setModelDetail: (state, action: PayloadAction<ModelDetail>) => {
      state.model = action.payload;
    },
    setModelDetailLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setModelDetailError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setModelDetailSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
    resetModelDetail: () => {
      return initialState;
    },
  },
});

export const { setModelDetail, setModelDetailLoading, setModelDetailError, setModelDetailSuccess, resetModelDetail } = modelDetailSlice.actions;

export default modelDetailSlice.reducer;
