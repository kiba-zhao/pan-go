import { SettingsName } from "./meta";
import { NetworkManageIcon, CircleCheckBig, CircleX } from "./Icon";
import { FormInputProps } from "./form";
import { TOCItem, TOCChapter, FieldsClassName } from "./MainBase";
import {
  useDeviceNetworkFields,
  useNetworkEnableMutation,
  useBroadcastEnableMutation,
} from "./ReactQuery";
import { withExtraState, ExtraType, type ExtraState } from "./ExtraBase";

import { newIPAddr } from "@/lib/utils";
import {
  type ChangeEvent,
  type ComponentProps,
  useMemo,
  useRef,
  useEffect,
} from "react";
import { useAppDispatch, useAppExtra } from "@/components/App/Context";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";
import {
  useController,
  useWatch,
  type FieldValues,
} from "@/components/App/Form";

import {
  List,
  ListItem,
  ListItemLabel,
  ListItemContent,
  ListItemSmall,
  ListItemText,
  ListItemButton,
  ListItemSwitch,
  ListItemMore,
} from "@/components/App/List";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
  InputGroupText,
  InputGroupButton,
  InputGroupTextarea,
} from "@/components/ui/input-group";
import { Badge } from "@/components/ui/badge";

export const DeviceNetworkTOCItem = () => {
  const { t } = useTranslation(SettingsName);
  return (
    <TOCItem
      chapter={TOCChapter.DeviceNetwork}
      text={t(`${I18nVariant.Main}.${TOCChapter.DeviceNetwork}.title`)}
      Avatar={NetworkManageIcon}
    />
  );
};

