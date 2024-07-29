"use client";

import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";

export default function SubscribePage() {
  const [walletAddress, setWalletAddress] = useState<string | null>(null);

  // Mock function to simulate fetching wallet address, replace this with actual logic
  useEffect(() => {
    const fetchWalletAddress = async () => {
      // Replace with actual logic to fetch wallet address
      const address = "0x1234...abcd"; // Example wallet address
      setWalletAddress(address);
    };
    fetchWalletAddress();
  }, []);

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
      <div className="flex flex-col items-center justify-between w-[706px] h-[469px] p-4 md:p-6 lg:p-8 bg-white rounded-lg shadow-lg mx-auto my-6">
        <h1
          className="text-center"
          style={{ fontFamily: '"Space Grotesk"', fontSize: '48px', fontWeight: 500, lineHeight: '43.2px', color: '#4C4C4C' }}
        >
          Become a lab.bio Subscriber
        </h1>
        <ul className="space-y-2 w-full">
          <li className="flex items-center" style={{ color: '#4C4C4C', fontFamily: '"Space Grotesk"', fontSize: '20px', fontWeight: 500, lineHeight: '21.1px', letterSpacing: '0.1px' }}>
            <span className="mr-2 text-green-500">✓</span>
            <span>Access x# of computation credits (about x number per x number)</span>
          </li>
          <li className="flex items-center" style={{ color: '#4C4C4C', fontFamily: '"Space Grotesk"', fontSize: '20px', fontWeight: 500, lineHeight: '21.1px', letterSpacing: '0.1px' }}>
            <span className="mr-2 text-green-500">✓</span>
            <span>Additional charges information Additional charges information</span>
          </li>
          <li className="flex items-center" style={{ color: '#4C4C4C', fontFamily: '"Space Grotesk"', fontSize: '20px', fontWeight: 500, lineHeight: '21.1px', letterSpacing: '0.1px' }}>
            <span className="mr-2 text-green-500">✓</span>
            <span>Additional charges information Additional charges information</span>
          </li>
          <li className="flex items-center" style={{ color: '#4C4C4C', fontFamily: '"Space Grotesk"', fontSize: '20px', fontWeight: 500, lineHeight: '21.1px', letterSpacing: '0.1px' }}>
            <span className="mr-2 text-green-500">✓</span>
            <span>Cancel subscription any time</span>
          </li>
        </ul>
        <div className="text-center">
          <p
            className="mt-4"
            style={{ fontFamily: '"Space Grotesk"', fontSize: '29px', fontWeight: 500, lineHeight: '30.595px', letterSpacing: '0.145px', color: '#4C4C4C' }}
          >
            X$ per month
          </p>
          <button
            className="flex items-center justify-center mt-4 w-[489px] h-[42px] px-[26.853px] py-[20.139px] rounded-[53.705px] bg-gradient-to-r from-[rgba(107,219,173,0.5)] to-[rgba(7,164,198,0.5)] shadow-lg"
            style={{ boxShadow: '67.131px 83.914px 167.828px 33.566px rgba(0, 0, 0, 0.25)' }}
          >
            Start Now
          </button>
        </div>
      </div>
    </div>
  );
}
