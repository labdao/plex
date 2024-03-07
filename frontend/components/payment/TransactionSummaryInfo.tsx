import { CreditCardIcon } from "lucide-react";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import {
  AppDispatch,
  selectTransactionsSummary,
  selectTransactionsSummaryError,
  selectTransactionsSummaryLoading,
  transactionsSummaryThunk,
} from "@/lib/redux";
import { cn } from "@/lib/utils";

import StripeCheckoutButton from "../payment/StripeCheckoutButton";

const InfoItem = ({ label, value }: { label: string; value: string | number | null }) => (
  <div className="flex items-center gap-1">
    <span className="font-mono text-xs font-bold uppercase">{label}:</span>
    <span className="text-sm">{value}</span>
  </div>
);

const TransactionSummaryInfo = ({ className }: { className?: string }) => {
  const dispatch = useDispatch<AppDispatch>();
  const transactionsSummary = useSelector(selectTransactionsSummary);
  const loading = useSelector(selectTransactionsSummaryLoading);
  const error = useSelector(selectTransactionsSummaryError);

  useEffect(() => {
    dispatch(transactionsSummaryThunk());
  }, [dispatch]);

  const { tokens, balance } = transactionsSummary;

  return (
    <div className={cn("flex flex-wrap justify-between gap-4 p-2 rounded-lg bg-primary/10", className)}>
      <div className="flex gap-4">
        <InfoItem label="Points" value={tokens} />
        <InfoItem label="Credits" value={balance} />
      </div>
      <StripeCheckoutButton variant="outline" size="sm">
        <CreditCardIcon size={20} className="mr-1" />
        Add Credits
      </StripeCheckoutButton>
    </div>
  );
};

export default TransactionSummaryInfo;
