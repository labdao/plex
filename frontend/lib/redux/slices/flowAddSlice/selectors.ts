import type { ReduxState } from '@/lib/redux'

export const selectFlowAddName = (state: ReduxState) => state.flowAdd.name
export const selectFlowAddTool = (state: ReduxState) => state.flowAdd.tool
export const selectFlowAddKwargs = (state: ReduxState) => state.flowAdd.kwargs
export const selectFlowAddLoading = (state: ReduxState) => state.flowAdd.loading
export const selectFlowAddError = (state: ReduxState) => state.flowAdd.error
export const selectFlowAddCid = (state: ReduxState) => state.flowAdd.cid
export const selectFlowAddSuccess = (state: ReduxState) => state.flowAdd.success
