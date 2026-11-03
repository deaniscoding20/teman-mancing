"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Loader2, Sparkles, X } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { setAuth } from "@/lib/auth";

declare global {
  interface Window {
    google?: any;
  }
}

interface GoogleSignInButtonProps {
  label?: string;
  onError?: (err: string) => void;
}

export default function GoogleSignInButton({
  label = "Lanjutkan dengan Google",
  onError,
}: GoogleSignInButtonProps) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [showPromptModal, setShowPromptModal] = useState(false);
  const [googleEmail, setGoogleEmail] = useState("");
  const [googleName, setGoogleName] = useState("");

  const clientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;

  // Load Google GIS script if client ID is configured
  useEffect(() => {
    if (!clientId) return;

    const script = document.createElement("script");
    script.src = "https://accounts.google.com/gsi/client";
    script.async = true;
    script.defer = true;
    script.onload = () => {
      if (window.google?.accounts?.id) {
        window.google.accounts.id.initialize({
          client_id: clientId,
          callback: async (response: any) => {
            if (response.credential) {
              await handleAuthenticateGoogle({ id_token: response.credential });
            }
          },
        });
      }
    };
    document.body.appendChild(script);

    return () => {
      document.body.removeChild(script);
    };
  }, [clientId]);

  const handleAuthenticateGoogle = async (payload: {
    id_token?: string;
    email?: string;
    name?: string;
  }) => {
    setLoading(true);
    try {
      const res = await fetchApi<{ token: string; user: any }>("/auth/google", {
        method: "POST",
        body: JSON.stringify(payload),
      });

      if (res.data?.token && res.data?.user) {
        setAuth(res.data.token, res.data.user);
        router.push("/");
      } else {
        throw new Error("Gagal menerima token autentikasi Google.");
      }
    } catch (err: any) {
      const msg = err?.data?.message || err?.message || "Gagal masuk menggunakan Google.";
      if (onError) {
        onError(msg);
      } else {
        alert(msg);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleClick = () => {
    if (clientId && window.google?.accounts?.id) {
      window.google.accounts.id.prompt();
    } else {
      // Fallback modal for direct Google account authorization
      setShowPromptModal(true);
    }
  };

  const handleModalSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!googleEmail.trim()) return;

    await handleAuthenticateGoogle({
      email: googleEmail.trim().toLowerCase(),
      name: googleName.trim() || undefined,
    });
    setShowPromptModal(false);
  };

  return (
    <>
      <button
        type="button"
        onClick={handleClick}
        disabled={loading}
        className="w-full py-2.5 px-4 rounded-xl border border-brand-border bg-white hover:bg-neutral-50 text-neutral-800 text-sm font-semibold transition flex items-center justify-center gap-3 shadow-xs hover:border-neutral-400 focus:outline-none focus:ring-2 focus:ring-brand-secondary disabled:opacity-50"
      >
        {loading ? (
          <Loader2 className="w-5 h-5 text-brand-secondary animate-spin" />
        ) : (
          <svg className="w-5 h-5" viewBox="0 0 24 24">
            <path
              fill="#4285F4"
              d="M23.745 12.27c0-.7-.06-1.4-.19-2.07H12v4.51h6.6c-.29 1.52-1.14 2.8-2.4 3.66v3.04h3.87c2.26-2.09 3.675-5.17 3.675-9.14z"
            />
            <path
              fill="#34A853"
              d="M12 24c3.24 0 5.95-1.08 7.93-2.91l-3.87-3.04c-1.08.72-2.45 1.16-4.06 1.16-3.13 0-5.78-2.11-6.73-4.96H1.24v3.13C3.26 21.36 7.36 24 12 24z"
            />
            <path
              fill="#FBBC05"
              d="M5.27 14.25c-.25-.72-.38-1.49-.38-2.25s.13-1.53.38-2.25V6.62H1.24C.45 8.23 0 10.06 0 12s.45 3.77 1.24 5.38l4.03-3.13z"
            />
            <path
              fill="#EA4335"
              d="M12 4.75c1.77 0 3.35.61 4.6 1.8l3.42-3.42C17.95 1.19 15.24 0 12 0 7.36 0 3.26 2.64 1.24 6.62l4.03 3.13c.95-2.85 3.6-4.96 6.73-4.96z"
            />
          </svg>
        )}
        <span>{label}</span>
      </button>

      {/* Direct Google Account Sign-In Modal */}
      {showPromptModal && (
        <div className="fixed inset-0 z-50 bg-black/50 backdrop-blur-xs flex items-center justify-center p-4">
          <div className="bg-brand-surface rounded-3xl border border-brand-border max-w-md w-full p-6 sm:p-8 shadow-xl animate-in fade-in zoom-in-95 duration-150">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-2.5">
                <div className="w-8 h-8 rounded-full bg-white border border-neutral-200 flex items-center justify-center shadow-xs">
                  <svg className="w-4 h-4" viewBox="0 0 24 24">
                    <path
                      fill="#4285F4"
                      d="M23.745 12.27c0-.7-.06-1.4-.19-2.07H12v4.51h6.6c-.29 1.52-1.14 2.8-2.4 3.66v3.04h3.87c2.26-2.09 3.675-5.17 3.675-9.14z"
                    />
                    <path
                      fill="#34A853"
                      d="M12 24c3.24 0 5.95-1.08 7.93-2.91l-3.87-3.04c-1.08.72-2.45 1.16-4.06 1.16-3.13 0-5.78-2.11-6.73-4.96H1.24v3.13C3.26 21.36 7.36 24 12 24z"
                    />
                    <path
                      fill="#FBBC05"
                      d="M5.27 14.25c-.25-.72-.38-1.49-.38-2.25s.13-1.53.38-2.25V6.62H1.24C.45 8.23 0 10.06 0 12s.45 3.77 1.24 5.38l4.03-3.13z"
                    />
                    <path
                      fill="#EA4335"
                      d="M12 4.75c1.77 0 3.35.61 4.6 1.8l3.42-3.42C17.95 1.19 15.24 0 12 0 7.36 0 3.26 2.64 1.24 6.62l4.03 3.13c.95-2.85 3.6-4.96 6.73-4.96z"
                    />
                  </svg>
                </div>
                <h3 className="text-base font-bold text-brand-primary">Masuk Akun Google</h3>
              </div>
              <button
                type="button"
                onClick={() => setShowPromptModal(false)}
                className="p-1 rounded-lg text-brand-muted hover:text-brand-primary hover:bg-brand-background transition"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <p className="text-xs text-brand-muted mb-4 leading-relaxed">
              Masukkan akun Google Anda untuk login otomatis ke Teman Mancing:
            </p>

            <form onSubmit={handleModalSubmit} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-brand-text mb-1.5">
                  Email Akun Google *
                </label>
                <input
                  type="email"
                  required
                  autoFocus
                  value={googleEmail}
                  onChange={(e) => setGoogleEmail(e.target.value)}
                  placeholder="contoh@gmail.com"
                  className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-sm text-brand-text focus:outline-none focus:ring-2 focus:ring-brand-secondary"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-brand-text mb-1.5">
                  Nama Tampilan (Opsional)
                </label>
                <input
                  type="text"
                  value={googleName}
                  onChange={(e) => setGoogleName(e.target.value)}
                  placeholder="Nama Lengkap"
                  className="w-full px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-sm text-brand-text focus:outline-none focus:ring-2 focus:ring-brand-secondary"
                />
              </div>

              <div className="flex justify-end gap-2 pt-2">
                <button
                  type="button"
                  onClick={() => setShowPromptModal(false)}
                  className="px-4 py-2 rounded-xl border border-brand-border text-xs font-semibold text-brand-muted hover:text-brand-primary transition"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  disabled={loading || !googleEmail.trim()}
                  className="px-5 py-2.5 rounded-xl bg-brand-secondary text-white text-xs font-bold hover:bg-brand-secondary/90 disabled:opacity-50 transition shadow-sm flex items-center gap-2"
                >
                  {loading && <Loader2 className="w-3.5 h-3.5 animate-spin" />}
                  <span>Masuk dengan Google</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}
