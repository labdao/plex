import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface File {
  Filename: string
  CID: string
  Tags: Tag[]
  S3URI: string
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

interface FileListSliceState {
  files: File[];
  pagination: Pagination;
  status: 'idle' | 'loading' | 'succeeded' | 'failed';
  error: string | null;
  success: boolean;
}

const initialState: FileListSliceState = {
  files: [],
  pagination: { currentPage: 1, totalPages: 0, pageSize: 50, totalCount: 0 },
  status: 'idle',
  error: null,
  success: false,
};

export const fileListSlice = createSlice({
  name: 'fileList',
  initialState,
  reducers: {
    setFileList: (state, action: PayloadAction<File[]>) => {
      state.files = action.payload;
    },
    setFileListPagination: (state, action: PayloadAction<Pagination>) => {
      state.pagination = action.payload;
    },
    setFileListLoading: (state, action: PayloadAction<boolean>) => {
      state.status = action.payload ? 'loading' : 'idle';
    },
    setFileListError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setFileListSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
      state.status = action.payload ? 'succeeded' : 'failed';
    },
  },
});

export const {
  setFileList,
  setFileListPagination,
  setFileListLoading,
  setFileListError,
  setFileListSuccess
} = fileListSlice.actions

export default fileListSlice.reducer
