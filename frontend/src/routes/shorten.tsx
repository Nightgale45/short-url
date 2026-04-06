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

  const hasResponse = result || error;

  return (
    <div className="flex items-start justify-center min-h-screen pt-24">
      <div className={`flex items-stretch gap-8 transition-all duration-300 ${hasResponse ? "justify-start" : "justify-center"}`}>
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
