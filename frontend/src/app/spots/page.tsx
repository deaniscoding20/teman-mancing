"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Compass, Plus, Search, MapPin, MessageSquare, Fish, Map as MapIcon, Grid, ExternalLink } from "lucide-react";
import { fetchApi, resolveMediaUrl } from "@/lib/api-client";
import { getStoredUser } from "@/lib/auth";
import { FishingSpotItem } from "@/types/api";

export default function SpotsPage() {
  const [spots, setSpots] = useState<FishingSpotItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [viewMode, setViewMode] = useState<"grid" | "map">("grid");
  const [selectedSpot, setSelectedSpot] = useState<FishingSpotItem | null>(null);
  const currentUser = getStoredUser();

  const loadSpots = (q: string = "") => {
    setLoading(true);
    const queryParam = q ? `?q=${encodeURIComponent(q)}` : "";
    fetchApi<FishingSpotItem[]>(`/spots${queryParam}`)
      .then((res) => {
        if (res.data) {
          setSpots(res.data);
        }
      })
      .catch(() => {
        setSpots([]);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    loadSpots();
  }, []);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    loadSpots(search);
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Header & Actions */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 className="text-3xl font-extrabold text-brand-primary flex items-center gap-2.5">
            <Compass className="w-8 h-8 text-brand-secondary" />
            Spot Mancing
          </h1>
          <p className="text-sm text-brand-muted mt-1">
            Jelajahi dan bagikan spot pemancingan terbaik di seluruh Indonesia
          </p>
        </div>

        {currentUser && (
          <Link
            href="/spots/create"
            className="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl font-bold text-white bg-brand-secondary hover:bg-brand-secondary/90 shadow-sm transition self-start md:self-auto"
          >
            <Plus className="w-4 h-4" />
            Tambah Spot Baru
          </Link>
        )}
      </div>

      {/* Search Bar & View Toggle */}
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4 mb-8">
        <form onSubmit={handleSearch} className="flex-1 max-w-xl">
          <div className="relative">
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Cari spot berdasarkan nama, alamat, atau kota..."
              className="w-full pl-11 pr-4 py-3 rounded-2xl border border-brand-border bg-brand-surface text-brand-text placeholder-brand-muted/70 text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary shadow-sm transition"
            />
            <Search className="w-5 h-5 text-brand-muted absolute left-4 top-1/2 -translate-y-1/2" />
          </div>
        </form>

        <div className="flex items-center gap-1.5 p-1 rounded-2xl border border-brand-border bg-brand-surface self-start sm:self-auto shrink-0">
          <button
            type="button"
            onClick={() => setViewMode("grid")}
            className={`flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-semibold transition ${
              viewMode === "grid"
                ? "bg-brand-primary text-white shadow-xs"
                : "text-brand-muted hover:text-brand-primary"
            }`}
          >
            <Grid className="w-4 h-4" />
            <span>Kartu</span>
          </button>
          <button
            type="button"
            onClick={() => {
              setViewMode("map");
              if (!selectedSpot && spots.length > 0) {
                setSelectedSpot(spots[0]);
              }
            }}
            className={`flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-semibold transition ${
              viewMode === "map"
                ? "bg-brand-primary text-white shadow-xs"
                : "text-brand-muted hover:text-brand-primary"
            }`}
          >
            <MapIcon className="w-4 h-4" />
            <span>Peta Lokasi</span>
          </button>
        </div>
      </div>

      {/* Spots Grid */}
      {loading ? (
        <div className="text-center py-20 text-brand-muted">Memuat daftar spot...</div>
      ) : spots.length === 0 ? (
        <div className="bg-brand-surface rounded-3xl border border-brand-border p-12 text-center max-w-lg mx-auto">
          <div className="w-14 h-14 rounded-2xl bg-brand-background flex items-center justify-center mx-auto mb-4 text-brand-muted">
            <Compass className="w-7 h-7" />
          </div>
          <h3 className="text-lg font-bold text-brand-primary mb-1">Belum Ada Spot</h3>
          <p className="text-sm text-brand-muted mb-6">
            Jadilah yang pertama membagikan spot mancing andalanmu di platform ini.
          </p>
          <Link
            href={currentUser ? "/spots/create" : "/login"}
            className="inline-flex items-center gap-2 px-6 py-2.5 rounded-xl font-semibold text-white bg-brand-primary hover:bg-brand-primary/90 text-sm transition"
          >
            <Plus className="w-4 h-4 text-brand-secondary" />
            Bagikan Spot Pertama
          </Link>
        </div>
      ) : viewMode === "map" ? (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Map Frame */}
          <div className="lg:col-span-2 bg-brand-surface rounded-3xl border border-brand-border overflow-hidden shadow-sm flex flex-col h-[520px]">
            <div className="p-4 bg-brand-surface border-b border-brand-border flex items-center justify-between text-xs">
              <span className="font-bold text-brand-primary flex items-center gap-1.5">
                <MapPin className="w-4 h-4 text-brand-secondary" />
                {selectedSpot ? selectedSpot.name : "Pilih spot di samping untuk melihat lokasi"}
              </span>
              {selectedSpot && (
                <a
                  href={`https://www.google.com/maps/search/?api=1&query=${selectedSpot.latitude},${selectedSpot.longitude}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1 font-semibold text-brand-secondary hover:underline"
                >
                  <span>Buka di Google Maps</span>
                  <ExternalLink className="w-3 h-3" />
                </a>
              )}
            </div>
            <div className="flex-1 w-full bg-brand-background">
              {selectedSpot ? (
                <iframe
                  title="Peta Spot Mancing"
                  width="100%"
                  height="100%"
                  frameBorder="0"
                  scrolling="no"
                  src={`https://www.openstreetmap.org/export/embed.html?bbox=${selectedSpot.longitude - 0.02}%2C${selectedSpot.latitude - 0.02}%2C${selectedSpot.longitude + 0.02}%2C${selectedSpot.latitude + 0.02}&layer=mapnik&marker=${selectedSpot.latitude}%2C${selectedSpot.longitude}`}
                  className="w-full h-full border-0"
                />
              ) : (
                <div className="h-full flex items-center justify-center text-xs text-brand-muted">
                  Klik salah satu spot untuk menampilkan titik koordinat di peta
                </div>
              )}
            </div>
          </div>

          {/* Spot Selection List */}
          <div className="space-y-3 h-[520px] overflow-y-auto pr-1">
            {spots.map((spot) => {
              const isSelected = selectedSpot?.id === spot.id;
              return (
                <div
                  key={spot.id}
                  onClick={() => setSelectedSpot(spot)}
                  className={`p-4 rounded-2xl border transition cursor-pointer flex flex-col gap-1.5 ${
                    isSelected
                      ? "bg-brand-surface border-brand-secondary shadow-sm ring-2 ring-brand-secondary/20"
                      : "bg-brand-surface border-brand-border hover:border-brand-secondary/40"
                  }`}
                >
                  <div className="flex items-center justify-between gap-2">
                    <h4 className="font-bold text-sm text-brand-primary line-clamp-1">{spot.name}</h4>
                    <span className="text-[10px] text-brand-secondary font-mono">
                      {spot.latitude?.toFixed(4)}, {spot.longitude?.toFixed(4)}
                    </span>
                  </div>
                  <p className="text-xs text-brand-muted line-clamp-1 flex items-center gap-1">
                    <MapPin className="w-3 h-3 text-brand-secondary shrink-0" />
                    <span>{spot.address}</span>
                  </p>
                  <div className="pt-2 mt-1 border-t border-brand-border flex items-center justify-between text-xs">
                    <span className="text-brand-muted text-[11px]">{spot.comments_count || 0} Ulasan</span>
                    <Link
                      href={`/spots/${spot.id}`}
                      className="text-brand-secondary font-semibold hover:underline"
                    >
                      Detail Spot &rarr;
                    </Link>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {spots.map((spot) => (
            <Link
              key={spot.id}
              href={`/spots/${spot.id}`}
              className="bg-brand-surface rounded-2xl border border-brand-border overflow-hidden hover:border-brand-secondary/40 shadow-sm transition flex flex-col group"
            >
              {/* Photo Banner */}
              <div className="h-48 bg-brand-background relative flex items-center justify-center overflow-hidden">
                {spot.primary_photo_url ? (
                  <img
                    src={resolveMediaUrl(spot.primary_photo_url)}
                    alt={spot.name}
                    className="w-full h-full object-cover group-hover:scale-105 transition duration-300"
                  />
                ) : (
                  <div className="flex flex-col items-center gap-1.5 text-brand-muted/70">
                    <Compass className="w-10 h-10" />
                    <span className="text-xs">Foto Spot</span>
                  </div>
                )}
              </div>

              {/* Card Body */}
              <div className="p-5 flex-1 flex flex-col">
                <h3 className="text-lg font-bold text-brand-primary group-hover:text-brand-secondary transition line-clamp-1 mb-1.5">
                  {spot.name}
                </h3>
                <p className="text-xs text-brand-muted flex items-center gap-1 mb-3">
                  <MapPin className="w-3.5 h-3.5 text-brand-secondary shrink-0" />
                  <span className="line-clamp-1">{spot.address}</span>
                </p>

                <p className="text-xs text-brand-muted line-clamp-2 leading-relaxed mb-4 flex-1">
                  {spot.description}
                </p>

                {/* Footer metadata */}
                <div className="pt-4 border-t border-brand-border flex items-center justify-between text-xs text-brand-muted">
                  <span className="flex items-center gap-1">
                    <Fish className="w-3.5 h-3.5 text-brand-accent" />
                    {spot.fish_types?.length ? spot.fish_types.slice(0, 2).join(", ") : "Aneka Ikan"}
                  </span>
                  <span className="flex items-center gap-1">
                    <MessageSquare className="w-3.5 h-3.5 text-brand-secondary" />
                    {spot.comments_count || 0} Ulasan
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
