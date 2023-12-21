import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { listDataFiles } from './asyncActions'
import {
  setDataFileList,
  setDataFileListError,
  setDataFileListPagination,
  setDataFileListSuccess
} from './slice'

export const dataFileListThunk = createAppAsyncThunk(
  'datafiles/listDataFiles',
  async (arg: Partial<{ page: number, pageSize: number, filters: Record<string, string | undefined> }> = { page: 1, pageSize: 50, filters: {} }, { dispatch }) => {
    const { page = 1, pageSize = 50, filters = {} } = arg; 
    try {
      const response = await listDataFiles({ page, pageSize, filters });
      if (response) {
        dispatch(setDataFileListSuccess(true));
        dispatch(setDataFileList(response.data));
        dispatch(setDataFileListPagination(response.pagination)); 
      } else {
        console.log('Failed to list DataFiles.', response);
        dispatch(setDataFileListError('Failed to list DataFiles.'));
      }
      return response;
    } catch (error: unknown) {
      console.log('Failed to list DataFiles.', error);
      if (error instanceof Error) {
        dispatch(setDataFileListError(error.message));
      } else {
        dispatch(setDataFileListError('Failed to list DataFiles.'));
      }
      return false;
    }
  }
)
