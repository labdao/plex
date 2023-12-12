import { createSlice, PayloadAction } from '@reduxjs/toolkit'

import { saveDataFilesAsync } from './thunks'

interface DataFileSliceState {
  filenames: string[]
  cids: string[]
  isLoading: boolean
  error: string | null
  isUploaded: boolean
}

const initialState: DataFileSliceState = {
  filenames: [],
  cids: [],
  isLoading: false,
  error: null,
  isUploaded: false,
}

export const dataFileAddSlice = createSlice({
  name: 'dataFile',
  initialState,
  reducers: {
    addFilename: (state, action: PayloadAction<string>) => {
      state.filenames.push(action.payload)
    },
    addCid: (state, action: PayloadAction<string>) => {
      state.cids.push(action.payload)
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
    resetDataFileSlice: (state) => {
      state.filenames = []
      state.cids = []
      state.isLoading = false
      state.error = null
      state.isUploaded = false
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(saveDataFilesAsync.pending, (state) => {
        state.isLoading = true
        state.error = null
      })
      .addCase(saveDataFilesAsync.fulfilled, (state, action) => {
        state.isLoading = false
        if (action.payload && Array.isArray(action.payload)) {
          for (const file of action.payload) {
            state.cids.push(file.cid)
            state.filenames.push(file.filename)
          }
        }
        state.isUploaded = true
      })
      .addCase(saveDataFilesAsync.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.error.message || 'An error occurred while saving data file.'
      })
  }
})

export const {
  addFilename,
  addCid,
  setDataFileError,
  startFileUploadDataSlice,
  endFileUploadDataSlice,
  setIsUploadedDataSlice,
} = dataFileAddSlice.actions

export default dataFileAddSlice.reducer
