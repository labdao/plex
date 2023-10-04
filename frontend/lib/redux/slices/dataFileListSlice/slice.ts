import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface DataFile {
  filename: string
  cid: string // Content Identifier in IPFS
}

interface DataFileListSliceState {
  dataFiles: DataFile[]
  loading: boolean
  error: string | null
  success: boolean
}

const initialState: DataFileListSliceState = {
  dataFiles: [],
  loading: false,
  error: null,
  success: false,
}

export const dataFileListSlice = createSlice({
  name: 'dataFileList',
  initialState,
  reducers: {
    setDataFileList: (state, action: PayloadAction<DataFile[]>) => {
      state.dataFiles = action.payload
    },
    setDataFileListLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload
    },
    setDataFileListError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload
    },
    setDataFileListSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload
    }
  }
})

export const {
  setDataFileList,
  setDataFileListLoading,
  setDataFileListError,
  setDataFileListSuccess
} = dataFileListSlice.actions

export default dataFileListSlice.reducer
