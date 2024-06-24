// redux/flowNames/slice.ts
import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface FlowName {
  ID: number;
  Name: string;
  StartTime: string
}

export interface CategorizedFlowNames {
  today: FlowName[];
  last7Days: FlowName[];
  last30Days: FlowName[];
  older: FlowName[];
}

interface FlowNamesSliceState {
  names: FlowName[];
  loading: boolean;
  error: string | null;
  success: boolean;
  categorizedFlowNames: {
    today: FlowName[];
    last7Days: FlowName[];
    last30Days: FlowName[];
    older: FlowName[];
  }
}

const initialState: FlowNamesSliceState = {
  names: [],
  loading: false,
  error: null,
  success: false,
  categorizedFlowNames: {
    today: [],
    last7Days: [],
    last30Days: [],
    older: []
  }
}

export const flowNamesSlice = createSlice({
  name: 'FlowNames',
  initialState,
  reducers: {
    setFlowNames: (state, action: PayloadAction<FlowName[]>) => {
      state.names = action.payload;
    },
    setFlowNamesLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setFlowNamesError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setFlowNamesSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
    setCategorizedFlowNames: (state, action: PayloadAction<CategorizedFlowNames>) => {
      state.categorizedFlowNames = action.payload;
    },
  },
});

export const { setFlowNames, setFlowNamesLoading, setFlowNamesError, setFlowNamesSuccess, setCategorizedFlowNames } = flowNamesSlice.actions;
export default flowNamesSlice.reducer;
