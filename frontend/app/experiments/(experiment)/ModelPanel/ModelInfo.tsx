"use client";

import { BookOpenIcon, FileJsonIcon, FileLineChart, GithubIcon } from "lucide-react";
import React from "react";

import { Button } from "@/components/ui/button";
import { ModelDetail } from "@/lib/redux";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";


interface ModelInfoProps {
  model: ModelDetail;
}

type OutputSummaryItem = {
  name: string;
  fileExtensions: string;
  filenames: string;
  multiple: boolean;
};

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

export default function ModelInfo({ model }: ModelInfoProps) {
  const { description, github, paper, outputs } = model.ModelJson;

  let outputSummaryInfo = { items: [] as OutputSummaryItem[] };
  for (const key in outputs) {
    outputSummaryInfo.items.push({
      name: key.replaceAll("_", " "),
      fileExtensions: outputs?.[key]?.glob?.map((glob: string) => glob.split(".").pop())?.join(", "),
      filenames: outputs?.[key]?.glob?.join(", "),
      multiple: outputs?.[key]?.type === "Array",
    });
  }

  return (
    <>
      <div className="text-left uppercase font-heading">About Model</div>
      <div>{renderDescriptionParagraphs(description)}</div>
      <div className="flex flex-wrap gap-2 mt-4">
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button asChild variant="outline" size="xs">
                <span>
                  <FileLineChart />
                  Example Result
                </span>
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <div className="text-xs">
                <p className="font-mono">Coming Soon!</p>
              </div>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
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
        {/* {model?.ID && (
          <Button asChild variant="outline" size="xs">
            <a href={`${model?.S3URI}/`} target="_blank">
              <FileJsonIcon /> Manifest
            </a>
          </Button>
        )} */}
      </div>
      {/*outputs && (
          <CardContent className="border-t">
            <div className="mb-2 uppercase font-heading">Expected Output</div>
            <div className="space-y-2 lowercase">
              {(outputSummaryInfo?.items || []).map((item, index) => (
                <div key={index}>
                  <div className="text-sm">{item.multiple ? <div>{item.fileExtensions} files</div> : <div>{item.filenames} file</div>}</div>
                  <div className="mr-3 text-xs text-muted-foreground">{item.name}</div>
                </div>
              ))}
            </div>
          </CardContent>
              )*/}
    </>
  );
}
