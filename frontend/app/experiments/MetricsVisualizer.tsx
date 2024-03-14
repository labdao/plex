"use client";

import backendUrl from "lib/backendUrl";
import { CircleDotDashedIcon, DownloadIcon, HelpCircleIcon } from "lucide-react";
import React, { useEffect, useState } from "react";
import { CartesianGrid, Cell, ResponsiveContainer, Scatter, ScatterChart, Tooltip as RechartsTooltip, XAxis, YAxis, ReferenceLine, ReferenceArea } from "recharts";

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
      const formattedData = { ...data };

      if (formattedData.plddt !== undefined) {
        formattedData.plddt = formattedData.plddt.toFixed(2);
      }
      if (formattedData.i_pae !== undefined) {
        formattedData.i_pae = (-formattedData.i_pae).toFixed(2);
      }
      const keysToShow = Object.keys(data).filter((key) => key !== "pdbFilePath" && key !== "checkpoint");

      return (
        <div className="max-w-xs p-3 text-xs bg-white border rounded shadow-md">
          {keysToShow.map((key) => (
            <p className="flex gap-1" key={key}>
              <span>{key}: </span>
              <strong className="max-w-full truncate">{formattedData[key]}</strong>
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
                        <strong>Stability Score:</strong> larger value indicates higher confidence in the predicted local structure
                      </p>
                      <p>
                        <strong>Affinity Score:</strong> larger value indicates higher confidence in the predicted interface residue distance
                      </p>
                      <p><strong>Note:</strong> designs which lie within the green square are likely high quality</p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
                <ResponsiveContainer width="100%" aspect={4 / 3}>
                  <ScatterChart margin={{ right: 10, bottom: 20 }}>
                    <CartesianGrid />
                    <XAxis 
                      type="number" 
                      dataKey="plddt" 
                      name="Stability Score (plddt)" 
                      domain={[0, 100]}
                      label={{ value: "Stability Score (plddt)", position: "insideBottom", offset: -10, dx: 0 }} />
                    <YAxis
                      type="number"
                      dataKey={(data) => -data.i_pae} // Invert i_pae values
                      name="Affinity Score (-i_pae)"
                      domain={[-40, 0]}
                      label={{ value: "Affinity Score (-i_pae)", angle: -90, position: "insideLeft", dx: 10, dy: 65 }}
                    />
                    {/* <Tooltip cursor={{ strokeDasharray: '3 3' }} /> */}
                    <RechartsTooltip content={<CustomTooltip />} cursor={{ strokeDasharray: "3 3" }} />
                    {/* plddt cutoff line */}
                    <ReferenceLine x={80} stroke="black" strokeDasharray="3 3" />
                    {/* -i_pae lower cutoff line */}
                    <ReferenceLine y={-10} stroke="black" strokeDasharray="3 3" />
                    {/* Color the top right quadrant */}
                    <ReferenceArea x1={80} y1={-10} strokeOpacity={0.3} fill="#6BDBAD" fillOpacity={0.3} />
                    <Scatter name="Checkpoints" data={plotData} fill="#000000" onClick={handlePointClick}>
                      {plotData?.map((entry: CheckpointChartData, index) => (
                        <Cell
                          key={`cell-${index}`}
                          strokeWidth={entry.pdbFilePath === activeCheckpointUrl ? 8 : 0}
                          stroke="#6BDBAD"
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
    </>
  );
}
