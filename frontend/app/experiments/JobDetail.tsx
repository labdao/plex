"use client";

import { getAccessToken } from "@privy-io/react-auth";
import { ColumnDef } from "@tanstack/react-table";
import backendUrl from "lib/backendUrl";
import { DownloadIcon } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { ScatterChart, Scatter, XAxis, YAxis, CartesianGrid, Tooltip } from 'recharts';

import MolstarComponent from "@/components/Molstar";
import { CopyToClipboard } from "@/components/shared/CopyToClipboard";
import { TruncatedString } from "@/components/shared/TruncatedString";
import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
import { DataTable } from "@/components/ui/data-table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { DataFile } from "@/lib/redux";

import LogViewer from "./LogViewer";

interface JobDetailProps {
  jobID: number;
}

interface CustomTooltipProps {
  active?: boolean;
  payload?: any[];
  label?: string;
}

interface CheckpointData {
  cycle: number;
  proposal: number;
  plddt: number;
  i_pae: number;
  dim1: number;
  dim2: number;
  pdbFilePath: string;
}

export interface JobDetail {
  ID: number | null;
  BacalhauJobID: string;
  JobUUID: string;
  State: string;
  Error: string;
  ToolID: string;
  FlowID: string;
  Inputs: {};
  InputFiles: DataFile[];
  OutputFiles: DataFile[];
  Status: string;
}

export default function JobDetail({ jobID }: JobDetailProps) {
  const [job, setJob] = useState({} as JobDetail);
  const [loading, setLoading] = useState(false);
  const [checkpoints, setCheckpoints] = useState([]);
  const [plotData, setPlotData] = useState([]);
  const [moleculeUrl, setMoleculeUrl] = useState('');
  const [activeTab, setActiveTab] = useState('dashboard');


  interface File {
    CID: string;
    Filename: string;
    Tags: Tag[];
  }

  interface Tag {
    Name: string;
    Type: string;
  }

  useEffect(() => {
    setLoading(true);
    const fetchData = async () => {
      try {
        const authToken = await getAccessToken(); // Get the access token
        const response = await fetch(`${backendUrl()}/jobs/${jobID}`, {
          headers: {
            Authorization: `Bearer ${authToken}`, // Include the authorization header
          },
        });

        if (!response.ok) {
          throw new Error(`HTTP error ${response.status}`);
        }

        const data = await response.json();
        console.log("Fetched job:", data);
        setJob(data);
      } catch (error) {
        console.error("Error fetching job:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();

    fetch(`${backendUrl()}/checkpoints/${jobID}`)
      .then(response => response.json())
      .then(data => {
        setCheckpoints(data);
      })
      .catch(error => console.error('Error fetching checkpoints:', error));

  }, [jobID]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const checkpointResponse = await fetch(`${backendUrl()}/checkpoints/${jobID}`);
        const checkpointData = await checkpointResponse.json();
        setCheckpoints(checkpointData);
  
        const plotDataResponse = await fetch(`${backendUrl()}/checkpoints/${jobID}/get-data`);
        const plotData = await plotDataResponse.json();
        setPlotData(plotData);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };
  
    if (job.State === "running") {
      fetchData();
      const intervalId = setInterval(fetchData, 5000);
  
      return () => clearInterval(intervalId);
    } else { //(job.State === "completed") {
      fetchData();
    }
  }, [jobID, job.State]);

  const handlePointClick = (data: CheckpointData) => {
    console.log('Clicked point data:', data);
    setMoleculeUrl(data.pdbFilePath);
    console.log("set molecule url:", data.pdbFilePath);
    // Switch to the visualize tab
    console.log(activeTab);
    setActiveTab('visualize');
    console.log(activeTab);
  };

  const CustomTooltip: React.FC<CustomTooltipProps> = ({ active, payload }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload; 
      const keysToShow = Object.keys(data).filter(key => key !== 'pdbFilePath');
  
      return (
        <div className="bg-white p-3 border rounded shadow-lg">
          {keysToShow.map((key) => (
            <p key={key}>
              {key}: {data[key]}
            </p>
          ))}
        </div>
      );
    }
  
    return null;
  };

  return (
    <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full @container ">
      <TabsList className="justify-start w-full px-6 pt-0 rounded-t-none">
        <TabsTrigger value="dashboard">Dashboard</TabsTrigger>
        <TabsTrigger value="visualize">Visualize</TabsTrigger>
        <TabsTrigger value="parameters">Parameters</TabsTrigger>
        <TabsTrigger value="outputs">Outputs</TabsTrigger>
        <TabsTrigger value="inputs">Inputs</TabsTrigger>
        <TabsTrigger value="logs">Logs</TabsTrigger>
      </TabsList>
      <TabsContent value="dashboard">
        <div style={{ display: 'flex', alignItems: 'center', marginTop: '15px', justifyContent: 'space-around' }}> 
        <div>
          <div className="font-heading" style={{ marginLeft: '70px', textAlign: 'center' }}>Metrics Space</div>
          <ScatterChart width={730} height={250} margin={{ top: 20, right: 20, bottom: 20, left: 20}}>
            <CartesianGrid />
            <XAxis type="number" dataKey="plddt" name="plddt" label={{ value: 'plddt', position: 'insideBottom', offset: -15, dx: 0 }}/>
            <YAxis type="number" dataKey="i_pae" name="i_pae" label={{ value: 'i_pae', angle: -90, position: 'insideLeft', dx: 0, dy: 15 }}/>
            {/* <Tooltip cursor={{ strokeDasharray: '3 3' }} /> */}
            <Tooltip content={<CustomTooltip />} cursor={{ strokeDasharray: '3 3' }} />
            <Scatter name="Checkpoints" data={plotData} fill="#8884d8" onClick={handlePointClick} />
          </ScatterChart>
          </div>
          <div style={{ maxWidth: '300px', fontSize: '0.75rem', color: 'grey', textAlign: 'left' }}>
            <p><b>plddt:</b> larger value indicates higher confidence in the predicted local structure</p>
            <p><b>i_pae:</b> larger value indicates higher confidence in the predicted interface residue distance</p>
            <p style = {{fontSize: '0.80rem', color: 'black'}}>Click on a datapoint to visualize</p>
          </div>
        </div>
        <CheckpointsList checkpoints={checkpoints} />
      </TabsContent>
      <TabsContent value="visualize">
        <MolstarComponent 
          moleculeUrl={moleculeUrl}
          customDataFormat="pdb" 
        />
      </TabsContent>
      <TabsContent value="parameters" className="px-6 pt-0">
        {Object.entries(job.Inputs || {}).map(([key, val]) => (
          <div key={key} className="flex justify-between py-1 text-base border-b last:border-none last:mb-3">
            <span className="text-muted-foreground/50">{key.replaceAll("_", " ")}</span>
            <span>{val ? <TruncatedString value={val.toString()} trimLength={10} /> : <span className="text-muted-foreground">n/a</span>}</span>
          </div>
        ))}
      </TabsContent>
      <TabsContent value="outputs">
        <FileList files={job.OutputFiles} />
      </TabsContent>
      <TabsContent value="inputs">
        <FileList files={job.InputFiles} />
      </TabsContent>
      <TabsContent value="logs">
        <div className="w-full">
          <LogViewer jobID={job.BacalhauJobID} />
        </div>
      </TabsContent>
    </Tabs>
  );
}

