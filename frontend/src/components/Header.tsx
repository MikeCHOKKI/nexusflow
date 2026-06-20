import { Bell } from "@phosphor-icons/react";
import SearchInput from "./ui/SearchInput";
import ThemeToggle from "./ui/ThemeToggle";

interface HeaderProps {
  title: string;
  search?: string;
  onSearch?: (v: string) => void;
}

export default function Header({ title, search, onSearch }: HeaderProps) {
  return (
    <header className="sticky top-0 z-30 h-16 bg-surface/80 backdrop-blur-lg border-b border-border flex items-center gap-4 px-8">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-sm min-w-0">
        <span className="text-text-muted">NexusFlow</span>
        <span className="text-text-muted">/</span>
        <span className="text-text-primary font-medium truncate">{title}</span>
      </div>

      {/* Spacer */}
      <div className="flex-1" />

      {/* Search */}
      {onSearch && (
        <div className="w-64">
          <SearchInput value={search ?? ""} onChange={onSearch} placeholder="Rechercher..." />
        </div>
      )}

      {/* Actions */}
      <div className="flex items-center gap-1">
        <ThemeToggle />

        <button className="relative p-2 rounded-lg text-text-muted hover:text-text-primary hover:bg-surface-hover transition-colors">
          <Bell size={18} />
          <span className="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-danger animate-[pulse-dot_2s_ease-in-out_infinite]" />
        </button>

        <div className="ml-2 pl-2 border-l border-border flex items-center gap-2">
          <div className="w-8 h-8 rounded-full bg-gradient-to-br from-accent to-accent-secondary flex items-center justify-center text-white text-xs font-semibold">
            AD
          </div>
          <div className="hidden sm:block">
            <p className="text-xs font-medium text-text-primary">Admin</p>
            <p className="text-[10px] text-text-muted">admin@nexusflow.io</p>
          </div>
        </div>
      </div>
    </header>
  );
}
