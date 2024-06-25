// redux/experimentNames/slice.ts
import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface ExperimentName {
  ID: number;
  Name: string;
  StartTime: string
}

export interface CategorizedExperimentNames {
  today: ExperimentName[];
  last7Days: ExperimentName[];
  last30Days: ExperimentName[];
  older: ExperimentName[];
}

interface ExperimentNamesSliceState {
  names: ExperimentName[];
  loading: boolean;
  error: string | null;
  success: boolean;
  categorizedExperimentNames: {
    today: ExperimentName[];
    last7Days: ExperimentName[];
    last30Days: ExperimentName[];
    older: ExperimentName[];
  }
}

const initialState: ExperimentNamesSliceState = {
  names: [],
  loading: false,
  error: null,
  success: false,
  categorizedExperimentNames: {
    today: [],
    last7Days: [],
    last30Days: [],
    older: []
  }
}

export const experimentNamesSlice = createSlice({
  name: 'ExperimentNames',
  initialState,
  reducers: {
    setExperimentNames: (state, action: PayloadAction<ExperimentName[]>) => {
      state.names = action.payload;
    },
    setExperimentNamesLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setExperimentNamesError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setExperimentNamesSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
    setCategorizedExperimentNames: (state, action: PayloadAction<CategorizedExperimentNames>) => {
      state.categorizedExperimentNames = action.payload;
    },
  },
});

export const { setExperimentNames, setExperimentNamesLoading, setExperimentNamesError, setExperimentNamesSuccess, setCategorizedExperimentNames } = experimentNamesSlice.actions;
export default experimentNamesSlice.reducer;
