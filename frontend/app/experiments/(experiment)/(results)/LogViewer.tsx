"use client";

import backendUrl from "lib/backendUrl";
import React, { useEffect, useRef, useState } from "react";

import { CopyToClipboard } from "@/components/shared/CopyToClipboard";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";

const LogViewer = ({ bacalhauJobID }: { bacalhauJobID: string }) => {
  const [logs, setLogs] = useState("");

  useEffect(() => {
    if (bacalhauJobID) {
      setLogs(`Connecting to stream with Bacalhau Job Id ${bacalhauJobID}`);

      let formattedBackendUrl = backendUrl().replace("http://", "").replace("https://", "");
      let wsProtocol = backendUrl().startsWith("https://") ? "wss" : "ws";

      console.log(formattedBackendUrl);
      const ws = new WebSocket(`${wsProtocol}://${formattedBackendUrl}/jobs/${bacalhauJobID}/logs`);

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
  }, [bacalhauJobID]);

  return (
    <div className="relative -mt-2">
      <CopyToClipboard string={logs} className="absolute z-10 top-2 right-6 bg-background" />
      <div className="flex flex-col-reverse w-full h-56 min-w-0 overflow-scroll no-scrollbar">
        <pre className="p-6 font-mono text-xs">{logs}</pre>
      </div>
    </div>
  );
};

export default LogViewer;
