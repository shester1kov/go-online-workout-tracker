import { SetStateAction } from "react";
import { Category } from "../models/category";
import { IExerciseFilters } from "../models/exerciseFilters";
import { useState, useEffect, useRef } from "react";

type Props = {
  filters: IExerciseFilters;
  onFilterChange: (newFilters: SetStateAction<IExerciseFilters>) => void;
  categories: Category[];
};

export function ExerciseFilters({
  filters,
  onFilterChange,
  categories,
}: Props) {
  const [searchValue, setSearchValue] = useState(filters.name);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleFilterChange = (partialFilters: Partial<IExerciseFilters>) => {
    onFilterChange((prev) => ({
      ...prev,
      ...partialFilters,
    }));
  };

  useEffect(() => {
    const handler = setTimeout(() => {
      onFilterChange((prev) => ({
        ...prev,
        name: searchValue,
      }));
    }, 500);

    return () => {
      clearTimeout(handler);
    };
  }, [searchValue]);

  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.focus();
      inputRef.current.setSelectionRange(
        searchValue.length,
        searchValue.length
      );
    }
  }, []);

  return (
    <div className="bg-white p-4 rounded-lg shadow mb-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label className="block mb-1">Название</label>
          <input
            type="text"
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
            className="w-full p-2 border rounded"
            placeholder="Поиск по названию"
          />
        </div>

        <div>
          <label className="block mb-1">Категория</label>
          <select
            value={filters.category || ""}
            onChange={(e) =>
              handleFilterChange({ category: e.target.value || null })
            }
            className="w-full p-2 border rounded"
          >
            <option value="">Все категории</option>
            {categories.map((category) => (
              <option key={category.id} value={category.id.toString()}>
                {category.name}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block mb-1">Сортировка</label>
          <select
            value={filters.order}
            onChange={(e) =>
              handleFilterChange({ order: e.target.value as "asc" | "desc" })
            }
            className="w-full p-2 border rounded"
          >
            <option value="asc">По возрастанию</option>
            <option value="desc">По убыванию</option>
          </select>
        </div>
      </div>
    </div>
  );
}
