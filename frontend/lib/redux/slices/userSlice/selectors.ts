import type { ReduxState } from '@/lib/redux'

export const selectWalletAddress = (state: ReduxState) => state.user.walletAddress
export const selectEmailAddress = (state: ReduxState) => state.user.emailAddress
export const selectIsLoggedIn = (state: ReduxState) => state.user.isLoggedIn
