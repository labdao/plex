import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";

import { saveFileToServer } from "./actions";
import { setCidDataSlice, setFileError, setFilenameDataSlice } from "./dataSlice";

interface FilePayload {
  file: File;
  metadata: { [key: string]: any };
  isPublic: boolean;
  handleSuccess: (cid: string) => void;
}

export const saveFileAsync = createAppAsyncThunk(
  "file/saveFile",
  async ({ file, metadata, isPublic, handleSuccess }: FilePayload, { dispatch }) => {
    try {
      const response = await saveFileToServer(file, metadata, isPublic);
      if (response.cid) {
        handleSuccess(response.cid);
      } else {
        dispatch(setFileError("Failed to save file."));
      }
      return response;
    } catch (error: unknown) {
      const errorMessage =
        typeof error === "object" && error !== null && "message" in error
          ? (error as { message?: string }).message
          : "An error occurred while saving file.";

      dispatch(setFileError(errorMessage || "An error occurred while saving file."));
    }
  }
);
