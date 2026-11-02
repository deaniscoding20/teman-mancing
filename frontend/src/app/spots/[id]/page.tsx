"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import {
  Compass,
  MapPin,
  Fish,
  Anchor,
  Sparkles,
  MessageSquare,
  Trash2,
  Edit3,
  ExternalLink,
  Send,
} from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { getStoredUser, getToken } from "@/lib/auth";

interface SpotDetail {
  id: string;
  user_id: string;
  name: string;
  description: string;
  address: string;
  latitude: number;
  longitude: number;
  fish_types?: { id: string; name: string }[];
  techniques?: { id: string; name: string }[];
  recommended_baits?: { id: string; name: string; description?: string }[];
  photos?: { id: string; object_key: string; is_primary: boolean }[];
  created_at: string;
}

interface SpotComment {
  id: string;
  user_id: string;
  user_name?: string;
  user_avatar?: string;
  comment: string;
  created_at: string;
}

export default function SpotDetailPage() {
  const params = useParams();
  const router = useRouter();
  const spotId = params?.id as string;

  const [spot, setSpot] = useState<SpotDetail | null>(null);
  const [comments, setComments] = useState<SpotComment[]>([]);
  const [newComment, setNewComment] = useState("");
  const [loading, setLoading] = useState(true);
  const [submittingComment, setSubmittingComment] = useState(false);

  const currentUser = getStoredUser();
  const token = getToken();

  useEffect(() => {
    if (!spotId) return;

    // Fetch Spot Details & Comments in parallel
    Promise.all([
      fetchApi<SpotDetail>(`/spots/${spotId}`),
      fetchApi<SpotComment[]>(`/spots/${spotId}/comments`),
    ])
      .then(([spotRes, commentsRes]) => {
        if (spotRes.data) setSpot(spotRes.data);
        if (commentsRes.data) setComments(commentsRes.data);
      })
      .catch((err) => {
        console.error("Failed to load spot:", err);
      })
      .finally(() => {
        setLoading(false);
      });
  }, [spotId]);

  const handlePostComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token || !newComment.trim()) return;

    setSubmittingComment(true);
    try {
      const res = await fetchApi<SpotComment>(`/spots/${spotId}/comments`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ comment: newComment }),
      });

      if (res.data) {
        setComments((prev) => [res.data!, ...prev]);
        setNewComment("");
      }
    } catch (err: any) {
      alert(err.message || "Gagal mengirim komentar.");
    } finally {
      setSubmittingComment(false);
    }
  };

  const handleDeleteComment = async (commentId: string) => {
    if (!token || !confirm("Yakin ingin menghapus komentar ini?")) return;

    try {
      await fetchApi(`/comments/${commentId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setComments((prev) => prev.filter((c) => c.id !== commentId));
    } catch (err: any) {
      alert(err.message || "Gagal menghapus komentar.");
    }
  };

  const handleDeleteSpot = async () => {
    if (!token || !confirm("Yakin ingin menghapus spot pemancingan ini?")) return;

    try {
      await fetchApi(`/spots/${spotId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      router.push("/spots");
    } catch (err: any) {
      alert(err.message || "Gagal menghapus spot.");
    }
  };

  if (loading) {
    return <div className="text-center py-20 text-brand-muted">Memuat informasi spot...</div>;
  }

  if (!spot) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-20 text-center">
        <h2 className="text-2xl font-bold text-brand-primary mb-2">Spot Tidak Ditemukan</h2>
        <Link href="/spots" className="text-sm text-brand-secondary hover:underline">
          &larr; Kembali ke daftar spot
        </Link>
      </div>
    );
  }

  const isOwner = currentUser && currentUser.id === spot.user_id;
  const directionsUrl = `https://www.google.com/maps/dir/?api=1&destination=${spot.latitude},${spot.longitude}`;

  return (
    <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Top Breadcrumb & Actions */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <Link href="/spots" className="text-sm font-semibold text-brand-secondary hover:underline">
          &larr; Kembali ke Semua Spot
        </Link>

        {isOwner && (
          <div className="flex items-center gap-2">
            <button
              onClick={handleDeleteSpot}
              className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-xl border border-red-200 text-red-600 bg-red-50 hover:bg-red-100 text-xs font-semibold transition"
            >
              <Trash2 className="w-3.5 h-3.5" />
              Hapus Spot
            </button>
          </div>
        )}
      </div>

      {/* Main Info Card */}
      <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-10 shadow-sm mb-8">
        <div className="flex flex-col gap-4">
          <div className="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-brand-secondary">
            <Compass className="w-4 h-4" /> Spot Pemancingan
          </div>
          <h1 className="text-3xl sm:text-4xl font-extrabold text-brand-primary">{spot.name}</h1>
          <p className="flex items-center gap-1.5 text-sm text-brand-muted">
            <MapPin className="w-4 h-4 text-brand-secondary shrink-0" />
            <span>{spot.address}</span>
          </p>

          <p className="mt-4 text-base text-brand-text leading-relaxed whitespace-pre-line">
            {spot.description}
          </p>
        </div>

        {/* Fish Types & Techniques Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8 pt-8 border-t border-brand-border">
          {/* Fish Types */}
          <div className="bg-brand-background rounded-2xl p-5 border border-brand-border">
            <h3 className="text-sm font-bold text-brand-primary flex items-center gap-2 mb-3">
              <Fish className="w-4 h-4 text-brand-accent" /> Target Jenis Ikan
            </h3>
            {spot.fish_types && spot.fish_types.length > 0 ? (
              <div className="flex flex-wrap gap-2">
                {spot.fish_types.map((f) => (
                  <span
                    key={f.id}
                    className="px-3 py-1 rounded-full bg-brand-surface text-brand-primary text-xs font-semibold border border-brand-border"
                  >
                    {f.name}
                  </span>
                ))}
              </div>
            ) : (
              <p className="text-xs text-brand-muted">Belum ada spesifikasi ikan.</p>
            )}
          </div>

          {/* Techniques */}
          <div className="bg-brand-background rounded-2xl p-5 border border-brand-border">
            <h3 className="text-sm font-bold text-brand-primary flex items-center gap-2 mb-3">
              <Anchor className="w-4 h-4 text-brand-secondary" /> Teknik Memancing Rekomendasi
            </h3>
            {spot.techniques && spot.techniques.length > 0 ? (
              <div className="flex flex-wrap gap-2">
                {spot.techniques.map((t) => (
                  <span
                    key={t.id}
                    className="px-3 py-1 rounded-full bg-brand-surface text-brand-primary text-xs font-semibold border border-brand-border"
                  >
                    {t.name}
                  </span>
                ))}
              </div>
            ) : (
              <p className="text-xs text-brand-muted">Belum ada spesifikasi teknik.</p>
            )}
          </div>
        </div>

        {/* Recommended Baits */}
        {spot.recommended_baits && spot.recommended_baits.length > 0 && (
          <div className="mt-6 p-5 rounded-2xl bg-amber-50/70 border border-amber-200">
            <h3 className="text-sm font-bold text-amber-900 flex items-center gap-2 mb-3">
              <Sparkles className="w-4 h-4 text-amber-600" /> Umpan Rekomendasi
            </h3>
            <ul className="space-y-2 text-xs text-amber-950">
              {spot.recommended_baits.map((b) => (
                <li key={b.id} className="flex flex-col sm:flex-row sm:items-baseline gap-1">
                  <span className="font-bold">&bull; {b.name}:</span>
                  <span className="text-amber-800">{b.description || "Umpan favorit"}</span>
                </li>
              ))}
            </ul>
          </div>
        )}

        {/* Map View & Directions Link */}
        <div className="mt-8 pt-8 border-t border-brand-border">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-bold text-brand-primary flex items-center gap-2">
              <MapPin className="w-4 h-4 text-brand-secondary" /> Peta Lokasi Spot
            </h3>
            <a
              href={directionsUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1.5 text-xs font-semibold text-brand-secondary hover:underline"
            >
              Petunjuk Arah Google Maps <ExternalLink className="w-3.5 h-3.5" />
            </a>
          </div>

          <div className="h-64 sm:h-80 w-full rounded-2xl overflow-hidden border border-brand-border relative bg-slate-100 flex items-center justify-center">
            {/* OpenStreetMap Embed Iframe */}
            <iframe
              title="Peta Lokasi"
              width="100%"
              height="100%"
              style={{ border: 0 }}
              loading="lazy"
              src={`https://www.openstreetmap.org/export/embed.html?bbox=${spot.longitude - 0.01}%2C${spot.latitude - 0.01}%2C${spot.longitude + 0.01}%2C${spot.latitude + 0.01}&layer=mapnik&marker=${spot.latitude}%2C${spot.longitude}`}
            />
          </div>
        </div>
      </div>

      {/* Comments Section */}
      <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-10 shadow-sm">
        <h2 className="text-xl font-bold text-brand-primary flex items-center gap-2.5 mb-6">
          <MessageSquare className="w-5 h-5 text-brand-secondary" />
          Ulasan & Diskusi Pemancing ({comments.length})
        </h2>

        {/* Comment input form */}
        {currentUser ? (
          <form onSubmit={handlePostComment} className="mb-8">
            <div className="flex flex-col gap-2.5">
              <textarea
                rows={3}
                required
                value={newComment}
                onChange={(e) => setNewComment(e.target.value)}
                placeholder="Bagikan pengalamanmu mancing di spot ini (kondisi air, hasil sambaran, umpan yang gacor)..."
                className="w-full p-4 rounded-2xl border border-brand-border bg-brand-background text-brand-text placeholder-brand-muted/70 text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary transition"
              />
              <button
                type="submit"
                disabled={submittingComment || !newComment.trim()}
                className="self-end inline-flex items-center gap-2 px-5 py-2.5 rounded-xl font-bold text-white bg-brand-primary hover:bg-brand-primary/90 disabled:opacity-50 text-xs shadow-sm transition"
              >
                <Send className="w-3.5 h-3.5 text-brand-secondary" />
                Kirim Komentar
              </button>
            </div>
          </form>
        ) : (
          <div className="mb-8 p-4 rounded-2xl bg-brand-background border border-brand-border text-center text-xs text-brand-muted">
            <Link href="/login" className="font-semibold text-brand-secondary hover:underline">
              Masuk
            </Link>{" "}
            untuk menulis ulasan atau berdiskusi di spot ini.
          </div>
        )}

        {/* Comments List */}
        {comments.length === 0 ? (
          <p className="text-center py-8 text-xs text-brand-muted">
            Belum ada komentar untuk spot ini. Mulai diskusi pertama!
          </p>
        ) : (
          <div className="space-y-4">
            {comments.map((c) => {
              const isCommentAuthor = currentUser && currentUser.id === c.user_id;
              return (
                <div
                  key={c.id}
                  className="p-4 rounded-2xl bg-brand-background border border-brand-border flex items-start justify-between gap-4"
                >
                  <div>
                    <div className="flex items-center gap-2 mb-1">
                      <span className="font-bold text-xs text-brand-primary">
                        {c.user_name || "Pemancing"}
                      </span>
                      <span className="text-[10px] text-brand-muted">
                        {new Date(c.created_at).toLocaleDateString("id-ID", {
                          day: "numeric",
                          month: "short",
                          year: "numeric",
                        })}
                      </span>
                    </div>
                    <p className="text-xs text-brand-text leading-relaxed whitespace-pre-line">
                      {c.comment}
                    </p>
                  </div>

                  {isCommentAuthor && (
                    <button
                      onClick={() => handleDeleteComment(c.id)}
                      className="p-1.5 text-brand-muted hover:text-red-500 rounded-lg transition"
                      title="Hapus Komentar"
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                    </button>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
