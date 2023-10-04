import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface Flow {
  cid: string
  walletAddress: string
}

interface FlowListSliceState {
  flows: Flow[]
  loading: boolean
  error: string | null
  success: boolean
}

const initialState: FlowListSliceState = {
  flows: [],
  loading: false,
  error: null,
  success: false,
}

export const flowListSlice = createSlice({
  name: 'FlowList',
  initialState,
  reducers: {
    setFlowList: (state, action: PayloadAction<Flow[]>) => {
      state.flows = action.payload
    },
    setFlowListLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload
    },
    setFlowListError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload
    },
    setFlowListSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload
    }
  }
})

export const {
  setFlowList,
  setFlowListLoading,
  setFlowListError,
  setFlowListSuccess
} = flowListSlice.actions

export default flowListSlice.reducer
