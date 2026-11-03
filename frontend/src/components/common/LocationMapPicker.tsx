"use client";

import { useState, useEffect, useRef } from "react";
import { MapPin, Search, Navigation, ExternalLink, Loader2 } from "lucide-react";

interface LocationMapPickerProps {
  address: string;
  latitude: string;
  longitude: string;
  onAddressChange: (address: string) => void;
  onCoordinateChange: (lat: string, lng: string) => void;
  label?: string;
}

interface GeocodingResult {
  place_id: number;
  display_name: string;
  lat: string;
  lon: string;
}

export default function LocationMapPicker({
  address,
  latitude,
  longitude,
  onAddressChange,
  onCoordinateChange,
  label = "Lokasi & Koordinat Peta",
}: LocationMapPickerProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [suggestions, setSuggestions] = useState<GeocodingResult[]>([]);
  const [searching, setSearching] = useState(false);
  const [locating, setLocating] = useState(false);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close suggestions dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setShowSuggestions(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  // Debounced search on Nominatim
  useEffect(() => {
    if (!searchQuery.trim() || searchQuery.trim().length < 3) {
      setSuggestions([]);
      return;
    }

    const timer = setTimeout(async () => {
      setSearching(true);
      try {
        const res = await fetch(
          `https://nominatim.openstreetmap.org/search?format=json&countrycodes=id&q=${encodeURIComponent(
            searchQuery.trim()
          )}&limit=5`,
          {
            headers: {
              "Accept-Language": "id,en",
            },
          }
        );
        if (res.ok) {
          const data: GeocodingResult[] = await res.json();
          setSuggestions(data);
          setShowSuggestions(true);
        }
      } catch (err) {
        console.warn("Geocoding fetch error:", err);
      } finally {
        setSearching(false);
      }
    }, 400);

    return () => clearTimeout(timer);
  }, [searchQuery]);

  const handleSelectSuggestion = (item: GeocodingResult) => {
    onAddressChange(item.display_name);
    onCoordinateChange(item.lat, item.lon);
    setShowSuggestions(false);
    setSearchQuery("");
  };

  const handleGetCurrentLocation = () => {
    if (!navigator.geolocation) {
      alert("Geolocation tidak didukung oleh browser Anda.");
      return;
    }

    setLocating(true);
    navigator.geolocation.getCurrentPosition(
      async (pos) => {
        const lat = pos.coords.latitude.toFixed(6);
        const lng = pos.coords.longitude.toFixed(6);
        onCoordinateChange(lat, lng);

        // Reverse geocoding to get human-readable address
        try {
          const res = await fetch(
            `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}`,
            {
              headers: { "Accept-Language": "id,en" },
            }
          );
          if (res.ok) {
            const data = await res.json();
            if (data.display_name) {
              onAddressChange(data.display_name);
            }
          }
        } catch (err) {
          console.warn("Reverse geocode error:", err);
        } finally {
          setLocating(false);
        }
      },
      (err) => {
        setLocating(false);
        alert(`Gagal mendeteksi lokasi GPS: ${err.message}`);
      },
      { enableHighAccuracy: true, timeout: 10000 }
    );
  };

  const parsedLat = parseFloat(latitude) || -6.2088;
  const parsedLng = parseFloat(longitude) || 106.8456;

  const googleMapsSearchUrl = `https://www.google.com/maps/search/?api=1&query=${parsedLat},${parsedLng}`;
  const googleMapsDirectionsUrl = `https://www.google.com/maps/dir/?api=1&destination=${parsedLat},${parsedLng}`;

  return (
    <div className="space-y-4">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <label className="block text-xs font-semibold text-brand-text uppercase tracking-wider">
          {label}
        </label>
        <button
          type="button"
          onClick={handleGetCurrentLocation}
          disabled={locating}
          className="inline-flex items-center gap-1.5 text-xs font-semibold text-brand-secondary hover:text-brand-secondary/80 disabled:opacity-50 transition cursor-pointer"
        >
          {locating ? (
            <Loader2 className="w-3.5 h-3.5 animate-spin" />
          ) : (
            <Navigation className="w-3.5 h-3.5" />
          )}
          <span>{locating ? "Mendeteksi GPS..." : "Gunakan Lokasi GPS Saya"}</span>
        </button>
      </div>

      {/* Search Autocomplete */}
      <div className="relative" ref={dropdownRef}>
        <div className="relative">
          <Search className="w-4 h-4 text-brand-muted absolute left-3.5 top-3" />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onFocus={() => suggestions.length > 0 && setShowSuggestions(true)}
            placeholder="Cari nama tempat / danau / pantai / kota di peta (misal: Pantai Ancol, Danau Sunter)..."
            className="w-full pl-10 pr-10 py-2.5 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
          />
          {searching && (
            <Loader2 className="w-4 h-4 text-brand-secondary animate-spin absolute right-3.5 top-3" />
          )}
        </div>

        {/* Suggestions Dropdown */}
        {showSuggestions && suggestions.length > 0 && (
          <div className="absolute z-30 left-0 right-0 mt-1 bg-brand-surface rounded-xl border border-brand-border shadow-lg max-h-60 overflow-y-auto divide-y divide-brand-border/60">
            {suggestions.map((item) => (
              <button
                type="button"
                key={item.place_id}
                onClick={() => handleSelectSuggestion(item)}
                className="w-full text-left p-3 hover:bg-brand-background flex items-start gap-2.5 transition text-xs"
              >
                <MapPin className="w-4 h-4 text-brand-secondary shrink-0 mt-0.5" />
                <span className="text-brand-text leading-snug line-clamp-2">
                  {item.display_name}
                </span>
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Alamat Terpilih */}
      <div>
        <label className="block text-xs font-semibold text-brand-muted mb-1">
          Alamat Lengkap / Petunjuk Jalan
        </label>
        <textarea
          rows={2}
          required
          value={address}
          onChange={(e) => onAddressChange(e.target.value)}
          placeholder="Alamat detail lokasi..."
          className="w-full px-4 py-2 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
        />
      </div>

      {/* Latitude & Longitude Inputs */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label className="block text-xs font-semibold text-brand-muted mb-1">Latitude</label>
          <input
            type="text"
            required
            value={latitude}
            onChange={(e) => onCoordinateChange(e.target.value, longitude)}
            placeholder="-6.2088"
            className="w-full px-4 py-2 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
          />
        </div>
        <div>
          <label className="block text-xs font-semibold text-brand-muted mb-1">Longitude</label>
          <input
            type="text"
            required
            value={longitude}
            onChange={(e) => onCoordinateChange(latitude, e.target.value)}
            placeholder="106.8456"
            className="w-full px-4 py-2 rounded-xl border border-brand-border bg-brand-background text-brand-text text-sm focus:outline-none focus:ring-2 focus:ring-brand-secondary"
          />
        </div>
      </div>

      {/* Interactive Map Preview with Google Maps Actions */}
      <div className="rounded-2xl border border-brand-border overflow-hidden bg-brand-background">
        <div className="h-52 w-full relative">
          <iframe
            title="Peta Lokasi"
            width="100%"
            height="100%"
            frameBorder="0"
            scrolling="no"
            src={`https://www.openstreetmap.org/export/embed.html?bbox=${parsedLng - 0.015}%2C${
              parsedLat - 0.015
            }%2C${parsedLng + 0.015}%2C${parsedLat + 0.015}&layer=mapnik&marker=${parsedLat}%2C${parsedLng}`}
            className="w-full h-full border-0"
          />
        </div>

        {/* Map toolbar & Google Maps external links */}
        <div className="p-3 bg-brand-surface border-t border-brand-border flex flex-wrap items-center justify-between gap-2 text-xs">
          <div className="flex items-center gap-1.5 text-brand-muted">
            <MapPin className="w-3.5 h-3.5 text-brand-secondary" />
            <span>
              Koordinat: {parsedLat.toFixed(5)}, {parsedLng.toFixed(5)}
            </span>
          </div>

          <div className="flex items-center gap-2">
            <a
              href={googleMapsSearchUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1 text-brand-secondary hover:underline font-semibold"
            >
              <span>Buka di Google Maps</span>
              <ExternalLink className="w-3 h-3" />
            </a>
            <span className="text-brand-border">|</span>
            <a
              href={googleMapsDirectionsUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1 text-brand-secondary hover:underline font-semibold"
            >
              <span>Rute Petunjuk Arah</span>
              <ExternalLink className="w-3 h-3" />
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}
