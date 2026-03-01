import type { ShortenRequest } from "@/models/shorten";

const basePath = "/api/v1/shorten";

export const shorten = async (request: ShortenRequest) =>
  fetch(basePath, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(request),
  });
