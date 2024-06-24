import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { ModelDetail } from "@/lib/redux";

interface ModelListSliceState {
  models: ModelDetail[];
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: ModelListSliceState = {
  models: [],
  loading: false,
  error: null,
  success: false,
};

export const modelListSlice = createSlice({
  name: "modelList",
  initialState,
  reducers: {
    setModelList: (state, action: PayloadAction<ModelDetail[]>) => {
      state.models = action.payload;
    },
    setModelListLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setModelListError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setModelListSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
    resetModelList: () => {
      return initialState;
    },
  },
});

export const { setModelList, setModelListLoading, setModelListError, setModelListSuccess, resetModelList } = modelListSlice.actions;

export default modelListSlice.reducer;
