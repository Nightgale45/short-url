import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";

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
    <form onSubmit={handleSubmit}>
      <div className="grid gap-2">
        <Label htmlFor="link">Link</Label>
        <Input
          id="link"
          type="text"
          placeholder="https://google.com"
          onChange={(e) => setInputValue(e.target.value)}
          required
        />
      </div>
      <Button
        type="submit"
        className="w-full mt-4"
        disabled={props.disableSubmit || inputValue.length == 0}
      >
        {props.disableSubmit ? "Shortening..." : "Submit"}
      </Button>
    </form>
  );
}

export default ShortenSubmit;
