"use client";

import { ChevronDownIcon } from "lucide-react";
import Link from "next/link";
import React from "react";

import { tasks } from "@/app/tasks/taskList";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";

export default function TasksMenu() {
  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger className="mr-2">
          <Button>
            Run Experiment <ChevronDownIcon size={18} className="ml-1" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" collisionPadding={10}>
          {tasks.map((task) => (
            <DropdownMenuItem asChild disabled={!task.available} key={task.slug}>
              <Link href={`/tasks/${task.slug}`} className="flex justify-between">
                {task.name}
                {!task.available && (
                  <Badge variant="secondary" className="ml-2">
                    Coming soon
                  </Badge>
                )}
              </Link>
            </DropdownMenuItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    </>
  );
}
