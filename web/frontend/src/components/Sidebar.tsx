import { Button, Divider } from "@nextui-org/react"
import { motion } from "framer-motion"
import { useState } from "react"
import { Link, useNavigate } from "react-router-dom"
import { FaBars, FaCode, FaHouse } from "react-icons/fa6"
import { useFetchSync } from "../hooks/useFetch"
import { CerberoService } from "../types/cerbero"
import SidebarItem from "./SidebarItem"

export default function Sidebar() {
  const [
    servicesResponse,
    ,
  ] = useFetchSync<CerberoService[]>("/api/services")

  const [isOpen, setIsOpen] = useState(false)
  const navigate = useNavigate()

  return (
    <motion.aside
      initial={{ width: "4rem" }}
      animate={{ width: isOpen ? "24rem" : "4rem" }}
      className="h-full flex flex-col items-center p-2 bg-default-100"
    >
      <div className="w-full flex items-center justify-center">
        {isOpen ?
          <Button onPress={() => setIsOpen(false)} variant="flat" className="w-full">
            <span>Close sidebar</span>
          </Button> :
          <Button isIconOnly={true} onPress={() => setIsOpen(true)} variant="flat">
            <FaBars/>
          </Button>}
      </div>
      <div className="h-full w-full flex flex-col items-center gap-2 p-4 m-2 border border-default-200 rounded-lg">
        <SidebarItem icon={<FaCode/>} isSidebarOpen={isOpen} link="/services" name="Services">
          {servicesResponse?.map(service => {
            return (
              <Link key={service.nfq} to={`/services/${service.nfq}`} className="px-4 py-1 rounded-lg bg-default-200 hover:opacity-75 border border-default-300">
                <li className="flex items-center gap-2">
                  <span className="text-sm">{service.name}</span>
                  <div className="flex items-center gap-2 ml-auto">
                    <span className="font-mono text-xs">nfq:{service.nfq}</span>
                    <span className="font-mono text-xs text-success">{service.port}</span>
                  </div>
                </li>
              </Link>
            )
          })}
        </SidebarItem>
        <Divider/>
      </div>
      <div className="w-full flex items-center justify-center">
        {isOpen ?
          <Button variant="bordered" className="w-full flex items-center gap-4" onPress={() => navigate("/")}>
            <FaHouse/>
            <span>Go to the landing page</span>
          </Button> :
          <Button isIconOnly={true} variant="bordered" onPress={() => navigate("/")}>
            <FaHouse/>
          </Button>}
      </div>
    </motion.aside>
  )
}
