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

interface Pagination {
  currentPage: number;
  totalPages: number;
  pageSize: number;
  totalCount: number;
}

interface DataFileListSliceState {
  dataFiles: DataFile[];
  pagination: Pagination;
  status: 'idle' | 'loading' | 'succeeded' | 'failed';
  error: string | null;
  success: boolean;
}

const initialState: DataFileListSliceState = {
  dataFiles: [],
  pagination: { currentPage: 1, totalPages: 0, pageSize: 50, totalCount: 0 },
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
    setDataFileListPagination: (state, action: PayloadAction<Pagination>) => {
      state.pagination = action.payload;
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
  setDataFileListPagination,
  setDataFileListLoading,
  setDataFileListError,
  setDataFileListSuccess
} = dataFileListSlice.actions

export default dataFileListSlice.reducer
