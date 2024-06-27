import type { ReduxState } from '@/lib/redux'

export const selectExperimentDetail = (state: ReduxState) => state.experimentDetail.experiment
export const selectExperimentDetailLoading = (state: ReduxState) => state.experimentDetail.loading
export const selectExperimentDetailSuccess = (state: ReduxState) => state.experimentDetail.success
export const selectExperimentDetailError = (state: ReduxState) => state.experimentDetail.error
