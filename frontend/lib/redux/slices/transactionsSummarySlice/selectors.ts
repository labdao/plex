import type { ReduxState } from "@/lib/redux";

export const selectTransactionsSummary = (state: ReduxState) => state.transactionsSummary.summary;
export const selectTransactionsSummaryLoading = (state: ReduxState) => state.transactionsSummary.loading;
export const selectTransactionsSummarySuccess = (state: ReduxState) => state.transactionsSummary.success;
export const selectTransactionsSummaryError = (state: ReduxState) => state.transactionsSummary.error;
