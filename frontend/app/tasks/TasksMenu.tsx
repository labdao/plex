"use client";

import { DropdownMenuContentProps } from "@radix-ui/react-dropdown-menu";
import Link from "next/link";
import React from "react";

import { tasks } from "@/app/tasks/taskList";
import { Badge } from "@/components/ui/badge";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";

interface TasksMenuProps {
  trigger: React.ReactNode;
  dropdownMenuContentProps?: DropdownMenuContentProps;
}

export default function TasksMenu({ dropdownMenuContentProps, trigger }: TasksMenuProps) {
  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>{trigger}</DropdownMenuTrigger>
        <DropdownMenuContent {...dropdownMenuContentProps} collisionPadding={10}>
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
