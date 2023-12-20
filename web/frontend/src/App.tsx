import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import Service from "./pages/Service"
import Services from "./pages/Services"

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="*"
          element={
            <Navigate to="/services"/>
          }
        />
        <Route
          path="/services"
          element={
            <Services/>
          }
        />
        <Route
          path="/services/:nfq"
          element={
            <Service/>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}
