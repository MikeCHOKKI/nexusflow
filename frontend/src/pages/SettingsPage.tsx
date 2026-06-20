import { useState } from "react";
import PageContainer from "../components/PageContainer";

export default function SettingsPage() {
  const [appName] = useState("NexusFlow");
  const [notifications, setNotifications] = useState(true);

  return (
    <PageContainer title="Paramètres">
      <div className="max-w-2xl space-y-6">
        {/* General */}
        <div className="bg-surface-elevated border border-border rounded-xl p-6">
          <h3 className="text-sm font-semibold text-text-primary mb-4">Général</h3>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-text-primary">Nom de l'application</p>
                <p className="text-xs text-text-muted">{appName} v1.0.0</p>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-text-primary">Notifications</p>
                <p className="text-xs text-text-muted">Recevoir les alertes système</p>
              </div>
              <button
                onClick={() => setNotifications(!notifications)}
                className={`relative w-10 h-5 rounded-full transition-colors ${
                  notifications ? "bg-accent" : "bg-surface-hover"
                }`}
              >
                <span
                  className={`absolute top-0.5 left-0.5 w-4 h-4 rounded-full bg-white transition-transform ${
                    notifications ? "translate-x-5" : "translate-x-0"
                  }`}
                />
              </button>
            </div>
          </div>
        </div>

        {/* Theme info */}
        <div className="bg-surface-elevated border border-border rounded-xl p-6">
          <h3 className="text-sm font-semibold text-text-primary mb-4">Apparence</h3>
          <p className="text-sm text-text-secondary">
            Utilise le bouton <span className="text-accent">☀️ / 🌙</span> dans l'en-tête pour basculer entre les
            thèmes clair et sombre.
          </p>
        </div>

        {/* API */}
        <div className="bg-surface-elevated border border-border rounded-xl p-6">
          <h3 className="text-sm font-semibold text-text-primary mb-4">API</h3>
          <div className="space-y-2">
            <div className="flex items-center justify-between py-2">
              <span className="text-sm text-text-primary">API Gateway</span>
              <code className="text-xs font-mono text-accent bg-accent/10 px-2 py-1 rounded">
                http://localhost:8080
              </code>
            </div>
            <div className="flex items-center justify-between py-2">
              <span className="text-sm text-text-primary">Version API</span>
              <code className="text-xs font-mono text-text-muted">v1</code>
            </div>
          </div>
        </div>
      </div>
    </PageContainer>
  );
}
