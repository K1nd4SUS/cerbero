import { Outlet } from "react-router-dom"
import Header from "../components/Header"
import Main from "./Main"

export default function ServicesLayout() {
  return (
    <>
      <Header/>
      <Main>
        <Outlet/>
      </Main>
    </>
  )
}
