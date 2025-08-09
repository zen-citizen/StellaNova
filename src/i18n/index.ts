import enLocale from "./locales/en.json";

type Messages = Record<string, Record<string, string>>;

export const messages: Messages = {
  en: enLocale,
};
