import type { ShortenRequest, ShortenResponse } from "@/models/shorten";

const basePath = "/api/v1/shorten";

export const shorten = async (request: ShortenRequest): Promise<ShortenResponse> => {
  const response = await fetch(basePath, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(request),
  });

  const body = await response.json() as ShortenResponse;

  if (!response.ok) {
    throw new Error(body.error_message ?? `Request failed: ${response.status}`);
  }

  return body;
};
