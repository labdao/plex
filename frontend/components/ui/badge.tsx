import { cva, type VariantProps } from "class-variance-authority";
import * as React from "react";

import { cn } from "@/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center rounded-full border  transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
  {
    variants: {
      variant: {
        default: "border-transparent bg-primary/20 text-primary-foreground hover:bg-primary/80",
        secondary: "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80",
        muted: "border-transparent bg-muted-foreground text-muted hover:bg-muted-foreground/80",
        destructive: "border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80",
        outline: "text-foreground",
      },
      size: {
        default: "px-3 py-0.5 text-xs",
        lg: "px-4 py-0.5 text-base",
      },
      status: {
        queued: "bg-gray-400/20",
        running: "bg-sky-500/20",
        completed: "bg-green-600/20",
        failed: "bg-destructive/20",
        "partial-failure": "bg-amber-500/20",
        unknown: "bg-gray-400/20",
      },
    },
    compoundVariants: [
      {
        variant: "outline",
        status: ["queued", "running", "completed", "failed", "partial-failure", "unknown"],
        class: "bg-transparent",
      },
      {
        variant: "outline",
        status: "queued",
        class: "border-gray-400 text-gray-400",
      },
      {
        variant: "outline",
        status: "running",
        class: "border-sky-500 text-sky-500",
      },
      {
        variant: "outline",
        status: "completed",
        class: "border-green-600 text-green-600",
      },
      {
        variant: "outline",
        status: "failed",
        class: "border-destructive text-destructive",
      },
      {
        variant: "outline",
        status: "partial-failure",
        class: "border-amber-500 text-amber-500",
      },
      {
        variant: "outline",
        status: "unknown",
        class: "border-gray-400 text-gray-400",
      },
    ],

    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
);

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement>, VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, size, status, ...props }: BadgeProps) {
  return <div className={cn(badgeVariants({ variant, size, status }), className)} {...props} />;
}

export { Badge, badgeVariants };
