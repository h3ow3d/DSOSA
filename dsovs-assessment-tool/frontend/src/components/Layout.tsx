import { Link, NavLink, Outlet } from "react-router-dom";

export default function Layout() {
  const navClass = ({ isActive }: { isActive: boolean }) =>
    `block px-3 py-2 rounded text-sm font-medium transition-colors ${
      isActive
        ? "bg-blue-700 text-white"
        : "text-gray-300 hover:bg-gray-700 hover:text-white"
    }`;

  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      <nav className="bg-gray-900 text-white shadow">
        <div className="max-w-7xl mx-auto px-4 flex items-center gap-6 h-14">
          <Link to="/dashboard" className="font-bold text-lg tracking-tight text-blue-400">
            DSOVS Assessment
          </Link>
          <NavLink to="/dashboard" className={navClass}>
            Dashboard
          </NavLink>
          <NavLink to="/projects" className={navClass}>
            Projects
          </NavLink>
        </div>
      </nav>
      <main className="flex-1 max-w-7xl mx-auto w-full px-4 py-6">
        <Outlet />
      </main>
    </div>
  );
}
