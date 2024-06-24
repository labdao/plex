import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface ModelAddSliceState {
    modelJson: string
    loading: boolean
    error: string | null
    success: boolean
}

const initialState: ModelAddSliceState = {
    modelJson: '',
    loading: false,
    error: null,
    success: false,
}

export const modelAddSlice = createSlice({
    name: 'modelAdd',
    initialState,
    reducers: {
        setAddModelJson: (state, action: PayloadAction<string>) => {
            state.modelJson = action.payload
        },
        setAddModelError: (state, action: PayloadAction<string | null>) => {
            state.error = action.payload
        },
        setAddModelLoading: (state, action: PayloadAction<boolean>) => {
            state.loading = action.payload
        },
        setAddModelSuccess: (state, action: PayloadAction<boolean>) => {
            state.success = action.payload
        }
    }
})

export const {
    setAddModelJson,
    setAddModelError,
    setAddModelLoading,
    setAddModelSuccess,
} = modelAddSlice.actions

export default modelAddSlice.reducer
