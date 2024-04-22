"use client";

import backendUrl from "lib/backendUrl";
import { HelpCircleIcon } from "lucide-react";
import React, { useContext, useEffect, useState } from "react";
import {
  CartesianGrid,
  Cell,
  ReferenceArea,
  ReferenceLine,
  ResponsiveContainer,
  Scatter,
  ScatterChart,
  Tooltip as RechartsTooltip,
  XAxis,
  YAxis,
} from "recharts";

import { PageLoader } from "@/components/shared/PageLoader";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import Molstar from "@/components/visualization/Molstar/index";
import { FlowDetail } from "@/lib/redux";
import { getAccessToken } from "@privy-io/react-auth";

import { aggregateJobStatus } from "../ExperimentStatus";
import { ExperimentUIContext } from "../ExperimentUIContext";

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
  jobUUID: string;
}

interface CheckpointData {
  filename: string;
  url: string;
}

export default function MetricsVisualizer({ flow }: { flow: FlowDetail }) {
  const [loading, setLoading] = useState(false);
  const [checkpoints, setCheckpoints] = useState([]);
  const [plotData, setPlotData] = useState([]);
  const { activeCheckpointUrl, setActiveCheckpointUrl, activeJobUUID, setActiveJobUUID } = useContext(ExperimentUIContext);

  const { status: flowStatus } = aggregateJobStatus(flow.Jobs || []);

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Currently unused
        //const checkpointResponse = await fetch(`${backendUrl()}/checkpoints/${flow.ID}`);
        //const checkpointData = await checkpointResponse.json();
        //setCheckpoints(checkpointData);

        let authToken;
        try {
          authToken = await getAccessToken()
        } catch (error) {
          console.log('Failed to get access token: ', error)
          throw new Error("Authentication failed");
        }
        const plotDataResponse = await fetch(`${backendUrl()}/checkpoints/${flow.ID}/get-data`, {
          method: 'Get',
          headers: {
            'Authorization': `Bearer ${authToken}`,
            'Content-Type': 'application/json',
          },
        });
        const plotData = await plotDataResponse.json();
        setPlotData(plotData);
      } catch (error) {
        console.error("Error fetching data:", error);
      } finally {
        setLoading(false);
      }
    };
    if (flow.ID) {
      setLoading(true);
      fetchData();
    } else {
      setPlotData([]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [flow]);

  useEffect(() => {
    if (plotData?.length > 0) {
      if (!activeCheckpointUrl) {
        console.log("Setting default active checkpoint url");
        setActiveCheckpointUrl((plotData?.[0] as CheckpointChartData)?.pdbFilePath);
      }
    } else {
      setActiveCheckpointUrl(undefined);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [plotData]);

  const handlePointClick = (entry: CheckpointChartData) => {
    setActiveCheckpointUrl(entry.pdbFilePath);
    setActiveJobUUID(entry.jobUUID);
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
              <span>{key === "i_pae" ? "-i_pae" : key}: </span>
              <strong className="max-w-full truncate">{formattedData[key]}</strong>
            </p>
          ))}
        </div>
      );
    }

    return null;
  };

  return (
    <div className="relative">
      {!plotData?.length && (
        <div className="absolute inset-0 z-20 flex items-center justify-center p-12 text-center text-muted-foreground bg-gray-50/80">
          <div>{loading ? <PageLoader /> : <p>Checkpoints will appear here as they complete.</p>}</div>
        </div>
      )}

      <div className="grid grid-cols-2">
        <div>
          <div>
            <div className="p-4 select-none">
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
                    <p>
                      <strong>Note:</strong> designs which lie within the green square are likely high quality
                    </p>
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
                    label={{ value: "Stability Score (plddt)", position: "insideBottom", offset: -10, dx: 0 }}
                  />
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
                  <Scatter name="Checkpoints" data={plotData} onClick={handlePointClick}>
                    {plotData?.map((entry: CheckpointChartData, index) => {
                      return (
                        <Cell
                          key={`cell-${index}-${activeCheckpointUrl}`}
                          strokeWidth={entry.pdbFilePath === activeCheckpointUrl ? 8 : 0}
                          stroke="#6BDBAD"
                          strokeOpacity={0.3}
                          paintOrder={"stroke"}
                          fill={entry.jobUUID === activeJobUUID ? "#000000" : "#959595"}
                        />
                      );
                    })}
                  </Scatter>
                </ScatterChart>
              </ResponsiveContainer>
            </div>
          </div>
        </div>
        <Molstar url={activeCheckpointUrl} />
      </div>
    </div>
  );
}