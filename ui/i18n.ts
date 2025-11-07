import i18n from "i18next";
import { initReactI18next } from "react-i18next";

i18n
  .use(initReactI18next)
  .init({
    fallbackLng: "en",
    lng: "en",
    debug: true,
    interpolation: {
      escapeValue: false,
    },
    resources: {
      en: {
        translation: require("../public/locales/en/translation.json"),
      },
      tr: {
        translation: require("../public/locales/tr/translation.json"),
      },
    },
  });

export default i18n;
