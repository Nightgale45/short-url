import ShortenSubmit from "@/components/shorten-submit";
import ShortenReponse from "@/components/shorten-reponse";
import ShortenError from "@/components/shorten-error";
import { shorten } from "@/services/shorten.service";
import { useState } from "react";
import type { ShortenRequest, ShortenResponse } from "@/models/shorten";

function Shorten() {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<ShortenResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (url: string) => {
    setLoading(true);
    setResult(null);
    setError(null);

    const request: ShortenRequest = { original_url: url };
    try {
      const data = await shorten(request);
      setResult(data);
    } catch (e) {
      setError(e instanceof Error ? e.message : "An unexpected error occurred");
    }

    setLoading(false);
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="flex items-stretch gap-8 justify-center transition-all duration-300">
        <div className="w-80">
          <ShortenSubmit onSubmit={handleSubmit} disableSubmit={loading} />
        </div>
        {result && (
          <div className="w-80">
            <ShortenReponse data={result} />
          </div>
        )}
        {error && (
          <div className="w-80">
            <ShortenError message={error} />
          </div>
        )}
      </div>
    </div>
  );
}

export default Shorten;
