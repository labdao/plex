import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk';

import { saveDataFileToServer } from './actions';
import { addCid, addFilename, setDataFileError } from './dataSlice';

interface DataFilePayload {
  files: File[],
  metadata: { [key: string]: any },
  handleSuccess: () => void
}

export const saveDataFilesAsync = createAppAsyncThunk(
  'dataFile/saveDataFiles',
  async ({ files, metadata, handleSuccess }: DataFilePayload, { dispatch }) => {
    const responses = [];
    for (const file of files) {
      try {
        const response = await saveDataFileToServer(file, metadata);
        console.log("Response:", response);
        if (response.cid) {
          dispatch(addCid(response.cid));
          dispatch(addFilename(response.filename));
          responses.push(response);
        } else {
          dispatch(setDataFileError(`Failed to save data file: ${file.name}.`));
        }
      } catch (error: unknown) {
        const errorMessage = (error as { message?: string }).message ?? `An error occurred while saving data file: ${file.name}.`;
        dispatch(setDataFileError(errorMessage));
      }
    }
    if (responses.length === files.length) {
      handleSuccess();
    }
    return responses;
  }
);