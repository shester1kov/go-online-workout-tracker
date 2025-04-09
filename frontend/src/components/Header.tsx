import { NavLink, useNavigate } from "react-router-dom"
import { useAuth } from "../hooks/useAuth"

export default function Header() {
    const { user, logout } = useAuth()
    const navigate = useNavigate()

    return (
        <header className="bg-white shadow-sm">
            <div className="max-w-6xl mx-auto px-4 py-3 flex justify-between items-center">
                <nav className="flex space-x-4">
                    <NavLink 
                        to='/'
                        className={
                            ({ isActive }) =>
                                `px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 ease-in-out ${
                                    isActive
                                        ? 'bg-blue-500 text-white'
                                        : 'text-gray-700 hover:bg-blue-100 hover:translate-y-[-2px] hover:shadow-sm'
                                }`
                        }
                    >
                        Главная
                    </NavLink>
                    <NavLink 
                        to='/about'
                        className={
                            ({ isActive }) =>
                                `px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 ease-in-out ${
                                    isActive
                                        ? 'bg-blue-500 text-white'
                                        : 'text-gray-700 hover:bg-blue-100 hover:translate-y-[-2px] hover:shadow-sm'
                                }`
                        }
                    >
                        О нас
                    </NavLink>
                    <NavLink 
                        to='/exercises'
                        className={
                            ({ isActive }) =>
                                `px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 ease-in-out ${
                                    isActive
                                        ? 'bg-blue-500 text-white'
                                        : 'text-gray-700 hover:bg-blue-100 hover:translate-y-[-2px] hover:shadow-sm'
                                }`
                        }
                    >
                        Упражнения
                    </NavLink>
                </nav>

                {user ? (
                    <div className="flex items-center space-x-2">
                        <button
                            onClick={() => navigate('/profile')}
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
                            to='/login'
                            className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
                        >
                            Войти
                        </NavLink>
                        <NavLink
                            to='/register'
                            className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
                        >
                            Регистрация
                        </NavLink>
                    </div>
                )}
            </div>
        </header>
    )

    
}