// redux/experimentNames/selectors.ts
import type { ReduxState } from '@/lib/redux';  // Correct the path as needed

export const selectExperimentNames = (state: ReduxState) => state.experimentNames.names;
export const selectExperimentNamesLoading = (state: ReduxState) => state.experimentNames.loading;
export const selectExperimentNamesError = (state: ReduxState) => state.experimentNames.error;
export const selectExperimentNamesSuccess = (state: ReduxState) => state.experimentNames.success;
export const selectCategorizedExperimentNames = (state: ReduxState) => state.experimentNames.categorizedExperimentNames;