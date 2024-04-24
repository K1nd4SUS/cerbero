import { Outlet } from "react-router-dom"

export default function BaseLayout() {
  return (
    <div className="absolute h-full w-full flex flex-col bg-default-50">
      <Outlet/>
    </div>
  )
}
