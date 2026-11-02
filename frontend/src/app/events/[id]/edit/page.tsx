"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { Calendar, Save, ArrowLeft, AlertCircle, Car, Ship } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { getStoredUser, getToken } from "@/lib/auth";

export default function EditEventPage() {
  const params = useParams();
  const router = useRouter();
  const eventId = params?.id as string;

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [destinationName, setDestinationName] = useState("");
  const [destinationAddress, setDestinationAddress] = useState("");
  const [meetingPoint, setMeetingPoint] = useState("");
  const [difficulty, setDifficulty] = useState("pemula");
  const [maxParticipants, setMaxParticipants] = useState("10");

  const [hasParkingFee, setHasParkingFee] = useState(false);
  const [parkingFee, setParkingFee] = useState("10000");
  const [requiresBoat, setRequiresBoat] = useState(false);
  const [boatRentalFee, setBoatRentalFee] = useState("150000");
  const [additionalCost, setAdditionalCost] = useState("");
  const [status, setStatus] = useState("open");

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const currentUser = getStoredUser();
  const token = getToken();

  useEffect(() => {
    if (!eventId) return;

    fetchApi<any>(`/events/${eventId}`)
      .then((res) => {
        if (res.data?.event) {
          const ev = res.data.event;
          if (currentUser && currentUser.id !== ev.creator_id) {
            router.push(`/events/${eventId}`);
            return;
          }
          setName(ev.name || "");
          setDescription(ev.description || "");
          setDestinationName(ev.destination_name || "");
          setDestinationAddress(ev.destination_address || "");
          setMeetingPoint(ev.meeting_point || "");
          setDifficulty(ev.difficulty || "pemula");
          setMaxParticipants(ev.max_participants ? ev.max_participants.toString() : "");
          setHasParkingFee(ev.has_parking_fee || false);
          setParkingFee(ev.parking_fee ? ev.parking_fee.toString() : "10000");
          setRequiresBoat(ev.requires_boat || false);
          setBoatRentalFee(ev.boat_rental_fee ? ev.boat_rental_fee.toString() : "150000");
          setAdditionalCost(ev.additional_cost || "");
          setStatus(ev.status || "open");
        }
      })
      .catch((err) => {
        setError(err.message || "Gagal memuat event.");
      })
      .finally(() => {
        setLoading(false);
      });
  }, [eventId, currentUser, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;

    const parsedMax = maxParticipants ? parseInt(maxParticipants, 10) : undefined;
    const parsedParking = hasParkingFee && parkingFee ? parseFloat(parkingFee) : undefined;
    const parsedBoat = requiresBoat && boatRentalFee ? parseFloat(boatRentalFee) : undefined;

    setSaving(true);
    setError(null);

    try {
      await fetchApi(`/events/${eventId}`, {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          name,
          description,
          destination_name: destinationName,
          destination_address: destinationAddress,
          meeting_point: meetingPoint,
          difficulty,
          max_participants: parsedMax,
          has_parking_fee: hasParkingFee,
          parking_fee: parsedParking,
          requires_boat: requiresBoat,
          boat_rental_fee: parsedBoat,
          additional_cost: additionalCost,
          status,
        }),
      });

      router.push(`/events/${eventId}`);
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Gagal memperbarui event.");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <div className="text-center py-20 text-brand-muted">Memuat agenda event...</div>;
  }

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <Link
        href={`/events/${eventId}`}
        className="inline-flex items-center gap-1.5 text-sm font-semibold text-brand-secondary hover:underline mb-6"
      >
        <ArrowLeft className="w-4 h-4" /> Kembali ke Detail Event
      </Link>

      <div className="bg-brand-surface rounded-3xl border border-brand-border p-6 sm:p-10 shadow-sm">
        <div className="mb-8">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand-secondary/10 text-brand-secondary text-xs font-semibold mb-2">
            <Calendar className="w-4 h-4" /> Manajemen Trip
          </div>
          <h1 className="text-2xl sm:text-3xl font-extrabold text-brand-primary">Edit Agenda Mancing</h1>
          <p className="text-sm text-brand-muted mt-1">Perbarui rincian trip, status, atau kuota peserta</p>
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
              Nama Event *
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
              Deskripsi Rundown *
            </label>
            <textarea
              rows={4}
              required
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full p-4 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Destinasi Tujuan *
              </label>
              <input
                type="text"
                required
                value={destinationName}
                onChange={(e) => setDestinationName(e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              />
            </div>
            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Titik Temu (Meeting Point) *
              </label>
              <input
                type="text"
                required
                value={meetingPoint}
                onChange={(e) => setMeetingPoint(e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              />
            </div>
          </div>

          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
              Alamat Lengkap Destinasi
            </label>
            <input
              type="text"
              value={destinationAddress}
              onChange={(e) => setDestinationAddress(e.target.value)}
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Tingkat Kesulitan
              </label>
              <select
                value={difficulty}
                onChange={(e) => setDifficulty(e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              >
                <option value="pemula">Pemula</option>
                <option value="menengah">Menengah</option>
                <option value="ahli">Ahli</option>
              </select>
            </div>

            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Status Event
              </label>
              <select
                value={status}
                onChange={(e) => setStatus(e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
              >
                <option value="open">Open (Terbuka)</option>
                <option value="closed">Closed (Ditutup)</option>
                <option value="completed">Completed (Selesai)</option>
                <option value="cancelled">Cancelled (Dibatalkan)</option>
              </select>
            </div>

            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-2">
                Batas Maksimal Kuota
              </label>
              <input
                type="number"
                min={1}
                value={maxParticipants}
                onChange={(e) => setMaxParticipants(e.target.value)}
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
            <span>{saving ? "Menyimpan..." : "Simpan Perubahan Event"}</span>
          </button>
        </form>
      </div>
    </div>
  );
}
