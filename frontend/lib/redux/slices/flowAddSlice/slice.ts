import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { ToolDetail } from "@/lib/redux";

export interface Kwargs {
  [key: string]: string[];
}

interface FlowAddSliceState {
  ID: number | null
  name: string
  tool: ToolDetail
  kwargs: Kwargs
  cid: string
  loading: boolean
  error: string | null
  success: boolean
}

const initialState: FlowAddSliceState = {
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
  cid: "",
  loading: false,
  error: null,
  success: false,
};

export const flowAddSlice = createSlice({
  name: "flowAdd",
  initialState,
  reducers: {
    setFlowAddName: (state, action: PayloadAction<string>) => {
      state.name = action.payload;
    },
    setFlowAddTool: (state, action: PayloadAction<ToolDetail>) => {
      state.tool = action.payload;
    },
    setFlowAddKwargs: (state, action: PayloadAction<Kwargs>) => {
      state.kwargs = action.payload;
    },
    setFlowAddCid: (state, action: PayloadAction<string>) => {
      state.cid = action.payload;
    },
    setFlowAddID: (state, action: PayloadAction<number | null>) => {
      state.ID = action.payload
    },
    setFlowAddError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setFlowAddLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setFlowAddSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
  },
});

export const {
  setFlowAddName,
  setFlowAddTool,
  setFlowAddCid,
  setFlowAddID,
  setFlowAddKwargs,
  setFlowAddError,
  setFlowAddLoading,
  setFlowAddSuccess,
} = flowAddSlice.actions

export default flowAddSlice.reducer;
