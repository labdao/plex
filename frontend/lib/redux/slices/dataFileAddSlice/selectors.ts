import type { ReduxState } from "@/lib/redux";

export const selectFilename = (state: ReduxState) => state.dataFileAdd.filename;
export const selectCID = (state: ReduxState) => state.dataFileAdd.cid;
export const selectDataFileError = (state: ReduxState) => state.dataFileAdd.error;
export const selectDataFileIsLoading = (state: ReduxState) => state.dataFileAdd.isLoading;
export const selectDateFileIsUploaded = (state: ReduxState) => state.dataFileAdd.isUploaded;
