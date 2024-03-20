import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface FlowUpdateSliceState {
    loading: boolean;
    error: string | null;
    success: boolean;
}

const initialState: FlowUpdateSliceState = {
    loading: false,
    error: null,
    success: false,
};

export const flowUpdateSlice = createSlice({
    name: "FlowUpdate",
    initialState,
    reducers: {
        setFlowUpdateLoading: (state, action: PayloadAction<boolean>) => {
            state.loading = action.payload;
        },
        setFlowUpdateError: (state, action: PayloadAction<string | null>) => {
            state.error = action.payload;
        },
        setFlowUpdateSuccess: (state, action: PayloadAction<boolean>) => {
            state.success = action.payload;
        },
    },
});

export const {
    setFlowUpdateLoading,
    setFlowUpdateError,
    setFlowUpdateSuccess,
} = flowUpdateSlice.actions;

export default flowUpdateSlice.reducer;