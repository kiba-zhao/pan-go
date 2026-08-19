import type { ComponentProps } from "react";
import { SearchIcon } from "./Icon";
import { SettingsName } from "./meta";

import { useTranslation, I18nVariant } from "@/components/App/I18Next";
import { useBrowser } from "@/components/App/Browser";
import { cn } from "@/lib/utils";

import {
  ListTitle,
  List,
  ListItem,
  ListItemText,
  ListItemVariant,
  ListItemLink,
  ListItemAvatar,
} from "@/components/App/List";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";

export enum TOCChapter {
  DeviceInfo = "deviceInfo",
  DeviceNetwork = "deviceNetwork",
  DeviceWeb = "deviceWeb",
  Appearance = "appearance",
  Language = "language",
  DeviceCluster = "deviceCluster",
}

const useTOCChapter = () => {
  const { window } = useBrowser() || {};
  const tocChapter = (window?.location?.hash || "").substring(1);
  return tocChapter;
};

type TOCItemProps = {
  chapter: TOCChapter;
  text: string;
  Avatar: ComponentProps<typeof ListItemAvatar>["as"];
} & Omit<ComponentProps<typeof ListItem>, "children" | "active">;
export const TOCItem = ({ chapter, text, Avatar, ...props }: TOCItemProps) => {
  const tocChapter = useTOCChapter();

  return (
    <ListItem
      {...props}
      active={tocChapter === chapter ? ListItemVariant.Primary : void 0}
    >
      <ListItemLink
        className="px-6"
        href={tocChapter === chapter ? "#" : `#${chapter}`}
      >
        <ListItemAvatar as={Avatar} />
        <ListItemText>{text}</ListItemText>
      </ListItemLink>
    </ListItem>
  );
};

export { List as TOCItemGroup };

export const FieldsSectionClassName = "max-w-3xl w-full pt-2";
type FieldsSectionProps = ComponentProps<"section"> & {
  chapter: TOCChapter;
};
export const FieldsSection = ({
  chapter,
  children,
  className,
  ...props
}: FieldsSectionProps) => {
  const { t } = useTranslation(SettingsName);

  const tocChapter = useTOCChapter();
  if (tocChapter !== chapter && tocChapter.length > 0) {
    return null;
  }

  return (
    <section className={cn(FieldsSectionClassName, className)} {...props}>
      <ListTitle id={chapter}>
        {t(`${I18nVariant.Main}.${chapter}.title`)}
      </ListTitle>
      {children}
    </section>
  );
};

export const FieldsClassName = "first:rounded-t-sm last:rounded-b-sm";

export const FieldsSearchFilter = () => {
  const { t } = useTranslation(SettingsName);
  return (
    <InputGroup>
      <InputGroupInput
        placeholder={t(`${I18nVariant.Main}.search.placeholder`)}
      />
      <InputGroupAddon>
        <SearchIcon />
      </InputGroupAddon>
    </InputGroup>
  );
};
