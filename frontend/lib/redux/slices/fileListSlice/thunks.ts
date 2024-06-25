import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { listFiles } from './asyncActions'
import {
  setFileList,
  setFileListError,
  setFileListPagination,
  setFileListSuccess
} from './slice'

export const fileListThunk = createAppAsyncThunk(
  'files/listFiles',
  async (arg: Partial<{ page: number, pageSize: number, filters: Record<string, string | undefined> }> = { page: 1, pageSize: 50, filters: {} }, { dispatch }) => {
    const { page = 1, pageSize = 50, filters = {} } = arg; 
    try {
      const response = await listFiles({ page, pageSize, filters });
      if (response) {
        dispatch(setFileListSuccess(true));
        dispatch(setFileList(response.data));
        dispatch(setFileListPagination(response.pagination)); 
      } else {
        console.log('Failed to list Files.', response);
        dispatch(setFileListError('Failed to list Files.'));
      }
      return response;
    } catch (error: unknown) {
      console.log('Failed to list Files.', error);
      if (error instanceof Error) {
        dispatch(setFileListError(error.message));
      } else {
        dispatch(setFileListError('Failed to list Files.'));
      }
      return false;
    }
  }
)
