import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { DataFile } from "../dataFileListSlice/slice";

export interface JobDetail {
  ID: number | null;
  State: string;
  Error: string;
  ToolID: string;
  ExperimentID: string;
  InputFiles: DataFile[];
  OutputFiles: DataFile[];
  Status: string;
  RayJobID: string;
  Model: any;
  Inputs: any;
}

export interface JobDetailSliceState {
  job: JobDetail;
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: JobDetailSliceState = {
  job: {
    ID: null,
    State: "",
    Error: "",
    ToolID: "",
    ExperimentID: "",
    InputFiles: [],
    OutputFiles: [],
    Status: "unknown",
    RayJobID: "",
    Model: {},
    Inputs: {},
  },
  loading: false,
  error: null,
  success: false,
};

export const jobDetailSlice = createSlice({
  name: "jobDetail",
  initialState,
  reducers: {
    setJobDetail: (state, action: PayloadAction<JobDetail>) => {
      state.job = action.payload;
    },
    setJobDetailLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setJobDetailError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setJobDetailSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
  },
});

export const { setJobDetail, setJobDetailLoading, setJobDetailError, setJobDetailSuccess } = jobDetailSlice.actions;

export default jobDetailSlice.reducer;
