"use client";

import { getAccessToken, usePrivy } from "@privy-io/react-auth";
import { Code2Icon, CreditCardIcon, DownloadIcon, FlaskRoundIcon, FolderIcon, Loader2Icon, LogOutIcon, User, UserCircleIcon } from "lucide-react";
import Link from "next/link";
import React, { useEffect, useState } from "react";

import { CopyToClipboard } from "@/components/shared/CopyToClipboard";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

import PrivyLoginButton from "../auth/PrivyLoginButton";
import StripeCheckoutButton from "../payment/StripeCheckoutButton";
import TransactionSummaryInfo from "../payment/TransactionSummaryInfo";
import { NavButton } from "./NavItem";
import backendUrl from "@/lib/backendUrl";

export default function UserMenu() {
  const { ready, authenticated, user, exportWallet, logout } = usePrivy();
  const walletAddress = user?.wallet?.address;
  const [isSubscribed, setIsSubscribed] = useState(false);

  const hasEmbeddedWallet =
    ready && authenticated && !!user?.linkedAccounts.find((account: any) => account.type === "wallet" && account.walletClient === "privy");

  const handleExportWallet = async () => {
    if (hasEmbeddedWallet) {
      exportWallet();
    }
  };

  const handleLogout = async () => {
    await logout();
  };

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
        setIsSubscribed(data.isSubscribed);
      }
    };

    if (authenticated) {
      checkSubscriptionStatus();
    }
  }, [authenticated]);

  if (!ready)
    return (
      <div className="flex items-center gap-2 px-3 py-2 text-sm text-muted-foreground">
        <Loader2Icon className="w-4 opacity-50 animate-spin" /> Authenticating...
      </div>
    );

  return (
    <>
      <PrivyLoginButton asChild>
        <span>
          <NavButton icon={<UserCircleIcon />} title="Log In" />
        </span>
      </PrivyLoginButton>
      {authenticated && (
        <>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <NavButton icon={<UserCircleIcon />} hasDropdown title={user?.email?.address || walletAddress} />
            </DropdownMenuTrigger>
            <DropdownMenuContent collisionPadding={10} side="right" align="start">
              {user?.email?.address && (
                <>
                  <DropdownMenuLabel>{user?.email?.address}</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                </>
              )}

              {walletAddress && (
                <>
                  <DropdownMenuLabel className="flex items-center truncate w-72">
                    Wallet: <em className="flex-grow ml-1 font-mono font-normal truncate">{walletAddress}</em>
                    <CopyToClipboard string={walletAddress} />
                  </DropdownMenuLabel>
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger className="w-full">
                        <DropdownMenuItem disabled={!hasEmbeddedWallet} onClick={handleExportWallet}>
                          <DownloadIcon size={20} className="mr-1" />
                          Export Wallet
                        </DropdownMenuItem>
                      </TooltipTrigger>
                      {!hasEmbeddedWallet && <TooltipContent>Export wallet only available for embedded wallets.</TooltipContent>}
                    </Tooltip>
                  </TooltipProvider>
                  <DropdownMenuSeparator />
                </>
              )}

              <DropdownMenuItem asChild>
                <Link href={isSubscribed ? "/subscription/manage" : "/subscribe"}>
                  <CreditCardIcon size={20} className="mr-1" />
                  {isSubscribed ? "Manage Subscription" : "Subscribe"}
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link href="/api">
                  <Code2Icon size={20} className="mr-1" />
                  API Keys
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link href="/experiments">
                  <FlaskRoundIcon size={20} className="mr-1" />
                  Experiments
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link href="/data">
                  <FolderIcon size={20} className="mr-1" />
                  Files
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem onClick={handleLogout}>
                <LogOutIcon size={20} className="mr-1" />
                Log Out
              </DropdownMenuItem>

              {/* <TransactionSummaryInfo className="mt-2" /> */}
            </DropdownMenuContent>
          </DropdownMenu>
        </>
      )}
    </>
  );
}
