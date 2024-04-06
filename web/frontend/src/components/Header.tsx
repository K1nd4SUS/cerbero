import { Button } from "@nextui-org/react"
import { useEffect } from "react"
import { FaCheck, FaFire, FaTriangleExclamation } from "react-icons/fa6"
import cerberoPng from "../assets/images/cerbero.png"
import { useFetch } from "../hooks/useFetch"

export default function Header() {
  const [
    firewall,
    triggerFwFetch,
    ,
  ] = useFetch<{ isConnected: boolean, isSynced: boolean }>()

  const [
    ,
    triggerFwUpdate,
    ,
  ] = useFetch()

  async function handleFwUpdate() {
    await triggerFwUpdate("/api/firewall", {
      method: "POST"
    })
  }

  useEffect(() => {
    void triggerFwFetch("/api/firewall")

    async function intervalCallback() {
      await triggerFwFetch("/api/firewall")
    }

    setInterval(intervalCallback, 5000)
  }, [])

  return (
    <header className="w-full flex items-center justify-center p-4 bg-default-50 shadow-2xl">
      <a href="/services" className="absolute">
        <img
          src={cerberoPng}
          alt="cerbero"
          className="h-8"
        />
      </a>
      <div className="ml-auto">
        {firewall?.isConnected ?
          firewall?.isSynced ?
            <Button color="success" variant="flat" disabled>
              <FaCheck/>
              <span className="font-bold">Firewall synced</span>
            </Button> :
            <Button onPress={() => void handleFwUpdate()} color="warning" variant="flat" className="ml-auto">
              <FaFire/>
              <span className="font-bold">Apply changes</span>
            </Button> :
          <Button color="danger" variant="flat" disabled>
            <FaTriangleExclamation/>
            <span className="font-bold">Firewall not connected</span>
          </Button>}
      </div>
    </header>
  )
}
