import { cn } from "@/lib/utils";

export function Card({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        "rounded-[1.5rem] border border-black/10 bg-white/90 p-6 shadow-sm dark:border-white/10 dark:bg-slate-900/85 dark:text-white",
        className,
      )}
      {...props}
    />
  );
}
