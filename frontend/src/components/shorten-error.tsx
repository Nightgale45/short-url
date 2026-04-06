import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";

interface ShortenErrorProps {
  message: string;
}

function ShortenError({ message }: ShortenErrorProps) {
  return (
    <Card className="w-full max-w-sm border-red-500 bg-red-50 h-full">
      <CardHeader>
        <CardTitle className="text-red-600">Error</CardTitle>
        <CardDescription className="text-red-500">
          Something went wrong
        </CardDescription>
      </CardHeader>
      <CardContent className="mb-3">
        <div className="gap-3">
          <div className="grid gap-2 mt-3">
            <p className="text-sm text-red-600 break-all">{message}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export default ShortenError;
