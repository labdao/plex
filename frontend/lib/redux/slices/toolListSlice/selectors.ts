import type { ReduxState } from '@/lib/redux'

export const selectToolList = (state: ReduxState) => state.toolList.tools
export const selectToolListLoading = (state: ReduxState) => state.toolList.loading
export const selectToolListSuccess = (state: ReduxState) => state.toolList.success
export const selectToolListError = (state: ReduxState) => state.toolList.error
