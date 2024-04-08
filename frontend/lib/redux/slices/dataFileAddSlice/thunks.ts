import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";

import { saveDataFileToServer } from "./actions";
import { setCidDataSlice, setDataFileError, setFilenameDataSlice } from "./dataSlice";

interface DataFilePayload {
  file: File;
  metadata: { [key: string]: any };
  isPublic: boolean;
  handleSuccess: (cid: string) => void;
}

export const saveDataFileAsync = createAppAsyncThunk(
  "dataFile/saveDataFile",
  async ({ file, metadata, isPublic, handleSuccess }: DataFilePayload, { dispatch }) => {
    try {
      const response = await saveDataFileToServer(file, metadata, isPublic);
      if (response.cid) {
        handleSuccess(response.cid);
      } else {
        dispatch(setDataFileError("Failed to save data file."));
      }
      return response;
    } catch (error: unknown) {
      const errorMessage =
        typeof error === "object" && error !== null && "message" in error
          ? (error as { message?: string }).message
          : "An error occurred while saving data file.";

      dispatch(setDataFileError(errorMessage || "An error occurred while saving data file."));
    }
  }
);
