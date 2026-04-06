import type { ShortenResponse } from "@/models/shorten";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { Label } from "@/components/ui/label";

interface ShortenReponseProps {
  data: ShortenResponse;
}

function ShortenReponse({ data }: ShortenReponseProps) {
  return (
    <Card className="w-full max-w-sm h-full">
      <CardHeader>
        <CardTitle>Result</CardTitle>
        <CardDescription>Your shortened URL is ready</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col gap-3">
          <div className="grid gap-2">
            <Label>Original URL</Label>
            <p className="text-sm break-all">{data.original_url}</p>
          </div>
          <div className="grid gap-2">
            <Label>Shortened URL</Label>
            <a
              href={data.shorten_url}
              className="text-sm underline break-all"
              target="_blank"
              rel="noopener noreferrer"
            >
              {data.shorten_url}
            </a>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export default ShortenReponse;
