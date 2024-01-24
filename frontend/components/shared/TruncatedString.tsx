import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

export function TruncatedString({ value, trimLength = 6 }: { value: string; trimLength?: number }) {
  if (!value) return null;
  if (value?.length) {
    if (value?.length < trimLength * 2) return value;
    return (
      <span>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger>{`${value.substring(0, trimLength)}...${value.substring(value.length - trimLength)}`}</TooltipTrigger>
            <TooltipContent>
              <span className="text-xs">{value}</span>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </span>
    );
  }
}
