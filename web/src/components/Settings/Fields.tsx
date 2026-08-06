import {
  List,
  ListItem,
  ListItemContent,
  ListItemLabel,
  ListItemSmall,
  ListItemMore,
  ListItemSwitch,
  ListTitle,
  ListItemText,
  ListItemButton,
  ListItemNavLink,
} from "@/components/App/List";
import { type ComponentProps, useRef, useEffect, useMemo } from "react";
import { useAppDispatch, useAppExtra } from "@/components/App/Context";
import { withExtraAction, ExtraType, ExtraState } from "./Extra";
import { TOCChapter } from "./TableOfContents";
import { useTranslation } from "@/components/App/I18Next";
import { SettingsName, MainI18nPrefix } from "./meta";
import {
  useInfoFields,
  useNetworkFields,
  useWebFields,
  useNetworkEnableMutation,
  useBroadcastEnableMutation,
  useWebEnableMutation,
  useLocalHostOnlyMutation,
} from "./ReactQuery";

const FieldsClassName = "first:rounded-t-sm last:rounded-b-sm";
const InfoKeyPrefix = `${MainI18nPrefix}.device-info`;
export const InfoFieldsTitle = ({
  className,
}: Pick<ComponentProps<typeof ListTitle>, "className">) => {
  const { t } = useTranslation(SettingsName, { keyPrefix: InfoKeyPrefix });

  return (
    <ListTitle id={TOCChapter.DeviceInfo} className={className}>
      {t("title")}
    </ListTitle>
  );
};

export const InfoFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { t } = useTranslation(SettingsName, { keyPrefix: InfoKeyPrefix });
  const { type, open } = useAppExtra<ExtraState<ExtraType>>({});

  const { data: infoFields } = useInfoFields({
    enabled: open !== true,
  });

  const [memoValue, memoMore] = useMemo(() => {
    const value = infoFields?.memo;
    if (!value || value.length < 1) {
      return [void 0, false];
    }
    const { index } = /\S/g.exec(value) || { index: -1 };
    if (index < 0) {
      return ["", false];
    }
    const lastIdx = value.indexOf("\n", index);
    const displayValue = value.substring(
      index,
      lastIdx < 0 || lastIdx === value.length ? void 0 : lastIdx,
    );

    if (lastIdx < 0 || lastIdx === value.length) {
      return [displayValue, false];
    }
    const { index: moreIdx } = /\S/g.exec(value.substring(lastIdx)) || {
      index: -1,
    };
    return [displayValue, moreIdx >= 0];
  }, [infoFields?.memo]);

  const nameButtonRef = useRef<HTMLButtonElement>(null);
  const memoButtonRef = useRef<HTMLButtonElement>(null);
  const dispatch = useAppDispatch();

  const handleNameEdit = () => {
    dispatch?.(withExtraAction(ExtraType.NameEdit));
  };

  const handleMemoEdit = () => {
    dispatch?.(withExtraAction(ExtraType.MemoEdit));
  };

  useEffect(() => {
    if (open) {
      return;
    }
    switch (type) {
      case ExtraType.NameEdit:
        nameButtonRef.current?.focus();
        break;
      case ExtraType.MemoEdit:
        memoButtonRef.current?.focus();
        break;
    }
  }, [type, open, nameButtonRef, memoButtonRef]);

  return (
    <List className={className}>
      <ListItem className={FieldsClassName}>
        <ListItemButton onClick={handleNameEdit} ref={nameButtonRef}>
          <ListItemContent>
            <ListItemText>{t("name")}</ListItemText>
            <ListItemSmall>{infoFields?.name || "-"}</ListItemSmall>
          </ListItemContent>
        </ListItemButton>
      </ListItem>
      <ListItem className={FieldsClassName}>
        <ListItemButton onClick={handleMemoEdit} ref={memoButtonRef}>
          <ListItemContent>
            <ListItemText>{t("memo")}</ListItemText>
            <ListItemSmall>
              {memoValue || "-"}
              {memoValue && memoMore ? "..." : ""}
            </ListItemSmall>
          </ListItemContent>
        </ListItemButton>
      </ListItem>
    </List>
  );
};

