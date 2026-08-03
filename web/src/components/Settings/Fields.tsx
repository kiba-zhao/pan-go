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
import { type ComponentProps, useRef, useEffect } from "react";
import { useAppDispatch, useAppExtra } from "@/components/App/Context";
import { switchExtra, ExtraType, ExtraState } from "./Extra";

const FieldsClassName = "first:rounded-t-sm last:rounded-b-sm";
export const InfoFieldsTitle = ({
  className,
}: Pick<ComponentProps<typeof ListTitle>, "className">) => {
  return <ListTitle className={className}>我的设备</ListTitle>;
};

export const InfoFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { type, open } = useAppExtra<ExtraState<ExtraType>>({});
  const nameButtonRef = useRef<HTMLButtonElement>(null);
  const memoButtonRef = useRef<HTMLButtonElement>(null);
  const dispatch = useAppDispatch();

  const handleNameEdit = () => {
    dispatch?.({
      extra: switchExtra(ExtraType.NameEdit),
    });
  };

  const handleMemoEdit = () => {
    dispatch?.({
      extra: switchExtra(ExtraType.MemoEdit),
    });
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
            <ListItemText>备注名</ListItemText>
            <ListItemSmall>desktop</ListItemSmall>
          </ListItemContent>
        </ListItemButton>
      </ListItem>
      <ListItem className={FieldsClassName}>
        <ListItemButton onClick={handleMemoEdit} ref={memoButtonRef}>
          <ListItemContent>
            <ListItemText>我的备忘</ListItemText>
            <ListItemSmall>无</ListItemSmall>
          </ListItemContent>
        </ListItemButton>
      </ListItem>
    </List>
  );
};

export const NetworkFieldsTitle = ({
  className,
}: Pick<ComponentProps<typeof ListTitle>, "className">) => {
  return <ListTitle className={className}>网络互联</ListTitle>;
};

export const NetworkFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { type, open } = useAppExtra<ExtraState<ExtraType>>({});
  const portButtonRef = useRef<HTMLButtonElement>(null);
  const publicAddrsButtonRef = useRef<HTMLButtonElement>(null);
  const broadcastAddrsButtonRef = useRef<HTMLButtonElement>(null);
  const dispatch = useAppDispatch();

  const handlePortEdit = () => {
    dispatch?.({
      extra: switchExtra(ExtraType.PeerPortEdit),
    });
  };

  const handlePublicAddrsEdit = () => {
    dispatch?.({
      extra: switchExtra(ExtraType.PublicAddrsEdit),
    });
  };

  const handleBroadcastAddrsEdit = () => {
    dispatch?.({
      extra: switchExtra(ExtraType.BroadcastAddrsEdit),
    });
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

  return (
    <List className={className}>
      <ListItem className={FieldsClassName}>
        <ListItemContent>
          <ListItemLabel htmlFor="network-fields-enabled">
            网络服务
          </ListItemLabel>
          <ListItemSwitch id="network-fields-enabled" />
        </ListItemContent>
      </ListItem>
      <ListItem className={FieldsClassName}>
        <ListItemButton onClick={handlePortEdit} ref={portButtonRef}>
          <ListItemContent>
            <ListItemText>网络端口</ListItemText>
            <ListItemSmall>9000</ListItemSmall>
          </ListItemContent>
        </ListItemButton>
      </ListItem>
      <ListItem className={FieldsClassName}>
        <ListItemContent>
          <ListItemLabel htmlFor="network-fields-broadcast-enabled">
            网络广播
          </ListItemLabel>
          <ListItemSwitch id="network-fields-broadcast-enabled" />
        </ListItemContent>
      </ListItem>
      <ListItem className={FieldsClassName}>
        <ListItemButton
          onClick={handleBroadcastAddrsEdit}
          ref={broadcastAddrsButtonRef}
        >
          <ListItemContent>
            <ListItemText>广播地址</ListItemText>
            <ListItemSmall>多播网络地址</ListItemSmall>
            <ListItemMore />
          </ListItemContent>
        </ListItemButton>
      </ListItem>
      <ListItem className={FieldsClassName}>
        <ListItemButton
          onClick={handlePublicAddrsEdit}
          ref={publicAddrsButtonRef}
        >
          <ListItemContent>
            <ListItemText>公共地址</ListItemText>
            <ListItemSmall>域名/主机名/IP地址 + 端口</ListItemSmall>
            <ListItemMore />
          </ListItemContent>
        </ListItemButton>
      </ListItem>
    </List>
  );
};

export const WebFieldsTitle = ({
  className,
}: Pick<ComponentProps<typeof ListTitle>, "className">) => {
  return <ListTitle className={className}>Web服务</ListTitle>;
};

export const WebFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  const { type, open } = useAppExtra<ExtraState<ExtraType>>({});
  const webPortButtonRef = useRef<HTMLButtonElement>(null);
  const dispatch = useAppDispatch();

  const handleWebPortEdit = () => {
    dispatch?.({
      extra: switchExtra(ExtraType.WebPortEdit),
    });
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

  return (
    <List className={className}>
      <ListItem className={FieldsClassName}>
        <ListItemContent>
          <ListItemLabel htmlFor="web-fields-enabled">开启服务</ListItemLabel>
          <ListItemSwitch id="web-fields-enabled" />
        </ListItemContent>
      </ListItem>
      <ListItem className={FieldsClassName}>
        <ListItemContent>
          <ListItemLabel htmlFor="web-fields-allowed">仅本机访问</ListItemLabel>
          <ListItemSwitch id="web-fields-allowed" />
        </ListItemContent>
      </ListItem>
      <ListItem className={FieldsClassName}>
        <ListItemButton onClick={handleWebPortEdit} ref={webPortButtonRef}>
          <ListItemContent>
            <ListItemText>Web端口</ListItemText>
            <ListItemSmall>9000</ListItemSmall>
          </ListItemContent>
        </ListItemButton>
      </ListItem>
    </List>
  );
};

export const AppearanceFieldsTitle = ({
  className,
}: Pick<ComponentProps<typeof ListTitle>, "className">) => {
  return <ListTitle className={className}>外观</ListTitle>;
};

export const AppearanceFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  return (
    <List className={className}>
      <ListItem className={FieldsClassName}>
        <ListItemText>主题模式</ListItemText>
        <ListItemSmall>跟随系统</ListItemSmall>
      </ListItem>
    </List>
  );
};

export const LanguageFieldsTitle = ({
  className,
}: Pick<ComponentProps<typeof ListTitle>, "className">) => {
  return <ListTitle className={className}>语言</ListTitle>;
};

export const LanguageFields = ({
  className,
}: Pick<ComponentProps<typeof List>, "className">) => {
  return (
    <List className={className}>
      <ListItem className={FieldsClassName}>
        <ListItemText>应用语言</ListItemText>
        <ListItemSmall>跟随系统</ListItemSmall>
      </ListItem>
    </List>
  );
};
