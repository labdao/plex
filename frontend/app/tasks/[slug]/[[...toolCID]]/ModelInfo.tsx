"use client";

import { BookOpenIcon, FileJsonIcon, FileLineChart, GithubIcon } from "lucide-react";
import React from "react";

import { Button } from "@/components/ui/button";
import { ToolDetail } from "@/lib/redux";

interface ModelInfoProps {
  tool: ToolDetail;
}

export default function ModelInfo({ tool }: ModelInfoProps) {
  const { description, github, paper, guide } = tool.ToolJson;

  const renderDescriptionParagraphs = (text: string) => {
    const paragraphs = text.split('\n');
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
          <ol className="list-decimal list-inside mt-2">
            {steps.map((step, index) => (
              <li key={index} className="mt-2 text-sm text-muted-foreground">
                {step.replace(/^\d+\. /, '')}
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

  return (
    <>
      {renderDescriptionParagraphs(description)}
      <div className="flex gap-2 mt-4">
        <Button asChild variant="outline" size="xs">
          <a href="/experiments/1" target="_blank">
            <FileLineChart />Example Result
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
    </>
  );
}