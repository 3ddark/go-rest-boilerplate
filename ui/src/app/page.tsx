"use client";

import React, { useState, useCallback } from "react";
import DataTable from "../components/DataTable/DataTable";
import { CustomColumnDef } from "../components/DataTable/DataTable";
import { getUsers, User, login2FA } from "../api/userService";
import { login } from "../api/authService";
import { useTranslation } from "react-i18next";

export default function HomePage() {
  const { t } = useTranslation();
  const [token, setToken] = useState<string | null>(null);
  const [userId, setUserId] = useState<number | null>(null);
  const [twoFactorRequired, setTwoFactorRequired] = useState(false);
  const [code, setCode] = useState("");

  const columns: CustomColumnDef<User>[] = [
    {
      accessorKey: "id",
      header: t("id"),
      size: 80,
      meta: { isSticky: true },
    },
    {
      accessorKey: "name",
      header: t("name_surname"),
      size: 180,
      meta: { isSticky: true },
    },
    {
      accessorKey: "email",
      header: t("email"),
      size: 200,
    },
    {
      accessorKey: "age",
      header: t("age"),
      size: 70,
    },
    {
      accessorKey: "city",
      header: t("city"),
      size: 150,
    },
    {
      accessorKey: "status",
      header: t("status"),
      size: 100,
      cell: ({ row }) => (
        <span
          className={`px-2 py-1 rounded-full text-xs font-semibold ${
            row.original.status === "Active"
              ? "bg-green-100 text-green-800"
              : row.original.status === "Passive"
              ? "bg-red-100 text-red-800"
              : "bg-yellow-100 text-yellow-800"
          }`}
        >
          {row.original.status}
        </span>
      ),
    },
    {
      id: "actions",
      header: t("actions"),
      size: 150,
      meta: { noExport: true },
      cell: ({ row }) => (
        <div className="flex gap-2">
          <button
            onClick={() => alert(t("edit") + ": " + row.original.name)}
            className="text-blue-600 hover:underline text-sm"
          >
            {t("edit")}
          </button>
          <button
            onClick={() => alert(t("delete") + ": " + row.original.name)}
            className="text-red-600 hover:underline text-sm"
          >
            {t("delete")}
          </button>
        </div>
      ),
    },
    { accessorKey: "col1", header: t("additional_info") + " 1", size: 120 },
    { accessorKey: "col2", header: t("additional_info") + " 2", size: 120 },
    { accessorKey: "col3", header: t("additional_info") + " 3", size: 120 },
  ];

  const handleLogin = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const email = formData.get("email") as string;
    const password = formData.get("password") as string;

    try {
      const response = await login(email, password);
      if (response.data.twoFactorRequired) {
        setUserId(response.data.userId);
        setTwoFactorRequired(true);
      } else {
        setToken(response.data.token);
      }
    } catch (error) {
      console.error("Login failed:", error);
    }
  };

  const handle2FALogin = async () => {
    if (userId) {
      try {
        const response = await login2FA(userId, code);
        setToken(response.data.token);
      } catch (error) {
        console.error("2FA login failed:", error);
      }
    }
  };

  const fetchFn = useCallback(
    (params: any) => {
      if (token) {
        return getUsers(params, token);
      }
      return Promise.resolve({ data: [], pageCount: 0, totalRowCount: 0 });
    },
    [token]
  );

  if (!token) {
    return (
      <main className="container mx-auto p-4 flex justify-center items-center h-screen">
        <div className="w-full max-w-xs">
          {!twoFactorRequired ? (
            <form
              onSubmit={handleLogin}
              className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
            >
              <div className="mb-4">
                <label
                  className="block text-gray-700 text-sm font-bold mb-2"
                  htmlFor="email"
                >
                  {t("email")}
                </label>
                <input
                  className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                  id="email"
                  type="email"
                  name="email"
                  placeholder={t("email")}
                />
              </div>
              <div className="mb-6">
                <label
                  className="block text-gray-700 text-sm font-bold mb-2"
                  htmlFor="password"
                >
                  {t("password")}
                </label>
                <input
                  className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
                  id="password"
                  type="password"
                  name="password"
                  placeholder="******************"
                />
              </div>
              <div className="flex items-center justify-between">
                <button
                  className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                  type="submit"
                >
                  {t("sign_in")}
                </button>
              </div>
            </form>
          ) : (
            <div className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
              <div className="mb-4">
                <label
                  className="block text-gray-700 text-sm font-bold mb-2"
                  htmlFor="code"
                >
                  {t("verification_code")}
                </label>
                <input
                  className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                  id="code"
                  type="text"
                  value={code}
                  onChange={(e) => setCode(e.target.value)}
                  placeholder={t("verification_code")}
                />
              </div>
              <div className="flex items-center justify-between">
                <button
                  onClick={handle2FALogin}
                  className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                  type="button"
                >
                  {t("verify")}
                </button>
              </div>
            </div>
          )}
        </div>
      </main>
    );
  }

  return (
    <main className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6 text-gray-800">
        {t("erp_data_list")}
      </h1>
      <DataTable<User>
        tableId="userListTable"
        columns={columns}
        fetchFn={fetchFn}
      />
    </main>
  );
}
