"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { Users, CheckCircle2, XCircle, ArrowLeft, AlertCircle, Ban } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { getStoredUser, getToken } from "@/lib/auth";

interface ParticipantReq {
  id: string;
  user_id: string;
  user_name?: string;
  role: string;
  status: string;
  joined_at: string;
}

export default function ManageEventPage() {
  const params = useParams();
  const router = useRouter();
  const eventId = params?.id as string;

  const [requests, setRequests] = useState<ParticipantReq[]>([]);
  const [participants, setParticipants] = useState<ParticipantReq[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);

  const currentUser = getStoredUser();
  const token = getToken();

  const loadData = () => {
    if (!token || !eventId) return;
    setLoading(true);

    Promise.all([
      fetchApi<ParticipantReq[]>(`/events/${eventId}/requests`, {
        headers: { Authorization: `Bearer ${token}` },
      }),
      fetchApi<any>(`/events/${eventId}`),
    ])
      .then(([reqRes, eventRes]) => {
        if (reqRes.data) setRequests(reqRes.data);
        if (eventRes.data?.participants) {
          setParticipants(eventRes.data.participants);
        }
      })
      .catch((err) => {
        console.error(err);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    loadData();
  }, [eventId, token]);

  const handleApprove = async (userId: string) => {
    if (!token) return;
    setActionLoading(true);
    try {
      await fetchApi(`/events/${eventId}/requests/${userId}/approve`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      });
      loadData();
    } catch (err: any) {
      alert(err.message || "Gagal menyetujui permohonan.");
    } finally {
      setActionLoading(false);
    }
  };

  const handleReject = async (userId: string) => {
    if (!token) return;
    setActionLoading(true);
    try {
      await fetchApi(`/events/${eventId}/requests/${userId}/reject`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      });
      loadData();
    } catch (err: any) {
      alert(err.message || "Gagal menolak permohonan.");
    } finally {
      setActionLoading(false);
    }
  };

  const handleCancelEvent = async () => {
    if (!token || !confirm("Yakin ingin membatalkan event ini?")) return;
    try {
      await fetchApi(`/events/${eventId}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });
      router.push("/events");
    } catch (err: any) {
      alert(err.message || "Gagal membatalkan event.");
    }
  };

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div className="flex items-center justify-between gap-4 mb-6">
        <Link
          href={`/events/${eventId}`}
          className="inline-flex items-center gap-1.5 text-sm font-semibold text-brand-secondary hover:underline"
        >
          <ArrowLeft className="w-4 h-4" /> Kembali ke Event
        </Link>

        <button
          onClick={handleCancelEvent}
          className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-xl border border-red-200 bg-red-50 text-red-600 hover:bg-red-100 text-xs font-semibold transition"
        >
          <Ban className="w-3.5 h-3.5" />
          Batalkan Event
        </button>
      </div>

      <div className="space-y-8">
        {/* Pending Requests */}
        <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-8 shadow-sm">
          <h2 className="text-xl font-bold text-brand-primary mb-1">
            Permohonan Bergabung ({requests.length})
          </h2>
          <p className="text-xs text-brand-muted mb-6">
            Daftar angler yang meminta izin untuk mengikuti event private ini
          </p>

          {loading ? (
            <div className="text-center py-6 text-xs text-brand-muted">Memuat data...</div>
          ) : requests.length === 0 ? (
            <p className="text-center py-6 text-xs text-brand-muted">
              Tidak ada permohonan bergabung yang pending.
            </p>
          ) : (
            <div className="space-y-3">
              {requests.map((r) => (
                <div
                  key={r.id}
                  className="p-4 rounded-2xl bg-brand-background border border-brand-border flex items-center justify-between gap-4"
                >
                  <div>
                    <h4 className="text-xs font-bold text-brand-primary">
                      {r.user_name || "Pemancing"}
                    </h4>
                    <span className="text-[10px] text-brand-muted">
                      Meminta izin pada{" "}
                      {new Date(r.joined_at).toLocaleDateString("id-ID", {
                        day: "numeric",
                        month: "short",
                        year: "numeric",
                      })}
                    </span>
                  </div>

                  <div className="flex items-center gap-2">
                    <button
                      onClick={() => handleApprove(r.user_id)}
                      disabled={actionLoading}
                      className="inline-flex items-center gap-1 px-3 py-1.5 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold shadow-xs transition"
                    >
                      <CheckCircle2 className="w-3.5 h-3.5" /> Setujui
                    </button>
                    <button
                      onClick={() => handleReject(r.user_id)}
                      disabled={actionLoading}
                      className="inline-flex items-center gap-1 px-3 py-1.5 rounded-xl border border-red-200 bg-red-50 hover:bg-red-100 text-red-600 text-xs font-semibold transition"
                    >
                      <XCircle className="w-3.5 h-3.5" /> Tolak
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* All Active Participants */}
        <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-8 shadow-sm">
          <h2 className="text-xl font-bold text-brand-primary mb-1">
            Daftar Peserta Aktif
          </h2>
          <p className="text-xs text-brand-muted mb-6">
            Seluruh peserta yang telah disetujui dan bergabung dalam trip
          </p>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
            {participants
              .filter((p) => p.status === "approved")
              .map((p) => (
                <div
                  key={p.id}
                  className="p-3.5 rounded-2xl bg-brand-background border border-brand-border flex items-center justify-between"
                >
                  <span className="text-xs font-bold text-brand-primary">
                    {p.user_name || "Peserta"}
                  </span>
                  <span className="text-[10px] uppercase font-bold text-brand-secondary">
                    {p.role}
                  </span>
                </div>
              ))}
          </div>
        </div>
      </div>
    </div>
  );
}
