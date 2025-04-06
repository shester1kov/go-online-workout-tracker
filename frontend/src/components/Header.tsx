import { Link } from "react-router-dom"
import { useAuth } from "../hooks/useAuth"

export default function Header() {
    const { user, logout } = useAuth()

    return (
        <header className="bg-white shadow-sm">
            <div className="max-w-6xl mx-auto px-4 py-3 flex justify-between items-center">
                <div className="flex space-x-4">
                    <Link 
                        to='/'
                        className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
                    >
                        Главная
                    </Link>
                    <Link 
                        to='/about'
                        className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
                    >
                        О нас
                    </Link>
                    <Link 
                        to='/exercises'
                        className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
                    >
                        Упражнения
                    </Link>
                </div>

                <div>
                    {user ? (
                        <button
                            onClick={logout}
                            className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
                        >
                            Выйти
                        </button>
                    ) : (
                        <>
                            <Link
                                to='/login'
                                className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
                            >
                                Войти
                            </Link>
                            <Link
                                to='/register'
                                className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
                            >
                                Регистрация
                            </Link>
                        </>
                    )}
                </div>
            </div>
        </header>
    )

    
}