import { useState } from "react";

type Props = {
  onCreate: (name: string, description: string) => Promise<void>;
  onClose: () => void;
};

export function AddCategoryModal({ onCreate, onClose }: Props) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handle = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      await onCreate(name, description);
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
        <h2 className="text-xl font-semibold">Новая категория</h2>
        {error && <p className="text-red-500">{error}</p>}
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Название категории"
          className="w-full p-2 border rounded"
          disabled={loading}
        />
        <input
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Описание категории"
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
            {loading ? "Сохраняем..." : "Создать"}
          </button>
        </div>
      </form>
    </div>
  );
}
