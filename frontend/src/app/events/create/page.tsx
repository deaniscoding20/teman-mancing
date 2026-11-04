"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Calendar, Plus, ArrowLeft, AlertCircle, Car, Ship } from "lucide-react";
import LocationMapPicker from "@/components/common/LocationMapPicker";
import { fetchApi } from "@/lib/api-client";
import { getToken } from "@/lib/auth";

export default function CreateEventPage() {
  const router = useRouter();

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [destinationName, setDestinationName] = useState("");
  const [destinationAddress, setDestinationAddress] = useState("");
  const [destinationLat, setDestinationLat] = useState("-6.2088");
  const [destinationLng, setDestinationLng] = useState("106.8456");
  const [eventDate, setEventDate] = useState("2026-12-15");
  const [startTime, setStartTime] = useState("06:00:00");
  const [endTime, setEndTime] = useState("12:00:00");
  const [meetingPoint, setMeetingPoint] = useState("");
  const [meetingAddress, setMeetingAddress] = useState("");
  const [meetingLat, setMeetingLat] = useState("-6.2088");
  const [meetingLng, setMeetingLng] = useState("106.8456");
  const [sameMeetingLocation, setSameMeetingLocation] = useState(false);
  const [locationNotes, setLocationNotes] = useState("");
  const [difficulty, setDifficulty] = useState("pemula");
  const [visibility, setVisibility] = useState<"public" | "private">("public");
  const [maxParticipants, setMaxParticipants] = useState("10");

  const [hasParkingFee, setHasParkingFee] = useState(false);
  const [parkingFee, setParkingFee] = useState("10000");
  const [requiresBoat, setRequiresBoat] = useState(false);
  const [boatRentalFee, setBoatRentalFee] = useState("150000");
  const [additionalCost, setAdditionalCost] = useState("Iuran umpan hidup bersama");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const token = getToken();
    if (!token) {
      setError("Anda harus masuk terlebih dahulu.");
      return;
    }

    const parsedMax = maxParticipants ? parseInt(maxParticipants, 10) : undefined;
    const parsedParking = hasParkingFee && parkingFee ? parseFloat(parkingFee) : undefined;
    const parsedBoat = requiresBoat && boatRentalFee ? parseFloat(boatRentalFee) : undefined;

    setLoading(true);

    try {
      const res = await fetchApi<{ id: string }>("/events", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          name,
          description,
          destination_name: destinationName,
          destination_address: destinationAddress,
          latitude: parseFloat(destinationLat) || -6.2088,
          longitude: parseFloat(destinationLng) || 106.8456,
          event_date: eventDate,
          start_time: startTime,
          end_time: endTime || undefined,
          meeting_point: sameMeetingLocation ? destinationName : meetingPoint,
          meeting_address: sameMeetingLocation ? destinationAddress : meetingAddress,
          meeting_latitude: sameMeetingLocation ? (parseFloat(destinationLat) || -6.2088) : (parseFloat(meetingLat) || -6.2088),
          meeting_longitude: sameMeetingLocation ? (parseFloat(destinationLng) || 106.8456) : (parseFloat(meetingLng) || 106.8456),
          same_meeting_location: sameMeetingLocation,
          location_notes: locationNotes.trim() || undefined,
          difficulty,
          visibility,
          max_participants: parsedMax,
          has_parking_fee: hasParkingFee,
          parking_fee: parsedParking,
          requires_boat: requiresBoat,
          boat_rental_fee: parsedBoat,
          additional_cost: additionalCost,
        }),
      });

      if (res.data?.id) {
        router.push(`/events/${res.data.id}`);
      } else {
        router.push("/events");
      }
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Gagal membuat agenda event.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Navigation */}
      <Link
        href="/events"
        className="inline-flex items-center gap-1.5 text-sm font-semibold text-brand-secondary hover:underline mb-6"
      >
        <ArrowLeft className="w-4 h-4" /> Kembali ke Daftar Janjian
      </Link>

      <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-10 shadow-sm">
        {/* Header */}
        <div className="mb-8">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand-secondary/10 text-brand-secondary text-xs font-semibold mb-2">
            <Calendar className="w-4 h-4" /> Janjian Mancing Bersama
          </div>
          <h1 className="text-2xl sm:text-3xl font-extrabold text-brand-primary">
            Buat Agenda Mancing Baru
          </h1>
          <p className="text-sm text-brand-muted mt-1">
            Ajak komunitas pemancing berkumpul dan menjelajahi spot idaman bersama
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
          {/* Nama Event */}
          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Nama Event / Agenda *
            </label>
            <input
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Contoh: Mancing Bareng Nila Babon Waduk Cirata"
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          {/* Deskripsi */}
          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Deskripsi Rencana & Rundown *
            </label>
            <textarea
              rows={4}
              required
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Tuliskan gambaran acara, target tangkapan, persiapan peralatan, dan aturan kenyamanan bersama..."
              className="w-full p-4 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          {/* Lokasi Destinasi Mancing */}
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Nama Destinasi / Spot Mancing *
              </label>
              <input
                type="text"
                required
                value={destinationName}
                onChange={(e) => setDestinationName(e.target.value)}
                placeholder="Contoh: Dermaga Apung Waduk Cirata"
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              />
            </div>

            <LocationMapPicker
              label="Pilih Lokasi Mancing di Peta Google Maps"
              address={destinationAddress}
              latitude={destinationLat}
              longitude={destinationLng}
              onAddressChange={(newAddr) => {
                setDestinationAddress(newAddr);
                if (!destinationName) {
                  setDestinationName(newAddr.split(",")[0] || newAddr);
                }
              }}
              onCoordinateChange={(lat, lng) => {
                setDestinationLat(lat);
                setDestinationLng(lng);
              }}
            />
          </div>

          {/* Opsi Lokasi Janjian / Meeting Point */}
          <div className="p-5 rounded-2xl bg-brand-background border border-brand-border space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-sm font-bold text-brand-primary">Titik Temu / Lokasi Janjian</h3>
                <p className="text-xs text-brand-muted">Tentukan tempat kumpul sebelum berangkat atau langsung di spot</p>
              </div>
              <label className="flex items-center gap-2 cursor-pointer select-none text-xs font-semibold text-brand-primary">
                <input
                  type="checkbox"
                  checked={sameMeetingLocation}
                  onChange={(e) => {
                    const checked = e.target.checked;
                    setSameMeetingLocation(checked);
                    if (checked) {
                      setMeetingPoint(destinationName || "Langsung di Lokasi Mancing");
                      setMeetingAddress(destinationAddress);
                      setMeetingLat(destinationLat);
                      setMeetingLng(destinationLng);
                    }
                  }}
                  className="w-4 h-4 rounded border-brand-border text-brand-secondary focus:ring-brand-secondary"
                />
                <span>Janjian di lokasi mancing yang sama</span>
              </label>
            </div>

            {!sameMeetingLocation ? (
              <div className="space-y-4 pt-2 border-t border-brand-border/60">
                <div>
                  <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                    Nama Titik Temu / Titik Kumpul *
                  </label>
                  <input
                    type="text"
                    required={!sameMeetingLocation}
                    value={meetingPoint}
                    onChange={(e) => setMeetingPoint(e.target.value)}
                    placeholder="Contoh: Parkiran Indomaret Rest Area KM 57"
                    className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-surface text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
                  />
                </div>

                <LocationMapPicker
                  label="Pilih Titik Temu di Peta Google Maps"
                  address={meetingAddress}
                  latitude={meetingLat}
                  longitude={meetingLng}
                  onAddressChange={(newAddr) => {
                    setMeetingAddress(newAddr);
                    if (!meetingPoint) {
                      setMeetingPoint(newAddr.split(",")[0] || newAddr);
                    }
                  }}
                  onCoordinateChange={(lat, lng) => {
                    setMeetingLat(lat);
                    setMeetingLng(lng);
                  }}
                />
              </div>
            ) : (
              <div className="p-3 bg-brand-surface rounded-xl border border-brand-border text-xs text-brand-muted flex items-center gap-2">
                <span className="font-semibold text-brand-secondary">✓ Titik kumpul:</span>
                <span>Peserta akan langsung bertemu di lokasi mancing ({destinationName || destinationAddress || "Sesuai titik peta destinasi"}).</span>
              </div>
            )}

            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
                Petunjuk / Catatan Tambahan Titik Kumpul (Opsional)
              </label>
              <textarea
                rows={2}
                value={locationNotes}
                onChange={(e) => setLocationNotes(e.target.value)}
                placeholder="Contoh: Kumpul di samping musholla dekat gerbang timur. Kontak HP koordinator: 0812xxxx"
                className="w-full p-3 rounded-xl border border-brand-border bg-brand-surface text-brand-text text-xs focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              />
            </div>
          </div>

          {/* Tanggal & Waktu */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Tanggal Pelaksanaan *
              </label>
              <input
                type="date"
                required
                value={eventDate}
                onChange={(e) => setEventDate(e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              />
            </div>

            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Jam Mulai *
              </label>
              <input
                type="time"
                required
                value={startTime}
                onChange={(e) => setStartTime(e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              />
            </div>

            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Jam Selesai (Opsional)
              </label>
              <input
                type="time"
                value={endTime}
                onChange={(e) => setEndTime(e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              />
            </div>
          </div>

          {/* Pengaturan Event */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-4 border-t border-brand-border">
            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Tingkat Kesulitan
              </label>
              <select
                value={difficulty}
                onChange={(e) => setDifficulty(e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              >
                <option value="pemula">Pemula (Santai)</option>
                <option value="menengah">Menengah (Tantangan)</option>
                <option value="ahli">Ahli (Ekstrem)</option>
              </select>
            </div>

            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Tipe Keanggotaan
              </label>
              <select
                value={visibility}
                onChange={(e) => setVisibility(e.target.value as "public" | "private")}
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              >
                <option value="public">Public (Langsung Gabung)</option>
                <option value="private">Private (Perlu Persetujuan Host)</option>
              </select>
            </div>

            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Maksimal Kuota Peserta
              </label>
              <input
                type="number"
                min={1}
                value={maxParticipants}
                onChange={(e) => setMaxParticipants(e.target.value)}
                placeholder="10"
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              />
            </div>
          </div>

          {/* Transparansi Biaya */}
          <div className="p-5 rounded-2xl bg-brand-background border border-brand-border space-y-4">
            <h3 className="text-xs font-bold uppercase tracking-wider text-brand-primary">
              Estimasi & Transparansi Biaya Bersama
            </h3>

            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
              <label className="flex items-center gap-2 text-sm text-brand-text font-medium cursor-pointer">
                <input
                  type="checkbox"
                  checked={hasParkingFee}
                  onChange={(e) => setHasParkingFee(e.target.checked)}
                  className="rounded border-brand-border text-brand-secondary focus:ring-brand-secondary"
                />
                <Car className="w-4 h-4 text-brand-secondary" /> Ada Biaya Parkir?
              </label>
              {hasParkingFee && (
                <input
                  type="number"
                  value={parkingFee}
                  onChange={(e) => setParkingFee(e.target.value)}
                  placeholder="Rp 10.000"
                  className="px-3 py-1.5 rounded-lg border border-brand-border bg-brand-surface text-xs w-40"
                />
              )}
            </div>

            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
              <label className="flex items-center gap-2 text-sm text-brand-text font-medium cursor-pointer">
                <input
                  type="checkbox"
                  checked={requiresBoat}
                  onChange={(e) => setRequiresBoat(e.target.checked)}
                  className="rounded border-brand-border text-brand-secondary focus:ring-brand-secondary"
                />
                <Ship className="w-4 h-4 text-brand-secondary" /> Perlu Sewa Perahu?
              </label>
              {requiresBoat && (
                <input
                  type="number"
                  value={boatRentalFee}
                  onChange={(e) => setBoatRentalFee(e.target.value)}
                  placeholder="Rp 150.000"
                  className="px-3 py-1.5 rounded-lg border border-brand-border bg-brand-surface text-xs w-40"
                />
              )}
            </div>

            <div>
              <label className="block text-xs font-medium text-brand-muted mb-1">
                Catatan Biaya Tambahan
              </label>
              <input
                type="text"
                value={additionalCost}
                onChange={(e) => setAdditionalCost(e.target.value)}
                placeholder="Contoh: Iuran makan siang bersama, biaya umpan hidup"
                className="w-full px-3 py-2 rounded-lg border border-brand-border bg-brand-surface text-xs"
              />
            </div>
          </div>

          {/* Submit Button */}
          <button
            type="submit"
            disabled={loading}
            className="mt-4 w-full py-3.5 px-4 rounded-xl font-bold text-white bg-brand-primary hover:bg-brand-primary/90 focus:outline-none focus:ring-2 focus:ring-brand-primary shadow-sm disabled:opacity-50 transition flex items-center justify-center gap-2"
          >
            {loading ? (
              <span>Mempublikasikan agenda...</span>
            ) : (
              <>
                <Plus className="w-4 h-4 text-brand-secondary" />
                <span>Publikasikan Janjian Mancing</span>
              </>
            )}
          </button>
        </form>
      </div>
    </div>
  );
}
