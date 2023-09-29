import type { ReduxState } from '@/lib/redux'

export const selectAddToolJson = (state: ReduxState) => state.toolAdd.toolJson
export const selectAddToolError = (state: ReduxState) => state.toolAdd.error
export const selectAddToolLoading = (state: ReduxState) => state.toolAdd.loading
export const selectAddToolSuccess = (state: ReduxState) => state.toolAdd.success
