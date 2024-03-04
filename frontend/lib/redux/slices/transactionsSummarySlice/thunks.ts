import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";

import { getTransactionsSummary } from "./asyncActions";
import { setTransactionsSummary, setTransactionsSummaryError, setTransactionsSummaryLoading, setTransactionsSummarySuccess } from "./slice";

export const transactionsSummaryThunk = createAppAsyncThunk("stripe/checkout", async (_, { dispatch }) => {
  dispatch(setTransactionsSummaryError(null));
  dispatch(setTransactionsSummaryLoading(true));
  try {
    const responseData = await getTransactionsSummary();
    dispatch(setTransactionsSummarySuccess(true));
    dispatch(setTransactionsSummary(responseData));
    dispatch(setTransactionsSummaryLoading(false));
    return responseData;
  } catch (error: unknown) {
    console.log("Problem getting checkout URL", error);
    if (error instanceof Error) {
      dispatch(setTransactionsSummaryError(error.message));
    } else {
      dispatch(setTransactionsSummaryError("Problem getting checkout URL"));
    }
    dispatch(setTransactionsSummaryLoading(false));
    return false;
  }
});
