"use client";

import { SiDiscord } from "@icons-pack/react-simple-icons";
import { usePrivy } from "@privy-io/react-auth";
import AddDataFileForm from "app/data/AddDataFileForm";
import TasksMenu from "app/tasks/TasksMenu";
import { BoxIcon, FlaskRoundIcon, FolderIcon, GithubIcon, UploadIcon } from "lucide-react";
import { SproutIcon } from "lucide-react";
import Link from "next/link";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { ExperimentStatus } from "@/app/experiments/(experiment)/ExperimentStatus";
import { NavButton } from "@/components/global/NavItem";
import { ScrollArea } from "@/components/ui/scroll-area";
import { AppDispatch, Flow, flowListThunk, selectCategorizedFlows, selectFlowList, selectFlowListLoading, selectUserIsAdmin } from "@/lib/redux";

import Logo from "./Logo";
import { NavLink } from "./NavItem";
import UserMenu from "./UserMenu";
import { Button } from "../ui/button";

const NavContent = (props: React.PropsWithChildren) => <div className="flex flex-col p-1 border-b border-border/50" {...props} />;

export default function Nav() {
  const { user } = usePrivy();
  const dispatch = useDispatch<AppDispatch>();
  const categorizedFlows = useSelector(selectCategorizedFlows);
  const flows = useSelector(selectFlowList);
  // const loading = useSelector(selectFlowListLoading);
  const walletAddress = user?.wallet?.address;
  const isAdmin = useSelector(selectUserIsAdmin);

  useEffect(() => {
    console.log("walletAddress", walletAddress);
    if (walletAddress) {
      console.log("dispatching flowListThunk");
      dispatch(flowListThunk(walletAddress));
    }
  }, [dispatch, walletAddress]);

  return (
    <nav className="sticky top-0 z-50 flex flex-col justify-between w-48 h-screen border-r shadow-lg border-border/50 shrink-0 bg-background">
      <Link href="/" className="flex items-center h-12 gap-2 p-2 text-lg font-bold uppercase font-heading whitespace-nowrap">
        <Logo className="w-auto h-6 text-primary" />
        Lab.Bio
        {isAdmin && <sup className="text-xs text-primary">Admin</sup>}
      </Link>
      <NavContent>
        <Button asChild color="primary" size="sm" className="w-full mb-2">
          <Link href="/experiments/new/protein-binder-design">
            <SproutIcon /> Design Molecule
          </Link>
        </Button>
      </NavContent>
      <div>
        <div className="p-2 font-mono text-xs font-bold text-muted-foreground opacity-70">Experiments</div>
      </div>
      <ScrollArea className="flex-grow border-b border-border/50">
        <div className="flex flex-col overflow-auto">
          {Object.keys(categorizedFlows).map((category) => {
            const flowsInCategory = categorizedFlows[category as keyof typeof categorizedFlows];
            if (flowsInCategory.length > 0) {
              return (
                <div key={category} className="px-2 mb-4">
                  <div className="p-3 text-xs font-bold text-muted-foreground opacity-70">
                    {category === "today" && "Today"}
                    {category === "last7Days" && "Previous 7 Days"}
                    {category === "last30Days" && "Previous 30 Days"}
                    {category === "older" && "Older"}
                  </div>
                  {flowsInCategory.map((flow: Flow) => (
                    <Link key={flow.ID} href={`/experiments/${flow.ID}`} legacyBehavior>
                      <a className="flex items-center w-full px-3 py-2 text-sm rounded-full hover:bg-muted text-muted-foreground hover:text-foreground">
                        {flow.Name}
                      </a>
                    </Link>
                  ))}
                </div>
              );
            }
            return null;
          })}
        </div>
      </ScrollArea>
      <div>
        <NavContent>
          <NavLink href="http://discord.gg/labdao" target="_blank" icon={<SiDiscord size={18} />}>
            Community
          </NavLink>
        </NavContent>
        <NavContent>
          <UserMenu />
        </NavContent>
      </div>
    </nav>
  );
}
