import { useTranslation } from "@/components/App/I18Next";
import {
  MonitorSmartphone,
  NetworkManageIcon,
  WebIcon,
  Palette,
  Languages,
} from "./Icon";
import { SettingsName } from "./meta";
import { type ComponentProps } from "react";
import { Separator } from "@/components/ui/separator";
import {
  List,
  ListItem,
  ListItemText,
  ListItemVariant,
  ListItemLink,
  ListItemAvatar,
} from "@/components/App/List";

import { MainI18nPrefix } from "./meta";
import { useBrowser } from "@/components/App/Browser";

export enum TOCChapter {
  DeviceInfo = "device-info",
  Network = "network",
  Web = "web",
  Appearance = "appearance",
  Language = "language",
}

export const useTOCChapter = () => {
  const { window } = useBrowser() || {};
  const tocChapter = (window?.location?.hash || "").substring(1);
  return tocChapter;
};

const TableOfContents = () => {
  const { t } = useTranslation(SettingsName);

  const tocChapter = useTOCChapter();

  return (
    <nav className="w-full flex flex-col gap-1">
      <List>
        <ListItem
          active={
            tocChapter === TOCChapter.DeviceInfo
              ? ListItemVariant.Primary
              : void 0
          }
        >
          <ItemLink
            href={
              tocChapter === TOCChapter.DeviceInfo
                ? "#"
                : `#${TOCChapter.DeviceInfo}`
            }
          >
            <ListItemAvatar as={MonitorSmartphone} />
            <ListItemText>
              {t("device-info.title", { keyPrefix: MainI18nPrefix })}
            </ListItemText>
          </ItemLink>
        </ListItem>
        <ListItem
          active={
            tocChapter === TOCChapter.Network ? ListItemVariant.Primary : void 0
          }
        >
          <ItemLink
            href={
              tocChapter === TOCChapter.Network ? "#" : `#${TOCChapter.Network}`
            }
          >
            <ListItemAvatar as={NetworkManageIcon} />
            <ListItemText>
              {t("network.title", { keyPrefix: MainI18nPrefix })}
            </ListItemText>
          </ItemLink>
        </ListItem>
        <ListItem
          active={
            tocChapter === TOCChapter.Web ? ListItemVariant.Primary : void 0
          }
        >
          <ItemLink
            href={tocChapter === TOCChapter.Web ? "#" : `#${TOCChapter.Web}`}
          >
            <ListItemAvatar as={WebIcon} />
            <ListItemText>
              {t("web.title", { keyPrefix: MainI18nPrefix })}
            </ListItemText>
          </ItemLink>
        </ListItem>
      </List>
      <Separator />
      <List>
        <ListItem
          active={
            tocChapter === TOCChapter.Appearance
              ? ListItemVariant.Primary
              : void 0
          }
        >
          <ItemLink
            href={
              tocChapter === TOCChapter.Appearance
                ? "#"
                : `#${TOCChapter.Appearance}`
            }
          >
            <ListItemAvatar as={Palette} />
            <ListItemText>
              {t("appearance.title", { keyPrefix: MainI18nPrefix })}
            </ListItemText>
          </ItemLink>
        </ListItem>
        <ListItem
          active={
            tocChapter === TOCChapter.Language
              ? ListItemVariant.Primary
              : void 0
          }
        >
          <ItemLink
            href={
              tocChapter === TOCChapter.Language
                ? "#"
                : `#${TOCChapter.Language}`
            }
          >
            <ListItemAvatar as={Languages} />
            <ListItemText>
              {t("languages.title", { keyPrefix: MainI18nPrefix })}
            </ListItemText>
          </ItemLink>
        </ListItem>
      </List>
    </nav>
  );
};

export default TableOfContents;

const ItemLink = ({
  className = "px-6",
  children,
  ...props
}: ComponentProps<typeof ListItemLink>) => {
  return (
    <ListItemLink {...props} className={className}>
      {children}
    </ListItemLink>
  );
};
