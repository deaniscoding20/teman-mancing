"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Calendar, Plus, Search, MapPin, Users, Shield, Clock } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { getStoredUser } from "@/lib/auth";
import { FishingEventItem } from "@/types/api";

export default function EventsPage() {
  const [events, setEvents] = useState<FishingEventItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const currentUser = getStoredUser();

  const loadEvents = (q: string = "") => {
    setLoading(true);
    const queryParam = q ? `?q=${encodeURIComponent(q)}` : "";
    fetchApi<FishingEventItem[]>(`/events${queryParam}`)
      .then((res) => {
        if (res.data) setEvents(res.data);
      })
      .catch(() => {
        setEvents([]);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    loadEvents();
  }, []);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    loadEvents(search);
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 className="text-3xl font-extrabold text-brand-primary flex items-center gap-2.5">
            <Calendar className="w-8 h-8 text-brand-secondary" />
            Janjian Mancing
          </h1>
          <p className="text-sm text-brand-muted mt-1">
            Temukan rekan memancing, ikuti trip seru bersama, dan nikmati petualangan baru
          </p>
        </div>

        {currentUser && (
          <Link
            href="/events/create"
            className="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl font-bold text-white bg-brand-secondary hover:bg-brand-secondary/90 shadow-sm transition self-start md:self-auto"
          >
            <Plus className="w-4 h-4" />
            Buat Janjian Mancing
          </Link>
        )}
      </div>

      {/* Search Bar */}
      <form onSubmit={handleSearch} className="mb-8">
        <div className="relative max-w-xl">
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari event berdasarkan nama kegiatan, destinasi, atau titik kumpul..."
            className="w-full pl-11 pr-4 py-3 rounded-2xl border border-brand-border bg-brand-surface text-brand-text placeholder-brand-muted/70 text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary shadow-sm transition"
          />
          <Search className="w-5 h-5 text-brand-muted absolute left-4 top-1/2 -translate-y-1/2" />
        </div>
      </form>

      {/* Events Listing */}
      {loading ? (
        <div className="text-center py-20 text-brand-muted">Memuat agenda janjian...</div>
      ) : events.length === 0 ? (
        <div className="bg-brand-surface rounded-3xl border border-brand-border p-12 text-center max-w-lg mx-auto">
          <div className="w-14 h-14 rounded-2xl bg-brand-background flex items-center justify-center mx-auto mb-4 text-brand-muted">
            <Calendar className="w-7 h-7" />
          </div>
          <h3 className="text-lg font-bold text-brand-primary mb-1">Belum Ada Janjian Mancing</h3>
          <p className="text-sm text-brand-muted mb-6">
            Rencanakan agenda trip mancing pertamamu dan ajak pemancing lain bergabung.
          </p>
          <Link
            href={currentUser ? "/events/create" : "/login"}
            className="inline-flex items-center gap-2 px-6 py-2.5 rounded-xl font-semibold text-white bg-brand-primary hover:bg-brand-primary/90 text-sm transition"
          >
            <Plus className="w-4 h-4 text-brand-secondary" />
            Buat Agenda Pertama
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {events.map((ev) => (
            <Link
              key={ev.id}
              href={`/events/${ev.id}`}
              className="bg-brand-surface rounded-2xl border border-brand-border overflow-hidden hover:border-brand-secondary/40 shadow-sm transition flex flex-col group"
            >
              <div className="p-6 flex-1 flex flex-col">
                <div className="flex items-center justify-between gap-2 mb-3">
                  <span
                    className={`px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${
                      ev.visibility === "private"
                        ? "bg-amber-100 text-amber-800"
                        : "bg-emerald-100 text-emerald-800"
                    }`}
                  >
                    {ev.visibility === "private" ? "Private Event" : "Public Event"}
                  </span>
                  <span className="text-xs text-brand-muted flex items-center gap-1">
                    <Shield className="w-3.5 h-3.5 text-brand-secondary" />
                    {ev.difficulty}
                  </span>
                </div>

                <h3 className="text-lg font-bold text-brand-primary group-hover:text-brand-secondary transition line-clamp-1 mb-2">
                  {ev.name}
                </h3>

                <p className="text-xs text-brand-muted flex items-center gap-1.5 mb-1.5">
                  <MapPin className="w-3.5 h-3.5 text-brand-secondary shrink-0" />
                  <span className="line-clamp-1">{ev.destination_name}</span>
                </p>

                <p className="text-xs text-brand-muted flex items-center gap-1.5 mb-4">
                  <Clock className="w-3.5 h-3.5 text-brand-accent shrink-0" />
                  <span>
                    {ev.event_date} &bull; {ev.start_time}
                  </span>
                </p>

                <p className="text-xs text-brand-muted line-clamp-2 leading-relaxed mb-4 flex-1">
                  {ev.description}
                </p>

                {/* Footer status */}
                <div className="pt-4 border-t border-brand-border flex items-center justify-between text-xs">
                  <span className="flex items-center gap-1 font-semibold text-brand-primary">
                    <Users className="w-3.5 h-3.5 text-brand-secondary" />
                    {ev.approved_participants}{" "}
                    {ev.max_participants ? `/ ${ev.max_participants}` : ""} Peserta
                  </span>
                  <span
                    className={`font-semibold capitalize ${
                      ev.status === "open" ? "text-emerald-600" : "text-brand-muted"
                    }`}
                  >
                    {ev.status}
                  </span>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
