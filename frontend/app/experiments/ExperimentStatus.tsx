import { Job } from "@/lib/redux";
interface JobStatusProps {
  status: string;
}

export function JobStatus({ status }: JobStatusProps) {
  return <></>;
}

interface ExperimentStatusProps {
  jobs: Job[];
}

export function ExperimentStatus({ jobs }: ExperimentStatusProps) {
  let status;
  const totalJobs = jobs.length;
  const queuedJobs = jobs.filter((job) => job.State === "queued").length;
  const runningJobs = jobs.filter((job) => job.State === "running").length;
  const failedJobs = jobs.filter((job) => job.State === "failed").length;
  const completedJobs = jobs.filter((job) => job.State === "completed").length;
  const pendingJobs = queuedJobs + runningJobs;

  const percentComplete = Math.round((completedJobs / totalJobs) * 100);

  if (totalJobs === completedJobs) {
    status = "completed";
  } else if (failedJobs === totalJobs) {
    status = "failed";
  } else if (queuedJobs === totalJobs) {
    status = "queued";
  } else if (runningJobs > 0) {
    status = "running";
  } else if (failedJobs > 0 && pendingJobs === 0) {
    status = "partial-failure";
  } else {
    status = "unknown";
  }

  return (
    <>
      {status} {percentComplete}
    </>
  );
}
