import { Routes, Route } from "react-router-dom";
import { Window } from "./components/Window";
import { List } from "./pages/List";
import { Form } from "./pages/Form";

function App() {
  return (
    <Window>
      <Routes>
        <Route path="/" element={<List />} />
        <Route path="/macro/create" element={<Form />} />
        <Route path="/macro/:id" element={<Form />} />
      </Routes>
    </Window>
  );
}

export default App;
