import type { ReduxState } from '@/lib/redux'

export const selectFlowDetail = (state: ReduxState) => state.flowDetail.flow
export const selectFlowDetailLoading = (state: ReduxState) => state.flowDetail.loading
export const selectFlowDetailSuccess = (state: ReduxState) => state.flowDetail.success
export const selectFlowDetailError = (state: ReduxState) => state.flowDetail.error
