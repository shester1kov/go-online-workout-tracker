import { useState } from "react";
import { useParams, Link } from "react-router-dom";
import { useWorkoutDetails } from "../hooks/useWorkoutDetails";
import { useExercises } from "../hooks/useExercises";
import { WorkoutExercise } from "../models/workoutExercise";
import { EditWorkoutExerciseForm } from "../components/EditWorkoutExerciseForm";
import { AddExerciseToWorkoutForm } from "../components/AddExerciseToWorkoutForm";
import jsPDF from "jspdf";
import html2canvas from "html2canvas";
import { API_URL } from "../config";

export default function WorkoutDetails() {
  const { id } = useParams<{ id: string }>();
  const workoutID = parseInt(id || "0");
  const { workout, exercises, loading, error, fetchExercises } =
    useWorkoutDetails(workoutID);
  const [editingExercise, setEditingExercise] =
    useState<WorkoutExercise | null>(null);
  const [addingExercise, setAddingExercise] = useState(false);
  const [operationError, setOperaionError] = useState<string | null>(null);
  const {
    exercises: allExercises,
    loading: loadingExercises,
    filters,
    setFilters,
    page,
    totalPages,
    setPage,
    limit,
    setLimit,
  } = useExercises();

  const handleAddExercise = async (data: {
    sets: number;
    reps: number;
    weight: number;
    notes: string;
    exercise_id: number;
  }) => {
    try {
      setOperaionError(null);

      const response = await fetch(
        `${API_URL}/workouts/${workoutID}/exercises`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify(data),
        }
      );

      if (!response.ok) throw new Error("Ошибка добавления");
      await fetchExercises();
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
    sets: number;
    reps: number;
    weight: number;
    notes: string;
    exercise_id: number;
  }) => {
    if (!editingExercise) return false;

    try {
      setOperaionError(null);

      const response = await fetch(
        `${API_URL}/workouts/${workoutID}/exercises/${editingExercise.id}`,
        {
          method: "PUT",
          credentials: "include",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(data),
        }
      );

      if (!response.ok) throw new Error("Ошибка обновления");
      await fetchExercises();
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

      const response = await fetch(
        `${API_URL}/workouts/${workoutID}/exercises/${exerciseID}`,
        {
          method: "DELETE",
          credentials: "include",
        }
      );

      if (!response.ok) throw new Error("Ошибка удаления упражнения");
      await fetchExercises();
    } catch (err) {
      setOperaionError(
        err instanceof Error ? err.message : "Неизвестная ошибка"
      );
    }
  };

  if (loading) return <div>Загрузка...</div>;
  if (error) return <div className="text-red-500">{error}</div>;
  if (!workout) return <div>Тренировка не найдена</div>;

  const dateOnly = workout.date.split("T")[0];

  const PdfContent = () => (
    <div
      id="pdf-export-content"
      style={{
        padding: "24px",
        maxWidth: "800px",
        margin: "0 auto",
        color: "#000",
        backgroundColor: "#fff",
        fontFamily: "Arial, sans-serif",
      }}
    >
      <h1
        style={{
          fontSize: "24px",
          fontWeight: "bold",
          marginBottom: "24px",
          textAlign: "center",
        }}
      >
        Тренировка от {dateOnly}
      </h1>

      {workout.notes && (
        <section style={{ marginBottom: "32px" }}>
          <h2
            style={{
              fontSize: "20px",
              fontWeight: "600",
              borderBottom: "1px solid #ccc",
              paddingBottom: "4px",
              marginBottom: "8px",
            }}
          >
            Описание
          </h2>
          <p style={{ fontSize: "16px", lineHeight: "1.5" }}>{workout.notes}</p>
        </section>
      )}

      <section>
        <h2
          style={{
            fontSize: "20px",
            fontWeight: "600",
            borderBottom: "1px solid #ccc",
            paddingBottom: "4px",
            marginBottom: "16px",
          }}
        >
          Упражнения
        </h2>

        {exercises.length === 0 ? (
          <p style={{ color: "#555" }}>Нет упражнений</p>
        ) : (
          <ul style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
            {exercises.map((ex) => (
              <li
                key={ex.id}
                style={{
                  border: "1px solid #ccc",
                  borderRadius: "8px",
                  padding: "16px",
                  backgroundColor: "#f9f9f9",
                }}
              >
                <h3
                  style={{
                    fontWeight: "bold",
                    fontSize: "18px",
                    marginBottom: "8px",
                  }}
                >
                  {ex.exercise?.name}
                </h3>
                <p style={{ marginBottom: "4px" }}>
                  <strong>Подходы:</strong> {ex.sets}
                </p>
                <p style={{ marginBottom: "4px" }}>
                  <strong>Повторения:</strong> {ex.reps}
                </p>
                <p style={{ marginBottom: "4px" }}>
                  <strong>Вес:</strong> {ex.weight} кг
                </p>
                {ex.notes && (
                  <p style={{ marginTop: "8px" }}>
                    <strong>Заметки:</strong> {ex.notes}
                  </p>
                )}
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );

  const handleExportPDF = async () => {
    // Clone the PDF content into a visible off-screen container
    const original = document.getElementById("pdf-export-content");
    if (!original) return;

    // Create a wrapper to hold cloned content
    const wrapper = document.createElement("div");
    wrapper.style.position = "absolute";
    wrapper.style.top = "-9999px";
    wrapper.style.left = "-9999px";
    document.body.appendChild(wrapper);

    // Clone node and append
    const clone = original.cloneNode(true) as HTMLElement;
    clone.style.display = "block";
    wrapper.appendChild(clone);

    try {
      const canvas = await html2canvas(clone, {
        scale: 2,
        backgroundColor: "#fff",
      });
      const imgData = canvas.toDataURL("image/png");

      // Determine PDF page size (A4 in px at 96dpi)
      const pdf = new jsPDF({
        orientation: "portrait",
        unit: "px",
        format: [canvas.width, canvas.height],
      });

      pdf.addImage(imgData, "PNG", 0, 0, canvas.width, canvas.height);
      pdf.save(`тренировка-${dateOnly}.pdf`);
    } catch (error) {
      console.error("Ошибка при создании PDF:", error);
      setOperaionError("Не удалось создать PDF. Попробуйте ещё раз.");
    } finally {
      // Clean up off-screen wrapper
      document.body.removeChild(wrapper);
    }
  };

  // ... return JSX with <div><PdfContent style={{display: 'none'}} /> ... and export button ...>

  return (
    <div className="container mx-auto p-4">
      <Link
        to="/workouts"
        className="text-blue-500 hover:underline mb-4 inline-block"
      >
        Назад к списку тренировок
      </Link>

      <h1 className="text-2xl font-bold mb-2">Тренировка {dateOnly}</h1>

      <div className="bg-white rounded-lg shadow p-4 mb-6">
        <h2 className="text-lg font-semibold mb-2">Описание</h2>

        <p className="text-gray-600">
          {workout.notes || <span className="text-gray-400">Нет описания</span>}
        </p>
      </div>

      {operationError && (
        <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4">
          {operationError}
        </div>
      )}

      <h2 className="text-xl font-bold mb-4">Упражнения</h2>

      <div className="hidden">
        <PdfContent />
      </div>

      <div className="flex space-x-0.5">
        <button
          onClick={() => setAddingExercise(true)}
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 mb-4"
        >
          Добавить упражнение
        </button>

        <button
          onClick={handleExportPDF}
          className="bg-blue-900 text-white px-4 py-2 rounded hover:bg-blue-950 mb-4"
        >
          Сохранить в PDF
        </button>
      </div>

      {exercises.length === 0 ? (
        <p className="text-gray-500">Нет упражнений в этой трениировке</p>
      ) : (
        <div className="grid gap-4">
          {exercises.map((exercise) => (
            <div
              key={exercise.id}
              className="bg-white rounded-lg shadow p-4 mb-4"
            >
              <h3 className="font-bold text-lg">{exercise.exercise?.name}</h3>

              <div className="flex space-x-2">
                <button
                  onClick={() => setEditingExercise(exercise)}
                  className="text-blue-500 hover:text-blue-700"
                >
                  Редактировать
                </button>

                <button
                  onClick={() => handleDeleteExercises(exercise.id)}
                  className="text-red-500 hover:text-red-700"
                >
                  Удалить
                </button>
              </div>

              <p className="text-gray-600">{exercise.notes}</p>

              <div className="mt-2 grid grid-cols-3 gap-2">
                <div>
                  <span className="text-sm text-gray-500">Подходы</span>
                  <p>{exercise.sets}</p>
                </div>
                <div>
                  <span className="text-sm text-gray-500">Повторения</span>
                  <p>{exercise.reps}</p>
                </div>
                <div>
                  <span className="text-sm text-gray-500">Вес (кг)</span>
                  <p>{exercise.weight}</p>
                </div>
              </div>
            </div>
          ))}
          {editingExercise && (
            <EditWorkoutExerciseForm
              exercise={editingExercise}
              onSave={handleUpdateExercise}
              onClose={() => setEditingExercise(null)}
              allExercises={allExercises}
              loading={loadingExercises}
              filters={filters}
              setFilters={setFilters}
              page={page}
              setPage={setPage}
              totalPages={totalPages}
              limit={limit}
              setLimit={setLimit}
            />
          )}
        </div>
      )}

      {addingExercise && (
        <AddExerciseToWorkoutForm
          onAdd={handleAddExercise}
          onClose={() => setAddingExercise(false)}
          allExercises={allExercises}
          loading={loadingExercises}
          filters={filters}
          setFilters={setFilters}
          page={page}
          setPage={setPage}
          totalPages={totalPages}
          limit={limit}
          setLimit={setLimit}
        />
      )}
    </div>
  );
}
