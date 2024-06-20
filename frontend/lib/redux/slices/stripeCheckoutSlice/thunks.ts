import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";

import { getCheckoutURL } from "./asyncActions";
import { setStripeCheckoutError, setStripeCheckoutLoading, setStripeCheckoutSuccess, setStripeCheckoutUrl } from "./slice";

export const stripeCheckoutThunk = createAppAsyncThunk("stripe/checkout", async (_, { dispatch }) => {
  dispatch(setStripeCheckoutError(null));
  dispatch(setStripeCheckoutLoading(true));
  try {
    // const responseData = await getCheckoutURL();
    // dispatch(setStripeCheckoutSuccess(true));
    // dispatch(setStripeCheckoutUrl(responseData));
    // dispatch(setStripeCheckoutLoading(false));
    // return responseData;
  } catch (error: unknown) {
    console.log("Problem getting checkout URL", error);
    if (error instanceof Error) {
      dispatch(setStripeCheckoutError(error.message));
    } else {
      dispatch(setStripeCheckoutError("Problem getting checkout URL"));
    }
    dispatch(setStripeCheckoutLoading(false));
    return false;
  }
});
