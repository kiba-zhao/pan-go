import { ClusterName } from "./meta";
import { ClusterInfoCard, ClusterAction } from "./Cluster";
import { DeviceCard, DeviceAction } from "./Device";
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

export const ClusterMain = () => {
  const i18nDispatch = useAppI18nDispatch();
  useEffect(() => {
    i18nDispatch?.(withAppI18nAction({ namespace: ClusterName }));
    return () => i18nDispatch?.(withAppI18nResetAction());
  }, [i18nDispatch]);

  const { t } = useTranslation(ClusterName);
  const dispatch = useAppDispatch();
  useEffect(() => {
    dispatch?.(withAppHeaderAction({ title: t(`${I18nVariant.Main}.title`) }));
    return () => dispatch?.(withResetAction());
  }, [dispatch, t]);
  return (
    <div className="p-3 @container flex gap-4">
      <ClusterSection />
      <DeviceSection />
    </div>
  );
};

const ClusterSection = () => (
  <section className="w-xs @max-3xl:hidden">
    <div className="py-2.5 pl-4 pr-3 flex items-center gap-2">
      <h3 className="text-sm text-muted-foreground grow">设备组</h3>
      <ClusterAction />
    </div>
    <ClusterInfoCard />
  </section>
);

const DeviceSection = () => (
  <section className="grow">
    <div className="py-2.5 pl-4 pr-3 flex items-center gap-2">
      <h3 className="text-sm text-muted-foreground grow">相关设备</h3>
      <DeviceAction />
    </div>
    <DeviceCard />
  </section>
);
