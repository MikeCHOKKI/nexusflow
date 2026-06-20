import type { OrderStatus, PaymentMethod, PaymentStatus, ProductStatus, UserRole, UserStatus } from "../../lib/types";

type BadgeColor = "cyan" | "green" | "yellow" | "red" | "blue" | "purple" | "slate" | "indigo" | "orange";

interface StatusBadgeProps {
  status: string;
  color?: BadgeColor;
}

const colorMap: Record<BadgeColor, { bg: string; text: string; dot: string }> = {
  green:  { bg: "bg-emerald-500/10", text: "text-emerald-400", dot: "bg-emerald-400" },
  yellow: { bg: "bg-amber-500/10",   text: "text-amber-400",  dot: "bg-amber-400" },
  red:    { bg: "bg-red-500/10",     text: "text-red-400",    dot: "bg-red-400" },
  blue:   { bg: "bg-blue-500/10",    text: "text-blue-400",   dot: "bg-blue-400" },
  purple: { bg: "bg-purple-500/10",  text: "text-purple-400", dot: "bg-purple-400" },
  cyan:   { bg: "bg-cyan-500/10",    text: "text-cyan-400",   dot: "bg-cyan-400" },
  slate:  { bg: "bg-slate-500/10",   text: "text-slate-400",  dot: "bg-slate-400" },
  indigo: { bg: "bg-indigo-500/10",  text: "text-indigo-400", dot: "bg-indigo-400" },
  orange: { bg: "bg-orange-500/10",  text: "text-orange-400", dot: "bg-orange-400" },
};

const orderStatusColor: Record<OrderStatus, BadgeColor> = {
  pending:    "yellow",
  confirmed:  "blue",
  processing: "indigo",
  shipped:    "purple",
  delivered:  "green",
  cancelled:  "red",
};

const paymentStatusColor: Record<PaymentStatus, BadgeColor> = {
  pending:   "yellow",
  completed: "green",
  failed:    "red",
  refunded:  "orange",
};

const productStatusColor: Record<ProductStatus, BadgeColor> = {
  active:   "green",
  inactive: "slate",
  draft:    "yellow",
};

const userRoleColor: Record<UserRole, BadgeColor> = {
  admin:   "purple",
  manager: "blue",
  viewer:  "slate",
};

const userStatusColor: Record<UserStatus, BadgeColor> = {
  active:   "green",
  inactive: "red",
};

const paymentMethodColor: Record<PaymentMethod, BadgeColor> = {
  wave:         "cyan",
  orange_money: "orange",
  mtn_money:    "yellow",
  visa:         "blue",
};

const labelMap: Record<string, string> = {
  pending:       "En attente",
  confirmed:     "Confirmé",
  processing:    "En cours",
  shipped:       "Expédié",
  delivered:     "Livré",
  cancelled:     "Annulé",
  completed:     "Complété",
  failed:        "Échoué",
  refunded:      "Remboursé",
  active:        "Actif",
  inactive:      "Inactif",
  draft:         "Brouillon",
  admin:         "Admin",
  manager:       "Manager",
  viewer:        "Lecteur",
  wave:          "Wave",
  orange_money:  "Orange Money",
  mtn_money:     "MTN Money",
  visa:          "Visa",
};

function getColor(status: string, color?: BadgeColor): BadgeColor {
  if (color) return color;
  return (orderStatusColor as Record<string, BadgeColor>)[status]
    ?? (paymentStatusColor as Record<string, BadgeColor>)[status]
    ?? (productStatusColor as Record<string, BadgeColor>)[status]
    ?? (userRoleColor as Record<string, BadgeColor>)[status]
    ?? (userStatusColor as Record<string, BadgeColor>)[status]
    ?? (paymentMethodColor as Record<string, BadgeColor>)[status]
    ?? "slate";
}

export default function StatusBadge({ status, color }: StatusBadgeProps) {
  const resolved = getColor(status, color);
  const { bg, text, dot } = colorMap[resolved];
  const label = labelMap[status] ?? status;

  return (
    <span className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium ${bg} ${text}`}>
      <span className={`w-1.5 h-1.5 rounded-full ${dot}`} />
      {label}
    </span>
  );
}
