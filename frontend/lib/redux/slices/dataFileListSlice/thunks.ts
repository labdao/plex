import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { listDataFiles } from './asyncActions'
import { setDataFileList,setDataFileListError, setDataFileListSuccess } from './slice'


export const dataFileListThunk = createAppAsyncThunk(
  'datafiles/listDataFiles',
  async (_, { dispatch }) => {
    try {
      const response = await listDataFiles()
      if (response) {
        dispatch(setDataFileListSuccess(true))
        dispatch(setDataFileList(response))
      } else {
        console.log('Failed to list DataFiles.', response)
        dispatch(setDataFileListError('Failed to list DataFiles.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to list DataFiles.', error)
      if (error instanceof Error) {
        dispatch(setDataFileListError(error.message))
      } else {
        dispatch(setDataFileListError('Failed to list DataFiles.'))
      }
      return false
    }
  }
)
