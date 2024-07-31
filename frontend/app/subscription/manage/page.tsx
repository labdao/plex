"use client";

import { usePrivy } from "@privy-io/react-auth";
import React, { useEffect, useState } from "react";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import { toast } from "sonner";
import { getAccessToken } from "@privy-io/react-auth";
import { AlertDialog, AlertDialogContent, AlertDialogOverlay, AlertDialogHeader, AlertDialogFooter, AlertDialogCancel } from "@/components/ui/alert-dialog";
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
  }

  export default function ManageSubscription() {
    const { user } = usePrivy();
    const walletAddress = user?.wallet?.address;
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [loading, setLoading] = useState(true);
    const [subscriptionDetails, setSubscriptionDetails] = useState<SubscriptionDetails | null>(null);
    const router = useRouter();
  
    const cancelSubscription = async () => {
      try {
        let authToken;
        try {
          authToken = await getAccessToken();
        } catch (error) {
          console.log("Failed to get access token: ", error);
          throw new Error("Authentication failed");
        }
        const response = await fetch(`${backendUrl()}/stripe/subscription/cancel`, {
          method: "POST",
          headers: {
            Authorization: `Bearer ${authToken}`, // Ensure you have access token
            "Content-Type": "application/json",
          },
        });
  
        if (!response.ok) {
          throw new Error("Failed to cancel subscription");
        }
  
        const data = await response.json();
        toast.success(data.message);
        router.replace("/subscribe"); // Redirect to the subscribe page after canceling
      } catch (error) {/stripe/subscription/check
        console.error("Error cancelling subscription", error);
        toast.error("Failed to cancel subscription");
      }
    };
  
    useEffect(() => {
      const fetchSubscriptionDetails = async () => {
        let authToken;
        try {
          authToken = await getAccessToken();
        } catch (error) {
          console.log("Failed to get access token: ", error);
          return;
        }
  
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
      };
  
      const checkSubscriptionStatus = async () => {
        let authToken;
        try {
          authToken = await getAccessToken();
        } catch (error) {
          console.log("Failed to get access token: ", error);
          return;
        }
  
        const response = await fetch(`${backendUrl()}/stripe/subscription/check-subscription`, {
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
      };
  
      checkSubscriptionStatus();
    }, [router]);
  
    if (loading) {
      return <div>Loading...</div>;
    }
  
    if (!walletAddress || !subscriptionDetails) {
      return <div>Loading...</div>;
    }

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
          <div className="w-96 p-6 bg-white rounded-lg shadow-lg">
            <h3 className="text-2xl font-bold font-heading text-black mb-4">Usage Details</h3>
            <div className="text-sm text-gray-600 space-y-4 font-heading">
            <div className="flex justify-between">
                   <span style={{ color: '#808080' }}>Included</span>
                   <span style={{ color: '#000000' }}>
                     {subscriptionDetails.used_credits} / {subscriptionDetails.included_credits} credits
                   </span>
                 </div>
                 <div className="flex justify-between">
                   <span style={{ color: '#808080' }}>Current cycle overage charges</span>
                   <span style={{ color: '#000000' }}>
                     {/* {subscriptionDetails.used_credits - subscriptionDetails.included_credits} credits /  */}
                     ${subscriptionDetails.overage_charge}
                   </span>
                 </div>
               </div>
             </div>
             <div className="w-96 p-6 bg-white rounded-lg shadow-lg font-heading">
               <h3 className="text-2xl font-bold text-black mb-4">Billing & Payment</h3>
               <div className="text-sm text-gray-600 space-y-4">
                 <div className="flex justify-between">
                   <span style={{ color: '#808080' }}>Cost until {subscriptionDetails.next_due}</span>
                   <span style={{ color: '#000000' }}>
                     ${subscriptionDetails.plan_amount}
                   </span>
                 </div>
                 <div className="flex justify-between">
                   <span style={{ color: '#808080' }}>Billing period</span>
                   <span style={{ color: '#000000' }}>{subscriptionDetails.plan_interval}</span>
                 </div>
                 <div className="flex justify-between">
                   <span style={{ color: '#808080' }}>Renewal date</span>
                   <span style={{ color: '#000000' }}>{subscriptionDetails.next_due}</span>
                 </div>
                 <div className="flex justify-between">
                     <span style={{ color: '#808080' }}>Renewal amount</span>
                     <span style={{ color: '#000000' }}>
                          $5 + overage charges{/* TODO: should be changed to the dynamic value */}
                     </span>
                 </div>
               </div>
               <div className="flex justify-between mt-6">
                 <button className="px-4 py-2 border rounded-md" style={{ borderColor: '#6BDBAD', color: '#6BDBAD' }}>Edit Billing</button>
                 <button
                   className="px-4 py-2 border rounded-md"
                   style={{ borderColor: '#000000', color: '#808080' }}
                   onClick={() => setIsDialogOpen(true)}
                 >
                   Cancel Subscription
                 </button>
               </div>
             </div>
           </div>
           <div className="w-96 p-6 bg-white rounded-lg shadow-lg">
           <div className="flex justify-between items-center mb-4">
                <h3 className="text-2xl font-bold font-heading text-black">Your Plan</h3>
                <div className="flex items-center border border-green-600 rounded-full px-2 py-1">
                    <span className="text-green-600 text-sm font-medium lowercase">{subscriptionDetails.status}</span>
                </div>
            </div>
             <div className="text-sm text-gray-600 space-y-4 font-heading" style={{ color: '#000000' }}>
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
            <br/>
             {subscriptionDetails.status === 'trialing' && (
            <div className="text-sm text-gray-600 mb-4 font-heading" style={{ color: '#808080' }}>
                <span>Your trial ends on {new Date(subscriptionDetails.current_period_end).toLocaleDateString()}. After the trial ends, your plan will continue with the selected subscription.</span>
            </div>
            )}
           </div>
         </div>
         <div className="absolute top-0 left-0 p-4 bg-white border-b border-gray-300 w-full flex justify-between items-center">
           <div className="text-gray-600 font-bold uppercase">Subscription/{walletAddress}</div>
         </div>
         <AlertDialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
           <AlertDialogOverlay />
           <AlertDialogContent>
             <AlertDialog>
               Are you sure you want to cancel your subscription? You can resubscribe at any time to continue enjoying our services.
             </AlertDialog>
             <AlertDialogFooter>
               <button onClick={() => setIsDialogOpen(false)} className="px-4 py-2 border rounded-md" style={{ borderColor: '#6BDBAD', color: '#6BDBAD' }}>
                 No, Keep Subscription
               </button>
               <button onClick={cancelSubscription} className="px-4 py-2 border rounded-md" style={{ borderColor: '#000000', color: '#808080' }}>
                 Yes, Cancel Subscription
               </button>
             </AlertDialogFooter>
           </AlertDialogContent>
         </AlertDialog>
       </div>
     );
   }