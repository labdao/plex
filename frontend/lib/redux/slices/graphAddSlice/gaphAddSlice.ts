import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface ToolInput {
    [key: string]: {
        type: string
        item: string
        glob: string[]
    }
}

interface GraphAddSliceState {
    toolCid: string
    toolName: string
    toolJson: {
        inputs: ToolInput
        [key: string]: any // optional if you have other properties
    }
    loadingTool: boolean
    errorLoadingTool: string | null
    successLoadingTool: boolean
}

const initialState: GraphAddSliceState = {
    toolCid: "",
    toolName: "",
    toolJson: { "inputs": {} },
    loadingTool: false,
    errorLoadingTool: null,
    successLoadingTool: false
}

export const graphAddSlice = createSlice({
    name: 'graphAdd',
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
