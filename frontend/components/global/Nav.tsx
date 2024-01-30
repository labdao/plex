"use client";

import AddDataFileForm from "app/data/AddDataFileForm";
import TasksMenu from "app/tasks/TasksMenu";
import { BoxIcon, FlaskRoundIcon, FolderIcon, GithubIcon, UploadIcon } from "lucide-react";
import { SproutIcon } from "lucide-react";
import Link from "next/link";
import React from "react";

import { NavButton } from "@/components/global/NavItem";
import { ScrollArea } from "@/components/ui/scroll-area";

import Logo from "./Logo";
import { NavLink } from "./NavItem";
import UserMenu from "./UserMenu";

const NavContent = (props: React.PropsWithChildren) => <div className="flex flex-col p-1 border-b border-border/50" {...props} />;

export default function Nav() {
  return (
    <nav className="sticky top-0 z-50 flex flex-col justify-between w-48 h-screen border-r shadow-lg border-border/50 shrink-0 bg-background">
      <ScrollArea>
        <Link
          href="/"
          className="flex items-center gap-2 p-2 text-lg font-bold uppercase border-b border-border/50 h-14 font-heading whitespace-nowrap"
        >
          <Logo className="w-auto h-6 text-primary" />
          Lab Exchange
        </Link>
        <NavContent>
          <NavLink href="/tasks" icon={<BoxIcon />} title="Tasks" />
        </NavContent>
        <NavContent>
          <div className="p-2 font-mono text-xs font-bold text-muted-foreground opacity-70">Personal</div>
          <NavLink href="/experiments" icon={<FlaskRoundIcon />} title="My Experiments" />
          <NavLink href="/data" icon={<FolderIcon />}>
            <>
              My Files&nbsp;<span className="opacity-70">(beta)</span>
            </>
          </NavLink>
          <TasksMenu
            trigger={<NavButton icon={<SproutIcon />} title="Run Experiment" hasDropdown />}
            dropdownMenuContentProps={{ side: "right", align: "start" }}
          />
          <AddDataFileForm trigger={<NavButton icon={<UploadIcon />} title="Upload Files" />} />
        </NavContent>
      </ScrollArea>
      <div>
        <NavContent>
          <NavLink href="https://github.com/labdao" target="_blank" icon={<GithubIcon />} title="GitHub" />
        </NavContent>
        <NavContent>
          <UserMenu />
        </NavContent>
      </div>
    </nav>
  );
}
