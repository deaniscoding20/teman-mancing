"use client";

import { useState, useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { Mail, CheckCircle2, AlertCircle, Loader2, ArrowRight, RefreshCw, Fish } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { setAuth } from "@/lib/auth";

function VerifyEmailContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const tokenParam = searchParams.get("token") || "";
  const emailParam = searchParams.get("email") || "";

  const [email, setEmail] = useState(emailParam);
  const [otp, setOtp] = useState("");
  const [loading, setLoading] = useState(false);
  const [resending, setResending] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [resendCooldown, setResendCooldown] = useState(0);

  // Auto-verify if token is provided in URL query
  useEffect(() => {
    if (tokenParam) {
      handleAutoVerify(tokenParam);
    }
  }, [tokenParam]);

  useEffect(() => {
    if (resendCooldown <= 0) return;
    const interval = setInterval(() => {
      setResendCooldown((prev) => prev - 1);
    }, 1000);
    return () => clearInterval(interval);
  }, [resendCooldown]);

  const handleAutoVerify = async (token: string) => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetchApi<{ token: string; user: any }>("/auth/verify-email", {
        method: "POST",
        body: JSON.stringify({ token }),
      });

      if (res.data?.token && res.data?.user) {
        setAuth(res.data.token, res.data.user);
      }
      setSuccess("Email Anda berhasil diverifikasi! Mengalihkan ke dashboard...");
      setTimeout(() => {
        router.push("/");
      }, 1500);
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Tautan verifikasi tidak valid atau telah kedaluwarsa.");
    } finally {
      setLoading(false);
    }
  };

  const handleManualVerify = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!otp.trim()) {
      setError("Silakan masukkan 6 digit kode OTP yang dikirim ke email Anda.");
      return;
    }

    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const res = await fetchApi<{ token: string; user: any }>("/auth/verify-email", {
        method: "POST",
        body: JSON.stringify({
          email: email.trim().toLowerCase(),
          otp: otp.trim(),
        }),
      });

      if (res.data?.token && res.data?.user) {
        setAuth(res.data.token, res.data.user);
      }
      setSuccess("Email Anda berhasil diverifikasi! Mengalihkan...");
      setTimeout(() => {
        router.push("/");
      }, 1500);
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Kode OTP salah atau telah kedaluwarsa.");
    } finally {
      setLoading(false);
    }
  };

  const handleResendOTP = async () => {
    if (!email.trim()) {
      setError("Masukkan email Anda terlebih dahulu untuk mengirim ulang kode.");
      return;
    }
    if (resendCooldown > 0) return;

    setResending(true);
    setError(null);

    try {
      await fetchApi("/auth/resend-verification", {
        method: "POST",
        body: JSON.stringify({ email: email.trim().toLowerCase() }),
      });
      setSuccess("Kode verifikasi baru telah berhasil dikirimkan ke email Anda.");
      setResendCooldown(60);
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Gagal mengirim ulang kode verifikasi.");
    } finally {
      setResending(false);
    }
  };

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md bg-brand-surface rounded-3xl border border-brand-border p-8 sm:p-10 shadow-sm">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="inline-flex w-14 h-14 rounded-2xl bg-sky-50 border border-sky-200 items-center justify-center text-sky-600 mb-3.5 shadow-xs">
            <Mail className="w-7 h-7" />
          </div>
          <h1 className="text-2xl font-extrabold text-brand-primary">Verifikasi Email</h1>
          <p className="text-xs sm:text-sm text-brand-muted mt-1.5">
            Kami telah mengirimkan 6 digit kode OTP verifikasi ke email Anda.
          </p>
        </div>

        {/* Alert Messages */}
        {error && (
          <div className="mb-6 p-4 rounded-2xl bg-red-50 border border-red-200 text-red-800 text-xs sm:text-sm flex items-start gap-2.5">
            <AlertCircle className="w-5 h-5 shrink-0 text-red-600 mt-0.5" />
            <span className="font-medium leading-relaxed">{error}</span>
          </div>
        )}

        {success && (
          <div className="mb-6 p-4 rounded-2xl bg-green-50 border border-green-200 text-green-800 text-xs sm:text-sm flex items-start gap-2.5">
            <CheckCircle2 className="w-5 h-5 shrink-0 text-green-600 mt-0.5" />
            <span className="font-semibold leading-relaxed">{success}</span>
          </div>
        )}

        {tokenParam && loading ? (
          <div className="text-center py-8">
            <Loader2 className="w-8 h-8 text-brand-secondary animate-spin mx-auto mb-3" />
            <p className="text-sm font-semibold text-brand-primary">Memverifikasi tautan token Anda...</p>
          </div>
        ) : (
          <form onSubmit={handleManualVerify} className="flex flex-col gap-4">
            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
                Alamat Email
              </label>
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="nama@email.com"
                className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary transition"
              />
            </div>

            <div>
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
                Kode Verifikasi (OTP 6 Digit)
              </label>
              <input
                type="text"
                required
                maxLength={6}
                value={otp}
                onChange={(e) => setOtp(e.target.value.replace(/\D/g, ""))}
                placeholder="123456"
                className="w-full px-4 py-3 rounded-xl border border-brand-border bg-brand-background text-brand-text text-center text-2xl font-mono font-bold tracking-widest placeholder-brand-muted/40 focus:outline-none focus:ring-2 focus:ring-brand-secondary transition"
              />
            </div>

            <button
              type="submit"
              disabled={loading || !otp.trim()}
              className="w-full mt-2 py-3 px-4 rounded-xl bg-brand-secondary hover:bg-brand-secondary/90 text-white font-bold text-sm transition shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
            >
              {loading ? (
                <Loader2 className="w-5 h-5 animate-spin" />
              ) : (
                <>
                  <span>Verifikasi Akun</span>
                  <ArrowRight className="w-4 h-4" />
                </>
              )}
            </button>

            <div className="pt-2 text-center">
              <button
                type="button"
                onClick={handleResendOTP}
                disabled={resending || resendCooldown > 0}
                className="text-xs font-semibold text-brand-primary hover:text-brand-secondary transition inline-flex items-center gap-1.5 disabled:opacity-50"
              >
                {resending ? (
                  <Loader2 className="w-3.5 h-3.5 animate-spin" />
                ) : (
                  <RefreshCw className="w-3.5 h-3.5" />
                )}
                <span>
                  {resendCooldown > 0
                    ? `Kirim ulang OTP (${resendCooldown}s)`
                    : "Kirim Ulang Kode OTP"}
                </span>
              </button>
            </div>
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

export default function VerifyEmailPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-[80vh] flex items-center justify-center">
          <Loader2 className="w-8 h-8 animate-spin text-brand-secondary" />
        </div>
      }
    >
      <VerifyEmailContent />
    </Suspense>
  );
}
