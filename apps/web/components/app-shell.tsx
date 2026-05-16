import Link from "next/link";
import { LayoutDashboard, Settings2, Tv, Wrench } from "lucide-react";

import { SignOutButton } from "@/components/sign-out-button";
import { ThemeToggle } from "@/components/theme-toggle";

const navItems = [
  { href: "/", label: "Dashboard", icon: LayoutDashboard },
  { href: "/plans", label: "Plans", icon: Tv },
  { href: "/settings/twitch", label: "Twitch", icon: Settings2 },
  { href: "/settings/bot", label: "Bot", icon: Settings2 },
  { href: "/settings/dev", label: "Dev", icon: Wrench },
];

export function AppShell({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="mx-auto flex min-h-screen max-w-[96rem] gap-6 px-4 py-6 sm:px-6">
      <aside className="sticky top-6 hidden h-[calc(100vh-3rem)] w-72 shrink-0 flex-col rounded-[2rem] border border-black/10 bg-white/85 p-5 backdrop-blur dark:border-white/10 dark:bg-slate-950/85 lg:flex">
        <div className="mb-8">
          <p className="text-xs uppercase tracking-[0.24em] text-black/40 dark:text-white/40">
            Cuepoint
          </p>
          <h1 className="mt-2 text-2xl font-semibold">Control Center</h1>
          <p className="mt-2 text-sm text-black/55 dark:text-white/55">
            Plan the stream, inspect runtime state, and prepare automation hooks.
          </p>
        </div>
        <nav className="space-y-2">
          {navItems.map((item) => {
            const Icon = item.icon;
            return (
              <Link
                key={item.href}
                href={item.href}
                className="flex items-center gap-3 rounded-2xl px-4 py-3 text-sm text-black/70 transition hover:bg-black/5 dark:text-white/75 dark:hover:bg-white/10"
              >
                <Icon className="h-4 w-4" />
                <span>{item.label}</span>
              </Link>
            );
          })}
        </nav>
        <div className="mt-auto space-y-3">
          <ThemeToggle />
          <SignOutButton />
        </div>
      </aside>
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="mb-6 flex items-center justify-between gap-4 rounded-[2rem] border border-black/10 bg-white/80 px-5 py-4 backdrop-blur dark:border-white/10 dark:bg-slate-950/75 lg:hidden">
          <div>
            <p className="text-xs uppercase tracking-[0.22em] text-black/45 dark:text-white/45">
              Cuepoint
            </p>
            <h1 className="text-xl font-semibold">Control Center</h1>
          </div>
          <div className="flex items-center gap-2">
            <ThemeToggle />
          </div>
        </header>
        <nav className="mb-4 flex gap-2 overflow-x-auto lg:hidden">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className="rounded-full border border-black/10 bg-white/80 px-4 py-2 text-sm text-black/70 dark:border-white/10 dark:bg-slate-950/75 dark:text-white/75"
            >
              {item.label}
            </Link>
          ))}
        </nav>
        <main className="flex-1">{children}</main>
      </div>
    </div>
  );
}
