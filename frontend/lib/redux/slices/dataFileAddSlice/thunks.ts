import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { saveDataFileToServer } from './actions'
import { setCidDataSlice, setFilenameDataSlice, setDataFileError } from './dataSlice'

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
        dispatch(setFilenameDataSlice(response.filename));
      } else {
        dispatch(setDataFileError('Failed to save data file.'))
      }
      return response;
    } catch (error: unknown) {
      const errorMessage = typeof error === 'object' && error !== null && 'message' in error
        ? (error as { message?: string }).message
        : 'An error occurred while saving data file.';

      dispatch(setDataFileError(errorMessage || 'An error occurred while saving data file.'));
    }
  }
)
