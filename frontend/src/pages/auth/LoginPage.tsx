import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { Button } from '../../components/ui/buttons/Button';
import { Input } from '../../components/ui/forms/Input';

export default function LoginPage() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('')
    const { login, isLoading } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await login(email, password);
            navigate('/exercises');
        } catch {
            setError('Неверный email или пароль');
        }
    };

    return (
        <div className='max-w-md mx-auto mt-10'>
            <h1 className='text-2xl font-bold mb-6'>Вход</h1>

            {error && (
                <div className='mb-4 p-2 bg-red-100 text-red-700 rounded'>
                    {error}
                </div>
            )}

            <form onSubmit={handleSubmit} className='space-y-4'>
                <Input
                    label="Email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                />

                <Input
                    label="Пароль"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                />

                <Button
                    type='submit'
                    variant='primary'
                    disabled={isLoading}
                    className="w-full"
                >
                    {isLoading ? 'Вход...' : 'Войти'}
                </Button>
            </form>
        </div> 
    );
}