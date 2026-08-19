import { SettingsName } from "./meta";
import { default as SettingsIcon } from "./Icon";
import { useEffect } from "react";

import {
  useAppDispatch,
  withResetAction,
  withAppHeaderAction,
} from "@/components/App/Context";

import {
  useTranslation,
  useAppI18nDispatch,
  withAppI18nAction,
  withAppI18nResetAction,
  I18nVariant,
} from "@/components/App/I18Next";
import { cn } from "@/lib/utils";

import {
  TOCChapter,
  TOCItemGroup,
  FieldsSection,
  FieldsSearchFilter,
  FieldsSectionClassName,
} from "./MainBase";
import { DeviceInfoTOCItem, DeviceInfoFields } from "./DeviceInfo";
import { DeviceNetworkTOCItem, DeviceNetworkFields } from "./DeviceNetwork";
import { DeviceWebTOCItem, DeviceWebFields } from "./DeviceWeb";
import { AppearanceTOCItem, AppearanceFields } from "./Appearance";
import { LanguageTOCItem, LanguageFields } from "./Language";
import { DeviceClusterTOCItem, DeviceClusterFields } from "./DeviceCluster";
import { Separator } from "@/components/ui/separator";

const SettingsMain = () => {
  const i18nDispatch = useAppI18nDispatch();
  useEffect(() => {
    i18nDispatch?.(withAppI18nAction({ namespace: SettingsName }));
    return () => i18nDispatch?.(withAppI18nResetAction());
  }, [i18nDispatch]);

  const { t } = useTranslation(SettingsName);
  const dispatch = useAppDispatch();
  useEffect(() => {
    dispatch?.(
      withAppHeaderAction({
        title: t(`${I18nVariant.Main}.title`),
      }),
    );
    return () => dispatch?.(withResetAction());
  }, [dispatch, t]);

  return (
    <div className="p-3 size-full @container">
      <TOCSection />
      <div className="@3xl:pl-68 flex flex-col items-center gap-1">
        <FilterSection />
        <DeviceInfoSection />
        <DeviceNetworkSection />
        <DeviceWebSection />
        <DeviceClusterSection />
        <AppearanceSection />
        <LanguageSection />
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

const FilterSection = () => {
  return (
    <section className={cn("pb-2", FieldsSectionClassName)}>
      <FieldsSearchFilter />
    </section>
  );
};

const DeviceInfoSection = () => (
  <FieldsSection chapter={TOCChapter.DeviceInfo}>
    <DeviceInfoFields className={PanelClassName} />
  </FieldsSection>
);

const DeviceNetworkSection = () => (
  <FieldsSection chapter={TOCChapter.DeviceNetwork}>
    <DeviceNetworkFields className={PanelClassName} />
  </FieldsSection>
);

const DeviceWebSection = () => (
  <FieldsSection chapter={TOCChapter.DeviceWeb}>
    <DeviceWebFields className={PanelClassName} />
  </FieldsSection>
);

const AppearanceSection = () => (
  <FieldsSection chapter={TOCChapter.Appearance}>
    <AppearanceFields className={PanelClassName} />
  </FieldsSection>
);

const LanguageSection = () => (
  <FieldsSection chapter={TOCChapter.Language}>
    <LanguageFields className={PanelClassName} />
  </FieldsSection>
);

const DeviceClusterSection = () => (
  <FieldsSection chapter={TOCChapter.DeviceCluster}>
    <DeviceClusterFields className={PanelClassName} />
  </FieldsSection>
);

const TOCHeader = () => {
  const { t } = useTranslation(SettingsName);
  return (
    <div className="px-6 flex flex-row gap-3 text-xl items-center align-bottom pb-3">
      <SettingsIcon size={28} className="h-7" />
      <h2 className="font-medium">{t(`${I18nVariant.Main}.subtitle`)}</h2>
    </div>
  );
};

const TableOfContents = () => (
  <nav className="w-full flex flex-col gap-1">
    <TOCItemGroup>
      <DeviceInfoTOCItem />
      <DeviceNetworkTOCItem />
      <DeviceWebTOCItem />
    </TOCItemGroup>
    <Separator />
    <TOCItemGroup>
      <DeviceClusterTOCItem />
    </TOCItemGroup>
    <Separator />
    <TOCItemGroup>
      <AppearanceTOCItem />
      <LanguageTOCItem />
    </TOCItemGroup>
  </nav>
);
