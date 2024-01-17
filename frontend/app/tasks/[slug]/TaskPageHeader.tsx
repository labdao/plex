"use client";

import { BookOpenIcon, GithubIcon, Loader2Icon } from "lucide-react";
import React from "react";

import { badgeVariants } from "@/components/ui/badge";
import { ToolDetail } from "@/lib/redux";

interface TaskPageHeaderProps {
  tool: ToolDetail;
  //@TODO will have a proper type once tasks are set up
  task: { name: string; slug: string; default_tool: { CID: string } };
  loading: boolean;
}

export default function TaskPageHeader({ tool, task, loading }: TaskPageHeaderProps) {
  const { author, name, description, github, paper } = tool.ToolJson;

  return (
    <div className="mb-6 border-b min-h-[11rem] border-b-border">
      <h1 className="mb-4 text-3xl font-heading">
        <span className="text-muted-foreground">
          {task.name}/{loading ? <Loader2Icon className="inline-block opacity-50 animate-spin" /> : <>{author || "unknown"}/</>}
        </span>
        {!loading && name}
      </h1>
      {!loading && (
        <>
          <p className="max-w-prose line-clamp-2 ">{description}</p>
          {(github || paper) && (
            <div className="flex gap-2 mt-4 ">
              {github && (
                <a href={github} target="_blank" className={badgeVariants({ variant: "references", size: "lg" })}>
                  github <GithubIcon className="ml-2" size={20} />
                </a>
              )}
              {paper && (
                <a href={paper} target="_blank" className={badgeVariants({ variant: "references", size: "lg" })}>
                  paper <BookOpenIcon className="ml-2" size={20} />
                </a>
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
}
