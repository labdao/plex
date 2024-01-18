import { Loader2Icon } from "lucide-react";

import Logo from "../global/Logo";

interface PageLoaderProps {
  variant?: "default" | "logo";
}

export function PageLoader({ variant = "default" }: PageLoaderProps) {
  return (
    <div className="flex justify-center w-full p-16 text-primary animate-pulse">
      {variant === "logo" ? <Logo className="w-auto h-16 animate-spin" /> : <Loader2Icon className="animate-spin" size={48} absoluteStrokeWidth />}
    </div>
  );
}
