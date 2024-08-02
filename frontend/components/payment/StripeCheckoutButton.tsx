import { Slot } from "@radix-ui/react-slot";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { toast } from "sonner";

import { ButtonProps } from "@/components/ui/button";
import { Button } from "@/components/ui/button";
import { AppDispatch, selectStripeCheckoutError, selectStripeCheckoutLoading, selectStripeCheckoutUrl, setStripeCheckoutError, setStripeCheckoutLoading, setStripeCheckoutSuccess, setStripeCheckoutUrl } from "@/lib/redux";
import { getCheckoutURL } from "@/lib/redux/slices/stripeCheckoutSlice/asyncActions";

import { PageLoader } from "../shared/PageLoader";
import { AlertDialog, AlertDialogContent } from "../ui/alert-dialog";

const StripeCheckoutButton = (props: ButtonProps) => {
  const Comp = props.asChild ? Slot : Button;

  const dispatch = useDispatch<AppDispatch>();
  const checkoutUrl = useSelector(selectStripeCheckoutUrl);
  const loading = useSelector(selectStripeCheckoutLoading);
  const error = useSelector(selectStripeCheckoutError);

  const handleCheckout = async () => {
    dispatch(setStripeCheckoutError(null));
    dispatch(setStripeCheckoutLoading(true));
    try {
      const checkoutResponse = await getCheckoutURL();
      if (checkoutResponse?.url) {
        dispatch(setStripeCheckoutSuccess(true));
        dispatch(setStripeCheckoutUrl({ url: checkoutResponse.url }));
        dispatch(setStripeCheckoutLoading(false));
        window.location.assign(checkoutResponse.url);
      } else {
        throw new Error("Failed to get checkout URL");
      }
    } catch (error: unknown) {
      console.log("Failed to initiate checkout", error);
      if (error instanceof Error) {
        dispatch(setStripeCheckoutError(error.message));
      } else {
        dispatch(setStripeCheckoutError("Failed to get checkout URL"));
      }
      dispatch(setStripeCheckoutLoading(false));
      toast.error("Failed to get checkout URL");
    }
  };

  useEffect(() => {
    if (error) {
      toast.error(error);
    }
  }, [error]);

  return (
    <>
      <AlertDialog open={loading}>
        <AlertDialogContent className="text-center">
          <PageLoader className="py-0" />
          <h4>Sending you to Stripe for subscription</h4>
        </AlertDialogContent>
      </AlertDialog>
      <Comp onClick={handleCheckout} {...props}>
        {props.children}
      </Comp>
    </>
  );
};

export default StripeCheckoutButton;
