"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Fish, UserPlus, AlertCircle, CheckCircle2 } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { setAuth } from "@/lib/auth";

export default function RegisterPage() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (password.length < 8) {
      setError("Password minimal harus 8 karakter");
      return;
    }

    setLoading(true);

    try {
      const res = await fetchApi<{ token: string; user: any }>("/auth/register", {
        method: "POST",
        body: JSON.stringify({ name, username, email, password }),
      });

      if (res.data?.token && res.data?.user) {
        setAuth(res.data.token, res.data.user);
        setSuccess(true);
        setTimeout(() => {
          router.push("/");
        }, 1500);
      }
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Registrasi gagal. Coba lagi.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[85vh] flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md bg-brand-surface rounded-2xl border border-brand-border p-8 shadow-sm">
        {/* Header */}
        <div className="text-center mb-6">
          <div className="inline-flex w-12 h-12 rounded-xl bg-brand-primary items-center justify-center text-white mb-3">
            <Fish className="w-7 h-7 text-brand-secondary" />
          </div>
          <h1 className="text-2xl font-bold text-brand-primary">Daftar Akun Baru</h1>
          <p className="text-sm text-brand-muted mt-1">
            Bergabunglah dengan komunitas pemancing Indonesia
          </p>
        </div>

        {/* Success feedback */}
        {success && (
          <div className="mb-6 p-4 rounded-xl bg-green-50 border border-green-200 flex items-center gap-3 text-green-800 text-sm">
            <CheckCircle2 className="w-5 h-5 shrink-0 text-green-600" />
            <span>Pendaftaran berhasil! Mengalihkan ke beranda...</span>
          </div>
        )}

        {/* Error notification */}
        {error && !success && (
          <div className="mb-6 p-4 rounded-xl bg-red-50 border border-red-200 flex items-start gap-3 text-red-700 text-sm">
            <AlertCircle className="w-5 h-5 shrink-0 mt-0.5" />
            <span>{error}</span>
          </div>
        )}

        {/* Form */}
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
              Nama Lengkap
            </label>
            <input
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Contoh: Budi Santoso"
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary text-sm transition"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
              Username
            </label>
            <input
              type="text"
              required
              value={username}
              onChange={(e) => setUsername(e.target.value.toLowerCase().replace(/[^a-z0-9_]/g, ""))}
              placeholder="budimancing"
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary text-sm transition"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
              Email
            </label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="nama@email.com"
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary text-sm transition"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider mb-1.5">
              Password (Min. 8 Karakter)
            </label>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary text-sm transition"
            />
          </div>

          <button
            type="submit"
            disabled={loading || success}
            className="mt-2 w-full py-3 px-4 rounded-xl font-bold text-white bg-brand-primary hover:bg-brand-primary/90 focus:outline-none focus:ring-2 focus:ring-brand-primary shadow-sm disabled:opacity-50 transition flex items-center justify-center gap-2"
          >
            {loading ? (
              <span>Mendaftarkan...</span>
            ) : (
              <>
                <UserPlus className="w-4 h-4 text-brand-secondary" />
                <span>Buat Akun Sekarang</span>
              </>
            )}
          </button>
        </form>

        {/* Footer */}
        <div className="mt-6 text-center text-xs text-brand-muted">
          Sudah punya akun?{" "}
          <Link href="/login" className="font-semibold text-brand-secondary hover:underline">
            Masuk di sini
          </Link>
        </div>
      </div>
    </div>
  );
}
