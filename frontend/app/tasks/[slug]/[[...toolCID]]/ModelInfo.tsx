"use client";

import { BookOpenIcon, FileJsonIcon, GithubIcon } from "lucide-react";
import React from "react";

import { Button } from "@/components/ui/button";
import { ToolDetail } from "@/lib/redux";

interface ModelInfoProps {
  tool: ToolDetail;
}

export default function ModelInfo({ tool }: ModelInfoProps) {
  const { description, github, paper } = tool.ToolJson;
  return (
    <>
      <p className="mt-4">{description}</p>
      <div className="flex gap-2 mt-4 ">
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
