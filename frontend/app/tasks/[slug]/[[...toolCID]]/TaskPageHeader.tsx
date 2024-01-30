"use client";

import { FileIcon, GithubIcon, Loader2Icon } from "lucide-react";
import React from "react";

import { Button } from "@/components/ui/button";
import { ToolDetail } from "@/lib/redux";

interface TaskPageHeaderProps {
  tool: ToolDetail;
  loading: boolean;
}

export default function TaskPageHeader({ tool, loading }: TaskPageHeaderProps) {
  const { author, name, description, github, paper } = tool.ToolJson;

  return (
    <div className="relative z-30 p-8 border-b rounded-b-lg shadow-md min-h-[14rem] border-b-border bg-background">
      <h1 className="mb-4 text-3xl font-heading">
        <span className="text-muted-foreground">
          {loading ? <Loader2Icon className="inline-block opacity-50 animate-spin" /> : <>{author || "unknown"}/</>}
        </span>
        {!loading && name}
      </h1>
      {!loading && (
        <>
          <p className="max-w-prose line-clamp-2 ">{description}</p>
          {(github || paper) && (
            <div className="flex gap-2 mt-4 ">
              {github && (
                <Button asChild variant="outline" size="sm">
                  <a href={github} target="_blank">
                    <GithubIcon /> GitHub
                  </a>
                </Button>
              )}
              {paper && (
                <Button asChild variant="outline" size="sm">
                  <a href={paper} target="_blank">
                    <FileIcon /> Reference
                  </a>
                </Button>
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
}
