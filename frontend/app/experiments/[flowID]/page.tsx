import ExperimentDetail from "../ExperimentDetail";

export default function ExperimentDetailPage({ params }: { params: { flowID: string } }) {
  return (
    <div className="max-w-6xl p-4">
      <ExperimentDetail experimentID={params.flowID} />
    </div>
  );
}
