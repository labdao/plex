import { CopyIcon } from "lucide-react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

interface CopyToClipboardProps {
  string: string;
  children: React.ReactNode;
}

export function CopyToClipboard({ string, children }: CopyToClipboardProps) {
  const [copied, setCopied] = useState(false);
  const copy = async () => {
    await navigator.clipboard.writeText(string);
    setCopied(true);
    setTimeout(() => {
      setCopied(false);
    }, 2000);
  };
  return (
    <div
      className="relative flex items-center gap-1 cursor-pointer group"
      onClick={() => {
        copy();
      }}
    >
      <span className="group-hover:underline">{children}</span>
      <TooltipProvider>
        <Tooltip open={copied}>
          <TooltipTrigger asChild>
            <Button size="icon" variant="ghost" className="p-1">
              <CopyIcon />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="right">Copied!</TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </div>
  );
}
