import type { ReactNode } from "react";

interface PageContainerProps {
  title?: string;
  actions?: ReactNode;
  children: ReactNode;
}

export default function PageContainer({ title, actions, children }: PageContainerProps) {
  return (
    <div className="p-8">
      {/* Page header */}
      {(title || actions) && (
        <div className="flex items-center justify-between mb-6">
          {title && <h1 className="text-xl font-semibold text-text-primary">{title}</h1>}
          {actions && <div className="flex items-center gap-2">{actions}</div>}
        </div>
      )}
      {children}
    </div>
  );
}
