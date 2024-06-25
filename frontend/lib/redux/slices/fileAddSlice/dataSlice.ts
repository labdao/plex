import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { saveFileAsync } from "./thunks";

interface FileSliceState {
  filename: string;
  cid: string; // Content Identifier in IPFS
  isLoading: boolean;
  error: string | null;
  isUploaded: boolean;
}

const initialState: FileSliceState = {
  filename: "",
  cid: "",
  isLoading: false,
  error: null,
  isUploaded: false,
};

export const fileAddSlice = createSlice({
  name: "file",
  initialState,
  reducers: {
    setFilenameDataSlice: (state, action: PayloadAction<string>) => {
      state.filename = action.payload;
    },
    setCidDataSlice: (state, action: PayloadAction<string>) => {
      state.cid = action.payload;
    },
    setFileError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    startFileUploadDataSlice: (state) => {
      state.isLoading = true;
    },
    endFileUploadDataSlice: (state) => {
      state.isLoading = false;
    },
    setIsUploadedDataSlice: (state, action: PayloadAction<boolean>) => {
      state.isUploaded = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(saveFileAsync.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(saveFileAsync.fulfilled, (state, action) => {
        state.isLoading = false;
        if (action.payload) {
          console.log("action.payload", action.payload);
          state.cid = action.payload.cid;
          state.filename = action.payload.filename;
        }
        state.isUploaded = true;
      })
      .addCase(saveFileAsync.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || "An error occurred while saving file.";
      });
  },
});

export const { setFilenameDataSlice, setCidDataSlice, setFileError, startFileUploadDataSlice, endFileUploadDataSlice, setIsUploadedDataSlice } =
  fileAddSlice.actions;

export default fileAddSlice.reducer;