const NetworkKeyPrefix = `${MainI18nPrefix}.network`;
export const NetworkFieldsTitle = ({
  className,
}: Pick<ComponentProps<typeof ListTitle>, "className">) => {
  const { t } = useTranslation(SettingsName, { keyPrefix: NetworkKeyPrefix });
  return (
    <ListTitle id={TOCChapter.Network} className={className}>
      {t("title")}
    </ListTitle>
  );
};

export const NetworkFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { t } = useTranslation(SettingsName, { keyPrefix: NetworkKeyPrefix });
  const { type, open } = useAppExtra<ExtraState<ExtraType>>({});

  const { data: networkFields } = useNetworkFields({
    enabled: open !== true,
  });

  const { mutate: saveNetworkEnabled, isPending: isNetworkEnabledPending } =
    useNetworkEnableMutation();
  const { mutate: saveBroadcastEnabled, isPending: isBroadcastEnabledPending } =
    useBroadcastEnableMutation();

  const portButtonRef = useRef<HTMLButtonElement>(null);
  const publicAddrsButtonRef = useRef<HTMLButtonElement>(null);
  const broadcastAddrsButtonRef = useRef<HTMLButtonElement>(null);
  const dispatch = useAppDispatch();

  const handlePortEdit = () => {
    dispatch?.(withExtraAction(ExtraType.PeerPortEdit));
  };

  const handlePublicAddrsEdit = () => {
    dispatch?.(withExtraAction(ExtraType.PublicAddrsEdit));
  };

  const handleBroadcastAddrsEdit = () => {
    dispatch?.(withExtraAction(ExtraType.BroadcastAddrsEdit));
  };

  useEffect(() => {
    if (open) {
      return;
    }
    switch (type) {
      case ExtraType.PeerPortEdit:
        portButtonRef.current?.focus();
        break;
      case ExtraType.PublicAddrsEdit:
        publicAddrsButtonRef.current?.focus();
        break;
      case ExtraType.BroadcastAddrsEdit:
        broadcastAddrsButtonRef.current?.focus();
        break;
    }
  }, [
    type,
    open,
    portButtonRef,
    publicAddrsButtonRef,
    broadcastAddrsButtonRef,
  ]);

  const disabled = isNetworkEnabledPending || networkFields?.enabled !== true;
  const broadcastDisabled =
    disabled ||
    isBroadcastEnabledPending ||
    networkFields?.broadcastEnabled !== true;

  const handleEnabledChange = (checked: boolean) => {
    saveNetworkEnabled({ enabled: checked });
  };

  const handleBroadcastEnabledChange = (checked: boolean) => {
    saveBroadcastEnabled({ broadcastEnabled: checked });
  };

  const [broadcastAddrsText, broadcastAddrsMore] = useMemo(() => {
    const addrs = networkFields?.broadcastAddrs;
    if (!addrs || addrs.length < 1) {
      return [void 0, false];
    }

    let idx = -1;
    let more = false;
    for (let i = 0; i < addrs.length; i++) {
      const addr = addrs[i];
      if (/^\s+$/g.test(addr) || addr.trim().length < 1) {
        continue;
      }
      if (idx < 0) {
        idx = i;
        continue;
      }
      more = true;
      break;
    }

    return [addrs[idx], more];
  }, [networkFields?.broadcastAddrs]);

  const [publicAddrsText, publicAddrsMore] = useMemo(() => {
    const addrs = networkFields?.publicAddrs;
    if (!addrs || addrs.length < 1) {
      return [void 0, false];
    }

    let idx = -1;
    let more = false;
    for (let i = 0; i < addrs.length; i++) {
      const addr = addrs[i];
      if (/^\s+$/g.test(addr) || addr.trim().length < 1) {
        continue;
      }
      if (idx < 0) {
        idx = i;
        continue;
      }
      more = true;
      break;
    }

    return [addrs[idx], more];
  }, [networkFields?.publicAddrs]);
  return (
    <List className={className}>
      <ListItem className={FieldsClassName} disabled={isNetworkEnabledPending}>
        <ListItemContent>
          <ListItemLabel htmlFor="network-fields-enabled">
            {t("enabled")}
          </ListItemLabel>
          <ListItemSwitch
            id="network-fields-enabled"
            checked={networkFields?.enabled === true}
            onCheckedChange={handleEnabledChange}
            disabled={isNetworkEnabledPending}
          />
        </ListItemContent>
      </ListItem>
      <ListItem className={FieldsClassName} disabled={disabled}>
        <ListItemButton
          disabled={disabled}
          onClick={handlePortEdit}
          ref={portButtonRef}
        >
          <ListItemContent>
            <ListItemText>{t("port")}</ListItemText>
            <ListItemSmall>{networkFields?.peerPort || "-"}</ListItemSmall>
          </ListItemContent>
        </ListItemButton>
      </ListItem>
      <ListItem
        className={FieldsClassName}
        disabled={disabled || isBroadcastEnabledPending}
      >
        <ListItemContent>
          <ListItemLabel htmlFor="network-fields-broadcast-enabled">
            {t("broadcast")}
          </ListItemLabel>
          <ListItemSwitch
            id="network-fields-broadcast-enabled"
            checked={networkFields?.broadcastEnabled === true}
            disabled={disabled || isBroadcastEnabledPending}
            onCheckedChange={handleBroadcastEnabledChange}
          />
        </ListItemContent>
      </ListItem>
      <ListItem className={FieldsClassName} disabled={broadcastDisabled}>
        <ListItemButton
          disabled={broadcastDisabled}
          onClick={handleBroadcastAddrsEdit}
          ref={broadcastAddrsButtonRef}
        >
          <ListItemContent>
            <ListItemText>{t("broadcast-addrs")}</ListItemText>
            <ListItemSmall>
              {broadcastAddrsText || "-"}
              {broadcastAddrsText && broadcastAddrsMore ? "..." : ""}
            </ListItemSmall>
            <ListItemMore />
          </ListItemContent>
        </ListItemButton>
      </ListItem>
      <ListItem className={FieldsClassName} disabled={disabled}>
        <ListItemButton
          disabled={disabled}
          onClick={handlePublicAddrsEdit}
          ref={publicAddrsButtonRef}
        >
          <ListItemContent>
            <ListItemText>{t("public-addrs")}</ListItemText>
            <ListItemSmall>
              {publicAddrsText || "-"}
              {publicAddrsText && publicAddrsMore ? "..." : ""}
            </ListItemSmall>
            <ListItemMore />
          </ListItemContent>
        </ListItemButton>
      </ListItem>
    </List>
  );
};

