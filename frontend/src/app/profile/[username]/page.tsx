"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { User, MapPin, Compass, Calendar, Edit3, Check, X } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { getStoredUser, getToken } from "@/lib/auth";
import { UserProfile } from "@/types/api";

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

  const currentUser = getStoredUser();
  const isOwner = currentUser && currentUser.username === username;

  useEffect(() => {
    if (!username) return;
    setLoading(true);

    fetchApi<UserProfile>(`/users/${username}`)
      .then((res) => {
        if (res.data) {
          setProfile(res.data);
          setEditName(res.data.name || "");
          setEditBio(res.data.bio || "");
          setEditLocation(res.data.location || "");
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
                  {isOwner && (
                    <button
                      onClick={() => setIsEditing(true)}
                      className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-xl border border-brand-border bg-brand-background text-brand-primary text-xs font-semibold hover:border-brand-secondary transition"
                    >
                      <Edit3 className="w-3.5 h-3.5 text-brand-secondary" />
                      Edit Profil
                    </button>
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

        {/* Counters / Stats */}
        <div className="grid grid-cols-3 gap-3 mt-8 pt-6 border-t border-brand-border text-center">
          <div className="p-3 rounded-2xl bg-brand-background">
            <span className="block text-xl font-extrabold text-brand-primary">
              {profile.spots_count ?? 0}
            </span>
            <span className="text-xs text-brand-muted flex items-center justify-center gap-1 mt-0.5">
              <Compass className="w-3.5 h-3.5 text-brand-secondary" /> Spot Dibagikan
            </span>
          </div>

          <div className="p-3 rounded-2xl bg-brand-background">
            <span className="block text-xl font-extrabold text-brand-primary">
              {profile.events_created_count ?? 0}
            </span>
            <span className="text-xs text-brand-muted flex items-center justify-center gap-1 mt-0.5">
              <Calendar className="w-3.5 h-3.5 text-brand-secondary" /> Event Dibuat
            </span>
          </div>

          <div className="p-3 rounded-2xl bg-brand-background">
            <span className="block text-xl font-extrabold text-brand-primary">
              {profile.events_joined_count ?? 0}
            </span>
            <span className="text-xs text-brand-muted flex items-center justify-center gap-1 mt-0.5">
              <Calendar className="w-3.5 h-3.5 text-brand-accent" /> Event Diikuti
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
