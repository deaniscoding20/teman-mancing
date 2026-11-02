"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Fish, LogIn, AlertCircle } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { setAuth } from "@/lib/auth";

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res = await fetchApi<{ token: string; user: any }>("/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });

      if (res.data?.token && res.data?.user) {
        setAuth(res.data.token, res.data.user);
        router.push("/");
      }
    } catch (err: any) {
      setError(err?.data?.message || err?.message || "Gagal masuk. Periksa kembali email dan password.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md bg-brand-surface rounded-2xl border border-brand-border p-8 shadow-sm">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="inline-flex w-12 h-12 rounded-xl bg-brand-primary items-center justify-center text-white mb-3">
            <Fish className="w-7 h-7 text-brand-secondary" />
          </div>
          <h1 className="text-2xl font-bold text-brand-primary">Masuk ke Teman Mancing</h1>
          <p className="text-sm text-brand-muted mt-1">
            Hubungkan kembali koneksi mancing dan agenda tripmu
          </p>
        </div>

        {/* Error notification */}
        {error && (
          <div className="mb-6 p-4 rounded-xl bg-red-50 border border-red-200 flex items-start gap-3 text-red-700 text-sm">
            <AlertCircle className="w-5 h-5 shrink-0 mt-0.5" />
            <span>{error}</span>
          </div>
        )}

        {/* Form */}
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
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
              Password
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
            disabled={loading}
            className="mt-2 w-full py-3 px-4 rounded-xl font-bold text-white bg-brand-primary hover:bg-brand-primary/90 focus:outline-none focus:ring-2 focus:ring-brand-primary shadow-sm disabled:opacity-50 transition flex items-center justify-center gap-2"
          >
            {loading ? (
              <span>Memproses...</span>
            ) : (
              <>
                <LogIn className="w-4 h-4 text-brand-secondary" />
                <span>Masuk Sekarang</span>
              </>
            )}
          </button>
        </form>

        {/* Footer */}
        <div className="mt-6 text-center text-xs text-brand-muted">
          Belum punya akun?{" "}
          <Link href="/register" className="font-semibold text-brand-secondary hover:underline">
            Daftar di sini
          </Link>
        </div>
      </div>
    </div>
  );
}
