import { Loader2Icon } from "lucide-react";

import { cn } from "@/lib/utils";

import Logo from "../global/Logo";

interface PageLoaderProps {
  variant?: "default" | "logo";
  className?: string;
}

export function PageLoader({ variant = "default", className }: PageLoaderProps) {
  return (
    <div className={cn("flex justify-center w-full p-16 text-primary animate-pulse", className)}>
      {variant === "logo" ? <Logo className="w-auto h-16 animate-spin" /> : <Loader2Icon className="animate-spin" size={48} absoluteStrokeWidth />}
    </div>
  );
}
