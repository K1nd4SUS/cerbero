import { Button, Link, Tooltip } from "@nextui-org/react"
import { ReactNode, useState } from "react"
import { FaCaretDown, FaCaretUp } from "react-icons/fa6"

export type SidebarItemProps = {
  children: ReactNode
  icon: ReactNode
  isSidebarOpen: boolean
  link: string
  name: string
}

export default function SidebarItem({ children, icon, isSidebarOpen, link, name }: SidebarItemProps) {
  const [isOpen, setIsOpen] = useState(false)

  if(!isSidebarOpen) {
    return (
      <Tooltip content="Services" size="sm" delay={1000}>
        <Button as={Link} href={link} isIconOnly={true} variant="bordered" size="sm">
          {icon}
        </Button>
      </Tooltip>
    )
  }

  return (
    <div className="w-full flex flex-col gap-2">
      <div className="flex items-center">
        <Link color="foreground" href={link} className="flex items-center gap-2">
          {icon}
          <span className="font-bold">{name}</span>
        </Link>
        {isOpen ?
          <Button isIconOnly={true} onPress={() => setIsOpen(false)} size="sm" className="ml-auto">
            <FaCaretUp/>
          </Button> :
          <Button isIconOnly={true} onPress={() => setIsOpen(true)} size="sm" className="ml-auto">
            <FaCaretDown/>
          </Button>}
      </div>
      {isOpen ?
        <ul className="flex flex-col gap-1">
          {children}
        </ul> : <></>}
    </div>
  )
}