// src/utils/tableStorage.ts
import {
  ColumnFiltersState,
  VisibilityState,
  SortingState,
  ColumnSizingState,
  PaginationState,
} from "@tanstack/react-table";

export interface TableUserPreferences {
  columnVisibility?: VisibilityState;
  columnFilters?: ColumnFiltersState;
  sorting?: SortingState;
  columnSizing?: ColumnSizingState;
  pagination?: PaginationState;
  globalFilter?: string;
}

const STORAGE_KEY_PREFIX = "erp_table_preferences_";

export function saveTablePreferences(
  tableId: string,
  preferences: TableUserPreferences
) {
  try {
    localStorage.setItem(
      STORAGE_KEY_PREFIX + tableId,
      JSON.stringify(preferences)
    );
    console.log(`Table preferences saved for ${tableId}`);
  } catch (error) {
    console.error(`Error saving table preferences for ${tableId}:`, error);
  }
}

export function loadTablePreferences(
  tableId: string
): TableUserPreferences | null {
  try {
    const item = localStorage.getItem(STORAGE_KEY_PREFIX + tableId);
    return item ? JSON.parse(item) : null;
  } catch (error) {
    console.error(`Error loading table preferences for ${tableId}:`, error);
    return null;
  }
}

export function clearTablePreferences(tableId: string) {
  try {
    localStorage.removeItem(STORAGE_KEY_PREFIX + tableId);
    console.log(`Table preferences cleared for ${tableId}`);
  } catch (error) {
    console.error(`Error clearing table preferences for ${tableId}:`, error);
  }
}
