"use client";

import { Moon, SunMedium } from "lucide-react";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";

type ThemeMode = "light" | "dark";

function applyTheme(mode: ThemeMode) {
  const root = document.documentElement;
  root.classList.toggle("dark", mode === "dark");
}

export function ThemeToggle() {
  const [mode, setMode] = useState<ThemeMode>("light");

  useEffect(() => {
    const stored = window.localStorage.getItem("cuepoint-theme");
    const nextMode =
      stored === "dark" || stored === "light"
        ? (stored as ThemeMode)
        : window.matchMedia("(prefers-color-scheme: dark)").matches
          ? "dark"
          : "light";
    setMode(nextMode);
    applyTheme(nextMode);
  }, []);

  const toggle = () => {
    const nextMode = mode === "dark" ? "light" : "dark";
    setMode(nextMode);
    applyTheme(nextMode);
    window.localStorage.setItem("cuepoint-theme", nextMode);
  };

  return (
    <Button
      type="button"
      variant="ghost"
      onClick={toggle}
      className="w-full justify-between rounded-2xl border border-black/10 px-4 py-3 text-sm dark:border-white/10 dark:text-white dark:hover:bg-white/10"
    >
      <span>{mode === "dark" ? "Dark mode" : "Light mode"}</span>
      {mode === "dark" ? <Moon className="h-4 w-4" /> : <SunMedium className="h-4 w-4" />}
    </Button>
  );
}
