import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface ToolSliceState {
    filename: string
    cid: string
    isLoading: boolean
    error: string | null
    isUploaded: boolean
}

const initialState: ToolSliceState = {
    filename: '',
    cid: '',
    isLoading: false,
    error: null,
    isUploaded: false,
}

export const toolAddSlice = createSlice({
    name: 'tool',
    initialState,
    reducers: {
        setFilename: (state, action: PayloadAction<string>) => {
            state.filename = action.payload
        },
        setCid: (state, action: PayloadAction<string>) => {
            state.cid = action.payload
        },
        setError: (state, action: PayloadAction<string | null>) => {
            state.error = action.payload
        },
        startFileUpload: (state) => {
            state.isLoading = true
        },
        endFileUpload: (state) => {
            state.isLoading = false
        },
        setIsUploaded: (state, action: PayloadAction<boolean>) => {
            state.isUploaded = action.payload
        },
    }
})

export const {
    setFilename,
    setCid,
    setError,
    startFileUpload,
    endFileUpload,
    setIsUploaded,
} = toolAddSlice.actions

export default toolAddSlice.reducer