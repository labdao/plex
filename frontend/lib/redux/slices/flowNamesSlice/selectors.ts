// redux/flowNames/selectors.ts
import type { ReduxState } from '@/lib/redux';  // Correct the path as needed

export const selectFlowNames = (state: ReduxState) => state.flowNames.names;
export const selectFlowNamesLoading = (state: ReduxState) => state.flowNames.loading;
export const selectFlowNamesError = (state: ReduxState) => state.flowNames.error;
export const selectFlowNamesSuccess = (state: ReduxState) => state.flowNames.success;
export const selectCategorizedFlowNames = (state: ReduxState) => state.flowNames.categorizedFlowNames;