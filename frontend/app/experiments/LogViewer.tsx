"use client";

import backendUrl from "lib/backendUrl";
import React, { useEffect, useState } from "react";
import { useSelector } from "react-redux";

import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { selectJobDetail } from "@/lib/redux";

const LogViewer = () => {
  const [logs, setLogs] = useState("");
  const job = useSelector(selectJobDetail);

  useEffect(() => {
    if (job.BacalhauJobID) {
      setLogs(`Connecting to stream with Bacalhau Job Id ${job.BacalhauJobID}`);

      let formattedBackendUrl = backendUrl().replace("http://", "").replace("https://", "");
      let wsProtocol = backendUrl().startsWith("https://") ? "wss" : "ws";

      console.log(formattedBackendUrl);
      const ws = new WebSocket(`${wsProtocol}://${formattedBackendUrl}/jobs/${job.BacalhauJobID}/logs`);

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
  }, [job]);

  return (
    <ScrollArea className="w-full h-56 min-w-0 ">
      <pre className="p-6 font-mono text-xs">{logs}</pre>
      <ScrollBar orientation="vertical" />
      <ScrollBar orientation="horizontal" />
    </ScrollArea>
  );
};

export default LogViewer;
