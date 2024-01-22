"use client";

import { usePrivy } from "@privy-io/react-auth";
import { DownloadIcon, Loader2Icon, LogOutIcon, User, UserCircleIcon } from "lucide-react";
import React from "react";

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
import { NavButton } from "./NavItem";

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
                  <DropdownMenuLabel className="truncate w-72">
                    Wallet: <em className="font-mono font-normal">{walletAddress}</em>
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
            </DropdownMenuContent>
          </DropdownMenu>
          <NavButton icon={<LogOutIcon />} title="Log Out" onClick={handleLogout} />
        </>
      )}
    </>
  );
}
