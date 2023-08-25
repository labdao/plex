import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { setError } from '@/lib/redux'
import { saveDataFileToServer } from './actions'

interface DataFilePayload {
  file: File,
  metadata: { [key: string]: any }
}

export const saveDataFileAsync = createAppAsyncThunk(
  'dataFile/saveDataFile',
  async ({ file, metadata }: DataFilePayload, { dispatch }) => {
    try {
      const response = await saveDataFileToServer(file, metadata)

      if (response.filename && response.cid) {
        // You can optionally store something in localStorage or do other operations here.
        // For instance: localStorage.setItem('filename', response.filename)

      } else {
        dispatch(setError('Failed to save data file.'))
      }
      return response
    } catch (error: unknown) {
      const errorMessage = typeof error === 'object' && error !== null && 'message' in error 
        ? (error as { message?: string }).message 
        : undefined;

      dispatch(setError(errorMessage || 'An error occurred while saving data file.'));
      return false;
    }
  }
)
