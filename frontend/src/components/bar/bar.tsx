import { Button } from "../ui/button";

function Bar() {
  return (
    <div className="header flex gap-1.5 justify-center relative p-2 rounded-md bg-slate-200">
      <Button
        variant="outline"
        className="hover:bg-neutral-100 hover:border-gray-100"
      >
        Shorten Url
      </Button>
      <Button variant="outline" className="hover:accent-slate-500">
        Analytics
      </Button>
    </div>
  );
}

export default Bar;
