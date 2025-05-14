import { useState } from "react";
import { useNutritions } from "../hooks/useNutrition";

export default function Nutrition() {
  const today = new Date().toISOString().slice(0, 10);
  const [date, setDate] = useState<string>(today);
  const { data, loading, error } = useNutritions(date);

  const totals = data.reduce(
    (acc, entry) => ({
      calories: acc.calories + entry.calories,
      protein: acc.protein + entry.protein,
      carbs: acc.carbs + entry.carbs,
      fat: acc.fat + entry.fat,
    }),
    { calories: 0, protein: 0, carbs: 0, fat: 0 }
  );

  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold mb-4">Питание за {date}</h1>
      <div className="mb-4">
        <input
          type="date"
          value={date}
          onChange={(e) => setDate(e.target.value)}
          className="border rounded p-2"
        />
      </div>
      {loading && <p>Загрузка...</p>}
      {error && <p className="text-red-500">Ошибка: {error}</p>}
      {!loading && !error && data !== null && (
        <div>
          <div className="bg-white rounded-lg shadow-md p-6 mb-8">
            <h3 className="text-xl font-semibold mb-4">Всего за день</h3>
            <div className="grid grid-cols-4 gap-4 text-center">
              <div className="bg-blue-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600">Калории</p>
                <p className="text-xl font-bold">{totals.calories}</p>
              </div>
              <div className="bg-green-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600">Белки (г)</p>
                <p className="text-xl font-bold">{totals.protein}</p>
              </div>
              <div className="bg-yellow-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600">Углеводы (г)</p>
                <p className="text-xl font-bold">{totals.carbs}</p>
              </div>
              <div className="bg-red-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600">Жиры (г)</p>
                <p className="text-xl font-bold">{totals.fat}</p>
              </div>
            </div>
          </div>

          <div className="border-4 border-white rounded-md shadow-md">
            <table className="min-w-full bg-white">
              <thead>
                <tr>
                  <th className="py-2 px-4 border-b">Продукт</th>
                  <th className="py-2 px-4 border-b">Калории</th>
                  <th className="py-2 px-4 border-b">Белки (г)</th>
                  <th className="py-2 px-4 border-b">Жиры (г)</th>
                  <th className="py-2 px-4 border-b">Углеводы (г)</th>
                </tr>
              </thead>
              <tbody>
                {data.map((entry) => (
                  <tr key={entry.id + entry.calories}>
                    <td className="py-2 px-4 border-b">{entry.food_name}</td>
                    <td className="py-2 px-4 border-b">{entry.calories}</td>
                    <td className="py-2 px-4 border-b">{entry.protein}</td>
                    <td className="py-2 px-4 border-b">{entry.fat}</td>
                    <td className="py-2 px-4 border-b">{entry.carbs}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}
