"use client";

import backendUrl from "lib/backendUrl";
import React, { useEffect, useState } from "react";

import { CopyToClipboard } from "@/components/shared/CopyToClipboard";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";

const LogViewer = ({ jobID }: { jobID: string }) => {
  const [logs, setLogs] = useState("");

  useEffect(() => {
    if (jobID) {
      setLogs(`Connecting to stream with Bacalhau Job Id ${jobID}`);

      let formattedBackendUrl = backendUrl().replace("http://", "").replace("https://", "");
      let wsProtocol = backendUrl().startsWith("https://") ? "wss" : "ws";

      console.log(formattedBackendUrl);
      const ws = new WebSocket(`${wsProtocol}://${formattedBackendUrl}/jobs/${jobID}/logs`);

      ws.onopen = () => {
        console.log("connected");
      };

      ws.onmessage = (event) => {
        // Handle incoming message
        console.log(event.data);
        setLogs((prevLogs) => `${prevLogs}\n${event.data}`);
      };

      ws.onclose = () => {
        console.log("disconnected");
      };

      return () => {
        ws.close();
      };
    } else {
      setLogs(`Logs will stream after job is sent to the Bacalhau network`);
    }
  }, [jobID]);

  return (
    <div className="relative">
      <CopyToClipboard string={logs} className="absolute z-10 right-6 bg-background" />
      <ScrollArea className="w-full h-56 min-w-0 ">
        <pre className="p-6 font-mono text-xs">{logs}</pre>
        <ScrollBar orientation="vertical" />
        <ScrollBar orientation="horizontal" />
      </ScrollArea>
    </div>
  );
};

export default LogViewer;
