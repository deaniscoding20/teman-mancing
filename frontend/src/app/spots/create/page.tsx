"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Compass, Plus, Trash2, ArrowLeft, AlertCircle } from "lucide-react";
import LocationMapPicker from "@/components/common/LocationMapPicker";
import { fetchApi } from "@/lib/api-client";
import { getToken } from "@/lib/auth";

export default function CreateSpotPage() {
  const router = useRouter();

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [address, setAddress] = useState("");
  const [latitude, setLatitude] = useState("-6.2088");
  const [longitude, setLongitude] = useState("106.8456");

  const [fishTypes, setFishTypes] = useState("Nila, Kakap Putih, Gabus");
  const [techniques, setTechniques] = useState("Casting, Dasaran");

  const [baits, setBaits] = useState<{ name: string; description: string }[]>([
    { name: "Cacing Tanah / Lur", description: "Sangat ampuh saat sore hari" },
  ]);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAddBait = () => {
    setBaits([...baits, { name: "", description: "" }]);
  };

  const handleRemoveBait = (index: number) => {
    setBaits(baits.filter((_, i) => i !== index));
  };

  const handleBaitChange = (index: number, field: "name" | "description", val: string) => {
    const updated = [...baits];
    updated[index][field] = val;
    setBaits(updated);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const token = getToken();
    if (!token) {
      setError("Anda harus masuk terlebih dahulu.");
      return;
    }

    const lat = parseFloat(latitude);
    const lng = parseFloat(longitude);
    if (isNaN(lat) || isNaN(lng)) {
      setError("Koordinat latitude dan longitude tidak valid.");
      return;
    }

    const parsedFishTypes = fishTypes
      .split(",")
      .map((s) => s.trim())
      .filter((s) => s.length > 0);

    const parsedTechniques = techniques
      .split(",")
      .map((s) => s.trim())
      .filter((s) => s.length > 0);

    const validBaits = baits.filter((b) => b.name.trim().length > 0);

    setLoading(true);

    try {
      const res = await fetchApi<{ id: string }>("/spots", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          name,
          description,
          address,
          latitude: lat,
          longitude: lng,
          fish_types: parsedFishTypes,
          techniques: parsedTechniques,
          recommended_baits: validBaits,
        }),
      });

      if (res.data?.id) {
        router.push(`/spots/${res.data.id}`);
      } else {
        router.push("/spots");
      }
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Gagal membuat spot mancing.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Navigation */}
      <Link
        href="/spots"
        className="inline-flex items-center gap-1.5 text-sm font-semibold text-brand-secondary hover:underline mb-6"
      >
        <ArrowLeft className="w-4 h-4" /> Kembali ke Daftar Spot
      </Link>

      <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-10 shadow-sm">
        {/* Header */}
        <div className="mb-8">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand-secondary/10 text-brand-secondary text-xs font-semibold mb-2">
            <Compass className="w-4 h-4" /> Kontribusi Komunitas
          </div>
          <h1 className="text-2xl sm:text-3xl font-extrabold text-brand-primary">
            Bagikan Spot Mancing Baru
          </h1>
          <p className="text-sm text-brand-muted mt-1">
            Bantu sesama angler menemukan titik lokasi terbaik untuk memancing
          </p>
        </div>

        {/* Error notification */}
        {error && (
          <div className="mb-6 p-4 rounded-xl bg-red-50 border border-red-200 flex flex-col gap-2 text-red-700 text-sm">
            <div className="flex items-start gap-3">
              <AlertCircle className="w-5 h-5 shrink-0 mt-0.5 text-red-600" />
              <span>{error}</span>
            </div>
            {error.toLowerCase().includes("verifikasi") && (
              <div className="mt-1 pl-8">
                <Link
                  href="/verify-email"
                  className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-brand-primary text-white text-xs font-bold rounded-lg hover:bg-brand-primary/90 transition shadow-xs"
                >
                  Verifikasi Email Sekarang →
                </Link>
              </div>
            )}
          </div>
        )}

        {/* Form */}
        <form onSubmit={handleSubmit} className="flex flex-col gap-6">
          {/* Nama Spot */}
          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Nama Spot Mancing *
            </label>
            <input
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Contoh: Dermaga Pantai Ancol atau Muara Gembong"
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          {/* Deskripsi */}
          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Deskripsi Kondisi & Akses *
            </label>
            <textarea
              rows={4}
              required
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Jelaskan kondisi spot, spot teduh/terik, waktu makan ikan paling ramai, atau kemudahan parkir..."
              className="w-full p-4 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          {/* Location Map Picker (Search, GPS, Map Preview & Google Maps) */}
          <LocationMapPicker
            address={address}
            latitude={latitude}
            longitude={longitude}
            onAddressChange={(newAddr) => setAddress(newAddr)}
            onCoordinateChange={(lat, lng) => {
              setLatitude(lat);
              setLongitude(lng);
            }}
          />

          {/* Jenis Ikan */}
          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Target Jenis Ikan (Pisahkan dengan koma)
            </label>
            <input
              type="text"
              value={fishTypes}
              onChange={(e) => setFishTypes(e.target.value)}
              placeholder="Contoh: Kakap Putih, Nila, Kerapu, Gabus"
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          {/* Teknik Memancing */}
          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Teknik Memancing (Pisahkan dengan koma)
            </label>
            <input
              type="text"
              value={techniques}
              onChange={(e) => setTechniques(e.target.value)}
              placeholder="Contoh: Casting, Dasaran, Jigging, Pelampung"
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          {/* Umpan Rekomendasi Dynamic List */}
          <div>
            <div className="flex items-center justify-between mb-2">
              <label className="text-xs font-semibold text-brand-text uppercase tracking-wider">
                Rekomendasi Racikan Umpan
              </label>
              <button
                type="button"
                onClick={handleAddBait}
                className="text-xs font-bold text-brand-secondary hover:underline flex items-center gap-1"
              >
                <Plus className="w-3.5 h-3.5" /> Tambah Umpan
              </button>
            </div>

            <div className="space-y-3">
              {baits.map((bait, idx) => (
                <div
                  key={idx}
                  className="flex items-center gap-3 p-3 rounded-2xl bg-brand-background border border-brand-border"
                >
                  <input
                    type="text"
                    value={bait.name}
                    onChange={(e) => handleBaitChange(idx, "name", e.target.value)}
                    placeholder="Nama Umpan (e.g. Udang Hidup)"
                    className="flex-1 px-3 py-1.5 rounded-lg border border-brand-border bg-brand-surface text-xs"
                  />
                  <input
                    type="text"
                    value={bait.description}
                    onChange={(e) => handleBaitChange(idx, "description", e.target.value)}
                    placeholder="Catatan / Cara Pasang"
                    className="flex-1 px-3 py-1.5 rounded-lg border border-brand-border bg-brand-surface text-xs"
                  />
                  {baits.length > 1 && (
                    <button
                      type="button"
                      onClick={() => handleRemoveBait(idx)}
                      className="text-brand-muted hover:text-red-500 p-1"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  )}
                </div>
              ))}
            </div>
          </div>

          {/* Submit Button */}
          <button
            type="submit"
            disabled={loading}
            className="mt-4 w-full py-3.5 px-4 rounded-xl font-bold text-white bg-brand-primary hover:bg-brand-primary/90 focus:outline-none focus:ring-2 focus:ring-brand-primary shadow-sm disabled:opacity-50 transition flex items-center justify-center gap-2"
          >
            {loading ? (
              <span>Menyimpan spot...</span>
            ) : (
              <>
                <Plus className="w-4 h-4 text-brand-secondary" />
                <span>Publikasikan Spot Mancing</span>
              </>
            )}
          </button>
        </form>
      </div>
    </div>
  );
}
