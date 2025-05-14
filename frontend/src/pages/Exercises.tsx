import { ExerciseFilters } from "../components/ExerciseFilters";
import { useExercises } from "../hooks/useExercises";
import { Exercise } from "../models/exercise";
import { useCategories } from "../hooks/useCategories";
import { useRoles } from "../hooks/useRoles";
import { useAuth } from "../hooks/useAuth";
import { hasRole } from "../utils/roleHelpers";
import { API_URL } from "../config";
import { useState } from "react";
import { AddExerciseModal } from "../components/AddExerciseModal";
import { EditExerciseModal } from "../components/EditExerciseModal";
import { Category } from "../models/category";
import { AddCategoryModal } from "../components/AddCategoryModal";
import { EditCategoryModal } from "../components/EditCategoryModal";

export default function Exercises() {
  const {
    exercises,
    loading,
    error,
    filters,
    setFilters,
    page,
    totalPages,
    setPage,
    limit,
    setLimit,
    refetch,
  } = useExercises();
  const {
    categories,
    loading: loadingCategories,
    error: errorCategories,
    createCategory,
    updateCategory,
    deleteCategory,
  } = useCategories();

  const { user, loading: loadingUser } = useAuth();
  const { roles } = useRoles(user?.id || 0);

  const [editingExercise, setEditingExercise] = useState<Exercise | null>(null);
  const [addingExercise, setAddingExercise] = useState(false);
  const [operationError, setOperaionError] = useState<string | null>(null);

  const isAdmin = user ? hasRole(roles, "admin") : false;
  const isModerator = user ? hasRole(roles, "moderator") : false;

  const [addingCategory, setAddingCategory] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);

  const handleAddExercise = async (data: {
    name: string;
    description: string;
    category_id: number;
  }) => {
    try {
      setOperaionError(null);

      const response = await fetch(`${API_URL}/exercises`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(data),
      });

      if (!response.ok) throw new Error("Ошибка добавления");
      refetch();
      setAddingExercise(false);
      return true;
    } catch (err) {
      setOperaionError(
        err instanceof Error ? err.message : "Неизвестная ошибка"
      );
      return false;
    }
  };

  const handleUpdateExercise = async (data: {
    name: string;
    description: string;
    category_id: number;
  }) => {
    if (!editingExercise) return false;

    try {
      setOperaionError(null);

      const response = await fetch(
        `${API_URL}/exercises/${editingExercise.id}`,
        {
          method: "PUT",
          credentials: "include",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(data),
        }
      );

      if (!response.ok) throw new Error("Ошибка обновления");
      refetch();
      return true;
    } catch (err) {
      setOperaionError(
        err instanceof Error ? err.message : "Неизвестная ошибка"
      );
      return false;
    }
  };

  const handleDeleteExercises = async (exerciseID: number) => {
    if (!window.confirm("Вы уверены что хотите удалить это упражнение?")) {
      return;
    }

    try {
      setOperaionError(null);

      const response = await fetch(`${API_URL}/exercises/${exerciseID}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (!response.ok) throw new Error("Ошибка удаления упражнения");
      refetch();
    } catch (err) {
      setOperaionError(
        err instanceof Error ? err.message : "Неизвестная ошибка"
      );
    }
  };

  if (loadingCategories || loadingUser) return <div>Загрузка...</div>;
  if (errorCategories)
    return <div className="text-red-500">{errorCategories}</div>;
  if (error) return <div className="text-red-500">{error}</div>;

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-6">Упражнения</h1>

      {(isAdmin || isModerator) && (
        <button
          onClick={() => setAddingExercise(true)}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 mb-4"
        >
          Добавить упражнение
        </button>
      )}

      {operationError && (
        <div className="mb-4 p-3 bg-red-100 text-red-700 rounded">
          {operationError}
        </div>
      )}

      <ExerciseFilters
        filters={filters}
        onFilterChange={setFilters}
        categories={categories}
      />

      {loading ? (
        <div className="text-center py-8">Загрузка упражнений...</div>
      ) : (
        <>
          <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
            {exercises.map((exercise) => (
              <ExerciseCard
                key={exercise.id}
                exercise={exercise}
                isAdmin={isAdmin}
                isModerator={isModerator}
                onEdit={() => setEditingExercise(exercise)}
                onDelete={() => handleDeleteExercises(exercise.id)}
              />
            ))}
          </div>

          <div className="flex flex-col md:flex-row justify-between items-center mt-6 space-y-4 md:space-y-0">
            <div className="flex items-center space-x-2">
              <label>На странице:</label>
              <select
                value={limit}
                onChange={(e) => setLimit(Number(e.target.value))}
                className="p-2 border rounded"
              >
                {[10, 20, 50, 100].map((opt) => (
                  <option key={opt} value={opt}>
                    {opt}
                  </option>
                ))}
              </select>
            </div>

            <div className="flex justify-center items-center space-x-4">
              <button
                onClick={() => setPage(page - 1)}
                disabled={page === 1}
                className="px-4 py-2 bg-gray-200 rounded disabled:opacity-50"
              >
                Назад
              </button>

              <span>
                Страница {page} из {totalPages}
              </span>

              <button
                onClick={() => setPage(page + 1)}
                disabled={page === totalPages}
                className="px-4 py-2 bg-gray-200 rounded disabled:opacity-50"
              >
                Далее
              </button>
            </div>
          </div>
        </>
      )}

      {addingExercise && (
        <AddExerciseModal
          categories={categories}
          onSubmit={handleAddExercise}
          onClose={() => setAddingExercise(false)}
        />
      )}

      {editingExercise && (
        <EditExerciseModal
          exercise={editingExercise}
          categories={categories}
          onSubmit={handleUpdateExercise}
          onClose={() => setEditingExercise(null)}
        />
      )}

      {isAdmin && (
        <div className="mt-8 p-4 bg-gray-50 rounded">
          <h2 className="text-lg font-semibold mb-4">Категории</h2>
          <button
            onClick={() => setAddingCategory(true)}
            className="mb-4 px-3 py-1 bg-blue-500 text-white rounded"
          >
            Добавить категорию
          </button>

          {categories.map((cat) => (
            <div
              key={cat.id}
              className="flex justify-between items-center mb-2"
            >
              <span>{cat.name}</span>
              <div className="space-x-2">
                <button
                  onClick={() => setEditingCategory(cat)}
                  className="text-blue-500"
                >
                  Редактировать
                </button>
                <button
                  onClick={() => {
                    if (window.confirm("Удалить категорию?")) {
                      deleteCategory(cat.id).catch(console.error);
                    }
                  }}
                  className="text-red-500"
                >
                  Удалить
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {addingCategory && (
        <AddCategoryModal
          onCreate={createCategory}
          onClose={() => setAddingCategory(false)}
        />
      )}
      {editingCategory && (
        <EditCategoryModal
          category={editingCategory}
          onUpdate={updateCategory}
          onClose={() => setEditingCategory(null)}
        />
      )}
    </div>
  );
}

function ExerciseCard({
  exercise,
  isAdmin,
  isModerator,
  onEdit,
  onDelete,
}: {
  exercise: Exercise;
  isAdmin: boolean;
  isModerator: boolean;
  onEdit: () => void;
  onDelete: () => void;
}) {
  return (
    <div className="bg-white rounded-lg shadow p-4">
      <h2 className="font-bold text-lg">{exercise.name}</h2>
      <p className="text-gray-600 mt-2">{exercise.description}</p>

      {(isAdmin || isModerator) && (
        <div>
          <button
            onClick={onEdit}
            className="p-1 text-blue-500 hover:text-blue-700"
            title="Редактировать"
          >
            Редактировать
          </button>
          <button
            onClick={onDelete}
            className="p-1 text-red-500 hover:text-red-700"
            title="Удалить"
          >
            Удалить
          </button>
        </div>
      )}
    </div>
  );
}
