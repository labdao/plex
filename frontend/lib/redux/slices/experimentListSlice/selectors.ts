import type { ReduxState } from '@/lib/redux'

export const selectExperimentList = (state: ReduxState) => state.experimentList.experiments
export const selectExperimentListLoading = (state: ReduxState) => state.experimentList.loading
export const selectExperimentListSuccess = (state: ReduxState) => state.experimentList.success
export const selectExperimentListError = (state: ReduxState) => state.experimentList.error
export const selectCategorizedExperiments = (state: ReduxState) => state.experimentList.categorizedExperiments
