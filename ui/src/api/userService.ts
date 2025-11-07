import {
  ColumnFiltersState,
  PaginationState,
  SortingState,
} from "@tanstack/react-table";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export interface User {
  id: number;
  name: string;
  age: number;
  city: string;
  status: "Active" | "Passive" | "Pending";
  email: string;
  col1?: string;
  col2?: string;
  col3?: string;
}

export interface GetUsersParams {
  pagination: PaginationState;
  sorting: SortingState;
  columnFilters: ColumnFiltersState;
  globalFilter: string | null;
}

export interface GetUsersResponse {
  data: User[];
  pageCount: number;
  totalRowCount: number;
}

export async function getUsers(
  params: GetUsersParams,
  token: string
): Promise<GetUsersResponse> {
  const queryParams = new URLSearchParams({
    page: (params.pagination.pageIndex + 1).toString(),
    limit: params.pagination.pageSize.toString(),
    sort: params.sorting.map((s) => `${s.id}:${s.desc ? "desc" : "asc"}`).join(","),
    filter: params.columnFilters.map((f) => `${f.id}:${f.value}`).join(","),
    globalFilter: params.globalFilter || "",
  });

  const response = await fetch(`${API_URL}/users?${queryParams}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  if (!response.ok) {
    throw new Error("Failed to fetch users");
  }

  const data = await response.json();
  return {
    data: data.data,
    pageCount: data.meta.last_page,
    totalRowCount: data.meta.total,
  };
}

export async function setup2FA(token: string) {
  const response = await fetch(`${API_URL}/users/2fa/setup`,
    {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  if (!response.ok) {
    throw new Error("Failed to setup 2FA");
  }

  return response.json();
}

export async function enable2FA(token: string, code: string) {
  const response = await fetch(`${API_URL}/users/2fa/enable`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ code }),
    }
  );

  if (!response.ok) {
    throw new Error("Failed to enable 2FA");
  }

  return response.json();
}

export async function disable2FA(token: string) {
  const response = await fetch(`${API_URL}/users/2fa/disable`,
    {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  if (!response.ok) {
    throw new Error("Failed to disable 2FA");
  }

  return response.json();
}

export async function login2FA(userId: number, code: string) {
  const response = await fetch(`${API_URL}/login/2fa`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ userId, code }),
    }
  );

  if (!response.ok) {
    throw new Error("Failed to login with 2FA");
  }

  return response.json();
}
