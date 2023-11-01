"use client";

import { useRouter } from "next/navigation";
import React from "react";

export function useAuth() {
  const user = localStorage.getItem("username");
  const walletAddress = localStorage.getItem("walletAddress");
  const router = useRouter();

  React.useEffect(() => {
    if (!user || !walletAddress) {
      router.push("/login");
    }
  }, [user, walletAddress, router]);
}
