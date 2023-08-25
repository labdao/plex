import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface DataFileSliceState {
  filename: string
  cid: string // Content Identifier in IPFS
  isLoading: boolean
  error: string | null
  isUploaded: boolean
}

const initialState: DataFileSliceState = {
  filename: '',
  cid: '',
  isLoading: false,
  error: null,
  isUploaded: false,
}

export const dataFileAddSlice = createSlice({
  name: 'dataFile',
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
} = dataFileAddSlice.actions

export default dataFileAddSlice.reducer
