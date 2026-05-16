import * as React from "react";

import { cn } from "@/lib/utils";

export function Textarea({
  className,
  ...props
}: React.TextareaHTMLAttributes<HTMLTextAreaElement>) {
  return (
    <textarea
      className={cn(
        "min-h-28 w-full rounded-xl border border-black/15 bg-white px-3 py-2 text-sm text-ink outline-none placeholder:text-black/35 focus:border-ember dark:border-white/10 dark:bg-slate-950 dark:text-white dark:placeholder:text-white/35 dark:focus:border-sky",
        className,
      )}
      {...props}
    />
  );
}