export const DeviceNetworkFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { t } = useTranslation(SettingsName);
  const { type, open } = useAppExtra<ExtraState>();

  const { data: networkFields } = useDeviceNetworkFields({
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
    dispatch?.(withExtraState({ type: ExtraType.PeerPortEdit }));
  };

  const handlePublicAddrsEdit = () => {
    dispatch?.(withExtraState({ type: ExtraType.PublicAddrsEdit }));
  };

  const handleBroadcastAddrsEdit = () => {
    dispatch?.(withExtraState({ type: ExtraType.BroadcastAddrsEdit }));
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
            {t(`${I18nVariant.Main}.${TOCChapter.DeviceNetwork}.enabled`)}
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
            <ListItemText>
              {t(`${I18nVariant.Main}.${TOCChapter.DeviceNetwork}.port`)}
            </ListItemText>
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
            {t(`${I18nVariant.Main}.${TOCChapter.DeviceNetwork}.broadcast`)}
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
            <ListItemText>
              {t(
                `${I18nVariant.Main}.${TOCChapter.DeviceNetwork}.broadcast-addrs`,
              )}
            </ListItemText>
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
            <ListItemText>
              {t(
                `${I18nVariant.Main}.${TOCChapter.DeviceNetwork}.public-addrs`,
              )}
            </ListItemText>
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

export const DeviceNetworkPortInput = <T extends FieldValues>({
  name,
  control,
  disabled,
  defaultValue,
  onChange,
  ...props
}: FormInputProps<T>) => {
  const { t } = useTranslation(SettingsName);

  const {
    field,
    fieldState: { invalid, isDirty },
    formState: { disabled: formDisabled },
  } = useController({
    ...props,
    name,
    control,
    rules: {
      required: true,
      min: 1,
      max: 65535,
      validate: validateDeviceNetworkPort,
      onChange: onChange,
    },
    disabled,
    defaultValue,
  });

  const handleReset = () => {
    field.onChange(defaultValue || "");
  };

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value.trim();
    if (typeof value === "string" && value.length < 1) {
      field.onChange(value);
      return;
    }
    if (validateDeviceNetworkPort(value)) field.onChange(Number(value));
  };

  return (
    <InputGroup>
      <InputGroupInput
        {...field}
        onChange={handleChange}
        type="number"
        placeholder={t(
          `${I18nVariant.Form}.deviceNetworkPortInput.placeholder`,
          { example: 9000 },
        )}
        autoComplete="off"
        aria-invalid={invalid}
      />
      <InputGroupAddon align="block-end">
        <InputGroupText className="text-xs text-muted-foreground">
          1 ~ 65535
        </InputGroupText>
        <InputGroupButton
          variant="default"
          size="sm"
          className="ml-auto max-md:hidden"
          onClick={handleReset}
          disabled={formDisabled || !isDirty}
        >
          {t(`${I18nVariant.Action}.reset`)}
        </InputGroupButton>
      </InputGroupAddon>
    </InputGroup>
  );
};

export const DeviceBroadcastAddrsTextInput = <T extends FieldValues>({
  name,
  control,
  disabled,
  defaultValue,
  onChange,
  ...props
}: FormInputProps<T>) => {
  const { t } = useTranslation(SettingsName);

  const {
    field,
    fieldState: { invalid, isDirty },
    formState: { disabled: formDisabled },
  } = useController({
    ...props,
    name,
    control,
    rules: {
      validate: validateBroadcastAddrText,
      onChange: onChange,
    },
    disabled,
    defaultValue,
  });

  const inputValue = useWatch({ name, control, disabled });
  const [validCount, invalidCount] = useMemo(
    () => countAddrsText(inputValue, parseBroadcastAddr),
    [inputValue],
  );

  const handleReset = () => {
    field.onChange(defaultValue || "");
  };

  return (
    <InputGroup>
      <InputGroupTextarea
        {...field}
        placeholder={t(
          `${I18nVariant.Form}.deviceBroadcastAddrsTextInput.placeholder`,
          {
            example: "224.0.0.1:9000",
          },
        )}
        className="min-h-29 max-h-29 scrollbar-thin"
        autoComplete="off"
        aria-invalid={invalid}
      />
      <InputGroupAddon align="block-end">
        <Badge variant="secondary">
          <CircleCheckBig data-icon="inline-start" /> {validCount}
        </Badge>
        <Badge
          variant="destructive"
          className={invalid && invalidCount > 0 ? void 0 : "hidden"}
        >
          <CircleX data-icon="inline-start" /> {invalidCount}
        </Badge>
        <InputGroupButton
          variant="default"
          size="sm"
          className="ml-auto max-md:hidden"
          onClick={handleReset}
          disabled={formDisabled || !isDirty}
        >
          {t(`${I18nVariant.Action}.reset`)}
        </InputGroupButton>
      </InputGroupAddon>
    </InputGroup>
  );
};

export const DevicePublicAddrsTextInput = <T extends FieldValues>({
  name,
  control,
  disabled,
  defaultValue,
  onChange,
  ...props
}: FormInputProps<T>) => {
  const { t } = useTranslation(SettingsName);

  const {
    field,
    fieldState: { invalid, isDirty },
    formState: { disabled: formDisabled },
  } = useController({
    ...props,
    name,
    control,
    rules: {
      validate: validatePublicAddrText,
      onChange: onChange,
    },
    disabled,
    defaultValue,
  });

  const inputValue = useWatch({ name, control, disabled });
  const [validCount, invalidCount] = useMemo(
    () => countAddrsText(inputValue, parsePublicAddr),
    [inputValue],
  );

  const handleReset = () => {
    field.onChange(defaultValue || "");
  };

  return (
    <InputGroup>
      <InputGroupTextarea
        {...field}
        placeholder={t(
          `${I18nVariant.Form}.devicePublicAddrsTextInput.placeholder`,
          {
            example: "www.example.com:9000",
          },
        )}
        className="min-h-29 max-h-29 scrollbar-thin"
        autoComplete="off"
        aria-invalid={invalid}
      />
      <InputGroupAddon align="block-end">
        <Badge variant="secondary">
          <CircleCheckBig data-icon="inline-start" /> {validCount}
        </Badge>
        <Badge
          variant="destructive"
          className={invalid && invalidCount > 0 ? void 0 : "hidden"}
        >
          <CircleX data-icon="inline-start" /> {invalidCount}
        </Badge>
        <InputGroupButton
          variant="default"
          size="sm"
          className="ml-auto max-md:hidden"
          onClick={handleReset}
          disabled={formDisabled || !isDirty}
        >
          {t(`${I18nVariant.Action}.reset`)}
        </InputGroupButton>
      </InputGroupAddon>
    </InputGroup>
  );
};

function validateDeviceNetworkPort(value: number | string) {
  const value_ = Number(value);
  if (isNaN(value_) || !Number.isInteger(value_)) {
    return false;
  }
  if (value_.toString() !== value.toString()) {
    return false;
  }
  return value_ >= 1 && value_ <= 65535;
}

function validateBroadcastAddrText(value: string) {
  const [_, invalidCount] = countAddrsText(value, parseBroadcastAddr);
  return invalidCount === 0;
}
export function parseBroadcastAddr(addr: string): undefined | string {
  if (/^\s+$/g.test(addr) || addr.trim().length === 0) {
    return;
  }
  addr = addr.replace(/\s+/g, "");
  const splitIdx = addr.lastIndexOf(":");
  if (splitIdx === -1) {
    return;
  }

  const portPart = addr.substring(splitIdx + 1);
  const port = Number(portPart);
  if (
    isNaN(port) ||
    !Number.isInteger(port) ||
    port < 1 ||
    port > 65535 ||
    port.toString() !== portPart
  ) {
    return;
  }

  const ipPart = newIPAddr(addr.substring(0, splitIdx));
  if (!ipPart || !ipPart.isMulticast()) {
    return;
  }

  return `${ipPart.correctForm()}:${portPart}`;
}

function countAddrsText(
  addrsText: string | undefined,
  parseFn: typeof parseBroadcastAddr,
) {
  const lines = addrsText?.split("\n") || [];
  if (lines.length === 0) {
    return [0, 0];
  }

  let validCount = 0;
  let invalidCount = 0;
  for (const line of lines) {
    if (/^\s+$/g.test(line) || line.trim().length === 0) {
      continue;
    }
    if (parseFn(line) !== void 0) {
      validCount++;
    } else {
      invalidCount++;
    }
  }

  return [validCount, invalidCount];
}

function validatePublicAddrText(value: string) {
  const [_, invalidCount] = countAddrsText(value, parsePublicAddr);
  return invalidCount === 0;
}

export function parsePublicAddr(addr: string): undefined | string {
  if (/^\s+$/g.test(addr) || addr.trim().length === 0) {
    return;
  }
  addr = addr.replace(/\s+/g, "");
  const splitIdx = addr.lastIndexOf(":");
  if (splitIdx === -1) {
    return;
  }

  const portPart = addr.substring(splitIdx + 1);
  const port = Number(portPart);
  if (
    isNaN(port) ||
    !Number.isInteger(port) ||
    port < 1 ||
    port > 65535 ||
    port.toString() !== portPart
  ) {
    return;
  }

  const ipPart = addr.substring(0, splitIdx);
  const ipAddr = newIPAddr(ipPart);
  if (ipAddr) {
    if (
      ipAddr.isBroadcast() ||
      ipAddr.isMulticast() ||
      ipAddr.isLinkLocal() ||
      ipAddr.isLoopback() ||
      ipAddr.isCGNAT() ||
      ipAddr.isUnspecified()
    )
      return;
    return `${ipAddr.correctForm()}:${portPart}`;
  }

  const url = URL.parse(`https://${addr}`);
  if (!url || url.host !== addr) {
    return;
  }
  return url.host;
}
