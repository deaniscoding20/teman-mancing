"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { User, MapPin, Compass, Calendar, Edit3, Check, X, MessageSquare, ArrowRight, Fish, PlusCircle } from "lucide-react";
import { fetchApi, resolveMediaUrl } from "@/lib/api-client";
import { getStoredUser, getToken } from "@/lib/auth";
import { UserProfile } from "@/types/api";

interface UserSpotItem {
  id: string;
  name: string;
  description: string;
  address: string;
  primary_photo_url?: string;
  comments_count?: number;
  created_at: string;
}

interface UserEventItem {
  id: string;
  name: string;
  description: string;
  destination_name: string;
  destination_address: string;
  event_date: string;
  start_time: string;
  approved_participants?: number;
  max_participants?: number;
  difficulty: string;
  status: string;
  created_at: string;
}

export default function ProfilePage() {
  const params = useParams();
  const username = params?.username as string;

  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [editName, setEditName] = useState("");
  const [editBio, setEditBio] = useState("");
  const [editLocation, setEditLocation] = useState("");
  const [saving, setSaving] = useState(false);

  // Activities state (Issue #17)
  const [activeTab, setActiveTab] = useState<"spots" | "created_events" | "joined_events">("spots");
  const [userSpots, setUserSpots] = useState<UserSpotItem[]>([]);
  const [createdEvents, setCreatedEvents] = useState<UserEventItem[]>([]);
  const [joinedEvents, setJoinedEvents] = useState<UserEventItem[]>([]);
  const [loadingActivities, setLoadingActivities] = useState(false);

  const currentUser = getStoredUser();
  const isOwner = currentUser && currentUser.username === username;

  useEffect(() => {
    if (!username) return;
    setLoading(true);

    fetchApi<UserProfile>(`/users/${username}`)
      .then((res) => {
        if (res.data) {
          const u = res.data;
          setProfile(u);
          setEditName(u.name || "");
          setEditBio(u.bio || "");
          setEditLocation(u.location || "");

          // Load user activities in parallel (Issue #17)
          setLoadingActivities(true);
          Promise.all([
            fetchApi<UserSpotItem[]>(`/spots?user_id=${u.id}&limit=20`),
            fetchApi<UserEventItem[]>(`/events?creator_id=${u.id}&limit=20`),
            fetchApi<UserEventItem[]>(`/events?participant_id=${u.id}&limit=20`),
          ])
            .then(([spotsRes, createdRes, joinedRes]) => {
              if (spotsRes.data) setUserSpots(spotsRes.data);
              if (createdRes.data) setCreatedEvents(createdRes.data);
              if (joinedRes.data) setJoinedEvents(joinedRes.data);
            })
            .catch((err) => {
              console.warn("Failed to load user activities:", err);
            })
            .finally(() => {
              setLoadingActivities(false);
            });
        }
      })
      .catch((err) => {
        setError(err.message || "Gagal memuat profil pengguna.");
      })
      .finally(() => {
        setLoading(false);
      });
  }, [username]);

  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = getToken();
    if (!token) return;

    setSaving(true);
    try {
      const res = await fetchApi<UserProfile>("/users/me", {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          name: editName,
          bio: editBio,
          location: editLocation,
        }),
      });

      if (res.data) {
        setProfile((prev) => prev ? { ...prev, ...res.data } : null);
        setIsEditing(false);
      }
    } catch (err: any) {
      alert(err.message || "Gagal menyimpan perubahan profil.");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-16 text-center text-brand-muted">
        Memuat profil...
      </div>
    );
  }

  if (error || !profile) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-16 text-center">
        <h2 className="text-xl font-bold text-brand-primary mb-2">Profil Tidak Ditemukan</h2>
        <p className="text-sm text-brand-muted">Pengguna dengan username @{username} tidak ada.</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Profile Header Card */}
      <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-8 shadow-sm">
        <div className="flex flex-col sm:flex-row items-center sm:items-start gap-6">
          {/* Avatar */}
          <div className="w-24 h-24 sm:w-28 sm:h-28 rounded-full bg-brand-primary/10 border-2 border-brand-secondary flex items-center justify-center text-brand-primary font-extrabold text-3xl shadow-sm shrink-0 overflow-hidden">
            {profile.avatar_url ? (
              <img src={profile.avatar_url} alt={profile.name} className="w-full h-full object-cover" />
            ) : (
              <User className="w-12 h-12 text-brand-secondary" />
            )}
          </div>

          {/* Details */}
          <div className="flex-1 text-center sm:text-left">
            {!isEditing ? (
              <>
                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                  <div>
                    <h1 className="text-2xl font-bold text-brand-primary">{profile.name}</h1>
                    <p className="text-sm text-brand-secondary font-medium">@{profile.username}</p>
                  </div>
                  {isOwner ? (
                    <button
                      onClick={() => setIsEditing(true)}
                      className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-xl border border-brand-border bg-brand-background text-brand-primary text-xs font-semibold hover:border-brand-secondary transition"
                    >
                      <Edit3 className="w-3.5 h-3.5 text-brand-secondary" />
                      Edit Profil
                    </button>
                  ) : (
                    <Link
                      href={`/chat?userId=${profile.id}`}
                      className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-xl bg-brand-secondary text-white text-xs font-semibold hover:bg-brand-secondary/90 transition shadow-sm"
                    >
                      <MessageSquare className="w-3.5 h-3.5" />
                      Kirim Pesan
                    </Link>
                  )}
                </div>

                <p className="mt-3 text-sm text-brand-muted leading-relaxed">
                  {profile.bio || "Belum ada bio yang dibagikan."}
                </p>

                {profile.location && (
                  <div className="mt-3 flex items-center justify-center sm:justify-start gap-1.5 text-xs text-brand-muted">
                    <MapPin className="w-3.5 h-3.5 text-brand-secondary" />
                    <span>{profile.location}</span>
                  </div>
                )}
              </>
            ) : (
              <form onSubmit={handleSaveProfile} className="flex flex-col gap-3 w-full">
                <div>
                  <label className="block text-xs font-semibold text-brand-text mb-1">Nama</label>
                  <input
                    type="text"
                    required
                    value={editName}
                    onChange={(e) => setEditName(e.target.value)}
                    className="w-full px-3 py-1.5 rounded-lg border border-brand-border bg-brand-background text-sm"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-brand-text mb-1">Bio</label>
                  <textarea
                    rows={2}
                    value={editBio}
                    onChange={(e) => setEditBio(e.target.value)}
                    placeholder="Ceritakan tentang hobi mancingmu..."
                    className="w-full px-3 py-1.5 rounded-lg border border-brand-border bg-brand-background text-sm"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-brand-text mb-1">Lokasi Domisili</label>
                  <input
                    type="text"
                    value={editLocation}
                    onChange={(e) => setEditLocation(e.target.value)}
                    placeholder="Contoh: Jakarta Barat"
                    className="w-full px-3 py-1.5 rounded-lg border border-brand-border bg-brand-background text-sm"
                  />
                </div>
                <div className="flex items-center gap-2 mt-1">
                  <button
                    type="submit"
                    disabled={saving}
                    className="inline-flex items-center gap-1 px-4 py-1.5 rounded-lg bg-brand-secondary text-white text-xs font-bold hover:bg-brand-secondary/90 transition"
                  >
                    <Check className="w-3.5 h-3.5" />
                    Simpan
                  </button>
                  <button
                    type="button"
                    onClick={() => setIsEditing(false)}
                    className="inline-flex items-center gap-1 px-3 py-1.5 rounded-lg border border-brand-border text-xs font-semibold hover:bg-brand-background transition"
                  >
                    <X className="w-3.5 h-3.5" />
                    Batal
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>

        {/* Counters / Stats Buttons to switch tabs */}
        <div className="grid grid-cols-3 gap-3 mt-8 pt-6 border-t border-brand-border text-center">
          <button
            onClick={() => setActiveTab("spots")}
            className={`p-3 rounded-2xl border transition text-center cursor-pointer ${
              activeTab === "spots"
                ? "bg-brand-primary/10 border-brand-secondary ring-2 ring-brand-secondary/20"
                : "bg-brand-background border-brand-border hover:border-brand-secondary/40"
            }`}
          >
            <span className="block text-xl font-extrabold text-brand-primary">
              {profile.spots_count ?? userSpots.length}
            </span>
            <span className="text-xs text-brand-muted flex items-center justify-center gap-1 mt-0.5 font-medium">
              <Compass className="w-3.5 h-3.5 text-brand-secondary" /> Spot Dibagikan
            </span>
          </button>

          <button
            onClick={() => setActiveTab("created_events")}
            className={`p-3 rounded-2xl border transition text-center cursor-pointer ${
              activeTab === "created_events"
                ? "bg-brand-primary/10 border-brand-secondary ring-2 ring-brand-secondary/20"
                : "bg-brand-background border-brand-border hover:border-brand-secondary/40"
            }`}
          >
            <span className="block text-xl font-extrabold text-brand-primary">
              {profile.events_created_count ?? createdEvents.length}
            </span>
            <span className="text-xs text-brand-muted flex items-center justify-center gap-1 mt-0.5 font-medium">
              <Calendar className="w-3.5 h-3.5 text-brand-secondary" /> Event Dibuat
            </span>
          </button>

          <button
            onClick={() => setActiveTab("joined_events")}
            className={`p-3 rounded-2xl border transition text-center cursor-pointer ${
              activeTab === "joined_events"
                ? "bg-brand-primary/10 border-brand-secondary ring-2 ring-brand-secondary/20"
                : "bg-brand-background border-brand-border hover:border-brand-secondary/40"
            }`}
          >
            <span className="block text-xl font-extrabold text-brand-primary">
              {profile.events_joined_count ?? joinedEvents.length}
            </span>
            <span className="text-xs text-brand-muted flex items-center justify-center gap-1 mt-0.5 font-medium">
              <Calendar className="w-3.5 h-3.5 text-brand-accent" /> Event Diikuti
            </span>
          </button>
        </div>
      </div>

      {/* Activity Content Tabs (Issue #17) */}
      <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-8 shadow-sm">
        {/* Tab Headers */}
        <div className="flex items-center gap-2 border-b border-brand-border pb-4 mb-6">
          <button
            onClick={() => setActiveTab("spots")}
            className={`px-4 py-2 rounded-xl text-xs sm:text-sm font-bold transition flex items-center gap-2 cursor-pointer ${
              activeTab === "spots"
                ? "bg-brand-primary text-white shadow-xs"
                : "text-brand-muted hover:text-brand-primary hover:bg-brand-background"
            }`}
          >
            <Compass className="w-4 h-4 text-brand-secondary" />
            <span>Spot Dibagikan ({userSpots.length})</span>
          </button>

          <button
            onClick={() => setActiveTab("created_events")}
            className={`px-4 py-2 rounded-xl text-xs sm:text-sm font-bold transition flex items-center gap-2 cursor-pointer ${
              activeTab === "created_events"
                ? "bg-brand-primary text-white shadow-xs"
                : "text-brand-muted hover:text-brand-primary hover:bg-brand-background"
            }`}
          >
            <Calendar className="w-4 h-4 text-brand-secondary" />
            <span>Event Dibuat ({createdEvents.length})</span>
          </button>

          <button
            onClick={() => setActiveTab("joined_events")}
            className={`px-4 py-2 rounded-xl text-xs sm:text-sm font-bold transition flex items-center gap-2 cursor-pointer ${
              activeTab === "joined_events"
                ? "bg-brand-primary text-white shadow-xs"
                : "text-brand-muted hover:text-brand-primary hover:bg-brand-background"
            }`}
          >
            <Calendar className="w-4 h-4 text-brand-accent" />
            <span>Event Diikuti ({joinedEvents.length})</span>
          </button>
        </div>

        {/* Tab 1: Spot Dibagikan */}
        {activeTab === "spots" && (
          <div>
            {loadingActivities ? (
              <div className="py-12 text-center text-xs text-brand-muted">Memuat daftar spot...</div>
            ) : userSpots.length === 0 ? (
              <div className="p-8 text-center bg-brand-background rounded-2xl border border-dashed border-brand-border">
                <Compass className="w-8 h-8 text-brand-muted/70 mx-auto mb-2" />
                <p className="text-sm font-bold text-brand-primary">Belum Ada Spot yang Dibagikan</p>
                <p className="text-xs text-brand-muted mt-1">
                  {isOwner
                    ? "Bagikan titik spot mancing andalanmu kepada komunitas!"
                    : "Angler ini belum pernah membagikan spot mancing."}
                </p>
                {isOwner && (
                  <Link
                    href="/spots/create"
                    className="mt-4 inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-brand-secondary text-white text-xs font-bold hover:bg-brand-secondary/90 transition shadow-xs"
                  >
                    <PlusCircle className="w-3.5 h-3.5" />
                    <span>Bagikan Spot Pertama</span>
                  </Link>
                )}
              </div>
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {userSpots.map((spot) => (
                  <Link
                    key={spot.id}
                    href={`/spots/${spot.id}`}
                    className="p-4 rounded-2xl bg-brand-background border border-brand-border hover:border-brand-secondary/40 transition flex gap-3.5 group"
                  >
                    <div className="w-20 h-20 rounded-xl bg-brand-primary/10 overflow-hidden shrink-0 flex items-center justify-center text-brand-muted">
                      {spot.primary_photo_url ? (
                        <img
                          src={resolveMediaUrl(spot.primary_photo_url)}
                          alt={spot.name}
                          className="w-full h-full object-cover group-hover:scale-105 transition"
                        />
                      ) : (
                        <Fish className="w-8 h-8 text-brand-secondary" />
                      )}
                    </div>
                    <div className="flex-1 min-w-0 flex flex-col justify-between">
                      <div>
                        <h4 className="text-sm font-bold text-brand-primary group-hover:text-brand-secondary transition truncate">
                          {spot.name}
                        </h4>
                        <p className="text-xs text-brand-muted truncate mt-1 flex items-center gap-1">
                          <MapPin className="w-3.5 h-3.5 text-brand-secondary shrink-0" />
                          <span>{spot.address}</span>
                        </p>
                      </div>
                      <div className="flex items-center justify-between text-[11px] text-brand-muted mt-2 pt-2 border-t border-brand-border/60">
                        <span>{spot.comments_count || 0} Ulasan</span>
                        <span className="text-brand-secondary font-bold inline-flex items-center gap-1">
                          Lihat Detail &rarr;
                        </span>
                      </div>
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Tab 2: Event Dibuat */}
        {activeTab === "created_events" && (
          <div>
            {loadingActivities ? (
              <div className="py-12 text-center text-xs text-brand-muted">Memuat daftar event...</div>
            ) : createdEvents.length === 0 ? (
              <div className="p-8 text-center bg-brand-background rounded-2xl border border-dashed border-brand-border">
                <Calendar className="w-8 h-8 text-brand-muted/70 mx-auto mb-2" />
                <p className="text-sm font-bold text-brand-primary">Belum Pernah Membuat Event</p>
                <p className="text-xs text-brand-muted mt-1">
                  {isOwner
                    ? "Inisiasi trip mancing bareng dan ajak teman-teman satu komunitas!"
                    : "Angler ini belum pernah membuat agenda janjian mancing."}
                </p>
                {isOwner && (
                  <Link
                    href="/events/create"
                    className="mt-4 inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-brand-secondary text-white text-xs font-bold hover:bg-brand-secondary/90 transition shadow-xs"
                  >
                    <PlusCircle className="w-3.5 h-3.5" />
                    <span>Buat Janjian Mancing</span>
                  </Link>
                )}
              </div>
            ) : (
              <div className="space-y-3">
                {createdEvents.map((ev) => (
                  <Link
                    key={ev.id}
                    href={`/events/${ev.id}`}
                    className="p-4 rounded-2xl bg-brand-background border border-brand-border hover:border-brand-secondary/40 transition flex flex-col sm:flex-row sm:items-center justify-between gap-3 group"
                  >
                    <div className="min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-brand-secondary/15 text-brand-secondary uppercase">
                          {ev.difficulty}
                        </span>
                        <span className="text-xs text-brand-muted">
                          {ev.event_date} • {ev.start_time}
                        </span>
                      </div>
                      <h4 className="text-sm font-bold text-brand-primary group-hover:text-brand-secondary transition truncate">
                        {ev.name}
                      </h4>
                      <p className="text-xs text-brand-muted truncate flex items-center gap-1 mt-0.5">
                        <MapPin className="w-3.5 h-3.5 text-brand-secondary shrink-0" />
                        <span>{ev.destination_name} ({ev.destination_address})</span>
                      </p>
                    </div>

                    <div className="flex items-center justify-between sm:justify-end gap-3 shrink-0 pt-2 sm:pt-0 border-t sm:border-0 border-brand-border/60">
                      <span className="text-xs font-semibold text-brand-muted">
                        {ev.approved_participants ?? 1} / {ev.max_participants ?? "∞"} Peserta
                      </span>
                      <span className="text-xs font-bold text-brand-secondary group-hover:underline">
                        Lihat Detail &rarr;
                      </span>
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Tab 3: Event Diikuti */}
        {activeTab === "joined_events" && (
          <div>
            {loadingActivities ? (
              <div className="py-12 text-center text-xs text-brand-muted">Memuat event diikuti...</div>
            ) : joinedEvents.length === 0 ? (
              <div className="p-8 text-center bg-brand-background rounded-2xl border border-dashed border-brand-border">
                <Calendar className="w-8 h-8 text-brand-muted/70 mx-auto mb-2" />
                <p className="text-sm font-bold text-brand-primary">Belum Mengikuti Event Apapun</p>
                <p className="text-xs text-brand-muted mt-1">
                  {isOwner
                    ? "Cari event seru dan bergabunglah dengan para angler di spot idaman."
                    : "Angler ini belum terdaftar pada agenda mancing manapun."}
                </p>
                <Link
                  href="/events"
                  className="mt-4 inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-brand-primary text-white text-xs font-semibold hover:bg-brand-primary/90 transition shadow-xs"
                >
                  <Compass className="w-3.5 h-3.5 text-brand-secondary" />
                  <span>Cari Janjian Mancing</span>
                </Link>
              </div>
            ) : (
              <div className="space-y-3">
                {joinedEvents.map((ev) => (
                  <Link
                    key={ev.id}
                    href={`/events/${ev.id}`}
                    className="p-4 rounded-2xl bg-brand-background border border-brand-border hover:border-brand-secondary/40 transition flex flex-col sm:flex-row sm:items-center justify-between gap-3 group"
                  >
                    <div className="min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-500/15 text-emerald-600 uppercase">
                          Peserta Terdaftar
                        </span>
                        <span className="text-xs text-brand-muted">
                          {ev.event_date} • {ev.start_time}
                        </span>
                      </div>
                      <h4 className="text-sm font-bold text-brand-primary group-hover:text-brand-secondary transition truncate">
                        {ev.name}
                      </h4>
                      <p className="text-xs text-brand-muted truncate flex items-center gap-1 mt-0.5">
                        <MapPin className="w-3.5 h-3.5 text-brand-secondary shrink-0" />
                        <span>{ev.destination_name}</span>
                      </p>
                    </div>

                    <div className="flex items-center justify-between sm:justify-end gap-3 shrink-0 pt-2 sm:pt-0 border-t sm:border-0 border-brand-border/60">
                      <span className="text-xs font-semibold text-brand-muted">
                        {ev.approved_participants ?? 1} Peserta
                      </span>
                      <span className="text-xs font-bold text-brand-secondary group-hover:underline">
                        Lihat Agenda &rarr;
                      </span>
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
