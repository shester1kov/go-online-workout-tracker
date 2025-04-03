import { useState } from "react";
import { Button } from "../../components/ui/buttons/Button";

export default function ExercisePage() {
    const [isLoading, setIsLoading] = useState(false);

    const loadExercises = () => {
        setIsLoading(true);
        console.log("Здесь будет запрос к апи");

        setTimeout(() => {
            setIsLoading(false);
        }, 1000);
    };

    return (
        <div className="container mx-auto px-4 py-8">
            <div className="bg-white rounded-lg shadow-md p-6">
                <h1 className="text-2xl font-bold mb-6">Упражнения</h1>

                <div className="border border-dashed border-gray-300 rounded-lg p-8 text-center">
                    <h2 className="text-xl font-semibold mb-2">Список упражнений</h2>
                    
                    <p className="text-gray-600 mb-6">
                        Список упражнений (после подключения к апи)
                    </p>
                    
                    <div className="space-y-4 mb-6">
                        <div className="border p-3 rounded-lg">
                            <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                            <div className="h-3 bg-gray-100 rounded w-1/2"></div>
                            <div className="border p-3 rounded-lg">
                                <div className="h-4 bg-gray-200 rounded w-2/3 mb-2"></div>
                                <div className="h-3 bg-gray-100 rounded w-1/3"></div>
                            </div>

                            <Button
                                onClick={loadExercises}
                                disabled={isLoading}
                                className="w-full md:w-auto"
                            >
                                {isLoading ? "Загрузка..." : "Загрузить упражнения"}
                            </Button>

                            <p className="text-sm text-gray-500 mt-4">
                                (Кнопка заглушка, иммитирует загрузку)
                            </p>
                        </div>
                    </div>
                </div>
            </div>

        </div>
    )
}