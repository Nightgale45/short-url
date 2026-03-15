import type { ShortenResponse } from "@/models/shorten";

interface ShortenReponseProps {
  data: ShortenResponse;
}

function ShortenReponse({ data }: ShortenReponseProps) {
  return (
    <div className="grid gap-2 mt-4">
      <p className="text-sm">
        Original URL: <span className="break-all">{data.original_url}</span>
      </p>
      <p className="text-sm">
        Shortened URL:{" "}
        <a href={data.shorten_url} className="underline break-all">
          {data.shorten_url}
        </a>
      </p>
    </div>
  );
}

export default ShortenReponse;
