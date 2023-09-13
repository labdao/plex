import type { ReduxState } from '@/lib/redux'

export const selectUsername = (state: ReduxState) => state.user.username
export const selectWalletAddress = (state: ReduxState) => state.user.walletAddress
export const selectUserFormError = (state: ReduxState) => state.user.error
export const selectUserFormIsLoading = (state: ReduxState) => state.user.isLoading
export const selectIsLoggedIn = (state: ReduxState) => state.user.isLoggedIn
