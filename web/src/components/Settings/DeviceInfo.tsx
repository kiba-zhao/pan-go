import { SettingsName } from "./meta";
import { MonitorSmartphone } from "./Icon";
import { TOCItem, TOCChapter, FieldsClassName } from "./MainBase";
import { useDeviceInfo } from "./ReactQuery";
import { withExtraState, ExtraType, type ExtraState } from "./ExtraBase";
import { FormInputProps } from "./form";

import {
  type ComponentProps,
  type ChangeEvent,
  useMemo,
  useRef,
  useEffect,
} from "react";

import {
  useWatch,
  useController,
  type FieldValues,
} from "@/components/App/Form";
import { useAppDispatch, useAppExtra } from "@/components/App/Context";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";
import {
  List,
  ListItem,
  ListItemContent,
  ListItemSmall,
  ListItemText,
  ListItemButton,
} from "@/components/App/List";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
  InputGroupText,
  InputGroupButton,
  InputGroupTextarea,
} from "@/components/ui/input-group";

export const DeviceInfoTOCItem = () => {
  const { t } = useTranslation(SettingsName);
  return (
    <TOCItem
      chapter={TOCChapter.DeviceInfo}
      text={t(`${I18nVariant.Main}.${TOCChapter.DeviceInfo}.title`)}
      Avatar={MonitorSmartphone}
    />
  );
};

export const DeviceInfoFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { t } = useTranslation(SettingsName);
  const { type, open } = useAppExtra<ExtraState>();

  const { data: infoFields } = useDeviceInfo({
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
    dispatch?.(withExtraState({ type: ExtraType.DeviceNameEdit }));
  };

  const handleMemoEdit = () => {
    dispatch?.(withExtraState({ type: ExtraType.DeviceMemoEdit }));
  };

  useEffect(() => {
    if (open) {
      return;
    }
    switch (type) {
      case ExtraType.DeviceNameEdit:
        nameButtonRef.current?.focus();
        break;
      case ExtraType.DeviceMemoEdit:
        memoButtonRef.current?.focus();
        break;
    }
  }, [type, open, nameButtonRef, memoButtonRef]);

  return (
    <List className={className}>
      <ListItem className={FieldsClassName}>
        <ListItemButton onClick={handleNameEdit} ref={nameButtonRef}>
          <ListItemContent>
            <ListItemText>
              {t(`${I18nVariant.Main}.${TOCChapter.DeviceInfo}.name`)}
            </ListItemText>
            <ListItemSmall>{infoFields?.name || "-"}</ListItemSmall>
          </ListItemContent>
        </ListItemButton>
      </ListItem>
      <ListItem className={FieldsClassName}>
        <ListItemButton onClick={handleMemoEdit} ref={memoButtonRef}>
          <ListItemContent>
            <ListItemText>
              {t(`${I18nVariant.Main}.${TOCChapter.DeviceInfo}.memo`)}
            </ListItemText>
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

export const DeviceNameInput = <T extends FieldValues>({
  name,
  control,
  disabled,
  onChange,
  defaultValue,
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
      maxLength: 32,
      pattern: /^\S+$/g,
      onChange: onChange,
    },
    defaultValue,
  });

  const nameValue = useWatch({ name, control, disabled });

  const handleReset = () => {
    field.onChange(defaultValue || "");
  };

  return (
    <InputGroup>
      <InputGroupInput
        {...field}
        type="text"
        placeholder={t(`${I18nVariant.Form}.deviceNameInput.placeholder`)}
        autoComplete="off"
        aria-invalid={invalid}
      />
      <InputGroupAddon align="block-end">
        <InputGroupText className="text-xs text-muted-foreground">
          {nameValue?.length || 0} / 32
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

export const DeviceMemoInput = <T extends FieldValues>({
  name,
  control,
  disabled,
  onChange,
  defaultValue,
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
      maxLength: 256,
      onChange: onChange,
    },
    defaultValue,
  });

  const memoValue = useWatch({ name, control, disabled });

  const handleReset = () => {
    field.onChange(defaultValue || "");
  };

  const handleChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    field.onChange(setMemoAs(event.target.value));
  };

  return (
    <InputGroup>
      <InputGroupTextarea
        {...field}
        onChange={handleChange}
        placeholder={t(`${I18nVariant.Form}.deviceMemoInput.placeholder`)}
        className="min-h-29 max-h-29 scrollbar-thin"
        autoComplete="off"
        aria-invalid={invalid}
      />
      <InputGroupAddon align="block-end">
        <InputGroupText className="text-xs text-muted-foreground">
          {memoValue?.length || 0} / 256
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

function setMemoAs(value?: string) {
  if (!value) {
    return "";
  }
  const lines = value
    .trim()
    .split("\n")
    .filter((_) => !/^\s+$/g.test(_))
    .map((_) => _.trim());
  return lines.join("\n");
}
