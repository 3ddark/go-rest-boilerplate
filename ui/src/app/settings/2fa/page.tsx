
"use client";

import React, { useState, useEffect } from 'react';
import { setup2FA, enable2FA, disable2FA } from '../../../api/userService';
import { useTranslation } from "react-i18next";

const TwoFactorAuthPage = () => {
  const { t } = useTranslation();
  const [token, setToken] = useState<string | null>(null);
  const [qrCode, setQrCode] = useState('');
  const [secret, setSecret] = useState('');
  const [code, setCode] = useState('');
  const [recoveryCodes, setRecoveryCodes] = useState<string[]>([]);
  const [isEnabled, setIsEnabled] = useState(false);

  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    if (storedToken) {
      setToken(storedToken);
    }
    // Fetch the current 2FA status for the user
    // and update the isEnabled state.
  }, []);

  const handleSetup = async () => {
    if (token) {
      try {
        const response = await setup2FA(token);
        setQrCode(response.data.qrCode);
        setSecret(response.data.secret);
      } catch (error) {
        console.error("Failed to setup 2FA:", error);
      }
    }
  };

  const handleEnable = async () => {
    if (token) {
      try {
        const response = await enable2FA(token, code);
        setRecoveryCodes(response.data);
        setIsEnabled(true);
        setQrCode('');
      } catch (error) {
        console.error("Failed to enable 2FA:", error);
      }
    }
  };

  const handleDisable = async () => {
    if (token) {
      try {
        await disable2FA(token);
        setIsEnabled(false);
        setRecoveryCodes([]);
      } catch (error) {
        console.error("Failed to disable 2FA:", error);
      }
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6 text-gray-800">{t("two_factor_authentication")}</h1>

      {!isEnabled ? (
        <div>
          <p className="mb-4">{t("2fa_disabled_message")}</p>
          <button
            onClick={handleSetup}
            className="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-md hover:bg-indigo-700 cursor-pointer transition-colors"
          >
            {t("enable_2fa")}
          </button>

          {qrCode && (
            <div className="mt-4">
              <p className="mb-2">{t("scan_qr_code")}</p>
              <img src={qrCode} alt="QR Code" />
              <p className="mt-2">{t("enter_secret_manually")}</p>
              <p className="font-mono bg-gray-100 p-2 rounded">{secret}</p>

              <div className="mt-4">
                <label htmlFor="code" className="block text-sm font-medium text-gray-700">{t("verification_code")}</label>
                <input
                  type="text"
                  id="code"
                  value={code}
                  onChange={(e) => setCode(e.target.value)}
                  className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                />
                <button
                  onClick={handleEnable}
                  className="mt-2 px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 cursor-pointer transition-colors"
                >
                  {t("verify_enable")}
                </button>
              </div>
            </div>
          )}
        </div>
      ) : (
        <div>
          <p className="mb-4">{t("2fa_enabled_message")}</p>
          <button
            onClick={handleDisable}
            className="px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700 cursor-pointer transition-colors"
          >
            {t("disable_2fa")}
          </button>

          {recoveryCodes.length > 0 && (
            <div className="mt-4">
              <h2 className="text-xl font-bold mb-2">{t("recovery_codes")}</h2>
              <p>{t("save_recovery_codes")}</p>
              <ul className="list-disc list-inside bg-gray-100 p-4 rounded">
                {recoveryCodes.map((c) => (
                  <li key={c} className="font-mono">{c}</li>
                ))}
              </ul>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default TwoFactorAuthPage;
