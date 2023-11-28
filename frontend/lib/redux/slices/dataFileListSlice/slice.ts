import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface DataFile {
  Filename: string
  CID: string
  Tags: Tag[]
}

interface Tag {
  Name: string;
  Type: string;
}

interface DataFileListSliceState {
  dataFiles: DataFile[];
  status: 'idle' | 'loading' | 'succeeded' | 'failed';
  error: string | null;
  success: boolean;
}

const initialState: DataFileListSliceState = {
  dataFiles: [],
  status: 'idle',
  error: null,
  success: false,
};

export const dataFileListSlice = createSlice({
  name: 'dataFileList',
  initialState,
  reducers: {
    setDataFileList: (state, action: PayloadAction<DataFile[]>) => {
      state.dataFiles = action.payload;
    },
    setDataFileListLoading: (state, action: PayloadAction<boolean>) => {
      state.status = action.payload ? 'loading' : 'idle';
    },
    setDataFileListError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setDataFileListSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
      state.status = action.payload ? 'succeeded' : 'failed';
    },
  },
});

export const {
  setDataFileList,
  setDataFileListLoading,
  setDataFileListError,
  setDataFileListSuccess
} = dataFileListSlice.actions

export default dataFileListSlice.reducer
