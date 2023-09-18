import type { ReduxState } from '@/lib/redux'

export const selectToolFilename = (state: ReduxState) => state.toolAdd.filename
export const selectToolCID = (state: ReduxState) => state.toolAdd.cid
export const selectToolError = (state: ReduxState) => state.toolAdd.error
export const selectToolIsLoading = (state: ReduxState) => state.toolAdd.isLoading
export const selectToolIsUploaded = (state: ReduxState) => state.toolAdd.isUploaded