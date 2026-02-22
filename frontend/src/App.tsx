import { BrowserRouter, Routes } from "react-router";
import Bar from "./components/bar/bar";

function App() {
  return (
    <div className="min-h-full px-[10%]">
            <BrowserRouter>
                <Routes>
                    {/* <Route path="/"></Route> */}
                    {/* <Route path></Route> */}
                </Routes>
            </BrowserRouter>
      <Bar></Bar>
    </div>
  );
}

export default App;
