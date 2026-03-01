import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { shorten } from "@/services/shorten.service";
import { useState } from "react";

function Shorten() {
  const [inputValue, setinputValue] = useState("");

  const handleSubmit = async () => {
    const response = await shorten({ original_url: inputValue });
    const data = await response.json();
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      <Card className="w-md">
        <CardHeader>
          <CardTitle>Link Shortener</CardTitle>
          <CardDescription>Enter a link to shorten</CardDescription>
        </CardHeader>
        <CardContent>
          <form>
            <div className="flex flex-col gap-6">
              <div className="grid gap-2">
                <Label htmlFor="email">Link</Label>
                <Input
                  id="link"
                  type="text"
                  placeholder="https://google.com"
                  onChange={(e) => setinputValue(e.target.value)}
                  required
                />
              </div>
            </div>
          </form>
        </CardContent>
        <CardFooter>
          <Button type="submit" className="w-full" onClick={handleSubmit}>
            Submit
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}

export default Shorten;
