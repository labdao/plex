import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { ToolDetail } from "@/lib/redux";

export interface Job {
  ID: number;
  BacalhauJobID: string;
  State: string;
  Error: string;
  Tool: ToolDetail;
  FlowId: string;
}

export interface FlowDetail {
  ID: number | null;
  CID: string;
  Jobs: Job[];
  Name: string;
  WalletAddress: string;
  StartTime: string;
  EndTime: string;
  Public: boolean;
}

interface FlowDetailSliceState {
  flow: FlowDetail;
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: FlowDetailSliceState = {
  flow: { ID: null, CID: "", Jobs: [], Name: "", WalletAddress: "", StartTime: "", EndTime: "", Public: false },
  loading: true,
  error: null,
  success: false,
};

export const flowDetailSlice = createSlice({
  name: "flowDetail",
  initialState,
  reducers: {
    setFlowDetail: (state, action: PayloadAction<FlowDetail>) => {
      state.flow = action.payload;
    },
    setFlowDetailLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setFlowDetailError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setFlowDetailSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
    setFlowDetailPublic: (state, action: PayloadAction<boolean>) => {
      state.flow.Public = action.payload;
    },
    resetFlowDetail: () => {
      return initialState;
    },
  },
});

export const { setFlowDetail, setFlowDetailLoading, setFlowDetailError, setFlowDetailPublic, setFlowDetailSuccess, resetFlowDetail } =
  flowDetailSlice.actions;

export default flowDetailSlice.reducer;
