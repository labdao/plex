import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { JobDetail } from "@/lib/redux";
import { cn } from "@/lib/utils";

export function aggregateJobStatus(jobs: JobDetail[]) {
  let status;
  let label;
  const totalJobs = jobs.length;
  const queuedJobs = jobs.filter((job) => job.Status === "queued").length;
  const runningJobs = jobs.filter((job) => job.Status === "running").length;
  const failedJobs = jobs.filter((job) => job.Status === "failed").length;
  const completedJobs = jobs.filter((job) => job.Status === "completed").length;

  //These statuses may be deprecated
  const errorJobs = jobs.filter((job) => job.Status === "error").length;
  const newJobs = jobs.filter((job) => job.Status === "new").length;

  const pendingJobs = queuedJobs + runningJobs + newJobs;

  if (totalJobs === 0) {
    status = "unknown";
    label = "Status unknown";
  } else if (totalJobs === completedJobs) {
    status = "completed";
    label = "All runs completed successfully";
  } else if (failedJobs === totalJobs || errorJobs === totalJobs) {
    status = "failed";
    label = "All runs failed";
  } else if (queuedJobs === totalJobs || newJobs === totalJobs) {
    status = "queued";
    label = "Waiting to start";
  } else if (pendingJobs > 0) {
    status = "running";
    label = "Running";
  } else if (failedJobs + errorJobs > 0) {
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
    pendingJobs,
    failedJobs,
    completedJobs,
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
        {(status === "failed" || status === "partial-failure" || status === "completed" || status === "unknown") && (
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
  const { status, label, pendingJobs, failedJobs, completedJobs } = aggregateJobStatus(jobs);
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
              {pendingJobs > 0 && <p>{pendingJobs} pending</p>}
              <p>{completedJobs} complete</p>
              <p>{failedJobs} failed</p>
            </div>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </>
  );
}
