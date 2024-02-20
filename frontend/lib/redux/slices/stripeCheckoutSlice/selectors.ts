import type { ReduxState } from "@/lib/redux";

export const selectStripeCheckoutUrl = (state: ReduxState) => state.stripeCheckout.url;
export const selectStripeCheckoutLoading = (state: ReduxState) => state.stripeCheckout.loading;
export const selectStripeCheckoutSuccess = (state: ReduxState) => state.stripeCheckout.success;
export const selectStripeCheckoutError = (state: ReduxState) => state.stripeCheckout.error;
