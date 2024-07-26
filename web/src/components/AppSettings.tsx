import { Title, useTranslate } from "react-admin";
import { useLocation, useNavigate } from "react-router-dom";

import Autocomplete from "@mui/material/Autocomplete";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardContent from "@mui/material/CardContent";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import Divider from "@mui/material/Divider";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { useMemo, useState } from "react";

import { FilePathInput, QRCode } from "./Custom";
import NotFound from "./NotFound";

import type { ReactNode } from "react";

const TabList = ["", "network-address"];
const TabUrlList = TabList.map((_) =>
  _.length > 0 ? `/app/settings/${_}` : `/app/settings`
);
const TabLabelList = TabList.map((_) =>
  _.length > 0 ? `custom.app/settings.${_}` : "custom.app/settings.summary"
);

const TabPanel = ({
  children,
  hidden,
}: {
  children?: ReactNode;
  hidden: boolean;
}) => (
  <Box hidden={hidden} padding={2}>
    {children}
  </Box>
);

export const AppSettings = () => {
  const location = useLocation();
  const value = useMemo(() => location.pathname, [location.pathname]);
  if (!TabUrlList.includes(value)) {
    return <NotFound />;
  }

  const t = useTranslate();
  const navigate = useNavigate();
  const onChange = (_: React.SyntheticEvent, newValue: string) => {
    navigate(newValue);
  };

  return (
    <>
      <Title title={t("custom.app/settings.name")} />
      <Paper>
        <Stack
          direction="row"
          spacing={1}
          alignItems="center"
          justifyContent={"space-between"}
        >
          <Tabs value={value} onChange={onChange}>
            {TabList.map((_, index) => (
              <Tab
                key={`custom.app/settings.${_}`}
                label={t(TabLabelList[index])}
                value={TabUrlList[index]}
              />
            ))}
          </Tabs>
          <Box></Box>
        </Stack>

        <TabPanel hidden={value !== TabUrlList[0]}>
          <AppSummarySettings />
        </TabPanel>
        <TabPanel hidden={value !== TabUrlList[1]}>
          <AppNetworkAddressSettings />
        </TabPanel>
      </Paper>
    </>
  );
};

const AppSummarySettings = () => {
  const t = useTranslate();
  return (
    <Stack
      padding={3}
      direction="row"
      spacing={5}
      alignItems="flex-start"
      justifyContent="flex-start"
      useFlexGap
      flexWrap="wrap"
    >
      <Stack spacing={2} alignItems="center" justifyContent={"space-between"}>
        <QRCode value="sample text" width={200} />
        <Stack
          direction="row"
          spacing={1}
          alignItems="center"
          justifyContent={"space-between"}
        >
          <Button variant="contained" size="small">
            {t("custom.button.copy")}
          </Button>
          <Button variant="contained" color="error" size="small">
            {t("custom.button.renew")}
          </Button>
        </Stack>
      </Stack>
      <Stack spacing={2} minWidth={200} maxWidth={760} width={"70%"}>
        <FilePathInput label={t("custom.app/settings.fields.rootPath")} />
        <TextField
          label={t("custom.app/settings.fields.nodeId")}
          fullWidth
          variant="filled"
          InputProps={{
            readOnly: true,
          }}
          defaultValue={"sample node id"}
        />
      </Stack>
    </Stack>
  );
};

type NewAddressButtonProps<T extends any> = {
  options: T[];
};

