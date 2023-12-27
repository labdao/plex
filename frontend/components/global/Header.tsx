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

import { Button } from "../ui/button";
import Logo from "./Logo";
import { NavLink } from "./NavLink";

const navItems = [
  {
    title: "Tasks",
    href: "/tasks",
  },
  {
    title: "Experiments",
    href: "/experiments",
  },
  {
    title: "Data",
    href: "/data",
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
    <nav className="flex items-center justify-between p-4 border-b bg-background">
      <Link href="/" className="flex items-center gap-4 text-lg font-bold uppercase font-heading">
        <Logo className="w-auto h-8" /> Lab Exchange
      </Link>
      {authenticated && (
        <>
          <div className="flex gap-8 ml-16 mr-auto">
            {navItems.map((navItem, idx) => (
              <NavLink key={idx} href={navItem.href} className="font-mono font-bold uppercase">
                {navItem.title}
              </NavLink>
            ))}
          </div>
          <Button asChild className="mr-4">
            <Link href="/tasks/protein-design">Run Experiment</Link>
          </Button>
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
