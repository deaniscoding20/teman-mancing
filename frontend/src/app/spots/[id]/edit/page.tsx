"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { Edit3, ArrowLeft, AlertCircle, Save } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { getStoredUser, getToken } from "@/lib/auth";

export default function EditSpotPage() {
  const params = useParams();
  const router = useRouter();
  const spotId = params?.id as string;

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [address, setAddress] = useState("");
  const [latitude, setLatitude] = useState("");
  const [longitude, setLongitude] = useState("");

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const currentUser = getStoredUser();
  const token = getToken();

  useEffect(() => {
    if (!spotId) return;

    fetchApi<any>(`/spots/${spotId}`)
      .then((res) => {
        if (res.data) {
          if (currentUser && currentUser.id !== res.data.user_id) {
            router.push(`/spots/${spotId}`);
            return;
          }
          setName(res.data.name || "");
          setDescription(res.data.description || "");
          setAddress(res.data.address || "");
          setLatitude(res.data.latitude?.toString() || "");
          setLongitude(res.data.longitude?.toString() || "");
        }
      })
      .catch((err) => {
        setError(err.message || "Gagal memuat spot.");
      })
      .finally(() => {
        setLoading(false);
      });
  }, [spotId, currentUser, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;

    const lat = parseFloat(latitude);
    const lng = parseFloat(longitude);
    if (isNaN(lat) || isNaN(lng)) {
      setError("Koordinat latitude dan longitude tidak valid.");
      return;
    }

    setSaving(true);
    setError(null);

    try {
      await fetchApi(`/spots/${spotId}`, {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          name,
          description,
          address,
          latitude: lat,
          longitude: lng,
        }),
      });

      router.push(`/spots/${spotId}`);
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Gagal memperbarui spot.");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <div className="text-center py-20 text-brand-muted">Memuat data spot...</div>;
  }

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <Link
        href={`/spots/${spotId}`}
        className="inline-flex items-center gap-1.5 text-sm font-semibold text-brand-secondary hover:underline mb-6"
      >
        <ArrowLeft className="w-4 h-4" /> Kembali ke Detail Spot
      </Link>

      <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-10 shadow-sm">
        <div className="mb-8">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand-secondary/10 text-brand-secondary text-xs font-semibold mb-2">
            <Edit3 className="w-4 h-4" /> Pengaturan Spot
          </div>
          <h1 className="text-2xl sm:text-3xl font-extrabold text-brand-primary">Edit Spot Mancing</h1>
          <p className="text-sm text-brand-muted mt-1">Perbarui informasi dan lokasi pemancingan</p>
        </div>

        {error && (
          <div className="mb-6 p-4 rounded-xl bg-red-50 border border-red-200 flex items-start gap-3 text-red-700 text-sm">
            <AlertCircle className="w-5 h-5 shrink-0 mt-0.5" />
            <span>{error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit} className="flex flex-col gap-6">
          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Nama Spot *
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
              Deskripsi Kondisi Spot *
            </label>
            <textarea
              rows={4}
              required
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full p-4 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Alamat Lengkap *
            </label>
            <input
              type="text"
              required
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Latitude *
              </label>
              <input
                type="text"
                required
                value={latitude}
                onChange={(e) => setLatitude(e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              />
            </div>
            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Longitude *
              </label>
              <input
                type="text"
                required
                value={longitude}
                onChange={(e) => setLongitude(e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={saving}
            className="mt-4 w-full py-3.5 px-4 rounded-xl font-bold text-white bg-brand-primary hover:bg-brand-primary/90 shadow-sm disabled:opacity-50 transition flex items-center justify-center gap-2"
          >
            <Save className="w-4 h-4 text-brand-secondary" />
            <span>{saving ? "Menyimpan perubahan..." : "Simpan Perubahan"}</span>
          </button>
        </form>
      </div>
    </div>
  );
}
