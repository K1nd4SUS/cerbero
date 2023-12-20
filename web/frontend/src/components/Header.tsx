import cerberoGif from "../assets/images/cerbero.gif"

export default function Header() {
  return (
    <header className="w-full flex items-center justify-center p-2 bg-kinda-primary shadow-2xl">
      <a href="/services">
        <img
          src={cerberoGif}
          alt="cerbero"
          className="h-12"
        />
      </a>
    </header>
  )
}
