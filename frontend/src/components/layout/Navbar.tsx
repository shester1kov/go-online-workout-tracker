import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext'
import { Button } from '../ui/buttons/Button';

export default function Navbar() {
    const { isAuthenticated, logout } = useAuth();
    const location = useLocation();

    const navLinks = [
        { path: '/', label: "Главная" },
        { path: '/exercises', label: 'Упражнения' },
        { path: '/about', label: 'О проекте' },
    ];

    return (
        <header className='bg-white shadow-sm'>
            <div className='container mx-auto px-4'>
                <nav className='flex items-center justify-between h-16'>
                    <div className='flex space-x-8'>
                    {navLinks.map((link) => (
                        <Link
                            key={link.path}
                            to={link.path}
                            className={`px-3 py-2 text-sm font-medium ${
                                location.pathname === link.path
                                    ? "text-blue-600 border-b-2 border-blue-600"
                                    : "text-gray-500 hover:text-gray-700"
                            }`}
                            >
                                {link.label}
                            </Link>
                    ))}
                    </div>

                    <div className='flex items-center space-x-4'>
                        {isAuthenticated ? (
                            <Button
                            onClick={logout}
                            >
                                Выйти
                            </Button>
                        ) : (
                            <>
                                <Link
                                    to='/login'
                                    className='px-4 py-2 text-sm text-gray-500 hover:text-gray-700'
                                >
                                    Вход
                                </Link>
                                <Link  
                                    to='/register'
                                    className='px-4 py-2 text-sm bg-blue-600 text-white rounded hover:bg-blue-700'
                                >
                                    Регистрация   
                                </Link>
                            </>
                        )}
                    </div>
                </nav>
            </div>
        </header>
    );
}