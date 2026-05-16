import * as React from "react";

import { cn } from "@/lib/utils";

export function Select({
  className,
  children,
  ...props
}: React.SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      className={cn(
        "w-full rounded-xl border border-black/15 bg-white px-3 py-2 text-sm text-ink outline-none focus:border-ember dark:border-white/10 dark:bg-slate-950 dark:text-white dark:focus:border-sky",
        className,
      )}
      {...props}
    >
      {children}
    </select>
  );
}
