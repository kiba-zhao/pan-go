import { useTranslation } from "./Context";

const CUSTOM_PREFIX = "custom";

type UseTranslationFunction = typeof useTranslation;
type UseTranslationNS = Parameters<UseTranslationFunction>[0];
type UseTranslationOptions = Parameters<UseTranslationFunction>[1];
export const useTranslate: UseTranslationFunction = (
  ns: UseTranslationNS,
  options: UseTranslationOptions
) => useTranslation(ns, options || { keyPrefix: CUSTOM_PREFIX });
