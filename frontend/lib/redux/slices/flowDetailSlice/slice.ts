import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { JobDetail, ToolDetail } from "@/lib/redux";

export interface FlowDetail {
  ID: number;
  Jobs: JobDetail[];
  Name: string;
  WalletAddress: string;
  StartTime: string;
  EndTime: string;
  Public: boolean;
  RecordCID: string;
}

interface FlowDetailSliceState {
  flow: FlowDetail;
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: FlowDetailSliceState = {
  flow: { 
    ID: 0, 
    Jobs: [], 
    Name: "", 
    WalletAddress: "", 
    StartTime: "", 
    EndTime: "", 
    Public: false,
    RecordCID: "",
  },
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
