import * as React from "react";

export function Switch(props: React.InputHTMLAttributes<HTMLInputElement>) {
  return <input type="checkbox" className="h-4 w-4 accent-pine" {...props} />;
}
