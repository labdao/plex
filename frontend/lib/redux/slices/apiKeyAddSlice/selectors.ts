import type { ReduxState } from "@/lib/redux"

export const selectApiKeyAdd = (state: ReduxState) => state.apiKeyAdd.key
export const selectApiKeyAddLoading = (state: ReduxState) => state.apiKeyAdd.loading
export const selectApiKeyAddError = (state: ReduxState) => state.apiKeyAdd.error
export const selectApiKeyAddSuccess = (state: ReduxState) => state.apiKeyAdd.success