const WebKeyPrefix = `${MainI18nPrefix}.web`;
export const WebFieldsTitle = ({
  className,
}: Pick<ComponentProps<typeof ListTitle>, "className">) => {
  const { t } = useTranslation(SettingsName, { keyPrefix: WebKeyPrefix });
  return (
    <ListTitle id={TOCChapter.Web} className={className}>
      {t("title")}
    </ListTitle>
  );
};

export const WebFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { t } = useTranslation(SettingsName, { keyPrefix: WebKeyPrefix });
  const { type, open } = useAppExtra<ExtraState<ExtraType>>({});

  const { data: webFields } = useWebFields({
    enabled: open !== true,
  });
  const { mutate: saveWebEnabled, isPending: isWebEnabledPending } =
    useWebEnableMutation();
  const { mutate: saveLocalHostOnly, isPending: isLocalHostOnlyPending } =
    useLocalHostOnlyMutation();

  const webPortButtonRef = useRef<HTMLButtonElement>(null);
  const dispatch = useAppDispatch();

  const handleWebPortEdit = () => {
    dispatch?.(withExtraAction(ExtraType.WebPortEdit));
  };

  useEffect(() => {
    if (open) {
      return;
    }
    switch (type) {
      case ExtraType.WebPortEdit:
        webPortButtonRef.current?.focus();
        break;
    }
  }, [type, open, webPortButtonRef]);

  const localhostOnlyDisabled =
    isWebEnabledPending || webFields?.webEnabled !== true;

  const handleEnabledChange = (checked: boolean) => {
    saveWebEnabled({ webEnabled: checked });
  };

  const handleLocalHostOnlyChange = (checked: boolean) => {
    saveLocalHostOnly({ localHostOnly: checked });
  };
  return (
    <List className={className}>
      <ListItem className={FieldsClassName} disabled={isWebEnabledPending}>
        <ListItemContent>
          <ListItemLabel htmlFor="web-fields-enabled">
            {t("enabled")}
          </ListItemLabel>
          <ListItemSwitch
            id="web-fields-enabled"
            checked={webFields?.webEnabled === true}
            onCheckedChange={handleEnabledChange}
            disabled={isWebEnabledPending}
          />
        </ListItemContent>
      </ListItem>
      <ListItem
        className={FieldsClassName}
        disabled={isLocalHostOnlyPending || localhostOnlyDisabled}
      >
        <ListItemContent>
          <ListItemLabel htmlFor="web-fields-allowed">
            {t("allowed")}
          </ListItemLabel>
          <ListItemSwitch
            id="web-fields-allowed"
            checked={webFields?.localHostOnly === true}
            disabled={isLocalHostOnlyPending || localhostOnlyDisabled}
            onCheckedChange={handleLocalHostOnlyChange}
          />
        </ListItemContent>
      </ListItem>
      <ListItem className={FieldsClassName} disabled={localhostOnlyDisabled}>
        <ListItemButton
          disabled={localhostOnlyDisabled}
          onClick={handleWebPortEdit}
          ref={webPortButtonRef}
        >
          <ListItemContent>
            <ListItemText>{t("port")}</ListItemText>
            <ListItemSmall>{webFields?.webPort || "-"}</ListItemSmall>
          </ListItemContent>
        </ListItemButton>
      </ListItem>
    </List>
  );
};

