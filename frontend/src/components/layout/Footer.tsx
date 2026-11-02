import Link from "next/link";
import { Fish, Heart } from "lucide-react";

export function Footer() {
  return (
    <footer className="w-full bg-brand-surface border-t border-brand-border mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        <div className="flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-brand-primary flex items-center justify-center text-white">
              <Fish className="w-5 h-5 text-brand-secondary" />
            </div>
            <div>
              <p className="font-bold text-brand-primary text-base">Teman Mancing</p>
              <p className="text-xs text-brand-muted">Cari Spot. Cari Teman. Gas Mancing!</p>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-6 text-sm text-brand-muted">
            <Link href="/spots" className="hover:text-brand-primary transition">
              Spot Mancing
            </Link>
            <Link href="/events" className="hover:text-brand-primary transition">
              Janjian Mancing
            </Link>
            <Link href="/chat" className="hover:text-brand-primary transition">
              Komunitas & Chat
            </Link>
          </div>

          <p className="text-xs text-brand-muted flex items-center gap-1">
            Built with <Heart className="w-3.5 h-3.5 text-red-500 fill-red-500" /> for Indonesian Anglers &copy; {new Date().getFullYear()}
          </p>
        </div>
      </div>
    </footer>
  );
}
