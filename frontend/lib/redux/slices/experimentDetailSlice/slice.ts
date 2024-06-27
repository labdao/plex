import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { JobDetail, ModelDetail } from "@/lib/redux";

export interface ExperimentDetail {
  ID: number;
  Jobs: JobDetail[];
  Name: string;
  WalletAddress: string;
  StartTime: string;
  EndTime: string;
  Public: boolean;
  RecordCID: string;
}

interface ExperimentDetailSliceState {
  experiment: ExperimentDetail;
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: ExperimentDetailSliceState = {
  experiment: { 
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

export const experimentDetailSlice = createSlice({
  name: "experimentDetail",
  initialState,
  reducers: {
    setExperimentDetail: (state, action: PayloadAction<ExperimentDetail>) => {
      state.experiment = action.payload;
    },
    setExperimentDetailLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setExperimentDetailError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setExperimentDetailSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
    setExperimentDetailPublic: (state, action: PayloadAction<boolean>) => {
      state.experiment.Public = action.payload;
    },
    resetExperimentDetail: () => {
      return initialState;
    },
  },
});

export const { setExperimentDetail, setExperimentDetailLoading, setExperimentDetailError, setExperimentDetailPublic, setExperimentDetailSuccess, resetExperimentDetail } =
  experimentDetailSlice.actions;

export default experimentDetailSlice.reducer;
