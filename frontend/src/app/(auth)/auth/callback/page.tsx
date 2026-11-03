"use client";

import { useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Loader2 } from "lucide-react";
import { setAuth } from "@/lib/auth";

function AuthCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const token = searchParams.get("token");
    const userStr = searchParams.get("user");
    const error = searchParams.get("error");

    if (error) {
      router.push(`/login?error=${encodeURIComponent(error)}`);
      return;
    }

    if (token) {
      let user = null;
      if (userStr) {
        try {
          user = JSON.parse(decodeURIComponent(userStr));
        } catch (e) {
          console.error("Failed to parse user payload:", e);
        }
      }
      setAuth(token, user);
      router.push("/");
    } else {
      router.push("/login?error=Token+autentikasi+tidak+ditemukan");
    }
  }, [searchParams, router]);

  return (
    <div className="min-h-[80vh] flex flex-col items-center justify-center gap-3">
      <Loader2 className="w-10 h-10 text-brand-secondary animate-spin" />
      <p className="text-sm font-semibold text-brand-primary">
        Menghubungkan sesi akun Anda...
      </p>
    </div>
  );
}

export default function AuthCallbackPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-[80vh] flex items-center justify-center">
          <Loader2 className="w-10 h-10 text-brand-secondary animate-spin" />
        </div>
      }
    >
      <AuthCallbackContent />
    </Suspense>
  );
}
