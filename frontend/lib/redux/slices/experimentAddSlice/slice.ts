import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { ModelDetail } from "@/lib/redux";

export interface Kwargs {
  [key: string]: string[];
}

interface ExperimentAddSliceState {
  ID: number | null
  name: string
  model: ModelDetail
  kwargs: Kwargs
  loading: boolean
  error: string | null
  success: boolean
}

const initialState: ExperimentAddSliceState = {
  ID: null,
  name: "",
  model: { 
    ID: "", 
    WalletAddress: "", 
    Name: "", 
    DefaultModel: false, 
    S3URI: "",
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
    } 
  },
  kwargs: {},
  loading: false,
  error: null,
  success: false,
};

export const experimentAddSlice = createSlice({
  name: "experimentAdd",
  initialState,
  reducers: {
    setExperimentAddName: (state, action: PayloadAction<string>) => {
      state.name = action.payload;
    },
    setExperimentAddModel: (state, action: PayloadAction<ModelDetail>) => {
      state.model = action.payload;
    },
    setExperimentAddKwargs: (state, action: PayloadAction<Kwargs>) => {
      state.kwargs = action.payload;
    },
    setExperimentAddID: (state, action: PayloadAction<number | null>) => {
      state.ID = action.payload
    },
    setExperimentAddError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setExperimentAddLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setExperimentAddSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
  },
});

export const {
  setExperimentAddName,
  setExperimentAddModel,
  setExperimentAddID,
  setExperimentAddKwargs,
  setExperimentAddError,
  setExperimentAddLoading,
  setExperimentAddSuccess,
} = experimentAddSlice.actions

export default experimentAddSlice.reducer;
