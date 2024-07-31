"use client";

import { usePrivy } from "@privy-io/react-auth";
import React, { useEffect, useState } from "react";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import { toast } from "sonner";
import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl";
import { useRouter } from "next/navigation";
import getPlanTemplate, { PlanDetail } from "lib/planTemplate";
import StripeCheckoutButton from "@/components/payment/StripeCheckoutButton";

interface PlanDetails {
  plan_name: string;
  plan_amount: number;
  plan_currency: string;
  plan_interval: string;
  included_credits: number;
  overage_charge: number;
}

export default function SubscribePage() {
  const { user } = usePrivy();
  const walletAddress = user?.wallet?.address;
  const [loading, setLoading] = useState(true);
  const [planDetails, setPlanDetails] = useState<PlanDetails | null>(null);
  const router = useRouter();

  useEffect(() => {
    const fetchPlanDetails = async () => {
      try {
        const authToken = await getAccessToken();
        const response = await fetch(`${backendUrl()}/stripe/plan-details`, {
          headers: {
            Authorization: `Bearer ${authToken}`,
            "Content-Type": "application/json",
          },
        });

        if (response.ok) {
          const data: PlanDetails = await response.json();
          setPlanDetails(data);
          setLoading(false);
        } else {
          console.error("Failed to fetch plan details. Response not OK.");
          setLoading(false);
        }
      } catch (error) {
        console.error("Failed to fetch plan details:", error);
        setLoading(false);
      }
    };

    const checkSubscriptionStatus = async () => {
      try {
        const authToken = await getAccessToken();

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
            fetchPlanDetails();
          }
        } else {
          console.error("Failed to check subscription status. Response not OK.");
          setLoading(false);
        }
      } catch (error) {
        console.error("Failed to check subscription status:", error);
        setLoading(false);
      }
    };

    if (walletAddress) {
      checkSubscriptionStatus();
    } else {
      console.log("Waiting for wallet address...");
    }
  }, [router, walletAddress]);

  if (loading || !walletAddress || !planDetails) {
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
      <div className="flex flex-col items-center justify-between w-[706px] h-[400px] p-4 bg-white rounded-lg shadow-lg mx-auto my-6">
        <h3 className="text-center font-heading" style={{ fontSize: '29px', lineHeight: '43.2px', letterSpacing: '0.5px', color: '#000000'}}>
          Become a lab.bio subscriber
        </h3>
        <div className="text-sm  text-gray-600 space-y-4 font-heading" style={{ fontSize: '16px', lineHeight: '28px', letterSpacing: '0.3px', color: '#000000' }}>
          {getPlanTemplate().details.map((detail: PlanDetail, index: number) => (
            <div key={index} className="flex items-start">
              <span className="mr-2 text-black">✓</span>
              <span>{detail.description
                .replace('{{includedCredits}}', planDetails.included_credits.toString())
                .replace('{{numMolecules}}', (planDetails.included_credits / 10).toString()) // Example calculation
                .replace('{{overageCharge}}', planDetails.overage_charge.toString())
              }</span>
            </div>
          ))}
        </div>
        <p className="mt-4 text-center font-heading" style={{ fontSize: '24px', lineHeight: '30px', letterSpacing: '0.14px', color: '#000000', fontWeight: '500' }}>
          ${planDetails.plan_amount} / month
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
