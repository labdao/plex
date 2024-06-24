import type { ReduxState } from '@/lib/redux'

export const selectAddModelJson = (state: ReduxState) => state.modelAdd.modelJson
export const selectAddModelError = (state: ReduxState) => state.modelAdd.error
export const selectAddModelLoading = (state: ReduxState) => state.modelAdd.loading
export const selectAddModelSuccess = (state: ReduxState) => state.modelAdd.success
