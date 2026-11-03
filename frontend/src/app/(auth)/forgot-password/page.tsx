"use client";

import { useState } from "react";
import Link from "next/link";
import { KeyRound, CheckCircle2, AlertCircle, Loader2, ArrowRight } from "lucide-react";
import { fetchApi } from "@/lib/api-client";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const cleanEmail = email.trim().toLowerCase();
    if (!cleanEmail) {
      setError("Silakan masukkan alamat email akun Anda.");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      await fetchApi("/auth/forgot-password", {
        method: "POST",
        body: JSON.stringify({ email: cleanEmail }),
      });
      setSubmitted(true);
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Gagal mengirimkan instruksi reset password.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md bg-brand-surface rounded-3xl border border-brand-border p-8 sm:p-10 shadow-sm">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="inline-flex w-14 h-14 rounded-2xl bg-amber-50 border border-amber-200 items-center justify-center text-amber-600 mb-3.5 shadow-xs">
            <KeyRound className="w-7 h-7" />
          </div>
          <h1 className="text-2xl font-extrabold text-brand-primary">Lupa Kata Sandi?</h1>
          <p className="text-xs sm:text-sm text-brand-muted mt-1.5">
            Masukkan email terdaftar Anda untuk menerima kode verifikasi atur ulang password.
          </p>
        </div>

        {error && (
          <div className="mb-6 p-4 rounded-2xl bg-red-50 border border-red-200 text-red-800 text-xs sm:text-sm flex items-start gap-2.5">
            <AlertCircle className="w-5 h-5 shrink-0 text-red-600 mt-0.5" />
            <span className="font-medium leading-relaxed">{error}</span>
          </div>
        )}

        {submitted ? (
          <div className="space-y-6">
            <div className="p-4 rounded-2xl bg-green-50 border border-green-200 text-green-800 text-xs sm:text-sm flex items-start gap-2.5">
              <CheckCircle2 className="w-5 h-5 shrink-0 text-green-600 mt-0.5" />
              <div>
                <p className="font-bold text-green-900 mb-1">Email Berhasil Dikirim!</p>
                <p className="leading-relaxed text-green-800/90">
                  Jika email <strong>{email}</strong> terdaftar di sistem, instruksi beserta kode verifikasi reset kata sandi telah kami kirimkan.
                </p>
              </div>
            </div>

            <Link
              href={`/reset-password?email=${encodeURIComponent(email)}`}
              className="w-full py-3 px-4 rounded-xl bg-brand-secondary hover:bg-brand-secondary/90 text-white font-bold text-sm transition shadow-sm flex items-center justify-center gap-2"
            >
              <span>Masukkan Kode OTP & Buat Password Baru</span>
              <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
                Alamat Email Akun
              </label>
              <input
                type="email"
                required
                autoFocus
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="nama@email.com"
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary transition"
              />
            </div>

            <button
              type="submit"
              disabled={loading || !email.trim()}
              className="w-full mt-2 py-3 px-4 rounded-xl bg-brand-secondary hover:bg-brand-secondary/90 text-white font-bold text-sm transition shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
            >
              {loading ? (
                <Loader2 className="w-5 h-5 animate-spin" />
              ) : (
                <>
                  <span>Kirim Kode Reset Password</span>
                  <ArrowRight className="w-4 h-4" />
                </>
              )}
            </button>
          </form>
        )}

        <div className="mt-8 pt-6 border-t border-brand-border/60 text-center">
          <Link href="/login" className="text-xs font-bold text-brand-secondary hover:underline">
            ← Kembali ke Halaman Masuk
          </Link>
        </div>
      </div>
    </div>
  );
}
