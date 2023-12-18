import { createAsyncThunk } from '@reduxjs/toolkit';

import { saveDataFileToServer } from './actions';
import { addCid, addFilename, setDataFileError } from './dataSlice';

interface DataFilePayload {
  files: File[],
  metadata: { [key: string]: any },
  handleSuccess: () => void
}

interface FileResponse {
  filename: string;
  cid: string;
}

interface ServerResponse {
  cids: string[];
}

export const saveDataFilesAsync = createAsyncThunk<FileResponse[], DataFilePayload>(
  'dataFile/saveDataFiles',
  async ({ files, metadata, handleSuccess }, { dispatch }) => {
    try {
      const serverResponse = await saveDataFileToServer(files, metadata);
      console.log("Server Response:", serverResponse);

      // @ts-ignore
      const fileResponses: FileResponse[] = serverResponse.cids.map((cid: string, index: number) => {
        return { cid: cid, filename: files[index].name };
      });

      fileResponses.forEach(fileResponse => {
        dispatch(addCid(fileResponse.cid));
        dispatch(addFilename(fileResponse.filename));
      });

      handleSuccess();
      return fileResponses;
    } catch (error: unknown) {
      const errorMessage = (error as { message?: string }).message ?? 'An error occurred while saving data files.';
      dispatch(setDataFileError(errorMessage));
      throw error;
    }
  }
);