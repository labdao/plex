"use client";

import { ExternalLinkIcon, PanelRightCloseIcon } from "lucide-react";
import Link from "next/link";

import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";

import ExperimentDetail from "../../ExperimentDetail";

export default function ExperimentDetailPage({ params }: { params: { flowID: string } }) {
  return (
    <ScrollArea className="relative z-50 flex flex-col h-screen min-w-0 shadow-xl bg-gray-50">
      <div className="flex items-center justify-between px-4 py-2 -mb-6 h-14 bg-background">
        <Button size="icon" variant="ghost" asChild>
          <a href={`/experiments/${params.flowID}`} target="_blank">
            <ExternalLinkIcon />
          </a>
        </Button>
        <Button size="icon" variant="ghost" asChild>
          <Link href="/experiments">
            <PanelRightCloseIcon />
          </Link>
        </Button>
      </div>
      <ExperimentDetail experimentID={params.flowID} />
    </ScrollArea>
  );
}
