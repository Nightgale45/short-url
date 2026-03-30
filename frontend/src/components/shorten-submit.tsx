import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import {
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { Card } from "@/components/ui/card";

interface ShortenSubmitProps {
  onSubmit: (url: string) => void;
  disableSubmit: boolean;
}

function ShortenSubmit(props: ShortenSubmitProps) {
  const [inputValue, setInputValue] = useState("");

  const handleSubmit = (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    props.onSubmit(inputValue); // sends the typed URL up to the parent
  };

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Shortener</CardTitle>
        <CardDescription>Enter a url to shorten</CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit}>
        <CardContent className="mb-3">
          <div className="flex flex-col gap-3">
            <div className="grid gap-2">
              <Label htmlFor="link">URL</Label>
              <Input
                id="link"
                type="text"
                placeholder="https://google.com"
                required
                onChange={(e) => setInputValue(e.target.value)}
              />
            </div>
          </div>
        </CardContent>
        <CardFooter className="flex-col gap-2">
          <Button
            type="submit"
            className="w-full"
            disabled={props.disableSubmit || inputValue.length == 0}
          >
            {props.disableSubmit ? "Shortening..." : "Submit"}
          </Button>
        </CardFooter>
      </form>
    </Card>
  );
}

export default ShortenSubmit;
