import { Share2 } from 'lucide-react';
import React, { useState } from 'react';

import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

const ExperimentShare = ({ experimentID }: { experimentID: string }) => {
  const [copied, setCopied] = useState(false);
  const currentPageLink = `${process.env.NEXT_PUBLIC_FRONTEND_URL}/experiments/${experimentID}`;

  const copyLinkToClipboard = async () => {
    await navigator.clipboard.writeText(currentPageLink);
    setCopied(true);
    setTimeout(() => {
      setCopied(false);
    }, 2000);
  };

  return (
    <TooltipProvider>
      <Tooltip open={copied}>
        <TooltipTrigger asChild>
          <Button variant="outline" className="text-sm" onClick={copyLinkToClipboard}>
            <Share2 className="w-4 h-4 mr-2" /> Share
          </Button>
        </TooltipTrigger>
        <TooltipContent side="right">Copied!</TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default ExperimentShare;