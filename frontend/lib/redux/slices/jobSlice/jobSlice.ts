import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface JobSliceState {
    selectedTool: string
    selectedDataFiles: string[]
    isLoading: boolean
    error: string | null
    isInitialized: boolean
}

const initialState: JobSliceState = {
    selectedTool: '',
    selectedDataFiles: [],
    isLoading: false,
    error: null,
    isInitialized: false,
}

export const jobSlice = createSlice({
    name: 'job',
    initialState,
    reducers: {
        setSelectedTool: (state, action: PayloadAction<string>) => {
            state.selectedTool = action.payload
        },
        setSelectedDataFiles: (state, action: PayloadAction<string[]>) => {
            state.selectedDataFiles = action.payload
        },
        setErrorJobSlice: (state, action: PayloadAction<string | null>) => {
            state.error = action.payload
        },
        startJobInitialization: (state) => {
            state.isLoading = true
        },
        endJobInitialization: (state) => {
            state.isLoading = false
        },
        setIsInitialized: (state, action: PayloadAction<boolean>) => {
            state.isInitialized = action.payload
        },
    }
})

export const {
    setSelectedTool,
    setSelectedDataFiles,
    setErrorJobSlice,
    startJobInitialization,
    endJobInitialization,
    setIsInitialized,
} = jobSlice.actions

export default jobSlice.reducer