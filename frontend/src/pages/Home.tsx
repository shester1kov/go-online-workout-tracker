export default function Home() {
    return (
        <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-4">
                Добро пожаловать!
            </h1>
            <p className="text-gray-600 mb-4">
                Это главная страница приложения. Здесь вы найдете основную информацию.
            </p>
            <div className="bg-blue-100 p-4 rounded-lg border border-blue-200" >
                <p className="text-blue-800">
                    Перейдите на страницу "О нас" используя навигацию выше.
                </p>
            </div>
        </div>
    )
}