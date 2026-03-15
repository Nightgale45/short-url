export interface ShortenRequest {
  original_url: string;
}

export interface ShortenResponse {
  original_url: string;
  shorten_url: string;
}
