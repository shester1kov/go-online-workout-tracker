import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { useState } from "react";
import { Menu, X } from "lucide-react";

export default function Header() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [isOpen, setIsOpen] = useState(false);

  const toggleMenu = () => setIsOpen(!isOpen);
  const closeMenu = () => setIsOpen(false);

  const navLinks = (
    <>
      <NavLink to="/" onClick={closeMenu} className={navLinkStyle}>
        Главная
      </NavLink>
      <NavLink to="/about" onClick={closeMenu} className={navLinkStyle}>
        О нас
      </NavLink>
      <NavLink to="/exercises" onClick={closeMenu} className={navLinkStyle}>
        Упражнения
      </NavLink>
      <NavLink to="/workouts" onClick={closeMenu} className={navLinkStyle}>
        Тренировки
      </NavLink>
      <NavLink to="/nutritions" onClick={closeMenu} className={navLinkStyle}>
        Питание
      </NavLink>
    </>
  );

  return (
    <header className="bg-white shadow-sm">
      <div className="max-w-6xl mx-auto px-4 py-3 flex justify-between items-center">
        <div className="md:hidden">
          <button onClick={toggleMenu} className="text-gray-700">
            {isOpen ? <X size={24} /> : <Menu size={24} />}
          </button>
        </div>

        <nav className="hidden md:flex space-x-4">
          <NavLink
            to="/"
            className={({ isActive }) =>
              `px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 ease-in-out ${
                isActive
                  ? "bg-blue-500 text-white"
                  : "text-gray-700 hover:bg-blue-100 hover:translate-y-[-2px] hover:shadow-sm"
              }`
            }
          >
            Главная
          </NavLink>
          <NavLink
            to="/about"
            className={({ isActive }) =>
              `px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 ease-in-out ${
                isActive
                  ? "bg-blue-500 text-white"
                  : "text-gray-700 hover:bg-blue-100 hover:translate-y-[-2px] hover:shadow-sm"
              }`
            }
          >
            О нас
          </NavLink>
          <NavLink
            to="/exercises"
            className={({ isActive }) =>
              `px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 ease-in-out ${
                isActive
                  ? "bg-blue-500 text-white"
                  : "text-gray-700 hover:bg-blue-100 hover:translate-y-[-2px] hover:shadow-sm"
              }`
            }
          >
            Упражнения
          </NavLink>
          <NavLink
            to="/workouts"
            className={({ isActive }) =>
              `px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 ease-in-out ${
                isActive
                  ? "bg-blue-500 text-white"
                  : "text-gray-700 hover:bg-blue-100 hover:translate-y-[-2px] hover:shadow-sm"
              }`
            }
          >
            Тренировки
          </NavLink>
          <NavLink
            to="/nutritions"
            className={({ isActive }) =>
              `px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 ease-in-out ${
                isActive
                  ? "bg-blue-500 text-white"
                  : "text-gray-700 hover:bg-blue-100 hover:translate-y-[-2px] hover:shadow-sm"
              }`
            }
          >
            Питание
          </NavLink>
        </nav>

        <div className="hidden md:flex items-center space-x-4">
          {user ? (
            <div className="flex items-center space-x-2">
              <button
                onClick={() => navigate("/profile")}
                className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
              >
                {user.username}
              </button>
              <button
                onClick={logout}
                className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
              >
                Выйти
              </button>
            </div>
          ) : (
            <div className="flex items-center space-x-2">
              <NavLink
                to="/login"
                className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
              >
                Войти
              </NavLink>
              <NavLink
                to="/register"
                className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
              >
                Регистрация
              </NavLink>
            </div>
          )}
        </div>
      </div>

      {isOpen && (
        <div className="md:hidden px-4 pb-4 space-y-2">
          <nav className="flex flex-col space-y-2">{navLinks}</nav>
          <div className="mt-2 border-t pt-2 flex flex-col space-y-2">
            {user ? (
              <>
                <button
                  onClick={() => {
                    navigate("/profile");
                    closeMenu();
                  }}
                  className={navButtonStyle}
                >
                  {user.username}
                </button>
                <button
                  onClick={() => {
                    logout();
                    closeMenu();
                  }}
                  className={navButtonStyle}
                >
                  Выйти
                </button>
              </>
            ) : (
              <>
                <NavLink
                  to="/login"
                  onClick={closeMenu}
                  className={navButtonStyle}
                >
                  Войти
                </NavLink>
                <NavLink
                  to="/register"
                  onClick={closeMenu}
                  className={navButtonStyle}
                >
                  Регистрация
                </NavLink>
              </>
            )}
          </div>
        </div>
      )}
    </header>
  );
}

const navButtonStyle =
  "text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium";

const navLinkStyle = ({ isActive }: { isActive: boolean }) =>
  `px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 ease-in-out ${
    isActive ? "bg-blue-500 text-white" : "text-gray-700 hover:bg-blue-100"
  }`;
