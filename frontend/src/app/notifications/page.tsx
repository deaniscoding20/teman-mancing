"use client";

import { useEffect, useState } from "react";
import { Bell, CheckCheck, CheckCircle2, MessageSquare, Users, AlertCircle } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { getStoredUser, getToken } from "@/lib/auth";

interface NotificationItem {
  id: string;
  type: string;
  title: string;
  message: string;
  reference_id?: string;
  is_read: boolean;
  created_at: string;
}

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<NotificationItem[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(true);

  const currentUser = getStoredUser();
  const token = getToken();

  const loadNotifications = () => {
    if (!token) return;
    setLoading(true);

    fetchApi<NotificationItem[]>("/notifications", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => {
        if (res.data) setNotifications(res.data);
        if (res.meta && (res.meta as any).unread_count !== undefined) {
          setUnreadCount((res.meta as any).unread_count);
        }
      })
      .catch((err) => {
        console.error("Failed to load notifications:", err);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    loadNotifications();
  }, [token]);

  const handleMarkAsRead = async (id: string) => {
    if (!token) return;

    try {
      await fetchApi(`/notifications/${id}/read`, {
        method: "PATCH",
        headers: { Authorization: `Bearer ${token}` },
      });
      setNotifications((prev) =>
        prev.map((n) => (n.id === id ? { ...n, is_read: true } : n))
      );
      setUnreadCount((prev) => Math.max(0, prev - 1));
    } catch (err: any) {
      console.error(err);
    }
  };

  const handleMarkAllAsRead = async () => {
    if (!token) return;

    try {
      await fetchApi("/notifications/read-all", {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      });
      setNotifications((prev) => prev.map((n) => ({ ...n, is_read: true })));
      setUnreadCount(0);
    } catch (err: any) {
      alert("Gagal menandai semua notifikasi.");
    }
  };

  if (!currentUser) {
    return (
      <div className="max-w-md mx-auto px-4 py-20 text-center">
        <h2 className="text-xl font-bold text-brand-primary mb-2">Masuk untuk Melihat Notifikasi</h2>
        <p className="text-sm text-brand-muted">
          Pemberitahuan seputar janjian mancing dan komentar akan muncul di sini.
        </p>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div className="flex items-center justify-between gap-4 mb-8">
        <div>
          <h1 className="text-2xl sm:text-3xl font-extrabold text-brand-primary flex items-center gap-2.5">
            <Bell className="w-7 h-7 text-brand-secondary" />
            Notifikasi Saya
          </h1>
          <p className="text-xs text-brand-muted mt-1">
            {unreadCount > 0
              ? `Terdapat ${unreadCount} pemberitahuan baru yang belum dibaca`
              : "Semua pemberitahuan sudah terbaca"}
          </p>
        </div>

        {unreadCount > 0 && (
          <button
            onClick={handleMarkAllAsRead}
            className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-xl border border-brand-border bg-brand-surface text-brand-primary text-xs font-semibold hover:border-brand-secondary transition"
          >
            <CheckCheck className="w-4 h-4 text-brand-secondary" />
            Tandai Semua Terbaca
          </button>
        )}
      </div>

      {loading ? (
        <div className="text-center py-20 text-brand-muted">Memuat pemberitahuan...</div>
      ) : notifications.length === 0 ? (
        <div className="bg-brand-surface rounded-3xl border border-brand-border p-12 text-center">
          <div className="w-12 h-12 rounded-2xl bg-brand-background flex items-center justify-center mx-auto mb-3 text-brand-muted">
            <Bell className="w-6 h-6" />
          </div>
          <h3 className="text-base font-bold text-brand-primary mb-1">Belum Ada Notifikasi</h3>
          <p className="text-xs text-brand-muted">
            Aktivitas seputar event dan spot mancing Anda akan muncul di halaman ini.
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {notifications.map((n) => (
            <div
              key={n.id}
              className={`p-4 rounded-2xl border transition flex items-start justify-between gap-4 ${
                n.is_read
                  ? "bg-brand-surface border-brand-border/70 opacity-80"
                  : "bg-white border-brand-secondary/40 shadow-xs"
              }`}
            >
              <div className="flex items-start gap-3.5">
                <div
                  className={`w-9 h-9 rounded-xl flex items-center justify-center shrink-0 mt-0.5 ${
                    n.is_read
                      ? "bg-brand-background text-brand-muted"
                      : "bg-brand-secondary/15 text-brand-secondary"
                  }`}
                >
                  {n.type.includes("comment") ? (
                    <MessageSquare className="w-4 h-4" />
                  ) : n.type.includes("join") ? (
                    <Users className="w-4 h-4" />
                  ) : (
                    <Bell className="w-4 h-4" />
                  )}
                </div>

                <div>
                  <h4 className="text-xs font-bold text-brand-primary mb-0.5">{n.title}</h4>
                  <p className="text-xs text-brand-text leading-relaxed">{n.message}</p>
                  <span className="text-[10px] text-brand-muted mt-1 inline-block">
                    {new Date(n.created_at).toLocaleDateString("id-ID", {
                      day: "numeric",
                      month: "short",
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </span>
                </div>
              </div>

              {!n.is_read && (
                <button
                  onClick={() => handleMarkAsRead(n.id)}
                  className="p-1 text-brand-muted hover:text-brand-secondary text-xs"
                  title="Tandai telah dibaca"
                >
                  <CheckCircle2 className="w-4 h-4" />
                </button>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
