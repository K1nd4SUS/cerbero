import { Button } from "@nextui-org/react"
import { motion } from "framer-motion"
import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { FaBars, FaHouse } from "react-icons/fa6"

export default function Sidebar() {
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
      <div className="h-full"></div>
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
