import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface ExperimentUpdateSliceState {
    loading: boolean;
    error: string | null;
    success: boolean;
}

const initialState: ExperimentUpdateSliceState = {
    loading: false,
    error: null,
    success: false,
};

export const experimentUpdateSlice = createSlice({
    name: "ExperimentUpdate",
    initialState,
    reducers: {
        setExperimentUpdateLoading: (state, action: PayloadAction<boolean>) => {
            state.loading = action.payload;
        },
        setExperimentUpdateError: (state, action: PayloadAction<string | null>) => {
            state.error = action.payload;
        },
        setExperimentUpdateSuccess: (state, action: PayloadAction<boolean>) => {
            state.success = action.payload;
        },
    },
});

export const {
    setExperimentUpdateLoading,
    setExperimentUpdateError,
    setExperimentUpdateSuccess,
} = experimentUpdateSlice.actions;

export default experimentUpdateSlice.reducer;