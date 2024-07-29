"use client";

import { usePrivy } from "@privy-io/react-auth";
import React from "react";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import { toast } from "sonner";
import { getAccessToken } from "@privy-io/react-auth";
import { AlertDialog, AlertDialogContent, AlertDialogOverlay, AlertDialogHeader, AlertDialogFooter, AlertDialogCancel } from "@/components/ui/alert-dialog";
import { useState } from "react";
import backendUrl from "lib/backendUrl";

export default function ManageSubscription() {
  const { user } = usePrivy();
  const walletAddress = user?.wallet?.address;
  const [isDialogOpen, setIsDialogOpen] = useState(false);

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
    } catch (error) {
    console.error("Error cancelling subscription", error);
    toast.error("Failed to cancel subscription");
    }
  };

  
  if (!walletAddress) {
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
                <span style={{ color: '#000000' }}>x credits / x credits</span>
              </div>
              <div className="flex justify-between">
                <span style={{ color: '#808080' }}>Current cycle overage charges</span>
                <span style={{ color: '#000000' }}>x credits / $xxxx</span>
              </div>
            </div>
          </div>
          <div className="w-96 p-6 bg-white rounded-lg shadow-lg font-heading">
            <h3 className="text-2xl font-bold text-black mb-4">Billing & Payment</h3>
            <div className="text-sm text-gray-600 space-y-4">
              <div className="flex justify-between">
                <span style={{ color: '#808080' }}>Cost</span>
                <span style={{ color: '#000000' }}>x$/mo</span>
              </div>
              <div className="flex justify-between">
                <span style={{ color: '#808080' }}>Billing period</span>
                <span style={{ color: '#000000' }}>Monthly</span>
              </div>
              <div className="flex justify-between">
                <span style={{ color: '#808080' }}>Renewal date</span>
                <span style={{ color: '#000000' }}>August 1</span>
              </div>
            </div>
            <div className="flex justify-between mt-6">
              <button className="px-4 py-2 border rounded-md" style={{ borderColor: '#6BDBAD', color: '#6BDBAD' }}>Edit Billing</button>
              <button
                className="px-4 py-2 border rounded-md"
                style={{ borderColor: '#000000', color: '#808080' }}
                onClick={cancelSubscription}
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
              <span className="text-green-600 text-sm font-medium lowercase">active</span>
            </div>
          </div>
          <div className="text-sm text-gray-600 space-y-4 font-heading" style={{ color: '#000000' }}>
            <div className="flex items-start">
              <span className="mr-2 text-black">✓</span>
              <span>Access x# of computation credits (about x number per x number)</span>
            </div>
            <div className="flex items-start">
              <span className="mr-2 text-black">✓</span>
              <span>Additional charges information Additional charges information</span>
            </div>
            <div className="flex items-start">
              <span className="mr-2 text-black">✓</span>
              <span>Additional charges information Additional charges information</span>
            </div>
            <div className="flex items-start">
              <span className="mr-2 text-black">✓</span>
              <span>Cancel subscription any time</span>
            </div>
          </div>
        </div>
      </div>
      <div className="absolute top-0 left-0 p-4 bg-white border-b border-gray-300 w-full flex justify-between items-center">
        <div className="text-gray-600 font-bold uppercase">Subscription/{walletAddress}</div>
      </div>
    <AlertDialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
    <AlertDialogOverlay />
    <AlertDialogContent>
    <AlertDialogHeader>
        <h2>Are you sure?</h2>
    </AlertDialogHeader>
    <AlertDialog>
        Do you really want to cancel your subscription? This action cannot be undone.
    </AlertDialog>
    <AlertDialogFooter>
        <button onClick={() => setIsDialogOpen(false)} className="px-4 py-2 border rounded-md" style={{ borderColor: '#000000', color: '#808080' }}>
        No, Keep Subscription
        </button>
        <button onClick={cancelSubscription} className="px-4 py-2 border rounded-md" style={{ borderColor: '#6BDBAD', color: '#6BDBAD' }}>
        Yes, Cancel Subscription
        </button>
    </AlertDialogFooter>
    </AlertDialogContent>
    </AlertDialog>
</div>
  );
}
