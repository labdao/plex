import type { ReduxState } from "@/lib/redux";

export const selectFilename = (state: ReduxState) => state.fileAdd.filename;
export const selectCID = (state: ReduxState) => state.fileAdd.id;
export const selectFileError = (state: ReduxState) => state.fileAdd.error;
export const selectFileIsLoading = (state: ReduxState) => state.fileAdd.isLoading;
export const selectDateFileIsUploaded = (state: ReduxState) => state.fileAdd.isUploaded;
