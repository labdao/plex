import type { ReduxState } from "@/lib/redux";

export const selectToolDetail = (state: ReduxState) => state.toolDetail.tool;
export const selectToolDetailLoading = (state: ReduxState) => state.toolDetail.loading;
export const selectToolDetailSuccess = (state: ReduxState) => state.toolDetail.success;
export const selectToolDetailError = (state: ReduxState) => state.toolDetail.error;
