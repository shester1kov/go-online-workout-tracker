import { Link, useNavigate } from "react-router-dom";
import { useState } from "react";
import { useAuth } from "../hooks/useAuth";


export default function Register() {
    const [email, setEmail] = useState('')
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const { register } = useAuth()
    const navigate = useNavigate()

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setError('')

        try {
            await register(email, password, username);
            navigate('/')
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Ошибка регистрации')
        }
    }


    return (
        <div className="max-w-md mx-auto mt-10 p-6 bg-white rounded-lg shadow-md">
            <h1 className="text-2xl font-bold mb-6 text-center">
                Регистрация
            </h1>

            {error && (
                <div className="mb-4 p-2 bg-red-100 text-red-700 rounded text-sm">
                    {error}
                </div>
            )}

            <form onSubmit={handleSubmit} className='space-y-4'>
                <div>
                    <label htmlFor='email' className='block mb-1 font-medium'>
                        Email
                    </label>
                    <input
                        id='email'
                        type='email'
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        className='w-full px-3 py-2 border rounded-md'
                        required
                    />
                </div>

                <div>
                    <label htmlFor='email' className='block mb-1 font-medium'>
                        Имя пользователя
                    </label>
                    <input
                        id='username'
                        type='text'
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        className='w-full px-3 py-2 border rounded-md'
                        required
                    />
                </div>

                <div>
                    <label htmlFor='password' className='block mb-1 font-medium'>
                        Пароль
                    </label>
                    <input
                        id='password'
                        type='password'
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className='w-full px-3 py-2 border rounded-md'
                        required
                    />
                </div>

                <button
                    type='submit'
                    className='w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700'
                >
                    Зарегистрироваться
                </button>
            </form>
            <div className="mt-4 text-center">
                Уже есть аккаунт?{' '}
                <Link
                    to='/login'
                    className="text-blue-600 hover:underline"
                >
                    Войдите 
                </Link>
            </div>
        </div>
    )
}