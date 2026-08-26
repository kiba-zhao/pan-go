import { SettingsName } from "./meta";
import { WebIcon } from "./Icon";
import { FormInputProps } from "./form";
import { TOCItem, TOCChapter, FieldsClassName } from "./MainBase";
import {
  useWebHost,
  useWebEnableMutation,
  useLocalHostOnlyMutation,
} from "./ReactQuery";
import { withExtraState, ExtraType, type ExtraState } from "./ExtraBase";

import {
  type ComponentProps,
  type ChangeEvent,
  useRef,
  useEffect,
} from "react";
import { useController, type FieldValues } from "@/components/App/Form";
import { useAppDispatch, useAppExtra } from "@/components/App/Context";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";

import {
  List,
  ListItem,
  ListItemLabel,
  ListItemContent,
  ListItemSmall,
  ListItemText,
  ListItemButton,
  ListItemSwitch,
} from "@/components/App/List";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
  InputGroupText,
  InputGroupButton,
} from "@/components/ui/input-group";

export const DeviceWebTOCItem = () => {
  const { t } = useTranslation(SettingsName);
  return (
    <TOCItem
      chapter={TOCChapter.DeviceWeb}
      text={t(`${I18nVariant.Main}.${TOCChapter.DeviceWeb}.title`)}
      Avatar={WebIcon}
    />
  );
};

export const DeviceWebFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { t } = useTranslation(SettingsName);
  const { type, open } = useAppExtra<ExtraState>();

  const { data: webFields } = useWebHost({
    enabled: open !== true,
  });
  const { mutate: saveWebEnabled, isPending: isWebEnabledPending } =
    useWebEnableMutation();
  const { mutate: saveLocalHostOnly, isPending: isLocalHostOnlyPending } =
    useLocalHostOnlyMutation();

  const webPortButtonRef = useRef<HTMLButtonElement>(null);
  const dispatch = useAppDispatch();

  const handleWebPortEdit = () => {
    dispatch?.(withExtraState({ type: ExtraType.WebPortEdit }));
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
            {t(`${I18nVariant.Main}.${TOCChapter.DeviceWeb}.enabled`)}
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
            {t(`${I18nVariant.Main}.${TOCChapter.DeviceWeb}.allowed`)}
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
            <ListItemText>
              {t(`${I18nVariant.Main}.${TOCChapter.DeviceWeb}.port`)}
            </ListItemText>
            <ListItemSmall>{webFields?.webPort || "-"}</ListItemSmall>
          </ListItemContent>
        </ListItemButton>
      </ListItem>
    </List>
  );
};

export const DeviceWebPortInput = <T extends FieldValues>({
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
      validate: validateDeviceWebPort,
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
    if (validateDeviceWebPort(value)) field.onChange(Number(value));
  };

  return (
    <InputGroup>
      <InputGroupInput
        {...field}
        onChange={handleChange}
        type="number"
        placeholder={t(`${I18nVariant.Form}.webPortInput.placeholder`, {
          example: 9000,
        })}
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

function validateDeviceWebPort(value: number | string) {
  const value_ = Number(value);
  if (isNaN(value_) || !Number.isInteger(value_)) {
    return false;
  }
  if (value_.toString() !== value.toString()) {
    return false;
  }
  return value_ >= 1 && value_ <= 65535;
}
