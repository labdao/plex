import React, { use, useContext, useEffect, useState } from "react";

import { TruncatedString } from "@/components/shared/TruncatedString";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import { ExperimentDetail } from "@/lib/redux";

import { ExperimentUIContext } from "../ExperimentUIContext";
import JobDetail from "./JobDetail";
import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "@/lib/backendUrl";

interface JobsAccordionProps {
  experiment: ExperimentDetail;
}

export default function JobsAccordion({ experiment }: JobsAccordionProps) {
  const { activeJobUUID, setActiveJobUUID } = useContext(ExperimentUIContext);
  const [jobs, setJobs] = useState(experiment.Jobs); // Initialize with initial job data


  useEffect(() => {
    if (!activeJobUUID) {
      setActiveJobUUID(jobs?.[0]?.RayJobID);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [jobs]);

  useEffect(() => {
    setJobs(experiment.Jobs);
  }, [experiment.Jobs]);

  // Polling to update job status and ID
  useEffect(() => {
    const fetchJobUpdates = async () => {
      try {
        const authToken = await getAccessToken();
        const nonFinalJobs = jobs.filter(
          (job) => !["failed", "succeeded", "stopped"].includes(job.JobStatus)
        );

        if (nonFinalJobs.length > 0) {
          const updatedJobs = await Promise.all(
            nonFinalJobs.map(async (job) => {
              const response = await fetch(`${backendUrl()}/jobs/${job.ID}`, {
                headers: {
                  Authorization: `Bearer ${authToken}`,
                },
              });
              const data = await response.json();
              return data;
            })
          );

          // Merge updated jobs back into the full jobs list
          const mergedJobs = jobs.map((job) =>
            updatedJobs.find((updatedJob) => updatedJob.ID === job.ID) || job
          );

          setJobs(mergedJobs);
        }
      } catch (error) {
        console.error("Error fetching job updates:", error);
      }
    };

    const intervalId = setInterval(() => {
      console.log("Polling for job updates");
      fetchJobUpdates();
    }, 10000); // Poll every 10 seconds
 
    return () => clearInterval(intervalId);
  }, [jobs]);


  return (
    <Accordion type="single" defaultValue={activeJobUUID} value={activeJobUUID} onValueChange={setActiveJobUUID} className="min-h-[600px]">
      {[...jobs]?.sort((a, b) => (a.ID || 0) - (b.ID || 0)).map((job, index) => {
        const validStates = ["queued", "processing", "pending", "running", "failed", "succeeded", "stopped"];
        const status = (validStates.includes(job.JobStatus) ? job.JobStatus : "unknown") as "queued" | "processing" | "pending" | "running" | "failed" | "succeeded" | "stopped" | "unknown";

        return (
          <AccordionItem value={job.RayJobID} className="border-0 [&[data-state=open]>div]:shadow-lg" key={job.ID}>
            <Card className="my-2 shadow-sm">
              <AccordionTrigger className="flex items-center justify-between w-full px-6 py-3 text-left hover:no-underline [&[data-state=open]]:bg-muted">
                <div className="flex items-center gap-2">
                  <div className="w-30">
                    <div>job {index + 1}</div>
                    <div className="flex gap-1 text-xs text-muted-foreground/70">
                      Job ID: {job.RayJobID ? <TruncatedString value={job.RayJobID} /> : "n/a"}
                    </div>
                  </div>
                  <Badge status={status} variant="outline">
                    {job.JobStatus}
                  </Badge>
                </div>
              </AccordionTrigger>
              <AccordionContent className="pb-0">
                <JobDetail jobID={job.ID} />
              </AccordionContent>
            </Card>
          </AccordionItem>
        );
      })}
    </Accordion>
  );
}
