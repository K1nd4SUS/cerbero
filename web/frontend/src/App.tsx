import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import BaseLayout from "./layouts/BaseLayout"
import ServicesLayout from "./layouts/ServicesLayout"
import Home from "./pages/Home"
import Service from "./pages/Service"
import Services from "./pages/Services"

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="*"
          element={
            <Navigate to="/"/>
          }
        />
        <Route
          path="/"
          element={
            <BaseLayout/>
          }
        >
          <Route
            index
            element={
              <Home/>
            }
          />
          <Route
            path="services"
            element={
              <ServicesLayout/>
            }
          >
            <Route
              index
              element={
                <Services/>
              }
            />
            <Route
              path=":nfq"
              element={
                <Service/>
              }
            />
          </Route>
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
