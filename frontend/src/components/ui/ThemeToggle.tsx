import { SunDim, MoonStars } from "@phosphor-icons/react";
import { useTheme } from "../../hooks/useTheme";

export default function ThemeToggle() {
  const { theme, toggle } = useTheme();

  return (
    <button
      onClick={toggle}
      className="p-2 rounded-lg text-text-muted hover:text-text-primary hover:bg-surface-hover transition-colors"
      aria-label="Basculer le thème"
    >
      {theme === "dark" ? <SunDim size={18} /> : <MoonStars size={18} />}
    </button>
  );
}