const NewAddressButton = <T extends any>({
  options,
}: NewAddressButtonProps<T>) => {
  const t = useTranslate();
  const [open, setOpen] = useState(false);
  const onOpen = () => setOpen(true);
  const onClose = () => setOpen(false);
  const onSubmit = () => {
    onClose();
  };
  return (
    <>
      <Button variant="contained" color="success" size="small" onClick={onOpen}>
        {t("custom.button.new")}
      </Button>
      <Dialog open={open} onClose={onClose}>
        <DialogTitle>
          <Stack direction="row" spacing={0.2} alignItems="flex-start">
            <Typography variant="h6">
              {t("custom.app/settings.network-address")}
            </Typography>
            <Typography
              variant="caption"
              bgcolor={"green"}
              height={0.5}
              paddingX={0.5}
              borderRadius={0.5}
            >
              {t("custom.button.new")}
            </Typography>
          </Stack>
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            请选择需要增加的网络地址类型进行提交，提交后将在下面对应的分类末尾添加默认地址项。请根据需要进行修改。
          </DialogContentText>
          <Autocomplete
            options={options}
            renderInput={(params) => (
              <TextField
                {...params}
                label={t("custom.app/settings.network-address")}
                variant="standard"
              />
            )}
          />
        </DialogContent>
        <DialogActions>
          <Button size="small" onClick={onClose}>
            {t("custom.button.cancel")}
          </Button>
          <Button size="small" onClick={onSubmit}>
            {t("custom.button.submit")}
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

const EditAddressButton = () => {
  const t = useTranslate();
  return (
    <Button variant="outlined" size="small">
      {t("custom.button.edit")}
    </Button>
  );
};

const RemoveAddressButton = () => {
  const t = useTranslate();
  return (
    <Button variant="contained" color="error" size="small">
      {t("custom.button.remove")}
    </Button>
  );
};

const APP_NETWORK_ADDRESS_TYPES = ["web", "node", "broadcast", "public"];
const AppNetworkAddressSettings = () => {
  const t = useTranslate();
  const addressTypes = useMemo(() => {
    return APP_NETWORK_ADDRESS_TYPES.map((_) => ({
      label: t(`custom.app/settings.${_}-address`),
      value: _,
    }));
  }, []);

  return (
    <Stack spacing={2}>
      <Stack
        direction="row"
        spacing={2}
        alignItems="center"
        justifyContent={"flex-end"}
      >
        <NewAddressButton options={addressTypes} />
        <Button variant="contained" color="error" size="small">
          {t("custom.button.save")}
        </Button>
        <Button variant="contained" size="small">
          {t("custom.button.reset")}
        </Button>
      </Stack>
      {addressTypes.map(({ label, value }) => (
        <AppAddressSettings
          key={value}
          avatar={value[0].toUpperCase()}
          title={label}
          addresses={["0.0.0.0:9000"]}
        >
          <EditAddressButton />
          <RemoveAddressButton />
        </AppAddressSettings>
      ))}
    </Stack>
  );
};

type AppAddressSettingsProps = {
  title: string;
  avatar: string;
  addresses: string[];
  children: ReactNode;
};
const AppAddressSettings = ({
  title,
  avatar,
  addresses,
  children,
}: AppAddressSettingsProps) => {
  return (
    <Stack spacing={2}>
      <Divider>{title}</Divider>
      <Stack
        direction="row"
        spacing={3}
        alignItems="center"
        useFlexGap
        flexWrap="wrap"
      >
        {addresses.map((address, i) => (
          <AppAddressSettingsItem key={i} avatar={avatar} address={address}>
            {children}
          </AppAddressSettingsItem>
        ))}
      </Stack>
    </Stack>
  );
};

type AppAddressSettingsItemProps = {
  avatar: string;
  address: string;
  children: ReactNode;
};
const AppAddressSettingsItem = ({
  avatar,
  address,
  children,
}: AppAddressSettingsItemProps) => {
  const t = useTranslate();

  const [ip, port] = useMemo(() => {
    const parts = address.split(":");
    const port = Number(parts.at(-1));
    const ip = parts.slice(0, -1).join(":");
    return [ip, port];
  }, [address]);

  return (
    <Card raised={true}>
      <CardContent>
        <Stack
          spacing={2}
          direction="row"
          alignItems="center"
          justifyContent="flex-start"
        >
          <Avatar
            variant="rounded"
            sx={{ width: 74, height: 74, fontSize: 48 }}
          >
            {avatar}
          </Avatar>
          <Stack>
            <TextField
              label="地址"
              value={ip}
              variant="outlined"
              inputProps={{ readOnly: true }}
            />
            <TextField
              label="端口"
              value={port}
              variant="outlined"
              inputProps={{ readOnly: true }}
            />
          </Stack>
        </Stack>
      </CardContent>
      <CardActions>
        <Stack
          direction="row"
          spacing={1}
          alignItems="center"
          justifyContent="flex-end"
          width="100%"
        >
          {children}
        </Stack>
      </CardActions>
    </Card>
  );
};
