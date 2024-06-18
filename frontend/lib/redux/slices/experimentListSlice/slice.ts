import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface Experiment {
  ID: number
  WalletAddress: string
  Name: string
  StartTime: string
  Public: boolean
}

export interface CategorizedExperiments {
  today: Experiment[];
  last7Days: Experiment[];
  last30Days: Experiment[];
  older: Experiment[];
}

interface ExperimentListSliceState {
  experiments: Experiment[]
  loading: boolean
  error: string | null
  success: boolean
  categorizedExperiments: {
    today: Experiment[];
    last7Days: Experiment[];
    last30Days: Experiment[];
    older: Experiment[];
  }
}

const initialState: ExperimentListSliceState = {
  experiments: [],
  loading: false,
  error: null,
  success: false,
  categorizedExperiments: {
    today: [],
    last7Days: [],
    last30Days: [],
    older: []
  }
}

export const experimentListSlice = createSlice({
  name: 'ExperimentList',
  initialState,
  reducers: {
    setExperimentList: (state, action: PayloadAction<Experiment[]>) => {
      state.experiments = action.payload
    },
    setExperimentListLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload
    },
    setExperimentListError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload
    },
    setExperimentListSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload
    },
    setCategorizedExperiments: (state, action: PayloadAction<CategorizedExperiments>) => {
      state.categorizedExperiments = action.payload
    },
  }
})

export const {
  setExperimentList,
  setExperimentListLoading,
  setExperimentListError,
  setExperimentListSuccess,
  setCategorizedExperiments,
} = experimentListSlice.actions

export default experimentListSlice.reducer
