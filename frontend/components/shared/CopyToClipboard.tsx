import { CopyIcon } from "lucide-react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";

interface CopyToClipboardProps {
  string: string;
  children?: React.ReactNode;
  className?: string;
}

export function CopyToClipboard({ string, children, className }: CopyToClipboardProps) {
  const [copied, setCopied] = useState(false);
  const copy = async () => {
    await navigator.clipboard.writeText(string);
    setCopied(true);
    setTimeout(() => {
      setCopied(false);
    }, 2000);
  };

  if (!string) return children;

  return (
    <div
      className={cn("relative cursor-pointer group", className)}
      onClick={() => {
        copy();
      }}
    >
      {children && <span className="group-hover:underline">{children}</span>}
      <TooltipProvider>
        <Tooltip open={copied}>
          <TooltipTrigger asChild>
            <Button size="icon" variant="ghost" className="inline-block p-1 text-muted-foreground">
              <CopyIcon />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="right">Copied!</TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </div>
  );
}
