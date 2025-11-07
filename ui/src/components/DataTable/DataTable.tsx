// src/components/DataTable/DataTable.tsx
import React from "react";
import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  ColumnDef,
  VisibilityState,
  ColumnFiltersState,
  SortingState,
  RowSelectionState,
  ColumnSizingState,
  PaginationState,
} from "@tanstack/react-table";
import ColumnFilter from "./ColumnFilter";
import { exportToCSV, exportToExcel } from "../../utils/exportData";
import {
  saveTablePreferences,
  loadTablePreferences,
  clearTablePreferences,
  TableUserPreferences,
} from "../../utils/tableStorage";

export type CustomColumnDef<TData extends object> = ColumnDef<TData> & {
  meta?: {
    isSticky?: boolean;
    noExport?: boolean;
  };
};

interface DataTableProps<TData extends object> {
  columns: CustomColumnDef<TData>[];
  onRowSelectionChange?: (selectedRows: TData[]) => void;
  tableId: string;
  fetchFn: (params: {
    pagination: PaginationState;
    sorting: SortingState;
    columnFilters: ColumnFiltersState;
    globalFilter: string;
  }) => Promise<{
    data: TData[];
    pageCount: number;
    totalRowCount: number;
  }>;
}

function DataTable<TData extends object>({
  columns,
  onRowSelectionChange,
  tableId,
  fetchFn,
}: DataTableProps<TData>) {
  const [sorting, setSorting] = React.useState<SortingState>([]);
  const [globalFilter, setGlobalFilter] = React.useState<string>("");
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    []
  );
  const [columnVisibility, setColumnVisibility] =
    React.useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = React.useState<RowSelectionState>({});
  const [columnSizing, setColumnSizing] = React.useState<ColumnSizingState>({});
  const [pagination, setPagination] = React.useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  });
  const [serverData, setServerData] = React.useState<TData[]>([]);
  const [serverPageCount, setServerPageCount] = React.useState(0);
  const [serverTotalRowCount, setServerTotalRowCount] = React.useState(0);
  const [isLoading, setIsLoading] = React.useState(true);

  const table = useReactTable({
    data: serverData,
    columns: [
      {
        id: "select",
        header: ({ table }) => (
          <input
            type="checkbox"
            checked={table.getIsAllRowsSelected()}
            indeterminate={table.getIsSomeRowsSelected()}
            onChange={table.getToggleAllRowsSelectedHandler()}
            className="form-checkbox h-4 w-4 text-indigo-600 transition duration-150 ease-in-out rounded-sm"
          />
        ),
        cell: ({ row }) => (
          <input
            type="checkbox"
            checked={row.getIsSelected()}
            disabled={!row.getCanSelect()}
            indeterminate={row.getIsSomeSelected()}
            onChange={row.getToggleSelectedHandler()}
            className="form-checkbox h-4 w-4 text-indigo-600 transition duration-150 ease-in-out rounded-sm"
          />
        ),
        enableSorting: false,
        enableHiding: false,
        enableColumnFilter: false,
        meta: { isSticky: true },
        size: 50,
      } as CustomColumnDef<TData>,
      ...columns,
    ],
    state: {
      sorting,
      globalFilter,
      columnFilters,
      columnVisibility,
      rowSelection,
      columnSizing,
      pagination,
    },
    onSortingChange: setSorting,
    onGlobalFilterChange: setGlobalFilter,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    onColumnSizingChange: setColumnSizing,
    onPaginationChange: setPagination,
    getCoreRowModel: getCoreRowModel(),
    manualPagination: true,
    manualSorting: true,
    manualFiltering: true,
    pageCount: serverPageCount,
    enableRowSelection: true,
    columnResizeMode: "onChange",
    debugTable: false,
  });

  React.useEffect(() => {
    const savedPreferences = loadTablePreferences(tableId);
    if (savedPreferences) {
      savedPreferences.columnVisibility &&
        setColumnVisibility(savedPreferences.columnVisibility);
      savedPreferences.columnFilters &&
        setColumnFilters(savedPreferences.columnFilters);
      savedPreferences.sorting && setSorting(savedPreferences.sorting);
      savedPreferences.columnSizing &&
        setColumnSizing(savedPreferences.columnSizing);
      savedPreferences.globalFilter &&
        setGlobalFilter(savedPreferences.globalFilter);
      savedPreferences.pagination && setPagination(savedPreferences.pagination);
    }
  }, [tableId]);

  const fetchData = React.useCallback(async () => {
    setIsLoading(true);
    setRowSelection({});
    try {
      const response = await fetchFn({
        pagination,
        sorting,
        columnFilters,
        globalFilter,
      });
      setServerData(response.data as TData[]);
      setServerPageCount(response.pageCount);
      setServerTotalRowCount(response.totalRowCount);
    } catch (error) {
      console.error("Veri çekilirken hata oluştu:", error);
    } finally {
      setIsLoading(false);
    }
  }, [pagination, sorting, columnFilters, globalFilter, fetchFn]);

  React.useEffect(() => {
    fetchData();
  }, [fetchData]);

  React.useEffect(() => {
    if (onRowSelectionChange) {
      const selectedRows = table
        .getSelectedRowModel()
        .rows.map((row) => row.original);
      onRowSelectionChange(selectedRows);
    }
  }, [rowSelection, onRowSelectionChange, serverData]);

  const handleSavePreferences = () => {
    const preferences: TableUserPreferences = {
      columnVisibility: table.getState().columnVisibility,
      columnFilters: table.getState().columnFilters,
      sorting: table.getState().sorting,
      columnSizing: table.getState().columnSizing,
      pagination: table.getState().pagination,
      globalFilter: table.getState().globalFilter,
    };
    saveTablePreferences(tableId, preferences);
  };

  const handleClearPreferences = () => {
    clearTablePreferences(tableId);
    setColumnVisibility({});
    setColumnFilters([]);
    setSorting([]);
    setColumnSizing({});
    setPagination({ pageIndex: 0, pageSize: 10 });
    setGlobalFilter("");
    table.resetRowSelection();
  };

  const stickyColumns = table
    .getVisibleFlatColumns()
    .filter((c) => c.columnDef.meta?.isSticky);
  const stickyColumnOffsets = stickyColumns.reduce((acc, col, i) => {
    const prevWidth =
      i > 0 ? acc[stickyColumns[i - 1].id] + stickyColumns[i - 1].getSize() : 0;
    acc[col.id] = prevWidth;
    return acc;
  }, {} as Record<string, number>);

  return (
    <div className="p-2">
      <div className="flex justify-between items-center mb-4 flex-wrap gap-2">
        <input
          type="text"
          value={globalFilter ?? ""}
          onChange={(e) => {
            setGlobalFilter(e.target.value);
            table.setPageIndex(0);
          }}
          placeholder="Genel Arama..."
          className="border p-2 w-full md:w-1/3 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
        />
        <div className="flex gap-2 flex-wrap">
          <button
            onClick={() =>
              exportToCSV(
                table.getRowModel().rows,
                table.getAllColumns(),
                `${tableId}.csv`
              )
            }
            className="px-4 py-2 bg-gray-600 text-white text-sm font-medium rounded-md hover:bg-gray-700 cursor-pointer transition-colors"
          >
            CSV Aktar
          </button>
          <button
            onClick={() =>
              exportToExcel(
                table.getRowModel().rows,
                table.getAllColumns(),
                `${tableId}.xlsx`
              )
            }
            className="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 cursor-pointer transition-colors"
          >
            Excel Aktar
          </button>
          <button
            onClick={handleSavePreferences}
            className="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-md hover:bg-indigo-700 cursor-pointer transition-colors"
          >
            Ayarları Kaydet
          </button>
          <button
            onClick={handleClearPreferences}
            className="px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700 cursor-pointer transition-colors"
          >
            Ayarları Sıfırla
          </button>
        </div>
      </div>

      <div className="mb-4 text-sm text-gray-700 flex justify-between">
        <span>
          {Object.keys(rowSelection).length > 0 &&
            `${Object.keys(rowSelection).length} satır seçildi.`}
        </span>
        <span>Toplam Kayıt: {serverTotalRowCount}</span>
      </div>

      <div className="overflow-x-auto relative">
        <table
          className="min-w-full divide-y divide-gray-200 shadow-md rounded-lg"
          style={{ width: table.getTotalSize() }}
        >
          <thead className="bg-gray-50">
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  const isSticky = header.column.columnDef.meta?.isSticky;
                  return (
                    <th
                      key={header.id}
                      colSpan={header.colSpan}
                      className={`relative px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider ${
                        isSticky ? "sticky-column sticky-column-header" : ""
                      }`}
                      style={{
                        width: header.getSize(),
                        left: isSticky
                          ? stickyColumnOffsets[header.column.id]
                          : undefined,
                        zIndex: isSticky ? 20 : undefined,
                      }}
                    >
                      {header.isPlaceholder ? null : (
                        <>
                          <div
                            className={
                              header.column.getCanSort()
                                ? "cursor-pointer select-none hover:opacity-70 transition-opacity"
                                : ""
                            }
                            onClick={header.column.getToggleSortingHandler()}
                          >
                            {flexRender(
                              header.column.columnDef.header,
                              header.getContext()
                            )}
                            {{ asc: " 🔼", desc: " 🔽" }[
                              header.column.getIsSorted() as string
                            ] ?? null}
                          </div>
                          {header.column.getCanFilter() && (
                            <ColumnFilter
                              column={header.column}
                              table={table}
                            />
                          )}
                          {header.column.getCanResize() && (
                            <div
                              onMouseDown={header.getResizeHandler()}
                              onTouchStart={header.getResizeHandler()}
                              className={`absolute top-0 right-0 h-full w-2 cursor-col-resize select-none touch-none ${
                                header.column.getIsResizing()
                                  ? "bg-indigo-500 opacity-70"
                                  : ""
                              }`}
                            />
                          )}
                        </>
                      )}
                    </th>
                  );
                })}
              </tr>
            ))}
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {isLoading ? (
              <tr>
                <td
                  colSpan={table.getAllColumns().length}
                  className="text-center p-4"
                >
                  Yükleniyor...
                </td>
              </tr>
            ) : table.getRowModel().rows.length === 0 ? (
              <tr>
                <td
                  colSpan={table.getAllColumns().length}
                  className="text-center p-4"
                >
                  Veri bulunamadı.
                </td>
              </tr>
            ) : (
              table.getRowModel().rows.map((row) => (
                <tr
                  key={row.id}
                  className={row.getIsSelected() ? "bg-indigo-50" : ""}
                >
                  {row.getVisibleCells().map((cell) => {
                    const isSticky = cell.column.columnDef.meta?.isSticky;
                    return (
                      <td
                        key={cell.id}
                        className={`px-6 py-4 whitespace-nowrap text-sm text-gray-900 ${
                          isSticky ? "sticky-column" : ""
                        }`}
                        style={{
                          width: cell.column.getSize(),
                          left: isSticky
                            ? stickyColumnOffsets[cell.column.id]
                            : undefined,
                          zIndex: isSticky ? 10 : undefined,
                        }}
                      >
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext()
                        )}
                      </td>
                    );
                  })}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
      <div className="flex items-center justify-between mt-4 flex-wrap gap-2">
        <div className="flex items-center gap-2">
          <button
            onClick={() => table.setPageIndex(0)}
            disabled={!table.getCanPreviousPage()}
            className="px-3 py-1 border rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer transition-colors"
          >
            {"<<"}
          </button>
          <button
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
            className="px-3 py-1 border rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer transition-colors"
          >
            {"<"}
          </button>
          <button
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
            className="px-3 py-1 border rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer transition-colors"
          >
            {">"}
          </button>
          <button
            onClick={() => table.setPageIndex(table.getPageCount() - 1)}
            disabled={!table.getCanNextPage()}
            className="px-3 py-1 border rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer transition-colors"
          >
            {">>"}
          </button>
          <span className="flex items-center gap-1">
            <div>Sayfa</div>
            <strong>
              {table.getState().pagination.pageIndex + 1} /{" "}
              {table.getPageCount()}
            </strong>
          </span>
        </div>
        <div className="flex items-center gap-2">
          <input
            type="number"
            defaultValue={table.getState().pagination.pageIndex + 1}
            onChange={(e) =>
              table.setPageIndex(
                e.target.value ? Number(e.target.value) - 1 : 0
              )
            }
            className="border p-1 rounded w-16"
          />
          <select
            value={table.getState().pagination.pageSize}
            onChange={(e) => table.setPageSize(Number(e.target.value))}
            className="p-1 border rounded"
          >
            {[10, 20, 30, 40, 50, 100].map((pageSize) => (
              <option key={pageSize} value={pageSize}>
                Göster {pageSize}
              </option>
            ))}
          </select>
        </div>
      </div>
    </div>
  );
}

export default DataTable;
