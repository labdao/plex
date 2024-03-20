import type { ReduxState } from '@/lib/redux'

export const selectFlowUpdateLoading = (state: ReduxState) => state.flowUpdate.loading;
export const selectFlowUpdateError = (state: ReduxState) => state.flowUpdate.error;
export const selectFlowUpdateSuccess = (state: ReduxState) => state.flowUpdate.success;
