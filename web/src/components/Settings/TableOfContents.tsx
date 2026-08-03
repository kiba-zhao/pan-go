import { useTranslation } from "@/components/App/I18Next";
import {
  MonitorSmartphone,
  NetworkManageIcon,
  WebIcon,
  Palette,
  Languages,
  Keyboard,
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

const TableOfContents = () => {
  const { t } = useTranslation(SettingsName);
  return (
    <nav className="w-full flex flex-col gap-1">
      <List>
        <ListItem active={ListItemVariant.Primary}>
          <ItemLink href="#1">
            <ListItemAvatar as={MonitorSmartphone} />
            <ListItemText>我的设备</ListItemText>
          </ItemLink>
        </ListItem>
        <ListItem>
          <ItemLink href="#2">
            <ListItemAvatar as={NetworkManageIcon} />
            <ListItemText>网络互联</ListItemText>
          </ItemLink>
        </ListItem>
        <ListItem>
          <ItemLink href="#3">
            <ListItemAvatar as={WebIcon} />
            <ListItemText>Web服务</ListItemText>
          </ItemLink>
        </ListItem>
      </List>
      <Separator />
      <List>
        <ListItem>
          <ItemLink href="#4">
            <ListItemAvatar as={Palette} />
            <ListItemText>外观</ListItemText>
          </ItemLink>
        </ListItem>
        <ListItem>
          <ItemLink href="#5">
            <ListItemAvatar as={Languages} />
            <ListItemText>语言</ListItemText>
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
