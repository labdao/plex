"use client";

import { usePrivy } from "@privy-io/react-auth";
import { useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { toast } from "sonner";
import backendUrl from "@/lib/backendUrl";
import { getAccessToken } from "@privy-io/react-auth";
import {
  AppDispatch,
  selectStripeCheckoutError,
  selectStripeCheckoutLoading,
} from "@/lib/redux";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import StripeCheckoutButton from "@/components/payment/StripeCheckoutButton";

export default function SubscribePage() {
  const dispatch = useDispatch<AppDispatch>();
  const error = useSelector(selectStripeCheckoutError);
  const { user } = usePrivy();
  const walletAddress = user?.wallet?.address;

  const router = useRouter();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkSubscriptionStatus = async () => {
      let authToken;
      try {
        authToken = await getAccessToken();
      } catch (error) {
        console.log("Failed to get access token: ", error);
        return;
      }

      const response = await fetch(`${backendUrl()}/stripe/subscription/check`, {
        headers: {
          Authorization: `Bearer ${authToken}`,
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        const data = await response.json();
        if (data.isSubscribed) {
          router.replace("/subscription/manage");
        } else {
          setLoading(false);
        }
      } else {
        setLoading(false);
      }
    };

    checkSubscriptionStatus();
  }, [router]);

  useEffect(() => {
    if (error) {
      toast.error(error);
    }
  }, [error]);

  if (!walletAddress || loading) {
    return <div>Loading...</div>;
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
        <h3 className="text-center font-heading" style={{ fontSize: '29px', lineHeight: '43.2px', letterSpacing: '0.5px', color: '#000000'}}>
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
          <StripeCheckoutButton color="primary" size="sm" className="w-full font-bold">
            Start now
          </StripeCheckoutButton>
        </div>
      </div>
    </div>
  );
}