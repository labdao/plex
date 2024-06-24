import type { ReduxState } from '@/lib/redux'

export const selectModelList = (state: ReduxState) => state.modelList.models
export const selectModelListLoading = (state: ReduxState) => state.modelList.loading
export const selectModelListSuccess = (state: ReduxState) => state.modelList.success
export const selectModelListError = (state: ReduxState) => state.modelList.error
