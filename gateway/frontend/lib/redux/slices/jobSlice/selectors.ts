import type { ReduxState } from '@/lib/redux'

export const selectSelectedTool = (state: ReduxState) => state.job.selectedTool
export const selectSelectedDataFiles = (state: ReduxState) => state.job.selectedDataFiles
export const selectJobError = (state: ReduxState) => state.job.error
export const selectJobIsLoading = (state: ReduxState) => state.job.isLoading
export const selectJobIsInitialized = (state: ReduxState) => state.job.isInitialized