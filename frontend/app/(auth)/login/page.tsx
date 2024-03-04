"use client";

import { usePrivy } from "@privy-io/react-auth";
import { LockIcon } from "lucide-react";
import { useRouter } from "next/navigation";

import PrivyLoginButton from "@/components/auth/PrivyLoginButton";
import { PageLoader } from "@/components/shared/PageLoader";
import { Button } from "@/components/ui/button";

export default function LoginPage() {
  const { ready, authenticated } = usePrivy();
  const router = useRouter();

  if (authenticated) {
    router.push("/");
  }

  return ready && !authenticated ? (
    <div>
      <div className="p-16 text-center">
        <LockIcon size={48} className="mx-auto mb-4" />
        <h1 className="mb-4 font-mono text-lg font-bold tracking-wide uppercase">Log In to Your Lab.Bio Account</h1>
        <PrivyLoginButton>Log In</PrivyLoginButton>
      </div>
    </div>
  ) : (
    <PageLoader variant="logo" />
  );
}
