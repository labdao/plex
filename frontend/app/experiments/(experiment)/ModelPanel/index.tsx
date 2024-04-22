"use client";

import { HelpCircleIcon, InfoIcon, PanelRightCloseIcon, PanelRightOpenIcon } from "lucide-react";
import React, { useContext, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import { ToolSelect } from "@/components/shared/ToolSelect";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { AppDispatch, selectToolDetail, selectToolDetailLoading, ToolDetail, toolDetailThunk } from "@/lib/redux";
import { cn } from "@/lib/utils";

import { ExperimentUIContext } from "../ExperimentUIContext";
import ModelGuide from "./ModelGuide";
import ModelInfo from "./ModelInfo";

interface ModelInfoProps {
  task?: {
    slug: string;
    name: string;
    available: boolean;
  };
  showSelect?: boolean;
  defaultOpen?: boolean;
}

export default function ModelPanel({ task, defaultOpen, showSelect }: ModelInfoProps) {
  const { modelPanelOpen, setModelPanelOpen } = useContext(ExperimentUIContext);
  const [activeTab, setActiveTab] = useState("info");
  const dispatch = useDispatch<AppDispatch>();
  const tool = useSelector(selectToolDetail);
  const toolDetailLoading = useSelector(selectToolDetailLoading);

  useEffect(() => {
    if (!toolDetailLoading) {
      setModelPanelOpen(Boolean(defaultOpen));
    }
  }, [toolDetailLoading, defaultOpen, setModelPanelOpen]);

  const handleToolChange = (value: any) => {
    dispatch(toolDetailThunk(value));
    setActiveTab("info");
  };

  const handleOpen = () => {
    setModelPanelOpen(!modelPanelOpen);
  };

  const handleTabChange = (tab: string) => {
    setActiveTab(tab);
    if (!open) {
      setModelPanelOpen(true);
    }
  };

  return (
    <Card
      className={cn(
        "transition-all lg:rounded-r-none m-2 lg:mx-0 lg:my-2 lg:sticky top-10 grow-0 h-[calc(100vh-7rem)] shrink-0 basis-14",
        modelPanelOpen && "basis-1/3"
      )}
    >
      <div className={cn("min-w-[26vw] flex flex-col h-full", modelPanelOpen && "opacity-1")}>
        <div className="flex items-center gap-3 p-3 border-b">
          <div>
            <Button size="icon" variant="ghost" className="" onClick={handleOpen}>
              {modelPanelOpen ? <PanelRightCloseIcon /> : <PanelRightOpenIcon />}
            </Button>
          </div>
          {showSelect ? (
            <ToolSelect onChange={handleToolChange} taskSlug={task?.slug} />
          ) : (
            <div className="ml-2 text-xl truncate font-heading">
              {tool.ToolJson?.author || "unknown"}/{tool.ToolJson?.name}
            </div>
          )}
        </div>

        <div className="flex grow">
          <div className="flex flex-col justify-start h-auto gap-2 p-3">
            <Button
              onClick={() => handleTabChange("info")}
              variant="ghost"
              size="icon"
              className={activeTab === "info" && modelPanelOpen ? "bg-muted" : undefined}
            >
              <InfoIcon size={48} />
            </Button>
            {tool.ToolJson?.guide && (
              <Button
                onClick={() => handleTabChange("guide")}
                variant="ghost"
                size="icon"
                className={activeTab === "guide" && modelPanelOpen ? "bg-muted" : undefined}
              >
                <HelpCircleIcon />
              </Button>
            )}
          </div>
          <div className="p-2 mt-2 grow">
            {tool.ToolJson?.description && activeTab === "info" && <ModelInfo tool={tool} />}
            {tool.ToolJson?.guide && activeTab === "guide" && <ModelGuide tool={tool} />}
          </div>
        </div>
      </div>
    </Card>
  );
}