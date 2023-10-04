import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { Tool } from '../toolListSlice/slice'

export interface Kwargs {
  [key: string]: string[];
}

interface FlowAddSliceState {
  tool: Tool,
  kwargs: Kwargs,
  loading: boolean
  error: string | null
  success: boolean
}

const initialState: FlowAddSliceState = {
  tool: { CID: '', WalletAddress: '', Name: '', ToolJson: { inputs: {} }},
  kwargs: {},
  loading: false,
  error: null,
  success: false,
}

export const flowAddSlice = createSlice({
  name: 'flowAdd',
  initialState,
  reducers: {
    setFlowAddTool: (state, action: PayloadAction<Tool>) => {
      state.tool = action.payload
    },
    setFlowAddKwargs: (state, action: PayloadAction<Kwargs>) => {
      state.kwargs = action.payload
    },
    setFlowAddError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload
    },
    setFlowAddLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload
    },
    setFlowAddSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload
    }
  }
})

export const {
  setFlowAddTool,
  setFlowAddKwargs,
  setFlowAddError,
  setFlowAddLoading,
  setFlowAddSuccess,
} = flowAddSlice.actions

export default flowAddSlice.reducer
