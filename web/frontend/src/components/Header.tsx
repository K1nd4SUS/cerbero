import cerberoPng from "../assets/images/cerbero.png"

export default function Header() {
  return (
    <header className="w-full flex items-center justify-center p-4 bg-default-50 shadow-2xl">
      <a href="/services">
        <img
          src={cerberoPng}
          alt="cerbero"
          className="h-8"
        />
      </a>
    </header>
  )
}
