import { createSlice, PayloadAction } from '@reduxjs/toolkit'

import { saveDataFileAsync } from './thunks'

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
    setFilenameDataSlice: (state, action: PayloadAction<string>) => {
      state.filename = action.payload
    },
    setCidDataSlice: (state, action: PayloadAction<string>) => {
      state.cid = action.payload
    },
    setDataFileError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload
    },
    startFileUploadDataSlice: (state) => {
      state.isLoading = true
    },
    endFileUploadDataSlice: (state) => {
      state.isLoading = false
    },
    setIsUploadedDataSlice: (state, action: PayloadAction<boolean>) => {
      state.isUploaded = action.payload
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(saveDataFileAsync.pending, (state) => {
        state.isLoading = true
        state.error = null
      })
      .addCase(saveDataFileAsync.fulfilled, (state, action) => {
        state.isLoading = false
        if (action.payload) {
          state.cid = action.payload.cid
          state.filename = action.payload.filename
        }
        state.isUploaded = true
      })
      .addCase(saveDataFileAsync.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.error.message || 'An error occurred while saving data file.'
      })
  }
})

export const {
  setFilenameDataSlice,
  setCidDataSlice,
  setDataFileError,
  startFileUploadDataSlice,
  endFileUploadDataSlice,
  setIsUploadedDataSlice,
} = dataFileAddSlice.actions

export default dataFileAddSlice.reducer
