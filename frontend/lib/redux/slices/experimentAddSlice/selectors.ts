import type { ReduxState } from '@/lib/redux'

export const selectExperimentAddName = (state: ReduxState) => state.experimentAdd.name
export const selectExperimentAddTool = (state: ReduxState) => state.experimentAdd.tool
export const selectExperimentAddKwargs = (state: ReduxState) => state.experimentAdd.kwargs
export const selectExperimentAddLoading = (state: ReduxState) => state.experimentAdd.loading
export const selectExperimentAddError = (state: ReduxState) => state.experimentAdd.error
export const selectExperimentAddID = (state: ReduxState) => state.experimentAdd.ID
export const selectExperimentAddSuccess = (state: ReduxState) => state.experimentAdd.success
