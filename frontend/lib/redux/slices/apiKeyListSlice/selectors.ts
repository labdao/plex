import type { ReduxState } from '@/lib/redux'

export const selectApiKeyList = (state: ReduxState) => state.apiKeyList.apiKeys
export const selectApiKeyListLoading = (state: ReduxState) => state.apiKeyList.loading
export const selectApiKeyListSuccess = (state: ReduxState) => state.apiKeyList.success
export const selectApiKeyListError = (state: ReduxState) => state.apiKeyList.error
