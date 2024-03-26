import React, { useContext, useEffect } from "react";

import { TruncatedString } from "@/components/shared/TruncatedString";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import { FlowDetail } from "@/lib/redux";

import JobDetail from "./JobDetail";
import { ActiveResultContext } from "./ActiveResultContext";

interface JobsAccordionProps {
  flow: FlowDetail;
}

export default function JobsAccordion({ flow }: JobsAccordionProps) {
  const { activeJobUUID, setActiveJobUUID } = useContext(ActiveResultContext);
  useEffect(() => {
    if (!activeJobUUID) {
      setActiveJobUUID(flow.Jobs?.[0]?.JobUUID);
    }
  }, [flow.Jobs]);

  return (
    <Accordion type="single" value={activeJobUUID} onValueChange={setActiveJobUUID} className="min-h-[800px]">
      {flow.Jobs?.map((job, index) => {
        const validStates = ["queued", "running", "failed", "completed"];
        const status = (validStates.includes(job.State) ? job.State : "unknown") as "queued" | "running" | "failed" | "completed" | "unknown";

        return (
          <AccordionItem value={job.JobUUID} className="border-0 [&[data-state=open]>div]:shadow-lg" key={job.ID}>
            <Card className="my-2 shadow-sm">
              <AccordionTrigger className="flex items-center justify-between w-full px-6 py-3 text-left hover:no-underline [&[data-state=open]]:bg-muted">
                <div className="flex items-center gap-2">
                  <div className="w-30">
                    <div>run {index + 1}</div>
                    <div className="flex gap-1 text-xs text-muted-foreground/70">
                      Job ID: {job.BacalhauJobID ? <TruncatedString value={job.BacalhauJobID} /> : "n/a"}
                    </div>
                  </div>
                  <Badge status={status} variant="outline">
                    {job.State}
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
