"use client";

import { useState, useRef, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  Fish,
  LogIn,
  AlertCircle,
  Eye,
  EyeOff,
  HelpCircle,
  CheckCircle2,
  Lock,
  User,
} from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { setAuth, getToken, getStoredUser } from "@/lib/auth";
import GoogleSignInButton from "@/components/auth/GoogleSignInButton";

export default function LoginPage() {
  const router = useRouter();

  // If user is already authenticated, redirect to home immediately
  useEffect(() => {
    const token = getToken();
    const user = getStoredUser();
    if (token && user) {
      router.replace("/");
    }
  }, [router]);
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [capsLockActive, setCapsLockActive] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [failedAttempts, setFailedAttempts] = useState(0);

  const passwordInputRef = useRef<HTMLInputElement>(null);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.getModifierState) {
      setCapsLockActive(e.getModifierState("CapsLock"));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const cleanIdent = identifier.trim();
    if (!cleanIdent) {
      setError("Silakan masukkan email atau username Anda.");
      return;
    }
    if (!password) {
      setError("Silakan masukkan password Anda.");
      return;
    }

    setLoading(true);

    try {
      const res = await fetchApi<{ token: string; user: any }>("/auth/login", {
        method: "POST",
        body: JSON.stringify({
          identifier: cleanIdent,
          email: cleanIdent.includes("@") ? cleanIdent : undefined,
          username: !cleanIdent.includes("@") ? cleanIdent : undefined,
          password,
        }),
      });

      if (res.data?.token && res.data?.user) {
        setAuth(res.data.token, res.data.user);
        router.push("/");
      } else {
        throw new Error("Token otorisasi tidak ditemukan.");
      }
    } catch (err: any) {
      setFailedAttempts((prev) => prev + 1);
      const backendMsg = err?.data?.message || err?.message;
      if (backendMsg && backendMsg.toLowerCase().includes("salah")) {
        setError("Email/username atau password yang Anda masukkan salah.");
      } else {
        setError(backendMsg || "Gagal masuk. Periksa kembali koneksi atau kredensial Anda.");
      }

      // Best practice: keep identifier, clear password and auto-focus for easy retry
      setPassword("");
      setTimeout(() => {
        passwordInputRef.current?.focus();
      }, 50);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md bg-brand-surface rounded-3xl border border-brand-border p-8 sm:p-10 shadow-sm">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="inline-flex w-14 h-14 rounded-2xl bg-brand-primary items-center justify-center text-white mb-3.5 shadow-sm">
            <Fish className="w-8 h-8 text-brand-secondary" />
          </div>
          <h1 className="text-2xl sm:text-3xl font-extrabold text-brand-primary">
            Masuk ke Teman Mancing
          </h1>
          <p className="text-xs sm:text-sm text-brand-muted mt-1.5">
            Cari Spot. Cari Teman. Gas Mancing!
          </p>
        </div>

        {/* Error notification with best-practice guidance */}
        {error && (
          <div className="mb-6 p-4 rounded-2xl bg-red-50 border border-red-200 text-red-800 text-xs sm:text-sm flex flex-col gap-2">
            <div className="flex items-start gap-2.5">
              <AlertCircle className="w-5 h-5 shrink-0 text-red-600 mt-0.5" />
              <span className="font-semibold leading-relaxed">{error}</span>
            </div>

            {error.toLowerCase().includes("verifikasi") && (
              <div className="mt-2 pt-2 border-t border-red-200">
                <Link
                  href={`/verify-email?email=${encodeURIComponent(identifier)}`}
                  className="w-full py-2 px-3 rounded-xl bg-brand-primary text-white text-xs font-bold hover:bg-brand-primary/90 transition flex items-center justify-center gap-1.5 shadow-xs"
                >
                  <span>Buka Halaman Verifikasi Email</span>
                  <Fish className="w-3.5 h-3.5 text-brand-secondary" />
                </Link>
              </div>
            )}

            {failedAttempts >= 1 && !error.toLowerCase().includes("verifikasi") && (
              <div className="mt-1 pt-2 border-t border-red-200/80 text-[11px] text-red-700 space-y-1">
                <p className="font-medium">Tips jika masih gagal masuk:</p>
                <ul className="list-disc list-inside space-y-0.5 text-red-600/90 pl-1">
                  <li>Pastikan penulisan email atau username sudah benar.</li>
                  <li>Periksa huruf besar/kecil (password bersifat *case-sensitive*).</li>
                  <li>Gunakan tombol ikon mata untuk melihat password yang Anda ketik.</li>
                </ul>
              </div>
            )}
          </div>
        )}

        {/* Form */}
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
              Email atau Username
            </label>
            <div className="relative">
              <input
                type="text"
                required
                autoFocus
                autoComplete="username"
                value={identifier}
                onChange={(e) => setIdentifier(e.target.value)}
                placeholder="nama@email.com atau username"
                className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary text-sm transition"
              />
              <User className="w-4 h-4 text-brand-muted absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            </div>
          </div>

          <div>
            <div className="flex items-center justify-between mb-1.5">
              <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider">
                Password
              </label>
              {capsLockActive && (
                <span className="text-[10px] font-bold text-amber-600 bg-amber-50 px-2 py-0.5 rounded border border-amber-200 animate-pulse">
                  Caps Lock Aktif
                </span>
              )}
            </div>

            <div className="relative">
              <input
                ref={passwordInputRef}
                type={showPassword ? "text" : "password"}
                required
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                onKeyDown={handleKeyDown}
                onKeyUp={handleKeyDown}
                placeholder="••••••••"
                className={`w-full pl-10 pr-10 py-2.5 rounded-xl border bg-brand-background text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary text-sm transition ${
                  error ? "border-red-300 ring-1 ring-red-300" : "border-brand-border"
                }`}
              />
              <Lock className="w-4 h-4 text-brand-muted absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                tabIndex={-1}
                className="absolute right-3.5 top-1/2 -translate-y-1/2 text-brand-muted hover:text-brand-primary transition"
                title={showPassword ? "Sembunyikan password" : "Tampilkan password"}
              >
                {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
              </button>
            </div>
            <div className="flex justify-end mt-1.5">
              <Link
                href="/forgot-password"
                className="text-xs font-semibold text-brand-secondary hover:underline"
              >
                Lupa kata sandi?
              </Link>
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="mt-3 w-full py-3.5 px-4 rounded-xl font-bold text-white bg-brand-primary hover:bg-brand-primary/90 focus:outline-none focus:ring-2 focus:ring-brand-primary shadow-sm disabled:opacity-50 transition flex items-center justify-center gap-2 text-sm"
          >
            {loading ? (
              <span className="flex items-center gap-2">
                <span className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                Memverifikasi kredensial...
              </span>
            ) : (
              <>
                <LogIn className="w-4 h-4 text-brand-secondary" />
                <span>Masuk Sekarang</span>
              </>
            )}
          </button>
        </form>

        {/* Divider */}
        <div className="relative my-6">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-brand-border" />
          </div>
          <div className="relative flex justify-center text-xs">
            <span className="px-3 bg-brand-surface text-brand-muted font-medium">atau masuk dengan</span>
          </div>
        </div>

        {/* Google Sign In */}
        <GoogleSignInButton
          label="Masuk dengan Google"
          onError={(err) => setError(err)}
        />

        {/* Footer Navigation */}
        <div className="mt-8 pt-6 border-t border-brand-border/80 text-center text-xs text-brand-muted space-y-2">
          <p>
            Belum punya akun?{" "}
            <Link href="/register" className="font-bold text-brand-secondary hover:underline">
              Daftar akun baru di sini
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
