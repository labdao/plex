import type { ReduxState } from '@/lib/redux'

export const selectFlowList = (state: ReduxState) => state.flowList.flows
export const selectFlowListLoading = (state: ReduxState) => state.flowList.loading
export const selectFlowListSuccess = (state: ReduxState) => state.flowList.success
export const selectFlowListError = (state: ReduxState) => state.flowList.error
