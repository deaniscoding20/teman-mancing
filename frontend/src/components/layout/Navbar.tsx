"use client";

import Link from "next/link";
import { Compass, Calendar, MessageSquare, Bell, User, Fish } from "lucide-react";

export function Navbar() {
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
        </div>
      </div>
    </header>
  );
}
