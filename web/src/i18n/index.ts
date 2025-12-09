import i18n from "i18next";
import { initReactI18next } from "react-i18next";

import hu from "./locales/hu.json";

i18n
    .use(initReactI18next)
    .init({
        resources: {
            hu: { translation: hu },
        },
        lng: "hu",
        fallbackLng: "hu",
        interpolation: {
            escapeValue: false,
        },
    });

export default i18n;
