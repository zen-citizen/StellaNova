import enLocale from "./locales/en.json";

export type enIds = keyof typeof enLocale;

export const messages = {
  en: enLocale,
};
