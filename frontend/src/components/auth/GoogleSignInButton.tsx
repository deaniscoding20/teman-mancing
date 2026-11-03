"use client";

import { useState } from "react";
import { Loader2 } from "lucide-react";

interface GoogleSignInButtonProps {
  label?: string;
  onError?: (err: string) => void;
}

export default function GoogleSignInButton({
  label = "Lanjutkan dengan Google",
  onError,
}: GoogleSignInButtonProps) {
  const [loading, setLoading] = useState(false);

  const handleClick = () => {
    setLoading(true);
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "https://temanmancing.my.id/api/v1";
    // Redirect directly to backend Google OAuth login route
    window.location.href = `${apiUrl.replace(/\/+$/, "")}/auth/google/login`;
  };

  return (
    <button
      type="button"
      onClick={handleClick}
      disabled={loading}
      className="w-full py-2.5 px-4 rounded-xl border border-brand-border bg-white hover:bg-neutral-50 text-neutral-800 text-sm font-semibold transition flex items-center justify-center gap-3 shadow-xs hover:border-neutral-400 focus:outline-none focus:ring-2 focus:ring-brand-secondary disabled:opacity-60 cursor-pointer"
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
  );
}
