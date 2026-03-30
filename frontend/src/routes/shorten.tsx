import ShortenSubmit from "@/components/shorten-submit";
import ShortenReponse from "@/components/shorten-reponse";
import { shorten } from "@/services/shorten.service";
import { useState } from "react";
import type { ShortenRequest, ShortenResponse } from "@/models/shorten";

function Shorten() {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<ShortenResponse | null>(null);

  const handleSubmit = async (url: string) => {
    setLoading(true);
    setResult(null); // clear previous result on new submission

    const request: ShortenRequest = { original_url: url };
    try {
      const data = await shorten(request);
      setResult(data);
    } catch (e) {
      console.error(e);
    }

    setLoading(false);
  };

  return (
    <div className="flex items-start justify-center min-h-screen pt-24">
      <div className={`flex gap-8 transition-all duration-300 ${result ? "justify-start" : "justify-center"}`}>
        <div className="w-80">
          <ShortenSubmit onSubmit={handleSubmit} disableSubmit={loading} />
        </div>
        {result && (
          <div className="w-80">
            <ShortenReponse data={result} />
          </div>
        )}
      </div>
    </div>
  );
}

export default Shorten;