const AppearanceKeyPrefix = `${MainI18nPrefix}.appearance`;
export const AppearanceFieldsTitle = ({
  className,
}: Pick<ComponentProps<typeof ListTitle>, "className">) => {
  const { t } = useTranslation(SettingsName, {
    keyPrefix: AppearanceKeyPrefix,
  });
  return (
    <ListTitle id={TOCChapter.Appearance} className={className}>
      {t("title")}
    </ListTitle>
  );
};

export const AppearanceFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { t } = useTranslation(SettingsName, {
    keyPrefix: AppearanceKeyPrefix,
  });
  return (
    <List className={className}>
      <ListItem className={FieldsClassName}>
        <ListItemText>{t("theme")}</ListItemText>
        <ListItemSmall>{t("system")}</ListItemSmall>
      </ListItem>
    </List>
  );
};

const LanguageKeyPrefix = `${MainI18nPrefix}.languages`;
export const LanguageFieldsTitle = ({
  className,
}: Pick<ComponentProps<typeof ListTitle>, "className">) => {
  const { t } = useTranslation(SettingsName, {
    keyPrefix: LanguageKeyPrefix,
  });
  return (
    <ListTitle id={TOCChapter.Language} className={className}>
      {t("title")}
    </ListTitle>
  );
};

export const LanguageFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { t } = useTranslation(SettingsName, {
    keyPrefix: LanguageKeyPrefix,
  });
  return (
    <List className={className}>
      <ListItem className={FieldsClassName}>
        <ListItemText>{t("lang")}</ListItemText>
        <ListItemSmall>{t("system")}</ListItemSmall>
      </ListItem>
    </List>
  );
};
