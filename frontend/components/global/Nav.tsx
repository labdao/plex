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
        <NavLink
          href="/tasks/protein-binder-design"
          icon={<SproutIcon />}
          title="Design Molecule"
          className="mb-3 hover:before:bg-opacity-80 [&>*]:hover:text-foreground relative z-0 [&>*]:text-primary bg-gradient-to-tr from-primary-light to-primary before:-z-10 before:bg-white before:block before:absolute before:inset-[1px] before:rounded-full"
        />
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
          <div className="flex items-center gap-3 px-3 py-1 text-xs text-muted-foreground/50 ">
            <svg xmlns="http://www.w3.org/2000/svg" width={15} height={16} fill="none">
              <path
                fill="currentColor"
                d="M0 8c0-2.8 1.1-4.7 3.3-5.6a12 12 0 0 0-1 5.6c0 2.4.3 4.2 1 5.5C1.1 12.5 0 10.7 0 8Zm7.3 6.3c1.6 0 2.9-.2 4-.7-.8 1.6-2.1 2.4-4 2.4s-3.2-.8-4-2.5c1 .6 2.4.8 4 .8ZM7.3 0c2 0 3.1.8 3.8 2.3-1-.4-2.3-.6-3.8-.6s-3 .2-4 .7a4.2 4.2 0 0 1 4-2.4ZM14 5.2h-2.2a9.1 9.1 0 0 0-.7-2.9c1.5.6 2.5 1.6 3 2.9Z"
              />
              <path fill="currentColor" d="M12 10.1h2.2a5 5 0 0 1-3 3.5c.4-1 .7-2 .7-3.5Z" />
            </svg>
            Powered by <br /> Convexity Labs
          </div>
        </NavContent>
        <NavContent>
          <NavLink href="http://discord.gg/labdao" target="_blank" icon={<SiDiscord size={18} />}>
            Community
          </NavLink>
        </NavContent>
        <NavContent>
          <div className="p-2 font-mono text-xs font-bold text-muted-foreground opacity-70">Personal</div>
          <NavLink href="/experiments" icon={<FlaskRoundIcon />} title="Experiments" />
          <NavLink href="/data" icon={<FolderIcon />}>
            <>Files</>
          </NavLink>
          <AddDataFileForm trigger={<NavButton icon={<UploadIcon />} title="Upload Files" />} />
          <UserMenu />
        </NavContent>
      </div>
    </nav>
  );
}
