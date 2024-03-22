"use client";

import { BookOpenIcon, FileJsonIcon, FileLineChart, GithubIcon, PanelRightCloseIcon, PanelRightOpenIcon } from "lucide-react";
import React, { useState } from "react";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { ToolDetail } from "@/lib/redux";
import { cn } from "@/lib/utils";

interface ModelInfoProps {
  tool: ToolDetail;
  selectComponent?: React.ReactNode;
  initialOpen?: boolean;
}

type OutputSummaryItem = {
  name: string;
  fileExtensions: string;
  fileNames: string;
  multiple: boolean;
};

export default function ModelInfo({ tool, initialOpen, selectComponent }: ModelInfoProps) {
  const [open, setOpen] = useState(initialOpen);

  const { description, github, paper, outputs } = tool.ToolJson;

  const renderDescriptionParagraphs = (text: string) => {
    const paragraphs = text.split("\n");
    const hasNumberedSteps = paragraphs.some((paragraph) => paragraph.match(/^\d+\. /));

    if (hasNumberedSteps) {
      const steps = paragraphs.filter((paragraph) => paragraph.match(/^\d+\. /));
      const nonStepParagraphs = paragraphs.filter((paragraph) => !paragraph.match(/^\d+\. /));

      return (
        <>
          {nonStepParagraphs.map((paragraph, index) => (
            <p key={index} className="mt-2">
              {paragraph}
            </p>
          ))}
          <ol className="mt-2 list-decimal list-inside">
            {steps.map((step, index) => (
              <li key={index} className="mt-2">
                {step.replace(/^\d+\. /, "")}
              </li>
            ))}
          </ol>
        </>
      );
    } else {
      return paragraphs.map((paragraph, index) => (
        <p key={index} className="mt-2">
          {paragraph}
        </p>
      ));
    }
  };

  let outputSummaryInfo = { items: [] as OutputSummaryItem[] };
  for (const key in outputs) {
    outputSummaryInfo.items.push({
      name: key.replaceAll("_", " "),
      fileExtensions: outputs?.[key]?.glob?.map((glob: string) => glob.split(".").pop())?.join(", "),
      fileNames: outputs?.[key]?.glob?.join(", "),
      multiple: outputs?.[key]?.type === "Array",
    });
  }

  return (
    <Card
      className={cn(
        "transition-all lg:rounded-r-none m-2 lg:mx-0 lg:my-2 lg:sticky top-0 grow-0 h-screen shrink-0 basis-24",
        open && "overflow-y-auto basis-1/3"
      )}
    >
      <div className="flex items-center p-4 pb-2">
        <Button size="sm" variant="ghost" className="-ml-4 text-base font-normal rounded-l-none font-heading" onClick={() => setOpen(!open)}>
          {open ? <PanelRightCloseIcon /> : <PanelRightOpenIcon />} Model
        </Button>
      </div>
      <div className={cn("transition-opacity opacity-0 min-w-[26vw]", open && "opacity-1")}>
        <CardContent className="pt-0">
          {selectComponent || (
            <div className="text-xl font-heading">
              {tool.ToolJson?.author || "unknown"}/{tool.ToolJson?.name}
            </div>
          )}
          <div>{renderDescriptionParagraphs(description)}</div>
          <div className="flex flex-wrap gap-2 mt-4">
            <Button asChild variant="outline" size="xs">
              <a href={process.env.NEXT_PUBLIC_DEMO_URL} target="_blank">
                <FileLineChart />
                Example Result
              </a>
            </Button>

            {github && (
              <Button asChild variant="outline" size="xs">
                <a href={github} target="_blank">
                  <GithubIcon /> GitHub
                </a>
              </Button>
            )}
            {paper && (
              <Button asChild variant="outline" size="xs">
                <a href={paper} target="_blank">
                  <BookOpenIcon /> Reference
                </a>
              </Button>
            )}
            {tool?.CID && (
              <Button asChild variant="outline" size="xs">
                <a href={`${process.env.NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT}${tool?.CID}/`} target="_blank">
                  <FileJsonIcon /> Manifest
                </a>
              </Button>
            )}
          </div>
        </CardContent>
        {outputs && (
          <CardContent className="border-t">
            <div className="mb-2 uppercase font-heading">Expected Output</div>
            <div className="space-y-2 lowercase">
              {(outputSummaryInfo?.items || []).map((item, index) => (
                <div key={index}>
                  <div className="text-sm">{item.multiple ? <div>{item.fileExtensions} files</div> : <div>{item.fileNames} file</div>}</div>
                  <div className="mr-3 text-xs text-muted-foreground">{item.name}</div>
                </div>
              ))}
            </div>
          </CardContent>
        )}
      </div>
    </Card>
  );
}
