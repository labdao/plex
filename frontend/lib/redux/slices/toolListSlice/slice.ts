import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { ToolDetail } from "@/lib/redux";

interface ToolListSliceState {
  models: ToolDetail[];
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: ToolListSliceState = {
  models: [],
  loading: false,
  error: null,
  success: false,
};

export const toolListSlice = createSlice({
  name: "toolList",
  initialState,
  reducers: {
    setToolList: (state, action: PayloadAction<ToolDetail[]>) => {
      state.models = action.payload;
    },
    setToolListLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setToolListError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setToolListSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
    resetToolList: () => {
      return initialState;
    },
  },
});

export const { setToolList, setToolListLoading, setToolListError, setToolListSuccess, resetToolList } = toolListSlice.actions;

export default toolListSlice.reducer;
