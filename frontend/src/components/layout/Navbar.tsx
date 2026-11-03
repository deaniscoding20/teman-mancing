"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Compass, Calendar, MessageSquare, Bell, User, Fish, LogOut, Settings } from "lucide-react";
import { getStoredUser, clearAuth } from "@/lib/auth";
import { fetchApi } from "@/lib/api-client";

export function Navbar() {
  const router = useRouter();
  const [currentUser, setCurrentUser] = useState<any>(null);

  useEffect(() => {
    const syncUser = () => {
      setCurrentUser(getStoredUser());
    };
    syncUser();

    window.addEventListener("auth-change", syncUser);
    window.addEventListener("storage", syncUser);
    return () => {
      window.removeEventListener("auth-change", syncUser);
      window.removeEventListener("storage", syncUser);
    };
  }, []);

  const handleLogout = async () => {
    try {
      await fetchApi("/auth/logout", { method: "POST" });
    } catch {
      // Ignore network errors on logout
    }
    clearAuth();
    router.push("/login");
  };

  return (
    <header className="sticky top-0 z-50 w-full bg-brand-surface/95 backdrop-blur border-b border-brand-border">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        {/* Brand */}
        <Link href="/" className="flex items-center gap-2.5 font-bold text-xl text-brand-primary">
          <div className="w-10 h-10 rounded-xl bg-brand-primary flex items-center justify-center text-white shadow-sm">
            <Fish className="w-6 h-6 text-brand-secondary" />
          </div>
          <div className="flex flex-col">
            <span className="leading-tight tracking-tight text-brand-primary font-extrabold text-lg">
              Teman Mancing
            </span>
            <span className="text-[10px] text-brand-muted tracking-wide font-normal -mt-0.5">
              Gas Mancing Bareng!
            </span>
          </div>
        </Link>

        {/* Navigation Items */}
        <nav className="hidden md:flex items-center gap-1 lg:gap-2">
          <Link
            href="/spots"
            className="flex items-center gap-1.5 px-3 py-2 rounded-lg text-sm font-medium text-brand-text hover:text-brand-secondary hover:bg-brand-background transition"
          >
            <Compass className="w-4 h-4 text-brand-muted" />
            <span>Spot Mancing</span>
          </Link>

          <Link
            href="/events"
            className="flex items-center gap-1.5 px-3 py-2 rounded-lg text-sm font-medium text-brand-text hover:text-brand-secondary hover:bg-brand-background transition"
          >
            <Calendar className="w-4 h-4 text-brand-muted" />
            <span>Janjian Mancing</span>
          </Link>

          <Link
            href="/chat"
            className="flex items-center gap-1.5 px-3 py-2 rounded-lg text-sm font-medium text-brand-text hover:text-brand-secondary hover:bg-brand-background transition"
          >
            <MessageSquare className="w-4 h-4 text-brand-muted" />
            <span>Chat</span>
          </Link>
        </nav>

        {/* Right Actions */}
        <div className="flex items-center gap-3">
          <Link
            href="/notifications"
            className="p-2 text-brand-text hover:text-brand-primary hover:bg-brand-background rounded-full transition relative"
            title="Notifikasi"
          >
            <Bell className="w-5 h-5" />
            <span className="absolute top-1 right-1 w-2 h-2 bg-brand-secondary rounded-full" />
          </Link>

          <div className="h-5 w-[1px] bg-brand-border hidden sm:block" />

          {currentUser ? (
            <div className="flex items-center gap-2">
              <Link
                href={`/profile/${currentUser.username}`}
                className="flex items-center gap-2 p-1.5 sm:px-3 sm:py-1.5 rounded-xl bg-brand-background hover:bg-brand-border/40 border border-brand-border transition"
              >
                <div className="w-7 h-7 rounded-full bg-brand-primary/10 flex items-center justify-center text-brand-secondary font-bold text-xs shrink-0">
                  <User className="w-4 h-4" />
                </div>
                <div className="hidden sm:flex flex-col text-left">
                  <span className="text-xs font-bold text-brand-primary leading-tight truncate max-w-[100px]">
                    {currentUser.name?.split(" ")[0] || "Angler"}
                  </span>
                  <span className="text-[10px] text-brand-muted leading-tight">
                    @{currentUser.username}
                  </span>
                </div>
              </Link>

              <Link
                href="/settings/profile"
                className="p-2 text-brand-muted hover:text-brand-primary rounded-xl transition"
                title="Pengaturan Profil"
              >
                <Settings className="w-4 h-4" />
              </Link>

              <button
                onClick={handleLogout}
                className="p-2 text-brand-muted hover:text-red-500 rounded-xl transition"
                title="Keluar"
              >
                <LogOut className="w-4 h-4" />
              </button>
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <Link
                href="/login"
                className="px-3.5 py-1.5 text-sm font-semibold text-brand-primary hover:text-brand-secondary transition"
              >
                Masuk
              </Link>
              <Link
                href="/register"
                className="px-4 py-2 text-sm font-semibold text-white bg-brand-secondary hover:bg-brand-secondary/90 rounded-lg shadow-sm transition"
              >
                Daftar
              </Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
