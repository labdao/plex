import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { JobDetail } from "@/lib/redux";
import { cn } from "@/lib/utils";

export function aggregateJobStatus(jobs: JobDetail[]) {
  let status;
  let label;
  const totalJobs = jobs.length;
  const queuedJobs = jobs.filter((job) => job.JobStatus === "queued").length;
  const processingJobs = jobs.filter((job) => job.JobStatus === "processing").length;
  const pendingJobs = jobs.filter((job) => job.JobStatus === "pending").length;
  const runningJobs = jobs.filter((job) => job.JobStatus === "running").length;
  const failedJobs = jobs.filter((job) => job.JobStatus === "failed").length;
  const succeededJobs = jobs.filter((job) => job.JobStatus === "succeeded").length;
  const stoppedJobs = jobs.filter((job) => job.JobStatus === "stopped").length;

  const progressingJobs = queuedJobs + runningJobs + processingJobs + pendingJobs;
  const unsuccessfulJobs = failedJobs + stoppedJobs;

  if (totalJobs === 0) {
    status = "unknown";
    label = "Status unknown";
  } else if (totalJobs === succeededJobs) {
    status = "succeeded";
    label = "All runs completed successfully";
  } else if (unsuccessfulJobs === totalJobs) {
    status = "failed";
    label = "All runs failed or stopped";
  } else if (queuedJobs === totalJobs || processingJobs === totalJobs || pendingJobs === totalJobs) {
    status = "queued";
    label = "Waiting to start";
  } else if (progressingJobs > 0) {
    status = "running";
    label = "Running";
  } else if (failedJobs + stoppedJobs > 0) {
    status = "partial-failure";
    label = "Some runs completed successfully";
  } else {
    status = "unknown";
    label = "Status unknown";
  }

  return {
    status,
    label,
    totalJobs,
    progressingJobs,
    unsuccessfulJobs,
    succeededJobs,
  };
}

interface JobStatusIconProps {
  size?: number;
  status: string;
}

export function JobStatusIcon({ size = 12, status }: JobStatusIconProps) {
  return (
    <span className={cn("relative z-0", (status === "queued" || status === "running") && "animate-pulse")}>
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 12 12" width={size} height={size} fill="none">
        {(status === "failed" || status === "partial-failure" || status === "succeeded" || status === "unknown") && (
          <circle cx={6} cy={6} r={6} fill={`url(#${status})`} />
        )}
        {status === "queued" && (
          <>
            <path fill="none" d="M0 6a6 6 0 0 1 12 0H0Z" />
            <path fill={`url(#${status})`} d="M12 6A6 6 0 1 1 0 6h12Z" />
          </>
        )}
        {status === "running" && <path fill={`url(#${status})`} d="M12 6a6 6 0 1 1-6-6v6h6Z" />}
        <defs>
          <linearGradient id="completed" x1={17.5} x2={0} y1={-2.5} y2={16} gradientUnits="userSpaceOnUse">
            <stop stopColor="#6BDBAD" />
            <stop offset={1} stopColor="#00B269" stopOpacity={0} />
          </linearGradient>
          <linearGradient id="failed" x1={17.5} x2={0} y1={-2.5} y2={16} gradientUnits="userSpaceOnUse">
            <stop stopColor="#ff5a5a" />
            <stop offset={1} stopColor="#FF0000" stopOpacity={0} />
          </linearGradient>
          <linearGradient id="partial-failure" x1={17.5} x2={0} y1={-2.5} y2={16} gradientUnits="userSpaceOnUse">
            <stop stopColor="#ffb255" />
            <stop offset={1} stopColor="#ff8c00" stopOpacity={0} />
          </linearGradient>
          <linearGradient id="queued" x1={17.5} x2={0} y1={-2.5} y2={16} gradientUnits="userSpaceOnUse">
            <stop stopColor="#c8c8c8" />
            <stop offset={1} stopColor="#A5A5A5" stopOpacity={0} />
          </linearGradient>
          <linearGradient id="unknown" x1={17.5} x2={0} y1={-2.5} y2={16} gradientUnits="userSpaceOnUse">
            <stop stopColor="#c8c8c8" />
            <stop offset={1} stopColor="#A5A5A5" stopOpacity={0} />
          </linearGradient>
          <linearGradient id="running" x1={17.5} x2={0} y1={-2.5} y2={16} gradientUnits="userSpaceOnUse">
            <stop stopColor="#8fd9e9" />
            <stop offset={1} stopColor="#07A4C6" stopOpacity={0} />
          </linearGradient>
        </defs>
      </svg>
    </span>
  );
}

interface ExperimentStatusProps {
  jobs: JobDetail[];
  className?: string;
}

export function ExperimentStatus({ jobs, className }: ExperimentStatusProps) {
  const { status, label, progressingJobs, unsuccessfulJobs, succeededJobs } = aggregateJobStatus(jobs);
  return (
    <>
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger asChild>
            <span className={className}>
              <JobStatusIcon status={status || "unknown"} />
            </span>
          </TooltipTrigger>
          <TooltipContent>
            <div className="text-xs">
              <p className="mb-1 font-mono uppercase">{label}</p>
              {progressingJobs > 0 && <p>{progressingJobs} pending</p>}
              <p>{succeededJobs} complete</p>
              <p>{unsuccessfulJobs} failed</p>
            </div>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </>
  );
}
