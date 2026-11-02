// Feature: Fishing Events (Janjian Mancing)
export interface EventFilterParams {
  page?: number;
  limit?: number;
  query?: string;
  date?: string;
  difficulty?: "pemula" | "menengah" | "ahli";
  visibility?: "public" | "private";
  status?: "open" | "closed" | "cancelled" | "completed";
}
