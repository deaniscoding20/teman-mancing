"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { AlertTriangle, Mail, RefreshCw, X, CheckCircle2 } from "lucide-react";
import { getCurrentUser } from "@/lib/auth";
import { fetchApi } from "@/lib/api-client";

export default function EmailVerificationBanner() {
  const [user, setUser] = useState<any>(null);
  const [dismissed, setDismissed] = useState(false);
  const [resending, setResending] = useState(false);
  const [resendStatus, setResendStatus] = useState<string | null>(null);

  useEffect(() => {
    const update = () => {
      const currentUser = getCurrentUser();
      setUser(currentUser);
    };
    update();
    window.addEventListener("auth-change", update);
    return () => window.removeEventListener("auth-change", update);
  }, []);

  if (!user || user.is_verified || dismissed) {
    return null;
  }

  const handleResend = async () => {
    setResending(true);
    setResendStatus(null);
    try {
      await fetchApi("/auth/resend-verification", {
        method: "POST",
        body: JSON.stringify({ email: user.email }),
      });
      setResendStatus("Kode OTP berhasil dikirim ulang ke email Anda!");
      setTimeout(() => setResendStatus(null), 6000);
    } catch (err: any) {
      setResendStatus(err?.data?.message || err?.message || "Gagal mengirim ulang email.");
      setTimeout(() => setResendStatus(null), 6000);
    } finally {
      setResending(false);
    }
  };

  return (
    <div className="bg-gradient-to-r from-amber-500 via-amber-600 to-orange-500 text-white shadow-sm transition-all duration-200">
      <div className="max-w-7xl mx-auto px-4 py-2.5 sm:px-6 lg:px-8 flex flex-col sm:flex-row items-center justify-between gap-2.5 text-xs sm:text-sm">
        <div className="flex items-center gap-2.5 text-center sm:text-left">
          <AlertTriangle className="w-5 h-5 shrink-0 text-amber-100 animate-pulse" />
          <span>
            <strong>Verifikasi Email Diperlukan:</strong> Cek inbox/spam di{" "}
            <span className="font-semibold underline decoration-white/50">{user.email}</span>{" "}
            untuk mengaktifkan fitur membuat spot, janjian mancing, dan kirim pesan.
          </span>
        </div>

        <div className="flex items-center gap-2 shrink-0">
          {resendStatus ? (
            <span className="bg-white/20 backdrop-blur-xs px-3 py-1 rounded-full text-xs font-semibold flex items-center gap-1.5">
              <CheckCircle2 className="w-3.5 h-3.5 text-green-200" />
              {resendStatus}
            </span>
          ) : (
            <button
              onClick={handleResend}
              disabled={resending}
              className="px-2.5 py-1 rounded-lg bg-white/15 hover:bg-white/25 text-white font-medium text-xs transition flex items-center gap-1 disabled:opacity-50"
            >
              <RefreshCw className={`w-3 h-3 ${resending ? "animate-spin" : ""}`} />
              <span>{resending ? "Mengirim..." : "Kirim Ulang"}</span>
            </button>
          )}

          <Link
            href={`/verify-email?email=${encodeURIComponent(user.email)}`}
            className="px-3.5 py-1 rounded-lg bg-white text-amber-800 font-bold text-xs hover:bg-neutral-100 transition shadow-xs"
          >
            Verifikasi Sekarang
          </Link>

          <button
            onClick={() => setDismissed(true)}
            className="p-1 rounded-lg text-white/80 hover:text-white hover:bg-white/10 transition"
            title="Tutup pemberitahuan"
          >
            <X className="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  );
}
