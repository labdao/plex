"use client";

import backendUrl from "lib/backendUrl";
import { CircleDotDashedIcon, DownloadIcon, HelpCircleIcon } from "lucide-react";
import React, { useEffect, useState } from "react";
import { CartesianGrid, Cell, ResponsiveContainer, Scatter, ScatterChart, Tooltip as RechartsTooltip, XAxis, YAxis } from "recharts";

import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import Molstar from "@/components/visualization/Molstar/index";
import { cn } from "@/lib/utils";

import { JobDetail } from "./JobDetail";

interface CustomTooltipProps {
  active?: boolean;
  payload?: any[];
  label?: string;
}

interface CheckpointChartData {
  cycle: number;
  proposal: number;
  plddt: number;
  i_pae: number;
  dim1: number;
  dim2: number;
  pdbFilePath: string;
}

interface CheckpointData {
  fileName: string;
  url: string;
}

export default function MetricsVisualizer({ job }: { job: JobDetail }) {
  const [loading, setLoading] = useState(false);
  const [checkpoints, setCheckpoints] = useState([]);
  const [plotData, setPlotData] = useState([]);
  const [activeCheckpointUrl, setActiveCheckpointUrl] = useState("");

  useEffect(() => {
    setLoading(true);
    const fetchData = async () => {
      try {
        const checkpointResponse = await fetch(`${backendUrl()}/checkpoints/${job.ID}`);
        const checkpointData = await checkpointResponse.json();
        setCheckpoints(checkpointData);

        const plotDataResponse = await fetch(`${backendUrl()}/checkpoints/${job.ID}/get-data`);
        const plotData = await plotDataResponse.json();
        setPlotData(plotData);
        if (activeCheckpointUrl === "") {
          setActiveCheckpointUrl(checkpointData?.[0]?.url);
        }
      } catch (error) {
        console.error("Error fetching data:", error);
      } finally {
        setLoading(false);
      }
    };
    if (job.ID) {
      fetchData();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [job]);

  const handlePointClick = (entry: CheckpointChartData) => {
    setActiveCheckpointUrl(entry.pdbFilePath);
  };

  const CustomTooltip: React.FC<CustomTooltipProps> = ({ active, payload }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      const keysToShow = Object.keys(data).filter((key) => key !== "pdbFilePath");

      return (
        <div className="max-w-xs p-3 text-xs bg-white border rounded shadow-md">
          {keysToShow.map((key) => (
            <p className="flex gap-1" key={key}>
              <span>{key}: </span>
              <strong className="max-w-full truncate">{data[key]}</strong>
            </p>
          ))}
        </div>
      );
    }

    return null;
  };
  if (!checkpoints?.length && !loading) {
    return (
      <div className="p-4 text-center text-muted-foreground">
        <CircleDotDashedIcon size={32} className={cn(job.State === "running" && "animate-spin", "mx-auto my-4")} absoluteStrokeWidth />
        <p>No checkpoints found yet. If this model has checkpoints, they&apos;ll appear here as they are reached.</p>
      </div>
    );
  }
  return (
    <>
      <div className="grid grid-cols-1 @sm:grid-cols-2">
        <div>
          <div>
            {!!checkpoints?.length && !!plotData?.length && (
              <div className="p-4 ">
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <div className="flex items-center justify-center gap-1 pl-10 mb-4 text-center font-heading">
                        Metrics Space <HelpCircleIcon className="text-muted-foreground" size={14} />
                      </div>
                    </TooltipTrigger>
                    <TooltipContent className="max-w-xs text-xs" side="bottom">
                      <p>
                        <strong>plddt:</strong> larger value indicates higher confidence in the predicted local structure
                      </p>
                      <p>
                        <strong>i_pae:</strong> larger value indicates higher confidence in the predicted interface residue distance
                      </p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
                <ResponsiveContainer width="100%" aspect={4 / 3}>
                  <ScatterChart margin={{ right: 10, bottom: 20 }}>
                    <CartesianGrid />
                    <XAxis type="number" dataKey="plddt" name="plddt" label={{ value: "plddt", position: "insideBottom", offset: -10, dx: 0 }} />
                    <YAxis
                      type="number"
                      dataKey="i_pae"
                      name="i_pae"
                      label={{ value: "i_pae", angle: -90, position: "insideLeft", dx: 15, dy: 15 }}
                    />
                    {/* <Tooltip cursor={{ strokeDasharray: '3 3' }} /> */}
                    <RechartsTooltip content={<CustomTooltip />} cursor={{ strokeDasharray: "3 3" }} />
                    <Scatter name="Checkpoints" data={plotData} fill="hsl(var(--accent))" onClick={handlePointClick}>
                      {plotData?.map((entry: CheckpointChartData, index) => (
                        <Cell
                          key={`cell-${index}`}
                          strokeWidth={entry.pdbFilePath === activeCheckpointUrl ? 8 : 0}
                          stroke="hsl(var(--primary))"
                          strokeOpacity={0.3}
                          paintOrder={"stroke"}
                        />
                      ))}
                    </Scatter>
                  </ScatterChart>
                </ResponsiveContainer>
              </div>
            )}
          </div>
        </div>
        {activeCheckpointUrl && <Molstar url={activeCheckpointUrl} />}
      </div>
      {/* <div>
        {checkpoints?.map((checkpoint: CheckpointData, index) => (
          <div
            key={index}
            className={cn(
              checkpoint.url === activeCheckpointUrl && "bg-primary/10",
              "flex items-center justify-between px-6 py-2 text-xs border-b border-border/50 last:border-none"
            )}
            onClick={() => setActiveCheckpointUrl(checkpoint?.url)}
          >
            <div className="truncate">{checkpoint.fileName}</div>
            <Button size="icon" variant="outline" asChild>
              <a href={checkpoint.url} download target="_blank" rel="noopener noreferrer">
                <DownloadIcon />
              </a>
            </Button>
          </div>
        ))}
      </div> */}
    </>
  );
}