function CheckpointsList({ checkpoints }: { checkpoints: Array<{ fileName: string, url: string }> }) {
  const safeCheckpoints = checkpoints || [];
  return (
    <div>
      {safeCheckpoints.length > 0 ? (
        safeCheckpoints.map((checkpoint, index) => (
          <div key={index} className="flex items-center justify-between px-6 py-2 text-xs border-b border-border/50 last:border-none">
            <div>
              <span className="text-accent">{checkpoint.fileName}</span>
            </div>
            <Button size="icon" variant="outline" asChild>
              <a href={checkpoint.url} download target="_blank" rel="noopener noreferrer">
                <DownloadIcon />
              </a>
            </Button>
          </div>
        ))
      ) : (
        <p>Checkpoints will start appearing soon... Please check back later</p>
      )}
    </div>
  );
}

function FileList({ files }: { files: DataFile[] }) {
  const handleDownload = async (file: DataFile) => {
    try {
      const authToken = await getAccessToken();
      const response = await fetch(`${backendUrl()}/datafiles/${file.CID}/download`, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });
      if (!response.ok) {
        throw new Error("Failed to download file");
        return;
      }
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = file.Filename || "download";
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (error) {
      console.error(error);
    }
  };
  return (
    <div>
      {!!files?.length ? (
        <>
          {files.map((file: DataFile) => (
            <div key={file.CID} className="flex items-center justify-between px-6 py-2 text-xs border-b border-border/50 last:border-none">
              <div>
                <a target="#" onClick={() => handleDownload(file)} className="text-accent" style={{ cursor: "pointer" }}>
                  <TruncatedString value={file.Filename} trimLength={30} />
                </a>
                <div className="opacity-70 text-muted-foreground">
                  <CopyToClipboard string={file.CID}>
                    cid: <TruncatedString value={file.CID} />
                  </CopyToClipboard>
                </div>
              </div>
              {/* @TODO: Add Filesize */}
              <Button size="icon" variant="outline" asChild>
                <a target="#" onClick={() => handleDownload(file)} style={{ cursor: "pointer" }}>
                  <DownloadIcon />
                </a>
              </Button>
            </div>
        ))}
        </>
      ) : (
        <div className="px-6 py-5">No files found.</div>
      )}
      </div>
  );
}
