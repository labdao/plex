import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface StripePayload {
  url: string | null;
}

export interface StripeCheckoutSliceState {
  url: string | null;
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: StripeCheckoutSliceState = {
  url: null,
  loading: false,
  error: null,
  success: false,
};

export const stripeCheckoutSlice = createSlice({
  name: "stripeCheckout",
  initialState,
  reducers: {
    setStripeCheckoutUrl: (state, action: PayloadAction<StripePayload>) => {
      state.url = action.payload?.url || null;
    },
    setStripeCheckoutLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setStripeCheckoutError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setStripeCheckoutSuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
  },
});

export const { setStripeCheckoutUrl, setStripeCheckoutLoading, setStripeCheckoutError, setStripeCheckoutSuccess } = stripeCheckoutSlice.actions;

export default stripeCheckoutSlice.reducer;
