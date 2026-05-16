export function FormField({
  label,
  hint,
  children,
  className = "",
}: {
  label: string;
  hint?: string;
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <label className={`space-y-2 ${className}`}>
      <span className="text-sm font-medium">{label}</span>
      {hint ? <p className="text-xs text-black/55 dark:text-white/55">{hint}</p> : null}
      {children}
    </label>
  );
}
