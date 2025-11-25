function Navbar() {
  return (
    <nav className="w-[50vw] mx-auto flex justify-between items-center py-4 px-4 bg-white">
      <div className="font-semibold text-4xl tracking-tight">VanGo</div>
      {/* Simulate authentication */}
      {false ? (
        <div className="text-base text-gray-700">John Doe</div>
      ) : (
        <div className="flex gap-3">
          <button className="px-4 py-2 bg-white border border-gray-300 text-gray-800 hover:bg-gray-50 transition-colors">
            Login
          </button>
          <button className="px-4 py-2 bg-blue-600 text-white hover:bg-blue-700 transition-colors">
            Signup
          </button>
        </div>
      )}
    </nav>
  )
}

export default Navbar;