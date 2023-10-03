import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface ToolAddSliceState {
    toolJson: string
    loading: boolean
    error: string | null
    success: boolean
}

const initialState: ToolAddSliceState = {
    toolJson: '',
    loading: false,
    error: null,
    success: false,
}

export const toolAddSlice = createSlice({
    name: 'toolAdd',
    initialState,
    reducers: {
        setAddToolJson: (state, action: PayloadAction<string>) => {
            state.toolJson = action.payload
        },
        setAddToolError: (state, action: PayloadAction<string | null>) => {
            state.error = action.payload
        },
        setAddToolLoading: (state, action: PayloadAction<boolean>) => {
            state.loading = action.payload
        },
        setAddToolSuccess: (state, action: PayloadAction<boolean>) => {
            state.success = action.payload
        }
    }
})

export const {
    setAddToolJson,
    setAddToolError,
    setAddToolLoading,
    setAddToolSuccess,
} = toolAddSlice.actions

export default toolAddSlice.reducer
