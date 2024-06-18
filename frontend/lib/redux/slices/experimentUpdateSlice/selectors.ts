import type { ReduxState } from '@/lib/redux'

export const selectExperimentUpdateLoading = (state: ReduxState) => state.experimentUpdate.loading;
export const selectExperimentUpdateError = (state: ReduxState) => state.experimentUpdate.error;
export const selectExperimentUpdateSuccess = (state: ReduxState) => state.experimentUpdate.success;
