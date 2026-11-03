"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { Compass, UserPlus } from "lucide-react";
import { getStoredUser, getToken } from "@/lib/auth";

export default function HomeCTAButton() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    const update = () => {
      const token = getToken();
      const user = getStoredUser();
      setIsLoggedIn(!!(token && user));
    };
    update();
    window.addEventListener("auth-change", update);
    return () => window.removeEventListener("auth-change", update);
  }, []);

  if (isLoggedIn) {
    return (
      <Link
        href="/spots"
        className="px-8 py-4 rounded-xl font-bold text-brand-primary bg-white hover:bg-white/95 shadow-md transition shrink-0 inline-flex items-center gap-2"
      >
        <Compass className="w-5 h-5 text-brand-secondary" />
        <span>Jelajahi Spot Sekarang</span>
      </Link>
    );
  }

  return (
    <Link
      href="/register"
      className="px-8 py-4 rounded-xl font-bold text-brand-primary bg-white hover:bg-white/95 shadow-md transition shrink-0 inline-flex items-center gap-2"
    >
      <UserPlus className="w-5 h-5 text-brand-secondary" />
      <span>Gabung Sekarang</span>
    </Link>
  );
}
