import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface TransactionsSummary {
  tokens: number | null;
  balance: number | null;
}

export interface TransactionsSummarySliceState {
  summary: TransactionsSummary;
  loading: boolean;
  error: string | null;
  success: boolean;
}

const initialState: TransactionsSummarySliceState = {
  summary: {
    tokens: null,
    balance: null,
  },
  loading: false,
  error: null,
  success: false,
};

export const transactionsSummarySlice = createSlice({
  name: "jobDetail",
  initialState,
  reducers: {
    setTransactionsSummary: (state, action: PayloadAction<TransactionsSummary>) => {
      state.summary = action.payload;
    },
    setTransactionsSummaryLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setTransactionsSummaryError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setTransactionsSummarySuccess: (state, action: PayloadAction<boolean>) => {
      state.success = action.payload;
    },
  },
});

export const { setTransactionsSummary, setTransactionsSummaryLoading, setTransactionsSummaryError, setTransactionsSummarySuccess } =
  transactionsSummarySlice.actions;

export default transactionsSummarySlice.reducer;
