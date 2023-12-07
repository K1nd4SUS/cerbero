import { BrowserRouter, Routes, Route } from "react-router-dom"
import Home from "./pages/Home"

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* TODO: Add 404 page */}
        <Route
          path="*"
          element={
            <Home/>
          }
        />
        <Route
          path="/"
          element={
            <Home/>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}
