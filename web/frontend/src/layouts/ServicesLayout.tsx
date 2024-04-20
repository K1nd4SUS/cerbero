import { Outlet } from "react-router-dom"
import Header from "../components/Header"
import Main from "./Main"
import Page from "./Page"

export default function ServicesLayout() {
  return (
    <Page>
      <Header/>
      <Main>
        <Outlet/>
      </Main>
    </Page>
  )
}

