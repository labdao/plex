import type { ReduxState } from '@/lib/redux'

export const selectJobDetail = (state: ReduxState) => state.jobDetail.job
export const selectJobDetailLoading = (state: ReduxState) => state.jobDetail.loading
export const selectJobDetailSuccess = (state: ReduxState) => state.jobDetail.success
export const selectJobDetailError = (state: ReduxState) => state.jobDetail.error
