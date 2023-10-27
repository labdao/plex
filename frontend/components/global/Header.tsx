"use client";

import { usePrivy } from "@privy-io/react-auth";
import { PrivyAuthContext } from "lib/PrivyContext";
import { selectIsLoggedIn, selectWalletAddress, setIsLoggedIn, setWalletAddress, useDispatch, useSelector } from "lib/redux";
import { DownloadIcon, User } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
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

import Logo from "./Logo";
import { NavLink } from "./NavLink";

const navItems = [
  {
    title: "Models",
    href: "/models",
  },
  {
    title: "Flows",
    href: "/flow/list",
  },
  {
    title: "Data",
    href: "/datafile/list",
  },
];

export default function Header() {
  const dispatch = useDispatch();
  const router = useRouter();
  const { ready, authenticated, user, exportWallet, logout } = usePrivy();
  const walletAddress = useSelector(selectWalletAddress);

  const hasEmbeddedWallet =
    ready && authenticated && !!user?.linkedAccounts.find((account: any) => account.type === "wallet" && account.walletClient === "privy");

  const handleExportWallet = async () => {
    if (hasEmbeddedWallet) {
      exportWallet();
    }
  };

  const handleLogout = async () => {
    logout();
    localStorage.removeItem("walletAddress");
    dispatch(setWalletAddress(""));
    dispatch(setIsLoggedIn(false));
    router.push("/login");
  };

  return (
    <nav className="flex items-center justify-between p-4 bg-background border-b">
      <Link href="/" className="flex items-center gap-4 font-bold uppercase text-lg">
        <Logo className="h-8 w-auto" /> Lab Exchange
      </Link>
      {authenticated && (
        <>
          <div className="flex gap-8 ml-16 mr-auto">
            {navItems.map((navItem, idx) => (
              <NavLink key={idx} href={navItem.href} className="font-medium uppercase">
                {navItem.title}
              </NavLink>
            ))}
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger>
              <User />
            </DropdownMenuTrigger>
            <DropdownMenuContent collisionPadding={10}>
              {user?.email?.address && (
                <>
                  <DropdownMenuLabel>{user?.email?.address}</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                </>
              )}

              {walletAddress && (
                <>
                  <DropdownMenuLabel>
                    Wallet: <em className="font-normal">{walletAddress}</em>
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
                </>
              )}
              <DropdownMenuItem onClick={handleLogout}>Log out</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </>
      )}
    </nav>
  );
}
