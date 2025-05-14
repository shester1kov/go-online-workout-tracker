import { useState } from "react";
import { Category } from "../models/category";

type Props = {
  category: Category;
  onUpdate: (id: number, name: string, description: string) => Promise<void>;
  onClose: () => void;
};

export function EditCategoryModal({ category, onUpdate, onClose }: Props) {
  const [name, setName] = useState(category.name);
  const [description, setDescription] = useState(category.description);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handle = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      await onUpdate(category.id, name, description);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center">
      <form
        onSubmit={handle}
        className="bg-white p-6 rounded shadow-lg w-full max-w-sm space-y-4"
      >
        <h2 className="text-xl font-semibold">Редактировать категорию</h2>
        {error && <p className="text-red-500">{error}</p>}
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="w-full p-2 border rounded"
          disabled={loading}
        />
        <input
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          className="w-full p-2 border rounded"
          disabled={loading}
        />
        <div className="flex justify-end space-x-2">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 border rounded"
            disabled={loading}
          >
            Отмена
          </button>
          <button
            type="submit"
            className="px-4 py-2 bg-blue-500 text-white rounded disabled:opacity-50"
            disabled={loading || !name.trim()}
          >
            {loading ? "Сохраняем..." : "Сохранить"}
          </button>
        </div>
      </form>
    </div>
  );
}
