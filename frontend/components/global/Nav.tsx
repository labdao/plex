"use client";

import { SiDiscord } from "@icons-pack/react-simple-icons";
import { usePrivy } from "@privy-io/react-auth";
import { PencilIcon, SproutIcon } from "lucide-react";
import Link from "next/link";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { ScrollArea } from "@/components/ui/scroll-area";
import { AppDispatch, Experiment, experimentListThunk, selectCategorizedExperiments, selectExperimentList, selectExperimentListLoading, selectUserIsAdmin } from "@/lib/redux";

import Logo from "./Logo";
import { NavLink } from "./NavItem";
import UserMenu from "./UserMenu";
import { Button } from "../ui/button";
import { cn } from "@/lib/utils";
import { useParams } from "next/navigation";
import { toast } from "sonner";
import { InlineEditExperiment } from "@/components/ui/inline-edit-experiment";

export default function Nav() {
  const { user } = usePrivy();
  const { experimentID } = useParams();
  const dispatch = useDispatch<AppDispatch>();
  const categorizedExperiments = useSelector(selectCategorizedExperiments);
  const experiments = useSelector(selectExperimentList);
  // const loading = useSelector(selectExperimentListLoading);
  const walletAddress = user?.wallet?.address;
  const isAdmin = useSelector(selectUserIsAdmin);

  useEffect(() => {
    console.log("walletAddress", walletAddress);
    if (walletAddress) {
      console.log("dispatching experimentListThunk");
      dispatch(experimentListThunk(walletAddress));
    }
  }, [dispatch, walletAddress]);

  return (
    <nav className="sticky top-0 z-50 flex flex-col justify-between w-48 h-screen border-r shadow-lg border-border/50 shrink-0 bg-background">
      <Link href="/" className="flex items-center h-12 gap-2 p-2 text-lg font-bold uppercase font-heading whitespace-nowrap">
        <Logo className="w-auto h-6 text-primary" />
        Lab.Bio
        {isAdmin && <sup className="text-xs text-primary">Admin</sup>}
      </Link>
      <div className="px-2 py-2">
        <Button asChild color="primary" size="sm" className="w-full">
          <Link href="/experiments/new/protein-binder-design">
            <SproutIcon /> Design Molecule
          </Link>
        </Button>
      </div>
      <ScrollArea className="flex-grow border-b border-border/50">
        <div className="w-48 p-2">
          {Object.keys(categorizedExperiments).map((category) => {
            const experimentsInCategory = categorizedExperiments[category as keyof typeof categorizedExperiments];
            if (experimentsInCategory.length > 0) {
              return (
                <div key={category} className="flex flex-col gap-1 mb-5">
                  <div className="px-3 mb-2 font-mono text-xs font-bold uppercase text-muted-foreground opacity-70">
                    {category === "today" && "Today"}
                    {category === "last7Days" && "Previous 7 Days"}
                    {category === "last30Days" && "Previous 30 Days"}
                    {category === "older" && "Older"}
                  </div>
                  {experimentsInCategory.map((experiment: Experiment) => (
                    <Link
                      key={experiment.ID}
                      href={`/experiments/${experiment.ID}`}
                      className={cn(
                        "relative group px-3 rounded-full py-2 text-sm truncate hover:bg-muted/50 text-muted-foreground hover:text-foreground",
                        experimentID === experiment.ID.toString() && "text-foreground bg-muted hover:bg-muted"
                      )}
                    >
                      <InlineEditExperiment experiment={experiment} />
                    </Link>
                  ))}
                </div>
              );
            }
            return null;
          })}
        </div>
      </ScrollArea>
      <div className="p-2">
        <NavLink href="http://discord.gg/labdao" target="_blank" icon={<SiDiscord size={18} />}>
          Community
        </NavLink>
        <UserMenu />
      </div>
    </nav>
  );
}
