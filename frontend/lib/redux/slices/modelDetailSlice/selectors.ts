import type { ReduxState } from "@/lib/redux";

export const selectModelDetail = (state: ReduxState) => state.modelDetail.model;
export const selectModelDetailLoading = (state: ReduxState) => state.modelDetail.loading;
export const selectModelDetailSuccess = (state: ReduxState) => state.modelDetail.success;
export const selectModelDetailError = (state: ReduxState) => state.modelDetail.error;
