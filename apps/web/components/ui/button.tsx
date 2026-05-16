import * as React from "react";

import { cn } from "@/lib/utils";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "default" | "outline" | "ghost";
};

export function Button({
  className,
  variant = "default",
  ...props
}: ButtonProps) {
  return (
    <button
      className={cn(
        "inline-flex items-center justify-center rounded-xl px-4 py-2 text-sm font-medium transition",
        variant === "default" &&
          "bg-ink text-sand hover:bg-black dark:bg-sky dark:text-slate-950 dark:hover:bg-sky/90",
        variant === "outline" &&
          "border border-ink/20 bg-white text-ink hover:bg-sand dark:border-white/10 dark:bg-slate-900 dark:text-white dark:hover:bg-slate-800",
        variant === "ghost" &&
          "bg-transparent text-ink hover:bg-black/5 dark:text-white dark:hover:bg-white/10",
        className,
      )}
      {...props}
    />
  );
}
