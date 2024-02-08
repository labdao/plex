"use client";

import { usePrivy } from "@privy-io/react-auth";
import { useDispatch } from "lib/redux";
import { Code2Icon, DownloadIcon, Loader2Icon, User } from "lucide-react";
import Link from "next/link";
import React from "react";

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

export default function UserMenu() {
  const { ready, authenticated, user, exportWallet, logout } = usePrivy();
  const walletAddress = user?.wallet?.address;

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

  if (!ready) return <Loader2Icon className="opacity-50 animate-spin" />;

  return (
    <>
      <PrivyLoginButton variant="ghost" />
      {authenticated && (
        <DropdownMenu>
          <DropdownMenuTrigger>
            <User />
          </DropdownMenuTrigger>
          <DropdownMenuContent collisionPadding={10} className="font-mono tracking-wider uppercase">
            {user?.email?.address && (
              <>
                <DropdownMenuLabel>{user?.email?.address}</DropdownMenuLabel>
                <DropdownMenuSeparator />
              </>
            )}

            {walletAddress && (
              <>
                <DropdownMenuLabel className="font-normal truncate w-72">
                  Wallet: <em>{walletAddress}</em>
                </DropdownMenuLabel>

                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger>
                      <DropdownMenuItem disabled={!hasEmbeddedWallet} onClick={handleExportWallet}>
                        <DownloadIcon size={20} className="mr-1" />
                        Export Wallet
                      </DropdownMenuItem>
                    </TooltipTrigger>
                    {!hasEmbeddedWallet && <TooltipContent>Export wallet only available for embedded wallets.</TooltipContent>}
                  </Tooltip>
                </TooltipProvider>

                <DropdownMenuSeparator />

                <Link href="/api">
                  <DropdownMenuItem>
                    <span className="flex items-center">
                      <Code2Icon size={20} className="mr-1" />
                      API Keys
                    </span>
                  </DropdownMenuItem>
                </Link>

                <DropdownMenuSeparator />
              </>
            )}
            <DropdownMenuItem onClick={handleLogout}>Log out</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )}
    </>
  );
}
