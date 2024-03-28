import { cn } from "@/lib/utils";
export default function PoweredByLogo({ className }: { className?: string }) {
  return (
    <div className={cn("flex items-center justify-center gap-3 p-3 text-xs text-muted-foreground/50 ", className)}>
      <svg xmlns="http://www.w3.org/2000/svg" width={15} height={16} fill="none">
        <path
          fill="currentColor"
          d="M0 8c0-2.8 1.1-4.7 3.3-5.6a12 12 0 0 0-1 5.6c0 2.4.3 4.2 1 5.5C1.1 12.5 0 10.7 0 8Zm7.3 6.3c1.6 0 2.9-.2 4-.7-.8 1.6-2.1 2.4-4 2.4s-3.2-.8-4-2.5c1 .6 2.4.8 4 .8ZM7.3 0c2 0 3.1.8 3.8 2.3-1-.4-2.3-.6-3.8-.6s-3 .2-4 .7a4.2 4.2 0 0 1 4-2.4ZM14 5.2h-2.2a9.1 9.1 0 0 0-.7-2.9c1.5.6 2.5 1.6 3 2.9Z"
        />
        <path fill="currentColor" d="M12 10.1h2.2a5 5 0 0 1-3 3.5c.4-1 .7-2 .7-3.5Z" />
      </svg>
      Powered by Convexity Labs
    </div>
  );
}
