import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import * as React from "react";

import { cn } from "@/lib/utils";

const buttonVariants = cva(
  "font-mono font-bold uppercase tracking-wider inline-flex items-center justify-center whitespace-nowrap rounded-full  ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-gradient-to-r from-secondary-blue to-primary-light text-primary-foreground hover:bg-primary/90",
        destructive: "text-destructive hover:bg-destructive/20",
        input: "rounded-md border border-input bg-background text-base font-normal tracking-normal lowercase font-body",
        outline: "border border-primary bg-transparent hover:bg-primary/10",
        secondary: "bg-secondary text-secondary-foreground hover:bg-secondary/80",
        ghost: "hover:bg-secondary hover:text-secondary-foreground",
        link: "text-accent underline-offset-4 hover:underline",
      },
      size: {
        default: "h-10 px-4 py-2 text-sm [&>svg]:mr-2",
        xs: "h-6  px-2 text-xs [&>svg]:mr-1",
        sm: "h-9  px-3 text-xs [&>svg]:mr-1",
        lg: "h-11  px-8 py-6 text-lg [&>svg]:mr-3",
        icon: "py-2 px-2 [&>svg]:m-0",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
);

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement>, VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(({ className, variant, size, type = "button", asChild = false, ...props }, ref) => {
  const Comp = asChild ? Slot : "button";
  return <Comp type={type} className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props} />;
});
Button.displayName = "Button";

export { Button, buttonVariants };
