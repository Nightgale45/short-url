import { BrowserRouter, Route, Routes } from "react-router";
import Shorten from "./routes/shorten";

function App() {
  return (
    <div className="min-h-full">
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Shorten />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
