import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { ToolDetail } from "@/lib/redux";

export interface Kwargs {
  [key: string]: string[];
}

interface ExperimentAddSliceState {
  ID: number | null
  name: string
  tool: ToolDetail
  kwargs: Kwargs
  loading: boolean
  error: string | null
  success: boolean
}

const initialState: ExperimentAddSliceState = {
  ID: null,
  name: "",
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
    setExperimentAddTool: (state, action: PayloadAction<ToolDetail>) => {
      state.tool = action.payload;
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
  setExperimentAddTool,
  setExperimentAddID,
  setExperimentAddKwargs,
  setExperimentAddError,
  setExperimentAddLoading,
  setExperimentAddSuccess,
} = experimentAddSlice.actions

export default experimentAddSlice.reducer;
