// src/components/DataTable/ColumnFilter.tsx
import React from "react";
import { Column, Table } from "@tanstack/react-table";

interface ColumnFilterProps<TData extends object> {
  column: Column<TData, unknown>;
  table: Table<TData>;
}

function ColumnFilter<TData extends object>({
  column,
  table,
}: ColumnFilterProps<TData>) {
  const columnFilterValue = column.getFilterValue();

  // Filtre değiştiğinde ilk sayfaya dön
  const onFilterChange = (value: string) => {
    column.setFilterValue(value);
    table.setPageIndex(0);
  };

  return (
    <input
      type="text"
      value={(columnFilterValue ?? "") as string}
      onChange={(e) => onFilterChange(e.target.value)}
      placeholder={`Ara...`}
      className="w-full border shadow rounded-sm text-xs mt-1 py-1 px-2 focus:ring-indigo-500 focus:border-indigo-500"
      onClick={(e) => e.stopPropagation()}
    />
  );
}

export default ColumnFilter;
