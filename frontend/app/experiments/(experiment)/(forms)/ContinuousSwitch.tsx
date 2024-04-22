import { Badge } from "@/components/ui/badge";
import { FormControl, FormDescription, FormLabel } from "@/components/ui/form";
import { Switch } from "@/components/ui/switch";

export default function ContinuousSwitch() {
  return (
    <div className="flex flex-row items-center gap-4">
      <FormControl>
        <Switch disabled />
      </FormControl>
      <div className="space-y-0.5">
        <FormLabel className="pt-0">
          continuous run <Badge>Coming soon!</Badge>
        </FormLabel>
        <FormDescription>continuously generate datapoints with the current sequence inputs. tURN OFF OR EDIT INPUT TO CANCEL.</FormDescription>
      </div>
    </div>
  );
}
