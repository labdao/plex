import { Loader2Icon } from "lucide-react";

export function PageLoader() {
  return (
    <div className="flex justify-center w-full p-16 text-primary">
      <Loader2Icon className="animate-spin" size={48} absoluteStrokeWidth />
    </div>
  );
}
