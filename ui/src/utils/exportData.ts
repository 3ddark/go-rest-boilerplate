// src/utils/exportData.ts
import * as XLSX from "xlsx";
import { saveAs } from "file-saver";
import { Row, Column } from "@tanstack/react-table";

export function exportToCSV<TData extends object>(
  rows: Row<TData>[],
  columns: Column<TData, unknown>[],
  filename: string = "data.csv"
) {
  if (!rows.length) {
    console.warn("Dışa aktarılacak veri bulunamadı.");
    return;
  }

  const exportableColumns = columns.filter(
    (col) => col.getIsVisible() && !col.columnDef.meta?.["noExport"]
  );
  const headers = exportableColumns.map((col) =>
    typeof col.columnDef.header === "string" ? col.columnDef.header : col.id
  );

  const csvRows = rows.map((row) => {
    return exportableColumns
      .map((col) => {
        const value = col.accessorFn
          ? col.accessorFn(row.original, row.index)
          : (row.original as any)[col.id];
        return `"${String(value).replace(/"/g, '""')}"`;
      })
      .join(",");
  });

  const csvContent = [headers.join(","), ...csvRows].join("\n");
  const blob = new Blob([`\uFEFF${csvContent}`], {
    type: "text/csv;charset=utf-8;",
  }); // BOM for Excel
  saveAs(blob, filename);
}

export function exportToExcel<TData extends object>(
  rows: Row<TData>[],
  columns: Column<TData, unknown>[],
  filename: string = "data.xlsx"
) {
  if (!rows.length) {
    console.warn("Dışa aktarılacak veri bulunamadı.");
    return;
  }

  const exportableColumns = columns.filter(
    (col) => col.getIsVisible() && !col.columnDef.meta?.["noExport"]
  );
  const headerNames = exportableColumns.map((col) =>
    typeof col.columnDef.header === "string" ? col.columnDef.header : col.id
  );

  const dataForExport = rows.map((row) => {
    let rowData: { [key: string]: any } = {};
    exportableColumns.forEach((col) => {
      const headerName =
        typeof col.columnDef.header === "string"
          ? col.columnDef.header
          : col.id;
      rowData[headerName] = col.accessorFn
        ? col.accessorFn(row.original, row.index)
        : (row.original as any)[col.id];
    });
    return rowData;
  });

  const worksheet = XLSX.utils.json_to_sheet(dataForExport);
  const workbook = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(workbook, worksheet, "Sayfa1");
  const excelBuffer = XLSX.write(workbook, { bookType: "xlsx", type: "array" });
  const blob = new Blob([excelBuffer], { type: "application/octet-stream" });
  saveAs(blob, filename);
}
