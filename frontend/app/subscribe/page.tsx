"use client";

import { usePrivy } from "@privy-io/react-auth";
import { useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { toast } from "sonner";
import { AppDispatch, selectStripeCheckoutError, selectStripeCheckoutLoading, selectStripeCheckoutUrl, stripeCheckoutThunk } from "@/lib/redux";
import { Button } from "@/components/ui/button";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import { PageLoader } from "@/components/shared/PageLoader";
import { AlertDialog, AlertDialogContent } from "@/components/ui/alert-dialog";

export default function SubscribePage() {
  const dispatch = useDispatch<AppDispatch>();
  const checkoutUrl = useSelector(selectStripeCheckoutUrl);
  const loading = useSelector(selectStripeCheckoutLoading);
  const error = useSelector(selectStripeCheckoutError);
  const { user } = usePrivy();
  const router = useRouter();
  const walletAddress = user?.wallet?.address;

  useEffect(() => {
    if (checkoutUrl) {
      window.location.assign(checkoutUrl);
    }
    if (error) {
      toast.error(error);
    }
  }, [checkoutUrl, error]);

  const handleCheckout = async () => {
    dispatch(stripeCheckoutThunk());
  };

  if (!walletAddress) {
    return <div>Loading...</div>; // Show a loading state while fetching wallet address
  }

  return (
    <div className="relative flex flex-col h-screen max-w-full grow">
      <Breadcrumbs
        items={[
          { name: "Subscribe", href: "/subscribe" },
          { name: walletAddress, href: `/subscribe/${walletAddress}` },
        ]}
        actions={null}
      />
      <div className="flex flex-col items-center justify-between w-[706px] h-[469px] p-4 bg-white rounded-lg shadow-lg mx-auto my-6">
        <h3 className="text-center font-heading" style={{ fontSize: '29px', lineHeight: '43.2px', letterSpacing: '0.5px', color: '#000000' }}>
          Become a lab.bio subscriber
        </h3>
        <ul className="space-y-4 w-full font-heading" style={{ fontSize: '16px', lineHeight: '28px', letterSpacing: '0.3px', color: '#000000' }}>
          <li className="flex items-center">
            <span className="mr-2 text-black">✓</span>
            <span>Access x# of computation credits (about x number per x number)</span>
          </li>
          <li className="flex items-center">
            <span className="mr-2 text-black">✓</span>
            <span>Additional charges information Additional charges information</span>
          </li>
          <li className="flex items-center">
            <span className="mr-2 text-black">✓</span>
            <span>Additional charges information Additional charges information</span>
          </li>
          <li className="flex items-center">
            <span className="mr-2 text-black">✓</span>
            <span>Cancel subscription any time</span>
          </li>
        </ul>
        <p className="mt-4 text-center font-heading" style={{ fontSize: '24px', lineHeight: '30px', letterSpacing: '0.14px', color: '#000000', fontWeight: '500' }}>
          X$ per month
        </p>
        <div className="px-2 py-2 w-full">
          <Button color="primary" size="sm" className="w-full font-bold" onClick={handleCheckout}>
            Start now
          </Button>
        </div>
      </div>
      <AlertDialog open={loading}>
        <AlertDialogContent className="text-center">
          <PageLoader className="py-0" />
          <h4>Sending you to Stripe for subscription</h4>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
