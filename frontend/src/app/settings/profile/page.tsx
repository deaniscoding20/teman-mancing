"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { User, Save, ArrowLeft, AlertCircle, CheckCircle2, Upload } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { getStoredUser, getToken, setAuth } from "@/lib/auth";

export default function SettingsProfilePage() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [bio, setBio] = useState("");
  const [location, setLocation] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [uploadingAvatar, setUploadingAvatar] = useState(false);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const currentUser = getStoredUser();
  const token = getToken();

  useEffect(() => {
    if (!token) {
      router.push("/login");
      return;
    }

    fetchApi<any>("/auth/me", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => {
        if (res.data?.User) {
          setName(res.data.User.name || "");
          setBio(res.data.User.bio || "");
          setLocation(res.data.User.location || "");
        }
      })
      .catch((err) => {
        setError(err.message || "Gagal memuat profil.");
      })
      .finally(() => {
        setLoading(false);
      });
  }, [token, router]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;

    setSaving(true);
    setStatusMessage(null);
    setError(null);

    try {
      const res = await fetchApi<any>("/users/me", {
        method: "PATCH",
        headers: { Authorization: `Bearer ${token}` },
        body: JSON.stringify({ name, bio, location }),
      });

      if (res.data) {
        setStatusMessage("Profil berhasil diperbarui!");
        if (currentUser) {
          setAuth(token, { ...currentUser, name, bio, location });
        }
      }
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Gagal memperbarui profil.");
    } finally {
      setSaving(false);
    }
  };

  const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !token) return;

    const formData = new FormData();
    formData.append("avatar", file);

    setUploadingAvatar(true);
    setError(null);

    try {
      await fetchApi("/users/me/avatar", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      });
      setStatusMessage("Foto avatar berhasil diupload!");
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Gagal mengupload avatar.");
    } finally {
      setUploadingAvatar(false);
    }
  };

  if (loading) {
    return <div className="text-center py-20 text-brand-muted">Memuat data akun...</div>;
  }

  return (
    <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <Link
        href={currentUser ? `/profile/${currentUser.username}` : "/"}
        className="inline-flex items-center gap-1.5 text-sm font-semibold text-brand-secondary hover:underline mb-6"
      >
        <ArrowLeft className="w-4 h-4" /> Kembali ke Profil
      </Link>

      <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-10 shadow-sm">
        <div className="mb-8">
          <h1 className="text-2xl sm:text-3xl font-extrabold text-brand-primary">Pengaturan Profil</h1>
          <p className="text-xs text-brand-muted mt-1">Perbarui biodata dan informasi publik akunmu</p>
        </div>

        {statusMessage && (
          <div className="mb-6 p-4 rounded-xl bg-green-50 border border-green-200 flex items-center gap-3 text-green-800 text-sm">
            <CheckCircle2 className="w-5 h-5 text-green-600 shrink-0" />
            <span>{statusMessage}</span>
          </div>
        )}

        {error && (
          <div className="mb-6 p-4 rounded-xl bg-red-50 border border-red-200 flex items-start gap-3 text-red-700 text-sm">
            <AlertCircle className="w-5 h-5 shrink-0 mt-0.5" />
            <span>{error}</span>
          </div>
        )}

        {/* Avatar Upload Field */}
        <div className="mb-8 p-5 rounded-2xl bg-brand-background border border-brand-border flex items-center justify-between gap-4">
          <div className="flex items-center gap-3.5">
            <div className="w-12 h-12 rounded-full bg-brand-secondary/15 flex items-center justify-center text-brand-secondary font-bold">
              <User className="w-6 h-6" />
            </div>
            <div>
              <h4 className="text-xs font-bold text-brand-primary">Foto Avatar</h4>
              <p className="text-[10px] text-brand-muted">Format: JPEG, PNG, WEBP (Maksimal 5MB)</p>
            </div>
          </div>

          <label className="cursor-pointer inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-xl border border-brand-border bg-brand-surface text-brand-primary text-xs font-semibold hover:border-brand-secondary transition">
            <Upload className="w-3.5 h-3.5 text-brand-secondary" />
            <span>{uploadingAvatar ? "Mengupload..." : "Ganti Foto"}</span>
            <input
              type="file"
              accept="image/jpeg,image/png,image/webp"
              onChange={handleAvatarChange}
              className="hidden"
            />
          </label>
        </div>

        <form onSubmit={handleSave} className="flex flex-col gap-5">
          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Nama Lengkap
            </label>
            <input
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Biodata Singkat (Bio)
            </label>
            <textarea
              rows={3}
              value={bio}
              onChange={(e) => setBio(e.target.value)}
              placeholder="Ceritakan tentang teknik mancing favoritmu..."
              className="w-full p-4 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Lokasi Domisili
            </label>
            <input
              type="text"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              placeholder="Contoh: Surabaya, Jawa Timur"
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          <button
            type="submit"
            disabled={saving}
            className="mt-3 w-full py-3.5 px-4 rounded-xl font-bold text-white bg-brand-primary hover:bg-brand-primary/90 shadow-sm disabled:opacity-50 transition flex items-center justify-center gap-2 text-sm"
          >
            <Save className="w-4 h-4 text-brand-secondary" />
            <span>{saving ? "Menyimpan perubahan..." : "Simpan Pengaturan"}</span>
          </button>
        </form>
      </div>
    </div>
  );
}
