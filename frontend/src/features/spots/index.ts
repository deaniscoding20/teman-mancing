// Feature: Fishing Spots
export interface SpotFilterParams {
  page?: number;
  limit?: number;
  query?: string;
  fish_type?: string;
  technique?: string;
  sort?: "newest" | "comments";
}
