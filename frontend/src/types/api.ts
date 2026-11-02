// Teman Mancing API TypeScript Types

export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  meta?: PaginationMeta;
}

export interface ApiFieldError {
  field: string;
  message: string;
}

export interface ApiErrorResponse {
  success: false;
  message: string;
  errors?: ApiFieldError[];
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total_items: number;
  total_pages: number;
}

export interface UserProfile {
  id: string;
  name: string;
  username: string;
  email: string;
  avatar_key?: string;
  avatar_url?: string;
  bio?: string;
  location?: string;
  spots_count?: number;
  events_created_count?: number;
  events_joined_count?: number;
  created_at: string;
}

export interface FishingSpotItem {
  id: string;
  user_id: string;
  user_name?: string;
  name: string;
  description: string;
  address: string;
  latitude: number;
  longitude: number;
  primary_photo_url?: string;
  fish_types?: string[];
  techniques?: string[];
  comments_count: number;
  created_at: string;
}

export interface FishingEventItem {
  id: string;
  creator_id: string;
  creator_name?: string;
  creator_avatar_url?: string;
  name: string;
  description: string;
  cover_image_url?: string;
  destination_name: string;
  event_date: string;
  start_time: string;
  meeting_point: string;
  difficulty: "pemula" | "menengah" | "ahli";
  visibility: "public" | "private";
  status: "open" | "closed" | "cancelled" | "completed";
  max_participants?: number;
  approved_participants: number;
}
