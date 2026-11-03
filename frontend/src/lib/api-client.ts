import { ApiResponse, ApiErrorResponse } from "@/types/api";
import { getToken } from "./auth";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export class ApiClientError extends Error {
  public status: number;
  public data: ApiErrorResponse;

  constructor(status: number, data: ApiErrorResponse) {
    super(data.message || `Request failed with status ${status}`);
    this.name = "ApiClientError";
    this.status = status;
    this.data = data;
  }
}

export async function fetchApi<T = any>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const url = `${API_BASE_URL}${endpoint.startsWith("/") ? endpoint : `/${endpoint}`}`;

  const headers = new Headers(options.headers || {});
  if (!headers.has("Content-Type") && !(options.body instanceof FormData)) {
    headers.set("Content-Type", "application/json");
  }

  // Auto attach Bearer token if not explicitly provided
  const token = getToken();
  if (token && !headers.has("Authorization")) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  let response: Response;
  try {
    response = await fetch(url, {
      credentials: "include",
      ...options,
      headers,
    });
  } catch (err: any) {
    throw new ApiClientError(0, {
      success: false,
      message: "Tidak dapat terhubung ke server. Pastikan backend aktif.",
    });
  }

  let data: any;
  try {
    data = await response.json();
  } catch (err) {
    throw new ApiClientError(response.status, {
      success: false,
      message: `Server mengembalikan response tidak valid (HTTP ${response.status})`,
    });
  }

  if (!response.ok || !data.success) {
    throw new ApiClientError(response.status, data);
  }

  return data as ApiResponse<T>;
}
