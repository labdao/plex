"use client";

import { usePrivy } from "@privy-io/react-auth";
import AddDataFileForm from "app/data/AddDataFileForm";
import TasksMenu from "app/tasks/TasksMenu";
import { BoxIcon, FlaskRoundIcon, FolderIcon, GithubIcon, UploadIcon } from "lucide-react";
import { SproutIcon } from "lucide-react";
import Link from "next/link";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { ExperimentStatus } from "@/app/experiments/ExperimentStatus";
import { NavButton } from "@/components/global/NavItem";
import { ScrollArea } from "@/components/ui/scroll-area";
import { AppDispatch, Flow, flowListThunk, selectCategorizedFlows, selectFlowList, selectFlowListLoading } from "@/lib/redux";

import Logo from "./Logo";
import { NavLink } from "./NavItem";
import UserMenu from "./UserMenu";

const NavContent = (props: React.PropsWithChildren) => <div className="flex flex-col p-1 border-b border-border/50" {...props} />;

export default function Nav() {
  const { user } = usePrivy();
  const dispatch = useDispatch<AppDispatch>();
  const categorizedFlows = useSelector(selectCategorizedFlows);
  const flows = useSelector(selectFlowList);
  // const loading = useSelector(selectFlowListLoading);
  const walletAddress = user?.wallet?.address;

  useEffect(() => {
    console.log('walletAddress', walletAddress);
    if (walletAddress) {
      console.log('dispatching flowListThunk')
      dispatch(flowListThunk(walletAddress));
    }
  }, [dispatch, walletAddress]);

  return (
    <nav className="sticky top-0 z-50 flex flex-col justify-between w-48 h-screen border-r shadow-lg border-border/50 shrink-0 bg-background">
      <Link
        href="/"
        className="flex items-center gap-2 p-2 text-lg font-bold uppercase border-b border-border/50 h-20 font-heading whitespace-nowrap"
      >
        <Logo className="w-auto h-6 text-primary" />
        Lab.Bio
      </Link>
      <div className="sticky top-14 bg-background border-b border-border/50 z-10">
        <div className="p-2 font-mono text-xs font-bold text-muted-foreground opacity-70">
          Experiments
        </div>
      </div>
      <ScrollArea className="flex-grow">
        <div className="flex flex-col overflow-auto">
          <NavContent>
            {Object.keys(categorizedFlows).map((category) => {
              const flowsInCategory = categorizedFlows[category as keyof typeof categorizedFlows];
              if (flowsInCategory.length > 0) {
                return (
                  <div key={category}>
                    <div className="p-2 mt-4 text-xs font-bold text-muted-foreground opacity-70">
                      {category === 'today' && 'Today'}
                      {category === 'last7Days' && 'Previous 7 Days'}
                      {category === 'last30Days' && 'Previous 30 Days'}
                      {category === 'older' && 'Older'}
                    </div>
                    {flowsInCategory.map((flow: Flow) => (
                      <Link key={flow.ID} href={`/experiments/${flow.ID}`} legacyBehavior>
                        <a className="w-full flex items-center text-sm px-3 py-2 hover:bg-muted rounded-full text-muted-foreground hover:text-foreground">
                          {flow.Name}
                        </a>
                      </Link>
                    ))}
                  </div>
                );
              }
              return null;
            })}
          </NavContent>
        </div>
      </ScrollArea>
      <div>
        {/* <NavContent>
          <NavLink href="https://github.com/labdao" target="_blank" icon={<GithubIcon />} title="GitHub" />
        </NavContent> */}
        <NavContent>
          <div className="p-2 font-mono text-xs font-bold text-muted-foreground opacity-70">Personal</div>
          <NavLink href="/experiments" icon={<FlaskRoundIcon />} title="Experiments" />
          <NavLink href="/tasks/protein-binder-design" icon={<SproutIcon />} title="Run Experiment" />
          <NavLink href="/data" icon={<FolderIcon />}>
            <>
              Files
            </>
          </NavLink>
          <AddDataFileForm trigger={<NavButton icon={<UploadIcon />} title="Upload Files" />} />
          <UserMenu />
        </NavContent>
      </div>
    </nav>
  );
}
