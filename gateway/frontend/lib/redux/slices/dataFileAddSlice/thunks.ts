import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { saveDataFileToServer } from './actions'
import { setCid, setFilename, setError } from './dataSlice'

interface DataFilePayload {
  file: File,
  metadata: { [key: string]: any }
}

export const saveDataFileAsync = createAppAsyncThunk(
  'dataFile/saveDataFile',
  async ({ file, metadata }: DataFilePayload, { dispatch }) => {
    try {
      const response = await saveDataFileToServer(file, metadata);
      console.log("Response:", response)
      if (response.filename) {
        dispatch(setFilename(response.filename));
      } else {
        dispatch(setError('Failed to save data file.'))
      }
      return response;
    } catch (error: unknown) {
      const errorMessage = typeof error === 'object' && error !== null && 'message' in error
        ? (error as { message?: string }).message
        : 'An error occurred while saving data file.';

      dispatch(setError(errorMessage || 'An error occurred while saving data file.'));
    }
  }
)
