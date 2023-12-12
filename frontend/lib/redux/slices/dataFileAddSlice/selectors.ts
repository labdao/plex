import type { ReduxState } from '@/lib/redux'

export const selectFilenames = (state: ReduxState) => state.dataFileAdd.filenames
export const selectCIDs = (state: ReduxState) => state.dataFileAdd.cids
export const selectDataFileError = (state: ReduxState) => state.dataFileAdd.error
export const selectDataFileIsLoading = (state: ReduxState) => state.dataFileAdd.isLoading
