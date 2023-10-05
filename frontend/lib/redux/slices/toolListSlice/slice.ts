import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface Tool {
  CID: string
  WalletAddress: string
  Name: string
  ToolJson: {inputs: {}}
}

interface ToolListSliceState {
  tools: Tool[]
  loading: boolean
  error: string | null
  success: boolean
}

const initialState: ToolListSliceState = {
  tools: [],
  loading: false,
  error: null,
  success: false,
}

export const toolListSlice = createSlice({
  name: 'toolList',
  initialState,
  reducers: {
    setToolList: (state, action: PayloadAction<Tool[]>) => {
      state.tools = action.payload
    },
    setToolListLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload
    },
    setToolListError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload
    },
    setToolListSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload
    }
  }
})

export const {
  setToolList,
  setToolListLoading,
  setToolListError,
  setToolListSuccess
} = toolListSlice.actions

export default toolListSlice.reducer
