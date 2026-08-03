import { useAppDispatch, resetToBlank } from "@/components/App/Context";
import { useTranslation } from "@/components/App/I18Next";
import { SettingsName } from "./meta";
import { useEffect } from "react";
import { cn } from "@/lib/utils";
import { default as SettingsIcon } from "./Icon";

import TableOfContents from "./TableOfContents";
import { SearchFilter } from "./Filters";
import {
  InfoFields,
  NetworkFields,
  InfoFieldsTitle,
  NetworkFieldsTitle,
  WebFieldsTitle,
  WebFields,
  AppearanceFieldsTitle,
  AppearanceFields,
  LanguageFieldsTitle,
  LanguageFields,
} from "./Fields";

const SettingsMain = () => {
  const dispatch = useAppDispatch();
  const { t } = useTranslation(SettingsName);

  useEffect(() => {
    dispatch?.({ header: { title: t(`main.title`) } });
    return () => dispatch?.(resetToBlank());
  }, [dispatch, t]);

  return (
    <div className="p-3 size-full @container">
      <TOCSection />
      <div className="@3xl:pl-68 flex flex-col items-center gap-1">
        <FilterSection />
        <InfoFieldsSection />
        <NetworkFieldsSection />
        <WebFieldsSection />
        <AppearanceFieldsSection />
        <LanguageFieldsSection />
      </div>
    </div>
  );
};

export default SettingsMain;

const PanelClassName =
  "rounded-sm bg-card text-card-foreground not-dark:border-border not-dark:border-1";
const TOCSection = () => {
  return (
    <section className={cn("fixed pt-3 w-66 @max-3xl:hidden", PanelClassName)}>
      <TOCHeader />
      <TableOfContents />
    </section>
  );
};

const ContentClassName = "max-w-3xl w-full pt-2";
const FilterSection = () => {
  return (
    <section className={cn("pb-2", ContentClassName)}>
      <SearchFilter />
    </section>
  );
};

const InfoFieldsSection = () => {
  return (
    <section className={ContentClassName}>
      <InfoFieldsTitle />
      <InfoFields className={PanelClassName} />
    </section>
  );
};

const NetworkFieldsSection = () => {
  return (
    <section className={ContentClassName}>
      <NetworkFieldsTitle />
      <NetworkFields className={PanelClassName} />
    </section>
  );
};

const WebFieldsSection = () => {
  return (
    <section className={ContentClassName}>
      <WebFieldsTitle />
      <WebFields className={PanelClassName} />
    </section>
  );
};

const AppearanceFieldsSection = () => {
  return (
    <section className={ContentClassName}>
      <AppearanceFieldsTitle />
      <AppearanceFields className={PanelClassName} />
    </section>
  );
};

const LanguageFieldsSection = () => {
  return (
    <section className={ContentClassName}>
      <LanguageFieldsTitle />
      <LanguageFields className={PanelClassName} />
    </section>
  );
};

const TOCHeader = () => {
  const { t } = useTranslation(SettingsName);
  return (
    <div className="px-6 flex flex-row gap-3 text-xl items-center align-bottom pb-3">
      <SettingsIcon size={28} className="h-7" />
      <h2 className="font-medium">{`${t(`main.subtitle`)}`}</h2>
    </div>
  );
};
