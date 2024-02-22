import { useLogin } from "@privy-io/react-auth";
import { usePrivy } from "@privy-io/react-auth";
import { Slot } from "@radix-ui/react-slot";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { toast } from "sonner";

import { ButtonProps } from "@/components/ui/button";
import { Button } from "@/components/ui/button";
import { AppDispatch, selectStripeCheckoutError, selectStripeCheckoutLoading, selectStripeCheckoutUrl, stripeCheckoutThunk } from "@/lib/redux";

import { PageLoader } from "../shared/PageLoader";
import { AlertDialog, AlertDialogContent } from "../ui/alert-dialog";

const StripeCheckoutButton = (props: ButtonProps) => {
  const Comp = props.asChild ? Slot : Button;

  const dispatch = useDispatch<AppDispatch>();
  const checkoutUrl = useSelector(selectStripeCheckoutUrl);
  const loading = useSelector(selectStripeCheckoutLoading);
  const error = useSelector(selectStripeCheckoutError);

  const handleCheckout = async () => {
    dispatch(stripeCheckoutThunk());
  };

  useEffect(() => {
    if (checkoutUrl) {
      window.location.assign(checkoutUrl);
    }
    if (error) {
      toast.error(error);
    }
  }, [checkoutUrl, error]);

  return (
    <span>
      <AlertDialog open={loading}>
        <AlertDialogContent className="text-center">
          <PageLoader className="py-0" />
          <h4>Sending you to Stripe to add credits</h4>
        </AlertDialogContent>
      </AlertDialog>
      <Comp onClick={handleCheckout} {...props}>
        {props.children}
      </Comp>
    </span>
  );
};

export default StripeCheckoutButton;
