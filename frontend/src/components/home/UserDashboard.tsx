"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import {
  Compass,
  Calendar,
  MessageSquare,
  PlusCircle,
  Clock,
  MapPin,
  ArrowRight,
  Fish,
  Sparkles,
  Users,
} from "lucide-react";
import { fetchApi, resolveMediaUrl } from "@/lib/api-client";
import { getStoredUser } from "@/lib/auth";

interface RecentSpotItem {
  id: string;
  name: string;
  address: string;
  primary_photo_url?: string;
  viewed_at?: string;
}

interface EventItem {
  id: string;
  name: string;
  destination_name: string;
  destination_address: string;
  event_date: string;
  start_time: string;
  approved_participants?: number;
  max_participants?: number;
  status: string;
}

export default function UserDashboard() {
  const [currentUser, setCurrentUser] = useState<any>(null);
  const [recentSpots, setRecentSpots] = useState<RecentSpotItem[]>([]);
  const [myEvents, setMyEvents] = useState<EventItem[]>([]);
  const [loadingEvents, setLoadingEvents] = useState(true);

  useEffect(() => {
    const user = getStoredUser();
    setCurrentUser(user);

    if (user?.id) {
      // 1. Load user-isolated recently viewed spots from storage
      try {
        const key = `recent_spots_${user.id}`;
        const raw = localStorage.getItem(key);
        if (raw) {
          setRecentSpots(JSON.parse(raw));
        }
      } catch (e) {
        // ignore
      }

      // 2. Load user's latest upcoming events (either joined or created)
      setLoadingEvents(true);
      Promise.all([
        fetchApi<EventItem[]>(`/events?participant_id=${user.id}&limit=4`),
        fetchApi<EventItem[]>(`/events?creator_id=${user.id}&limit=4`),
      ])
        .then(([joinedRes, createdRes]) => {
          const joined = joinedRes.data || [];
          const created = createdRes.data || [];
          // Deduplicate by event id
          const map = new Map<string, EventItem>();
          [...joined, ...created].forEach((ev) => map.set(ev.id, ev));
          setMyEvents(Array.from(map.values()).slice(0, 4));
        })
        .catch((err) => {
          console.warn("Could not load user events:", err);
        })
        .finally(() => {
          setLoadingEvents(false);
        });
    }
  }, []);

  if (!currentUser) return null;

  return (
    <div className="flex flex-col gap-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full py-4 sm:py-6">
      {/* Welcome Banner */}
      <div className="bg-gradient-to-r from-brand-primary via-[#0f3d53] to-brand-secondary/80 rounded-3xl p-6 sm:p-8 text-white shadow-md relative overflow-hidden">
        <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div className="space-y-2 max-w-2xl">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/15 backdrop-blur-xs text-xs font-semibold text-brand-secondary">
              <Sparkles className="w-3.5 h-3.5" /> Ruang Angler
            </div>
            <h1 className="text-2xl sm:text-3xl font-extrabold tracking-tight">
              Halo, {currentUser.name?.split(" ")[0] || "Angler"}! Siap Mancing Hari Ini?
            </h1>
            <p className="text-white/80 text-sm leading-relaxed">
              Temukan spot potensial di sekitarmu, buat agenda janjian mancing bareng komunitas, atau obrolkan umpan gacor terbaru.
            </p>
          </div>

          <div className="flex flex-wrap sm:flex-nowrap items-center gap-2.5 shrink-0">
            <Link
              href="/spots/create"
              className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-brand-secondary text-white text-xs font-bold hover:bg-brand-secondary/90 transition shadow-sm"
            >
              <PlusCircle className="w-4 h-4" />
              <span>Bagi Spot Baru</span>
            </Link>
            <Link
              href="/events/create"
              className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-white/20 hover:bg-white/30 text-white text-xs font-bold backdrop-blur-xs transition"
            >
              <Calendar className="w-4 h-4 text-brand-secondary" />
              <span>Buat Janjian Mancing</span>
            </Link>
          </div>
        </div>
      </div>

      {/* Quick Navigation Cards */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <Link
          href="/spots"
          className="p-5 rounded-2xl bg-brand-surface border border-brand-border hover:border-brand-secondary/40 hover:shadow-xs transition flex flex-col gap-3 group"
        >
          <div className="w-10 h-10 rounded-xl bg-brand-primary/10 text-brand-secondary flex items-center justify-center group-hover:scale-105 transition">
            <Compass className="w-5 h-5" />
          </div>
          <div>
            <h3 className="text-sm font-bold text-brand-primary group-hover:text-brand-secondary transition">
              Jelajahi Spot
            </h3>
            <p className="text-[11px] text-brand-muted mt-0.5">Katalog spot di seluruh Indonesia</p>
          </div>
        </Link>

        <Link
          href="/events"
          className="p-5 rounded-2xl bg-brand-surface border border-brand-border hover:border-brand-secondary/40 hover:shadow-xs transition flex flex-col gap-3 group"
        >
          <div className="w-10 h-10 rounded-xl bg-emerald-500/10 text-emerald-600 flex items-center justify-center group-hover:scale-105 transition">
            <Calendar className="w-5 h-5" />
          </div>
          <div>
            <h3 className="text-sm font-bold text-brand-primary group-hover:text-emerald-600 transition">
              Janjian Mancing
            </h3>
            <p className="text-[11px] text-brand-muted mt-0.5">Ikuti trip bareng angler lain</p>
          </div>
        </Link>

        <Link
          href="/chat"
          className="p-5 rounded-2xl bg-brand-surface border border-brand-border hover:border-brand-secondary/40 hover:shadow-xs transition flex flex-col gap-3 group"
        >
          <div className="w-10 h-10 rounded-xl bg-sky-500/10 text-sky-600 flex items-center justify-center group-hover:scale-105 transition">
            <MessageSquare className="w-5 h-5" />
          </div>
          <div>
            <h3 className="text-sm font-bold text-brand-primary group-hover:text-sky-600 transition">
              Ruang Chat
            </h3>
            <p className="text-[11px] text-brand-muted mt-0.5">Diskusi personal dan grup event</p>
          </div>
        </Link>

        <Link
          href={`/profile/${currentUser.username}`}
          className="p-5 rounded-2xl bg-brand-surface border border-brand-border hover:border-brand-secondary/40 hover:shadow-xs transition flex flex-col gap-3 group"
        >
          <div className="w-10 h-10 rounded-xl bg-amber-500/10 text-amber-600 flex items-center justify-center group-hover:scale-105 transition">
            <Users className="w-5 h-5" />
          </div>
          <div>
            <h3 className="text-sm font-bold text-brand-primary group-hover:text-amber-600 transition">
              Aktivitas Profil
            </h3>
            <p className="text-[11px] text-brand-muted mt-0.5">Spot & event yang kamu miliki</p>
          </div>
        </Link>
      </div>

      {/* Grid: Spot Terakhir Dilihat & Event Saya */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Section 1: Spot Terakhir Dilihat */}
        <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 shadow-xs flex flex-col">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-base font-bold text-brand-primary flex items-center gap-2">
              <Clock className="w-4 h-4 text-brand-secondary" />
              <span>Spot Terakhir Dilihat</span>
            </h2>
            <Link href="/spots" className="text-xs font-semibold text-brand-secondary hover:underline">
              Semua Spot &rarr;
            </Link>
          </div>

          {recentSpots.length === 0 ? (
            <div className="flex-1 flex flex-col items-center justify-center text-center p-8 bg-brand-background rounded-2xl border border-dashed border-brand-border">
              <div className="w-12 h-12 rounded-full bg-brand-primary/5 flex items-center justify-center text-brand-muted mb-3">
                <Compass className="w-6 h-6 text-brand-muted/70" />
              </div>
              <p className="text-sm font-bold text-brand-primary">Belum Ada Riwayat Spot</p>
              <p className="text-xs text-brand-muted mt-1 max-w-xs">
                Spot mancing yang kamu buka akan tersimpan di sini agar mudah kamu temukan kembali.
              </p>
              <Link
                href="/spots"
                className="mt-4 px-4 py-2 rounded-xl bg-brand-primary text-white text-xs font-semibold hover:bg-brand-primary/90 transition shadow-xs"
              >
                Mulai Jelajahi Spot
              </Link>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {recentSpots.map((spot) => (
                <Link
                  key={spot.id}
                  href={`/spots/${spot.id}`}
                  className="p-3 rounded-2xl bg-brand-background border border-brand-border hover:border-brand-secondary/40 transition flex gap-3 group"
                >
                  <div className="w-16 h-16 rounded-xl bg-brand-primary/10 overflow-hidden shrink-0 flex items-center justify-center text-brand-muted">
                    {spot.primary_photo_url ? (
                      <img
                        src={resolveMediaUrl(spot.primary_photo_url)}
                        alt={spot.name}
                        className="w-full h-full object-cover group-hover:scale-105 transition"
                      />
                    ) : (
                      <Fish className="w-6 h-6 text-brand-secondary" />
                    )}
                  </div>
                  <div className="flex-1 min-w-0 flex flex-col justify-center">
                    <h4 className="text-xs font-bold text-brand-primary group-hover:text-brand-secondary transition truncate">
                      {spot.name}
                    </h4>
                    <p className="text-[11px] text-brand-muted truncate mt-0.5 flex items-center gap-1">
                      <MapPin className="w-3 h-3 text-brand-secondary shrink-0" />
                      <span>{spot.address}</span>
                    </p>
                    <span className="text-[10px] text-brand-secondary font-semibold mt-1 inline-flex items-center gap-0.5">
                      Lihat Spot <ArrowRight className="w-2.5 h-2.5" />
                    </span>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>

        {/* Section 2: Agenda Janjian Mancing Saya */}
        <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 shadow-xs flex flex-col">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-base font-bold text-brand-primary flex items-center gap-2">
              <Calendar className="w-4 h-4 text-brand-secondary" />
              <span>Agenda Janjian Mancing Saya</span>
            </h2>
            <Link href="/events" className="text-xs font-semibold text-brand-secondary hover:underline">
              Cari Agenda &rarr;
            </Link>
          </div>

          {loadingEvents ? (
            <div className="flex-1 flex items-center justify-center p-8 text-xs text-brand-muted">
              Memuat agenda mancing...
            </div>
          ) : myEvents.length === 0 ? (
            <div className="flex-1 flex flex-col items-center justify-center text-center p-8 bg-brand-background rounded-2xl border border-dashed border-brand-border">
              <div className="w-12 h-12 rounded-full bg-emerald-500/10 flex items-center justify-center text-emerald-600 mb-3">
                <Calendar className="w-6 h-6" />
              </div>
              <p className="text-sm font-bold text-brand-primary">Belum Mengikuti Janjian</p>
              <p className="text-xs text-brand-muted mt-1 max-w-xs">
                Ayo cari rekan mancing atau buat trip bareng angler se-wilayahmu sekarang!
              </p>
              <div className="mt-4 flex items-center gap-2">
                <Link
                  href="/events"
                  className="px-3.5 py-1.5 rounded-xl border border-brand-border bg-brand-surface text-brand-primary text-xs font-semibold hover:border-brand-secondary transition"
                >
                  Cari Trip
                </Link>
                <Link
                  href="/events/create"
                  className="px-3.5 py-1.5 rounded-xl bg-brand-secondary text-white text-xs font-bold hover:bg-brand-secondary/90 transition shadow-xs"
                >
                  Buat Trip Baru
                </Link>
              </div>
            </div>
          ) : (
            <div className="space-y-3">
              {myEvents.map((ev) => (
                <Link
                  key={ev.id}
                  href={`/events/${ev.id}`}
                  className="p-3.5 rounded-2xl bg-brand-background border border-brand-border hover:border-brand-secondary/40 transition flex items-center justify-between gap-3 group"
                >
                  <div className="min-w-0">
                    <span className="text-[10px] font-bold text-brand-secondary uppercase tracking-wider">
                      {ev.event_date} • {ev.start_time}
                    </span>
                    <h4 className="text-xs font-bold text-brand-primary group-hover:text-brand-secondary transition truncate mt-0.5">
                      {ev.name}
                    </h4>
                    <p className="text-[11px] text-brand-muted truncate flex items-center gap-1 mt-0.5">
                      <MapPin className="w-3 h-3 text-brand-secondary shrink-0" />
                      <span>{ev.destination_name}</span>
                    </p>
                  </div>

                  <div className="text-right shrink-0">
                    <span className="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-brand-secondary/15 text-brand-primary">
                      {ev.approved_participants ?? 1} Peserta
                    </span>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
