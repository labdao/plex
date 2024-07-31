"use client";

import { usePrivy } from "@privy-io/react-auth";
import React, { useEffect, useState } from "react";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import { toast } from "sonner";
import { getAccessToken } from "@privy-io/react-auth";
import { AlertDialog, AlertDialogContent, AlertDialogOverlay, AlertDialogFooter } from "@/components/ui/alert-dialog";
import backendUrl from "lib/backendUrl";
import { useRouter } from "next/navigation";
import getPlanTemplate, { PlanDetail, PlanTemplate } from "lib/planTemplate";

interface SubscriptionDetails {
  plan_name: string;
  plan_amount: number;
  plan_currency: string;
  plan_interval: string;
  current_period_start: string;
  current_period_end: string;
  next_due: string;
  status: string;
  included_credits: number;
  used_credits: number;
  overage_charge: number;
  cancel_at_period_end: boolean;
}

export default function ManageSubscription() {
  const { user } = usePrivy();
  const walletAddress = user?.wallet?.address;
  const [loading, setLoading] = useState(true);
  const [subscriptionDetails, setSubscriptionDetails] = useState<SubscriptionDetails | null>(null);
  const router = useRouter();

  const handleSubscriptionAction = async () => {
    try {
      let authToken = await getAccessToken();

      const response = await fetch(`${backendUrl()}/stripe/billing-portal`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${authToken}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          returnURL: window.location.href, // Return to current page after managing billing
        }),
      });

      const data = await response.json();
      if (data.url) {
        window.location.href = data.url; // Redirect to the Stripe billing portal
      } else {
        throw new Error("Failed to create billing portal session");
      }
    } catch (error) {
      console.error("Error redirecting to billing portal", error);
      toast.error("Failed to open billing portal");
    }
  };

  useEffect(() => {
    const fetchSubscriptionDetails = async () => {
      try {
        let authToken = await getAccessToken();

        const response = await fetch(`${backendUrl()}/stripe/subscription`, {
          headers: {
            Authorization: `Bearer ${authToken}`,
            "Content-Type": "application/json",
          },
        });

        if (response.ok) {
          const data: SubscriptionDetails = await response.json();
          setSubscriptionDetails(data);
          setLoading(false);
        } else {
          setLoading(false);
        }
      } catch (error) {
        console.error("Failed to fetch subscription details:", error);
      }
    };

    const checkSubscriptionStatus = async () => {
      try {
        let authToken = await getAccessToken();

        const response = await fetch(`${backendUrl()}/stripe/subscription/check`, {
          headers: {
            Authorization: `Bearer ${authToken}`,
            "Content-Type": "application/json",
          },
        });

        if (response.ok) {
          const data = await response.json();
          if (!data.isSubscribed) {
            router.replace("/subscribe");
          } else {
            fetchSubscriptionDetails();
          }
        } else {
          setLoading(false);
        }
      } catch (error) {
        console.error("Failed to check subscription status:", error);
        setLoading(false);
      }
    };

    checkSubscriptionStatus();
  }, [router]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!walletAddress || !subscriptionDetails) {
    return <div>Loading...</div>;
  }

  const showRenewalInfo = !subscriptionDetails.cancel_at_period_end;
  const isTrialingButNotRenewing = subscriptionDetails.status === 'trialing' && subscriptionDetails.cancel_at_period_end;

  return (
    <div className="relative flex flex-col h-screen max-w-full grow">
      <Breadcrumbs
        items={[
          { name: "subscription/manage", href: "/subscription/manage" },
          { name: walletAddress, href: `/subscription/${walletAddress}` },
        ]}
        actions={null}
      />
      <div className="flex justify-between space-x-6 mx-auto my-6">
        <div className="flex flex-col space-y-6">
          <div className="w-[450px] p-6 bg-white rounded-lg shadow-lg">
            <h3 className=" font-heading text-black">Usage Details</h3>
            <br/>
            <div className="text-sm text-gray-600 space-y-1.5 font-mono">
              <div className="flex justify-between ">
                <span style={{ color: '#808080' }}>Used / Included credits</span>
                <span style={{ color: '#000000' }}>
                  {subscriptionDetails.used_credits} / {subscriptionDetails.included_credits} credits
                </span>
              </div>
              <div className="flex justify-between">
                <span style={{ color: '#808080' }}>Current cycle overage charges</span>
                <span style={{ color: '#000000' }}>
                  ${subscriptionDetails.overage_charge}
                </span>
              </div>
            </div>
          </div>
          <div className="w-[450px] p-6 bg-white rounded-lg shadow-lg font-heading">
            <h3 className="font-heading text-black">Billing & Payment</h3><br/>
            <div className="text-sm text-gray-600 space-y-1.5 font-mono">
              <div className="flex justify-between">
                <span style={{ color: '#808080' }}>Cost until {subscriptionDetails.next_due}</span>
                <span style={{ color: '#000000' }}>
                  {subscriptionDetails.status === 'trialing' ? '$0' : `$${subscriptionDetails.plan_amount + Math.max(0, (subscriptionDetails.used_credits - subscriptionDetails.included_credits) * subscriptionDetails.overage_charge)}`}
                </span>
              </div>
              <div className="flex justify-between">
                <span style={{ color: '#808080' }}>Billing period</span>
                <span style={{ color: '#000000' }}>{subscriptionDetails.plan_interval}</span>
              </div>
              {showRenewalInfo ? (
                <>
                  <div className="flex justify-between">
                    <span style={{ color: '#808080' }}>Plan renews on</span>
                    <span style={{ color: '#000000' }}>{subscriptionDetails.next_due}</span>
                  </div>
                  <div className="flex justify-between">
                    <span style={{ color: '#808080' }}>Plan amount</span>
                    <span style={{ color: '#000000' }}>
                      $5 + overage charges{/* TODO: should be changed to the dynamic value */}
                    </span>
                  </div>
                </>
              ) : (
                <div className="flex justify-between">
                  <span style={{ color: '#808080' }}>Auto renewal</span>
                  <span style={{ color: '#FF0000' }}>Off</span>
                </div>
              )}
            </div>
            <div className="flex justify-between mt-6 space-x-4">
              <button
                className="px-3 py-1 border rounded-md font-mono text-black uppercase tracking-wider"
                style={{ borderColor: '#6BDBAD', fontSize: '12px'}}
                onClick={handleSubscriptionAction}
              >
                {showRenewalInfo ? 'Edit Billing' : 'Auto Renew'}
              </button> 
              <button
                className={`px-3 py-1 border rounded-md font-mono uppercase tracking-wider ${
                  showRenewalInfo ? 'text-gray-500' : 'text-black'
                }`}
                style={{ fontSize: '12px', borderColor: showRenewalInfo ? '#000000' : '#6BDBAD' }}
                onClick={handleSubscriptionAction}
              >
                {showRenewalInfo ? 'Cancel Subscription' : 'Subscribe Again'}
              </button>
            </div>

          </div>
        </div>
        <div className="w-96 p-6 bg-white rounded-lg shadow-lg">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-xl font-bold font-heading text-black">Your Plan</h3>
            <div className="flex items-center border border-green-600 rounded-full px-2 py-1">
              <span className="text-green-600 text-sm font-medium lowercase">{subscriptionDetails.status}</span>
            </div>
          </div>
          <div className="text-sm text-gray-600 space-y-4 font-mono" style={{ color: '#000000' }}>
            {getPlanTemplate().details.map((detail: PlanDetail, index: number) => (
              <div key={index} className="flex items-start">
                <span className="mr-2 text-black">✓</span>
                <span>{detail.description
                  .replace('{{includedCredits}}', subscriptionDetails?.included_credits.toString() || '0')
                  .replace('{{numMolecules}}', ((subscriptionDetails?.included_credits || 0) / 10).toString()) // Example calculation
                  .replace('{{overageCharge}}', subscriptionDetails?.overage_charge.toString() || '0.00')
                }</span>
              </div>
            ))}
          </div>
          <br />
            <div className="text-sm text-gray-600 mb-4 font-mono" style={{ color: '#808080' }}>
              <span>Your trial ends on {new Date(subscriptionDetails.current_period_end).toISOString().split('T')[0]}. After the trial ends, your plan will {isTrialingButNotRenewing ? 'not renew.' : 'continue with the selected subscription.'}</span>
            </div>
        </div>
      </div>
      <div className="absolute top-0 left-0 p-4 bg-white border-b border-gray-300 w-full flex justify-between items-center">
        <div className="text-gray-600 font-bold uppercase">Subscription/{walletAddress}</div>
      </div>
    </div>
  );
}
