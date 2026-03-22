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
    <>
      <div className="submit-content">
        <ShortenSubmit onSubmit={handleSubmit} disableSubmit={loading} />
        {/* Only render ShortenReponse once data is returned */}
        {result && <ShortenReponse data={result} />}
      </div>
    </>
  );
}

export default Shorten;
