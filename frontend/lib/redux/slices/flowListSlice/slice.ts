import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface Flow {
  ID: number
  CID: string
  WalletAddress: string
  Name: string
  StartTime: string
}

export interface CategorizedFlows {
  today: Flow[];
  last7Days: Flow[];
  last30Days: Flow[];
  older: Flow[];
}

interface FlowListSliceState {
  flows: Flow[]
  loading: boolean
  error: string | null
  success: boolean
  categorizedFlows: {
    today: Flow[];
    last7Days: Flow[];
    last30Days: Flow[];
    older: Flow[];
  }
}

const initialState: FlowListSliceState = {
  flows: [],
  loading: false,
  error: null,
  success: false,
  categorizedFlows: {
    today: [],
    last7Days: [],
    last30Days: [],
    older: []
  }
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
    },
    setCategorizedFlows: (state, action: PayloadAction<CategorizedFlows>) => {
      state.categorizedFlows = action.payload
    },
  }
})

export const {
  setFlowList,
  setFlowListLoading,
  setFlowListError,
  setFlowListSuccess,
  setCategorizedFlows,
} = flowListSlice.actions

export default flowListSlice.reducer
