"use client";

import { usePrivy } from "@privy-io/react-auth";
import AddDataFileForm from "app/data/AddDataFileForm";
import TasksMenu from "app/tasks/TasksMenu";
import dayjs from "dayjs";
import isBetween from "dayjs/plugin/isBetween";
import isToday from "dayjs/plugin/isToday";
import { BoxIcon, FlaskRoundIcon, FolderIcon, GithubIcon, UploadIcon } from "lucide-react";
import { SproutIcon } from "lucide-react";
import Link from "next/link";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { ExperimentStatus } from "@/app/experiments/ExperimentStatus";
import { NavButton } from "@/components/global/NavItem";
import { ScrollArea } from "@/components/ui/scroll-area";
import { AppDispatch, Flow, flowListThunk, selectFlowList, selectFlowListLoading } from "@/lib/redux";

import Logo from "./Logo";
import { NavLink } from "./NavItem";
import UserMenu from "./UserMenu";

dayjs.extend(isToday);
dayjs.extend(isBetween);

const NavContent = (props: React.PropsWithChildren) => <div className="flex flex-col p-1 border-b border-border/50" {...props} />;

export default function Nav() {
  const { user } = usePrivy();
  const dispatch = useDispatch<AppDispatch>();
  const flows = useSelector(selectFlowList);
  const loading = useSelector(selectFlowListLoading);
  const walletAddress = user?.wallet?.address;

  useEffect(() => {
    if (walletAddress) {
      dispatch(flowListThunk(walletAddress));
    }
  }, [dispatch, walletAddress]);

  const categorizeFlows = (flows: Flow[]) => {
    const today = dayjs();
    const categories = {
      today: [] as Flow[],
      last7Days: [] as Flow[],
      last30Days: [] as Flow[],
      older: [] as Flow[],
    };

    flows.forEach((flow: Flow) => {
      const start = dayjs(flow.StartTime);
      if (start.isToday()) {
        categories.today.push(flow);
      } else if (start.isBetween(today.subtract(7, "day"), today)) {
        categories.last7Days.push(flow);
      } else if (start.isBetween(today.subtract(30, "day"), today)) {
        categories.last30Days.push(flow);
      } else {
        categories.older.push(flow);
      }
    });

    return categories;
  }

  const categorizedFlows = categorizeFlows(flows);

  return (
    <nav className="sticky top-0 z-50 flex flex-col justify-between w-48 h-screen border-r shadow-lg border-border/50 shrink-0 bg-background">
        <Link
          href="/"
          className="flex items-center gap-2 p-2 text-lg font-bold uppercase border-b border-border/50 h-14 font-heading whitespace-nowrap"
        >
          <Logo className="w-auto h-6 text-primary" />
          Lab.Bio
        </Link>
        <div className="p-3 font-mono text-xs font-bold text-muted-foreground opacity-70 sticky top-0 bg-background">
          My Experiments
        </div>
        <ScrollArea className="flex-grow">
          <div className="flex flex-col overflow-auto">
            <NavContent>
              {Object.keys(categorizedFlows).map((category) => {
                const categoryKey = category as keyof typeof categorizedFlows;
                const flowsInCategory = categorizedFlows[categoryKey];
                if (flowsInCategory.length > 0) {
                  return (
                    <div key={category}>
                      <div className="p-2 mt-4 text-xs font-bold text-muted-foreground opacity-70">
                        {category === 'today' && 'Today'}
                        {category === 'last7Days' && 'Previous 7 Days'}
                        {category === 'last30Days' && 'Previous 30 Days'}
                        {category === 'older' && 'Older'}
                      </div>
                      {flowsInCategory.map((flow) => (
                        <NavLink key={flow.ID} href={`/experiments/${flow.ID}`} title={flow.Name} />
                      ))}
                    </div>
                  );
                }
                return null;
              })}
            </NavContent>
          </div>
        </ScrollArea>
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
