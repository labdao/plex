"use client";

import { ChevronRightIcon, SproutIcon } from "lucide-react";
import Link from "next/link";
import React from "react";

import { tasks } from "@/app/tasks/taskList";
import { Badge } from "@/components/ui/badge";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";

import { NavButton } from "./NavItem";

export default function TasksMenu() {
  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <NavButton icon={<SproutIcon />} title="Run Experiment" hasDropdown />
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start" side="right" collisionPadding={10}>
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
