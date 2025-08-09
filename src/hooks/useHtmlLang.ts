import { useEffect, useState } from "react";

export const useHtmlLang = (): string => {
  const [lang, setLang] = useState<string>("en");
  useEffect(() => {
    const initialLang = document.documentElement.lang || "en";
    setLang(initialLang);

    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (
          mutation.type === "attributes" &&
          mutation.attributeName === "lang"
        ) {
          const newLang = document.documentElement.lang || "en";
          setLang(newLang);
        }
      });
    });

    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ["lang"],
    });

    return () => {
      observer.disconnect();
    };
  }, []);

  return lang;
};
