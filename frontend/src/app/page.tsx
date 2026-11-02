import Link from "next/link";
import { Compass, Calendar, Users, MapPin, ArrowRight, ShieldCheck, Waves } from "lucide-react";

export default function HomePage() {
  return (
    <div className="flex flex-col gap-16 py-10">
      {/* Hero Section */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full text-center flex flex-col items-center pt-8 sm:pt-14 pb-8">
        <div className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-brand-secondary/15 text-brand-primary text-xs font-semibold mb-6">
          <Waves className="w-4 h-4 text-brand-secondary" />
          Komunitas Pemancing Terbesar di Indonesia
        </div>

        <h1 className="text-4xl sm:text-5xl lg:text-6xl font-extrabold text-brand-primary tracking-tight max-w-4xl leading-[1.15]">
          Cari Spot. Cari Teman. <br className="hidden sm:inline" />
          <span className="text-brand-secondary">Gas Mancing!</span>
        </h1>

        <p className="mt-5 text-lg sm:text-xl text-brand-muted max-w-2xl font-normal leading-relaxed">
          Temukan titik spot memancing rahasia, ikuti agenda janjian mancing bersama rekan komunitas, dan nikmati serunya petualangan mancing tanpa ribet.
        </p>

        <div className="mt-8 flex flex-col sm:flex-row items-center gap-3.5 w-full sm:w-auto">
          <Link
            href="/spots"
            className="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-7 py-3.5 rounded-xl font-semibold text-white bg-brand-primary hover:bg-brand-primary/90 shadow-md transition"
          >
            <Compass className="w-5 h-5 text-brand-secondary" />
            Jelajahi Spot Mancing
          </Link>
          <Link
            href="/events"
            className="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-7 py-3.5 rounded-xl font-semibold text-brand-primary bg-brand-surface border border-brand-border hover:bg-brand-border/40 shadow-sm transition"
          >
            <Calendar className="w-5 h-5 text-brand-secondary" />
            Ikuti Janjian Mancing
          </Link>
        </div>
      </section>

      {/* 3 Core Pillars */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 sm:gap-8">
          {/* Card 1 */}
          <div className="bg-brand-surface rounded-2xl p-7 border border-brand-border shadow-sm flex flex-col hover:border-brand-secondary/40 transition">
            <div className="w-12 h-12 rounded-xl bg-brand-secondary/10 flex items-center justify-center text-brand-secondary mb-5">
              <MapPin className="w-6 h-6" />
            </div>
            <h3 className="text-xl font-bold text-brand-primary mb-2">Spot Mancing Interaktif</h3>
            <p className="text-brand-muted text-sm leading-relaxed mb-6 flex-1">
              Peta interaktif Leaflet OSM dengan info teknik memancing, jenis ikan target, rekomendasi racikan umpan, dan ulasan sesama angler.
            </p>
            <Link
              href="/spots"
              className="inline-flex items-center gap-1.5 text-sm font-semibold text-brand-secondary hover:underline mt-auto"
            >
              Lihat Spot Terdekat <ArrowRight className="w-4 h-4" />
            </Link>
          </div>

          {/* Card 2 */}
          <div className="bg-brand-surface rounded-2xl p-7 border border-brand-border shadow-sm flex flex-col hover:border-brand-secondary/40 transition">
            <div className="w-12 h-12 rounded-xl bg-brand-accent/15 flex items-center justify-center text-brand-accent mb-5">
              <Calendar className="w-6 h-6 text-brand-accent" />
            </div>
            <h3 className="text-xl font-bold text-brand-primary mb-2">Janjian Mancing Bareng</h3>
            <p className="text-brand-muted text-sm leading-relaxed mb-6 flex-1">
              Buat agenda trip mancing bersama (public atau private). Dilengkapi transparansi biaya sewa kapal, parkir, dan kuota peserta aman.
            </p>
            <Link
              href="/events"
              className="inline-flex items-center gap-1.5 text-sm font-semibold text-brand-secondary hover:underline mt-auto"
            >
              Cari Agenda Event <ArrowRight className="w-4 h-4" />
            </Link>
          </div>

          {/* Card 3 */}
          <div className="bg-brand-surface rounded-2xl p-7 border border-brand-border shadow-sm flex flex-col hover:border-brand-secondary/40 transition">
            <div className="w-12 h-12 rounded-xl bg-brand-primary/10 flex items-center justify-center text-brand-primary mb-5">
              <Users className="w-6 h-6" />
            </div>
            <h3 className="text-xl font-bold text-brand-primary mb-2">Realtime Chat & Komunitas</h3>
            <p className="text-brand-muted text-sm leading-relaxed mb-6 flex-1">
              Otomatis terhubung ke grup chat room untuk setiap event yang kamu ikuti, plus pesan langsung personal chat dengan teman mancing barumu.
            </p>
            <Link
              href="/chat"
              className="inline-flex items-center gap-1.5 text-sm font-semibold text-brand-secondary hover:underline mt-auto"
            >
              Buka Chat Room <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
        </div>
      </section>

      {/* Safety & Community Banner */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full">
        <div className="bg-gradient-to-r from-brand-primary to-[#1c4d66] rounded-3xl p-8 sm:p-12 text-white flex flex-col md:flex-row items-center justify-between gap-8 shadow-lg">
          <div className="flex flex-col gap-3 max-w-xl text-center md:text-left">
            <div className="inline-flex items-center justify-center md:justify-start gap-2 text-brand-secondary text-sm font-semibold">
              <ShieldCheck className="w-5 h-5" /> Terverifikasi & Ramah Pemula
            </div>
            <h2 className="text-2xl sm:text-3xl font-bold">Siap Menemukan Rekan Mancing Pertamamu?</h2>
            <p className="text-white/80 text-sm sm:text-base">
              Daftar gratis sekarang dan bergabunglah dengan ribuan pemancing dari seluruh penjuru Indonesia.
            </p>
          </div>
          <Link
            href="/register"
            className="px-8 py-4 rounded-xl font-bold text-brand-primary bg-white hover:bg-white/95 shadow-md transition shrink-0"
          >
            Gabung Sekarang
          </Link>
        </div>
      </section>
    </div>
  );
}
