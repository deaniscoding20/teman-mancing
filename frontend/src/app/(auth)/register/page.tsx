"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  Fish,
  UserPlus,
  AlertCircle,
  CheckCircle2,
  Eye,
  EyeOff,
  Lock,
  User,
  Mail,
  AtSign,
  Check,
  X,
} from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { setAuth, getToken, getStoredUser } from "@/lib/auth";
import GoogleSignInButton from "@/components/auth/GoogleSignInButton";

export default function RegisterPage() {
  const router = useRouter();

  // If user is already authenticated, redirect to home immediately
  useEffect(() => {
    const token = getToken();
    const user = getStoredUser();
    if (token && user) {
      router.replace("/");
    }
  }, [router]);
  const [name, setName] = useState("");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  // Password verification checks
  const isLengthValid = password.length >= 8 && password.length <= 72;
  const hasLetter = /[a-zA-Z]/.test(password);
  const hasNumber = /[0-9]/.test(password);
  const isPasswordStrong = isLengthValid && hasLetter && hasNumber;
  const isPasswordMatch = confirmPassword.length > 0 && password === confirmPassword;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!isLengthValid) {
      setError("Password harus memiliki panjang minimal 8 karakter.");
      return;
    }

    if (password !== confirmPassword) {
      setError("Konfirmasi password tidak cocok dengan password yang Anda masukkan.");
      return;
    }

    if (username.length < 3) {
      setError("Username minimal 3 karakter.");
      return;
    }

    setLoading(true);

    try {
      const res = await fetchApi<{ token: string; user: any }>("/auth/register", {
        method: "POST",
        body: JSON.stringify({
          name: name.trim(),
          username: username.trim().toLowerCase(),
          email: email.trim().toLowerCase(),
          password,
        }),
      });

      if (res.data?.token && res.data?.user) {
        setAuth(res.data.token, res.data.user);
        setSuccess(true);
        setTimeout(() => {
          router.push("/?registered=true");
        }, 1200);
      }
    } catch (err: any) {
      setError(
        err?.data?.message || err?.message || "Registrasi gagal. Silakan coba kembali."
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[85vh] flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md bg-brand-surface rounded-3xl border border-brand-border p-8 sm:p-10 shadow-sm">
        {/* Header */}
        <div className="text-center mb-6">
          <div className="inline-flex w-14 h-14 rounded-2xl bg-brand-primary items-center justify-center text-white mb-3.5 shadow-sm">
            <Fish className="w-8 h-8 text-brand-secondary" />
          </div>
          <h1 className="text-2xl sm:text-3xl font-extrabold text-brand-primary">Daftar Akun Baru</h1>
          <p className="text-xs sm:text-sm text-brand-muted mt-1.5">
            Bergabunglah dengan ribuan pemancing dari seluruh Indonesia
          </p>
        </div>

        {/* Success feedback */}
        {success && (
          <div className="mb-6 p-4 rounded-2xl bg-green-50 border border-green-200 flex items-center gap-3 text-green-800 text-sm">
            <CheckCircle2 className="w-5 h-5 shrink-0 text-green-600" />
            <span className="font-semibold">Pendaftaran berhasil! Mengalihkan ke beranda...</span>
          </div>
        )}

        {/* Error notification */}
        {error && !success && (
          <div className="mb-6 p-4 rounded-2xl bg-red-50 border border-red-200 flex items-start gap-3 text-red-700 text-xs sm:text-sm">
            <AlertCircle className="w-5 h-5 shrink-0 mt-0.5 text-red-600" />
            <span className="font-semibold">{error}</span>
          </div>
        )}

        {/* Form */}
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
              Nama Lengkap
            </label>
            <div className="relative">
              <input
                type="text"
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Contoh: Budi Santoso"
                className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary text-sm transition"
              />
              <User className="w-4 h-4 text-brand-muted absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            </div>
          </div>

          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
              Username
            </label>
            <div className="relative">
              <input
                type="text"
                required
                value={username}
                onChange={(e) => setUsername(e.target.value.toLowerCase().replace(/[^a-z0-9_]/g, ""))}
                placeholder="budimancing"
                className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary text-sm transition"
              />
              <AtSign className="w-4 h-4 text-brand-muted absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            </div>
          </div>

          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
              Alamat Email
            </label>
            <div className="relative">
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="nama@email.com"
                className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary text-sm transition"
              />
              <Mail className="w-4 h-4 text-brand-muted absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            </div>
          </div>

          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
              Password
            </label>
            <div className="relative">
              <input
                type={showPassword ? "text" : "password"}
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Minimal 8 karakter"
                className="w-full pl-10 pr-10 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary text-sm transition"
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
          </div>

          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
              Konfirmasi Password
            </label>
            <div className="relative">
              <input
                type={showConfirmPassword ? "text" : "password"}
                required
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Ulangi password di atas"
                className={`w-full pl-10 pr-10 py-2.5 rounded-xl border bg-brand-background text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary text-sm transition ${
                  confirmPassword && !isPasswordMatch
                    ? "border-red-300 ring-1 ring-red-300"
                    : confirmPassword && isPasswordMatch
                    ? "border-green-400"
                    : "border-brand-border"
                }`}
              />
              <Lock className="w-4 h-4 text-brand-muted absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
              <button
                type="button"
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                tabIndex={-1}
                className="absolute right-3.5 top-1/2 -translate-y-1/2 text-brand-muted hover:text-brand-primary transition"
                title={showConfirmPassword ? "Sembunyikan password" : "Tampilkan password"}
              >
                {showConfirmPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
              </button>
            </div>
          </div>

          {/* Realtime Password Verification Checklist */}
          {password.length > 0 && (
            <div className="p-3.5 rounded-2xl bg-brand-background border border-brand-border text-xs space-y-1.5">
              <div className="flex items-center gap-2">
                {isLengthValid ? (
                  <Check className="w-3.5 h-3.5 text-green-600 shrink-0" />
                ) : (
                  <X className="w-3.5 h-3.5 text-red-500 shrink-0" />
                )}
                <span className={isLengthValid ? "text-green-700 font-medium" : "text-brand-muted"}>
                  Minimal 8 karakter (maksimal 72)
                </span>
              </div>

              <div className="flex items-center gap-2">
                {hasLetter && hasNumber ? (
                  <Check className="w-3.5 h-3.5 text-green-600 shrink-0" />
                ) : (
                  <X className="w-3.5 h-3.5 text-brand-muted shrink-0" />
                )}
                <span className={hasLetter && hasNumber ? "text-green-700 font-medium" : "text-brand-muted"}>
                  Kombinasi huruf dan angka
                </span>
              </div>

              {confirmPassword.length > 0 && (
                <div className="flex items-center gap-2">
                  {isPasswordMatch ? (
                    <Check className="w-3.5 h-3.5 text-green-600 shrink-0" />
                  ) : (
                    <X className="w-3.5 h-3.5 text-red-500 shrink-0" />
                  )}
                  <span className={isPasswordMatch ? "text-green-700 font-medium" : "text-red-600"}>
                    {isPasswordMatch ? "Konfirmasi password cocok" : "Konfirmasi password belum cocok"}
                  </span>
                </div>
              )}
            </div>
          )}

          <button
            type="submit"
            disabled={loading || success || (confirmPassword.length > 0 && !isPasswordMatch)}
            className="mt-3 w-full py-3.5 px-4 rounded-xl font-bold text-white bg-brand-primary hover:bg-brand-primary/90 focus:outline-none focus:ring-2 focus:ring-brand-primary shadow-sm disabled:opacity-50 transition flex items-center justify-center gap-2 text-sm"
          >
            {loading ? (
              <span className="flex items-center gap-2">
                <span className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                Mendaftarkan akun...
              </span>
            ) : (
              <>
                <UserPlus className="w-4 h-4 text-brand-secondary" />
                <span>Buat Akun Sekarang</span>
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
            <span className="px-3 bg-brand-surface text-brand-muted font-medium">atau daftar dengan</span>
          </div>
        </div>

        {/* Google Sign In */}
        <GoogleSignInButton
          label="Daftar dengan Google"
          onError={(err) => setError(err)}
        />

        {/* Footer */}
        <div className="mt-8 pt-6 border-t border-brand-border/80 text-center text-xs text-brand-muted">
          Sudah punya akun?{" "}
          <Link href="/login" className="font-bold text-brand-secondary hover:underline">
            Masuk di sini
          </Link>
        </div>
      </div>
    </div>
  );
}
