"use client";

import { BookOpenIcon, GithubIcon } from "lucide-react";
import React from "react";

import { badgeVariants } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { ToolDetail } from "@/lib/redux";

interface TaskPageHeaderProps {
  tool: ToolDetail;
  //@TODO will have a proper type once tasks are set up
  task: { name: string; slug: string; default_tool: { CID: string } };
}

export default function TaskPageHeader({ tool, task }: TaskPageHeaderProps) {
  const { author, name, description, github, paper } = tool.ToolJson;

  return (
    <>
      <h1 className="mb-4 text-3xl font-heading">
        <span className="text-muted-foreground">
          <span className="lowercase">{task.name}</span>/{author || "unknown"}/
        </span>
        {name}
      </h1>
      <p className="max-w-prose line-clamp-2 min-h-[3em]">{description}</p>
      {(github || paper) && (
        <div className="flex gap-2 mt-4">
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
      <Separator className="my-6" />
    </>
  );
}
