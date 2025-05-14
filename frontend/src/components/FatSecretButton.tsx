import { useState, useEffect } from 'react';
import { API_URL } from '../config';

export default function FatSecretButton() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [jQueryLoaded, setjQueryLoaded] = useState(false);

  useEffect(() => {
    const loadjQuery = async () => {
      try {

        const script = document.createElement('script');
        script.src = 'https://code.jquery.com/jquery-3.6.0.min.js';
        script.async = true;


        script.onload = () => {
          setjQueryLoaded(true);
        };

        document.head.appendChild(script);
      } catch (err) {
        console.error('Ошибка при загрузке jQuery:', err);
      }
    };

    loadjQuery();
  }, []);

  const handleConnect = async () => {
    try {
      setLoading(true);
      setError('');
      
      if (!jQueryLoaded) {
        throw new Error('jQuery не загружен');
      }

      window.open(`${API_URL}/connect/fatsecret`, "_blank");
    } catch  {
      setError('Ошибка при открытии окна авторизации');
    } finally {
      setLoading(false);
    }
  };


  return (
    <div>
      <button
        onClick={handleConnect}
        disabled={loading}
        className="bg-blue-500 text-white px-4 py-2 rounded"
      >
        {loading ? 'Подключение...' : 'Подключить FatSecret'}
      </button>
      {error && <p className="text-red-500 mt-2">{error}</p>}
    </div>
  );
}


// import { useState } from 'react';

// export default function FatSecretButton() {
//   const [loading, setLoading] = useState(false);
//   const [error, setError] = useState('');

//   const handleConnect = async () => {
//     try {
//       setLoading(true);
//       setError('');

//       window.open("http://localhost:8080/api/v1/connect/fatsecret", "_blank");

//     //   window.location.href = 'http://localhost:8080/api/v1/connect/fatsecret'

//     //   window.open(
//     //     'http://localhost:8080/api/v1/connect/fatsecret',
//     //     '_blank',
//     //     'width=600,height=600'
//     //   );

//     // await fetch('http://localhost:8080/api/v1/connect/fatsecret', {
//     //     method: "GET",
//     //     credentials: 'include'
//     // })

//     } catch {
//       setError('Ошибка при открытии окна авторизации');
//     } finally {
//       setLoading(false);
//     }
//   };