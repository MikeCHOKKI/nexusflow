import { type ReactNode } from "react";
import { motion } from "motion/react";
import { TrendUp, TrendDown } from "@phosphor-icons/react";

interface StatsCardProps {
  icon: ReactNode;
  label: string;
  value: string;
  trend?: number;
  index?: number;
  onClick?: () => void;
}

export default function StatsCard({ icon, label, value, trend, index = 0, onClick }: StatsCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.35, delay: index * 0.08, ease: "easeOut" }}
      className={`bg-surface-elevated border border-border rounded-xl p-5 ${onClick ? "cursor-pointer hover:bg-surface-hover" : ""} transition-colors`}
      onClick={onClick}
    >
      <div className="flex items-start justify-between mb-3">
        <div className="p-2 rounded-lg bg-accent/10 text-accent">
          {icon}
        </div>
        {trend !== undefined && trend !== 0 && (
          <div className={`flex items-center gap-1 text-xs font-medium ${trend > 0 ? "text-success" : "text-danger"}`}>
            {trend > 0 ? <TrendUp size={14} weight="fill" /> : <TrendDown size={14} weight="fill" />}
            <span>{Math.abs(trend)}%</span>
          </div>
        )}
      </div>
      <p className="text-2xl font-semibold text-text-primary tracking-tight mb-0.5">
        {value}
      </p>
      <p className="text-xs text-text-muted">{label}</p>
    </motion.div>
  );
}
