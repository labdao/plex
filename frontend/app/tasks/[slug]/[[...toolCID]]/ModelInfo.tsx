"use client";

import { BookOpenIcon, FileJsonIcon, FileLineChart, GithubIcon } from "lucide-react";
import React from "react";

import { Button } from "@/components/ui/button";
import { CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { ToolDetail } from "@/lib/redux";

interface ModelInfoProps {
  tool: ToolDetail;
}

type OutputSummaryItem = {
  name: string;
  fileExtensions: string;
  fileNames: string;
  multiple: boolean;
};

export default function ModelInfo({ tool }: ModelInfoProps) {
  const { description, github, paper, guide, outputs } = tool.ToolJson;

  const renderDescriptionParagraphs = (text: string) => {
    const paragraphs = text.split("\n");
    const hasNumberedSteps = paragraphs.some((paragraph) => paragraph.match(/^\d+\. /));

    if (hasNumberedSteps) {
      const steps = paragraphs.filter((paragraph) => paragraph.match(/^\d+\. /));
      const nonStepParagraphs = paragraphs.filter((paragraph) => !paragraph.match(/^\d+\. /));

      return (
        <>
          {nonStepParagraphs.map((paragraph, index) => (
            <p key={index} className="mt-2 text-sm text-muted-foreground">
              {paragraph}
            </p>
          ))}
          <ol className="mt-2 list-decimal list-inside">
            {steps.map((step, index) => (
              <li key={index} className="mt-2 text-sm text-muted-foreground">
                {step.replace(/^\d+\. /, "")}
              </li>
            ))}
          </ol>
        </>
      );
    } else {
      return paragraphs.map((paragraph, index) => (
        <p key={index} className="mt-2 text-sm text-muted-foreground">
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

  console.log(outputSummaryInfo);

  return (
    <>
      {renderDescriptionParagraphs(description)}
      <div className="flex gap-2 mt-4">
        <Button asChild variant="outline" size="xs">
          <a href={process.env.NEXT_PUBLIC_DEMO_URL} target="_blank">
            <FileLineChart />
            Example Result
          </a>
        </Button>
      </div>
      <div className="flex gap-2 mt-4">
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
      {outputs && (
        <>
          <Separator className="my-2" />
          <CardContent>
            <div className="mb-4 font-mono text-sm font-bold uppercase">Expected Output</div>
            <div className="space-y-2 lowercase">
              {(outputSummaryInfo?.items || []).map((item, index) => (
                <div key={index}>
                  {item.multiple ? <div>{item.fileExtensions} files</div> : <div>{item.fileNames} file</div>}
                  <div className="mr-3 text-xs text-muted-foreground">{item.name}</div>
                </div>
              ))}
            </div>
          </CardContent>
        </>
      )}
    </>
  );
}
