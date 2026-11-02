"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import {
  Calendar,
  MapPin,
  Clock,
  Users,
  Shield,
  DollarSign,
  Ship,
  Car,
  CheckCircle2,
  AlertCircle,
  MessageSquare,
  LogOut,
  Settings,
} from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { getStoredUser, getToken } from "@/lib/auth";

interface EventDetail {
  id: string;
  creator_id: string;
  creator_name?: string;
  creator_avatar?: string;
  name: string;
  description: string;
  destination_name: string;
  destination_address: string;
  latitude: number;
  longitude: number;
  event_date: string;
  start_time: string;
  end_time?: string;
  meeting_point: string;
  difficulty: string;
  has_parking_fee: boolean;
  parking_fee?: number;
  requires_boat: boolean;
  boat_rental_fee?: number;
  additional_cost?: string;
  visibility: "public" | "private";
  max_participants?: number;
  approved_participants: number;
  status: string;
  created_at: string;
}

interface ParticipantItem {
  id: string;
  event_id: string;
  user_id: string;
  user_name?: string;
  role: string;
  status: string;
  joined_at: string;
}

export default function EventDetailPage() {
  const params = useParams();
  const router = useRouter();
  const eventId = params?.id as string;

  const [event, setEvent] = useState<EventDetail | null>(null);
  const [participants, setParticipants] = useState<ParticipantItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);
  const [actionMessage, setActionMessage] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);

  const currentUser = getStoredUser();
  const token = getToken();

  const loadEvent = () => {
    if (!eventId) return;
    setLoading(true);

    fetchApi<{ event: EventDetail; participants: ParticipantItem[] }>(`/events/${eventId}`)
      .then((res) => {
        if (res.data) {
          setEvent(res.data.event);
          setParticipants(res.data.participants || []);
        }
      })
      .catch((err) => {
        console.error("Failed to load event:", err);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    loadEvent();
  }, [eventId]);

  const handleJoin = async () => {
    if (!token) {
      router.push("/login");
      return;
    }

    setActionLoading(true);
    setActionError(null);
    setActionMessage(null);

    try {
      const res = await fetchApi(`/events/${eventId}/join`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setActionMessage(res.message || "Berhasil bergabung!");
      loadEvent();
    } catch (err: any) {
      setActionError(err?.data?.message || err?.message || "Gagal bergabung ke event.");
    } finally {
      setActionLoading(false);
    }
  };

  const handleLeave = async () => {
    if (!token || !confirm("Yakin ingin membatalkan keikutsertaan dalam event ini?")) return;

    setActionLoading(true);
    try {
      await fetchApi(`/events/${eventId}/leave`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      loadEvent();
    } catch (err: any) {
      alert(err.message || "Gagal meninggalkan event.");
    } finally {
      setActionLoading(false);
    }
  };

  if (loading) {
    return <div className="text-center py-20 text-brand-muted">Memuat agenda event...</div>;
  }

  if (!event) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-20 text-center">
        <h2 className="text-2xl font-bold text-brand-primary mb-2">Event Tidak Ditemukan</h2>
        <Link href="/events" className="text-sm text-brand-secondary hover:underline">
          &larr; Kembali ke daftar event
        </Link>
      </div>
    );
  }

  const isCreator = currentUser && currentUser.id === event.creator_id;
  const myParticipation = currentUser
    ? participants.find((p) => p.user_id === currentUser.id)
    : null;

  return (
    <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Top breadcrumb */}
      <Link href="/events" className="text-sm font-semibold text-brand-secondary hover:underline mb-6 inline-block">
        &larr; Kembali ke Semua Janjian Mancing
      </Link>

      {/* Action Notification */}
      {actionMessage && (
        <div className="mb-6 p-4 rounded-2xl bg-green-50 border border-green-200 flex items-center gap-3 text-green-800 text-sm">
          <CheckCircle2 className="w-5 h-5 text-green-600 shrink-0" />
          <span>{actionMessage}</span>
        </div>
      )}
      {actionError && (
        <div className="mb-6 p-4 rounded-2xl bg-red-50 border border-red-200 flex items-center gap-3 text-red-700 text-sm">
          <AlertCircle className="w-5 h-5 text-red-600 shrink-0" />
          <span>{actionError}</span>
        </div>
      )}

      {/* Main Event Card */}
      <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-10 shadow-sm mb-8">
        <div className="flex flex-col gap-4">
          <div className="flex flex-wrap items-center gap-2">
            <span
              className={`px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider ${
                event.visibility === "private"
                  ? "bg-amber-100 text-amber-800"
                  : "bg-emerald-100 text-emerald-800"
              }`}
            >
              {event.visibility === "private" ? "Private Event" : "Public Event"}
            </span>

            <span className="px-3 py-1 rounded-full bg-brand-background text-brand-muted text-xs font-semibold flex items-center gap-1.5 border border-brand-border">
              <Shield className="w-3.5 h-3.5 text-brand-secondary" />
              Tingkat: {event.difficulty}
            </span>
          </div>

          <h1 className="text-3xl sm:text-4xl font-extrabold text-brand-primary">{event.name}</h1>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-2 text-sm text-brand-muted">
            <div className="flex items-center gap-2">
              <Calendar className="w-4 h-4 text-brand-secondary shrink-0" />
              <span>
                {event.event_date} &bull; Pukul {event.start_time} {event.end_time ? `- ${event.end_time}` : ""}
              </span>
            </div>

            <div className="flex items-center gap-2">
              <MapPin className="w-4 h-4 text-brand-secondary shrink-0" />
              <span>{event.destination_name}</span>
            </div>
          </div>

          <p className="mt-4 text-base text-brand-text leading-relaxed whitespace-pre-line">
            {event.description}
          </p>
        </div>

        {/* Meeting point */}
        <div className="mt-8 p-5 rounded-2xl bg-brand-background border border-brand-border">
          <h3 className="text-xs font-bold uppercase tracking-wider text-brand-primary mb-2 flex items-center gap-1.5">
            <MapPin className="w-4 h-4 text-brand-secondary" /> Titik Kumpul (Meeting Point)
          </h3>
          <p className="text-sm font-semibold text-brand-text">{event.meeting_point}</p>
          <p className="text-xs text-brand-muted mt-1">{event.destination_address}</p>
        </div>

        {/* Cost Breakdown (README Section 8.2 & 8.3) */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mt-6">
          <div className="p-4 rounded-2xl bg-brand-background border border-brand-border">
            <span className="text-xs text-brand-muted flex items-center gap-1.5 mb-1">
              <Car className="w-4 h-4 text-brand-secondary" /> Biaya Parkir
            </span>
            <span className="text-sm font-bold text-brand-primary">
              {event.has_parking_fee && event.parking_fee
                ? `Rp ${event.parking_fee.toLocaleString("id-ID")}`
                : "Gratis / Tidak ada"}
            </span>
          </div>

          <div className="p-4 rounded-2xl bg-brand-background border border-brand-border">
            <span className="text-xs text-brand-muted flex items-center gap-1.5 mb-1">
              <Ship className="w-4 h-4 text-brand-secondary" /> Sewa Kapal/Perahu
            </span>
            <span className="text-sm font-bold text-brand-primary">
              {event.requires_boat && event.boat_rental_fee
                ? `Rp ${event.boat_rental_fee.toLocaleString("id-ID")}`
                : "Tidak butuh perahu"}
            </span>
          </div>

          <div className="p-4 rounded-2xl bg-brand-background border border-brand-border">
            <span className="text-xs text-brand-muted flex items-center gap-1.5 mb-1">
              <DollarSign className="w-4 h-4 text-brand-accent" /> Biaya Lainnya
            </span>
            <span className="text-xs font-semibold text-brand-text">
              {event.additional_cost || "Tidak ada biaya ekstra"}
            </span>
          </div>
        </div>

        {/* Join Actions */}
        <div className="mt-8 pt-8 border-t border-brand-border flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2 text-sm text-brand-primary font-semibold">
            <Users className="w-5 h-5 text-brand-secondary" />
            <span>
              {event.approved_participants}{" "}
              {event.max_participants ? `/ ${event.max_participants}` : ""} Pemancing Bergabung
            </span>
          </div>

          <div className="flex items-center gap-3 w-full sm:w-auto">
            {isCreator ? (
              <span className="text-xs font-bold text-brand-secondary px-4 py-2 rounded-xl bg-brand-background border border-brand-border">
                Anda adalah Pembuat Event ini
              </span>
            ) : myParticipation?.status === "approved" ? (
              <div className="flex items-center gap-2">
                <span className="text-xs font-bold text-emerald-700 bg-emerald-50 px-3 py-1.5 rounded-lg border border-emerald-200 flex items-center gap-1">
                  <CheckCircle2 className="w-3.5 h-3.5" /> Terdaftar
                </span>
                <Link
                  href="/chat"
                  className="px-4 py-2 rounded-xl text-xs font-bold text-white bg-brand-primary hover:bg-brand-primary/90 flex items-center gap-1.5 shadow-sm"
                >
                  <MessageSquare className="w-3.5 h-3.5 text-brand-secondary" />
                  Grup Chat
                </Link>
                <button
                  onClick={handleLeave}
                  disabled={actionLoading}
                  className="p-2 text-brand-muted hover:text-red-500 rounded-xl border border-brand-border hover:bg-red-50 transition"
                  title="Batalkan keikutsertaan"
                >
                  <LogOut className="w-4 h-4" />
                </button>
              </div>
            ) : myParticipation?.status === "pending" ? (
              <span className="text-xs font-bold text-amber-800 bg-amber-50 px-4 py-2 rounded-xl border border-amber-200">
                Menunggu Persetujuan Creator
              </span>
            ) : (
              <button
                onClick={handleJoin}
                disabled={actionLoading}
                className="w-full sm:w-auto px-8 py-3 rounded-xl font-bold text-white bg-brand-secondary hover:bg-brand-secondary/90 shadow-md transition flex items-center justify-center gap-2"
              >
                <Users className="w-4 h-4" />
                {event.visibility === "private" ? "Minta Gabung (Request Join)" : "Gabung Janjian Sekarang"}
              </button>
            )}
          </div>
        </div>
      </div>

      {/* Participants Card */}
      <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-10 shadow-sm">
        <h2 className="text-xl font-bold text-brand-primary flex items-center gap-2.5 mb-6">
          <Users className="w-5 h-5 text-brand-secondary" />
          Daftar Peserta Terdaftar ({participants.filter((p) => p.status === "approved").length})
        </h2>

        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
          {participants
            .filter((p) => p.status === "approved")
            .map((p) => (
              <div
                key={p.id}
                className="p-3.5 rounded-2xl bg-brand-background border border-brand-border flex items-center justify-between"
              >
                <span className="text-xs font-bold text-brand-primary">
                  {p.user_name || "Pemancing"}
                </span>
                <span className="text-[10px] uppercase font-bold text-brand-secondary">
                  {p.role}
                </span>
              </div>
            ))}
        </div>
      </div>
    </div>
  );
}
