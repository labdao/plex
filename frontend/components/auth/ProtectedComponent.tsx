"use client";

import { usePrivy } from "@privy-io/react-auth";
import { LockIcon } from "lucide-react";
import { useRouter } from "next/navigation";
import { ReactNode } from "react";

import { Card, CardContent } from "@/components/ui/card";

import Logo from "../global/Logo";
import { PageLoader } from "../shared/PageLoader";
import PrivyLoginButton from "./PrivyLoginButton";

type ProtectedComponentProps = {
  children: ReactNode;
  method: "overlay" | "hide" | "redirect";
  message?: string;
};

const ProtectedComponent = ({ children, method = "overlay", message }: ProtectedComponentProps) => {
  const router = useRouter();
  const { ready, authenticated } = usePrivy();
  if (!authenticated && method === "redirect") {
    router.push("/login");
  }

  if (!ready) return <PageLoader variant="logo" />;

  if (authenticated) return children;

  if (!authenticated && method !== "redirect")
    return (
      <div className="relative">
        <div className="sticky z-40 -mt-2 border-t top-[-1px] inset-x-6 inset-y-12">
          <Card className="rounded-t-none">
            <CardContent className="flex items-center justify-between gap-4">
              <span className="font-mono font-bold tracking-wide uppercase">
                <LockIcon size={16} absoluteStrokeWidth className="inline-block mr-2" />
                {message}
              </span>
              <PrivyLoginButton>Log In</PrivyLoginButton>
            </CardContent>
          </Card>
        </div>
        {method === "overlay" && <div className="pointer-events-none select-none opacity-30">{children}</div>}
      </div>
    );
};

export default ProtectedComponent;
