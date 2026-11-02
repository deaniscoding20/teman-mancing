import type { Metadata } from "next";
import "./globals.css";
import { Navbar } from "@/components/layout/Navbar";
import { Footer } from "@/components/layout/Footer";

export const metadata: Metadata = {
  title: "Teman Mancing — Cari Spot. Cari Teman. Gas Mancing!",
  description: "Platform komunitas digital pemancing Indonesia. Temukan spot baru, adakan janjian mancing bersama, dan diskusikan pengalaman terbaikmu.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="id">
      <body className="flex flex-col min-h-screen bg-brand-background text-brand-text antialiased">
        <Navbar />
        <main className="flex-1">
          {children}
        </main>
        <Footer />
      </body>
    </html>
  );
}
