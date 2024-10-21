import { Title, useTranslate } from "react-admin";

import AdminPanelSettingsIcon from "@mui/icons-material/AdminPanelSettings";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardContent from "@mui/material/CardContent";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import FormControl from "@mui/material/FormControl";
import FormControlLabel from "@mui/material/FormControlLabel";
import FormLabel from "@mui/material/FormLabel";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Switch from "@mui/material/Switch";
import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";

import type { AppSettingsFields } from "../API";
import { useAPI } from "../API";
import {
  Dialog,
  DialogConfirmActions,
  DialogConfirmContent,
  DialogSubmitActions,
} from "./Feedback/Dialog";
import { QRCodeDownloadButton } from "./QRCode/Base";
import { NodeQRCode } from "./QRCode/Node";

import {
  useIsFetching,
  useMutation,
  usePrefetchQuery,
  useQueryClient,
  useSuspenseQuery,
} from "@tanstack/react-query";
import type { ChangeEvent, ReactNode } from "react";
import {
  Fragment,
  Suspense,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";
import {
  Controller,
  FormProvider,
  useFieldArray,
  useForm,
  useFormContext,
} from "react-hook-form";

export const AppSettingsIcon = AdminPanelSettingsIcon;
export const AppSettingsRoutePath = "/app/settings";

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

export const APP_SETTINGS_QUERY_KEY = ["app-settings"];
export const AppSettings = () => {
  const t = useTranslate();
  const [tab, setTab] = useState("summary");
  const onChange = (_: React.SyntheticEvent, newValue: string) => {
    setTab(newValue);
  };

  const api = useAPI();
  usePrefetchQuery({
    queryKey: APP_SETTINGS_QUERY_KEY,
    queryFn: api?.selectAllAppSettings,
  });

  return (
    <Suspense fallback={<AppSettingsLoading />}>
      <Title title={t("custom.app/settings.name")} />
      <Paper>
        <Tabs value={tab} onChange={onChange}>
          <Tab label={t("custom.app/settings.summary")} value={"summary"} />
          <Tab
            label={t("custom.app/settings.web-address")}
            value={"web-address"}
          />
          <Tab
            label={t("custom.app/settings.node-address")}
            value={"node-address"}
          />
          <Tab
            label={t("custom.app/settings.broadcast-address")}
            value={"broadcast-address"}
          />
          <Tab
            label={t("custom.app/settings.public-address")}
            value={"public-address"}
          />
        </Tabs>
        <TabPanel hidden={tab !== "summary"}>
          <AppSummarySettings />
        </TabPanel>
        <TabPanel hidden={tab !== "web-address"}>
          <WebAddressSettings />
        </TabPanel>
        <TabPanel hidden={tab !== "node-address"}>
          <NodeAddressSettings />
        </TabPanel>
        <TabPanel hidden={tab !== "broadcast-address"}>
          <BroadcastAddressSettings />
        </TabPanel>
        <TabPanel hidden={tab !== "public-address"}>
          <PublicAddressSettings />
        </TabPanel>
      </Paper>
    </Suspense>
  );
};

const RefreshButton = () => {
  const t = useTranslate();

  const isFetching = useIsFetching({
    queryKey: APP_SETTINGS_QUERY_KEY,
  });

  const queryClient = useQueryClient();
  const onClick = async () => {
    if (isFetching) return;
    await queryClient.refetchQueries({
      queryKey: APP_SETTINGS_QUERY_KEY,
      type: "active",
    });
  };
  return (
    <Button variant="contained" size="small" onClick={onClick}>
      {t("custom.button.refresh")}
    </Button>
  );
};

const SaveButton = () => {
  const t = useTranslate();

  const [open, setOpen] = useState(false);
  const onOpen = () => setOpen(true);
  const onClose = () => setOpen(false);

  const api = useAPI();
  const { mutate } = useMutation({
    mutationFn: api?.saveAppSettings,
  });
  const onSave = useCallback(
    async (fields: AppSettingsFields) => {
      return await mutate(fields);
    },
    [api]
  );
  const { handleSubmit } = useFormContext();

  const onConfirm = async (confirm: boolean) => {
    if (confirm) {
      await handleSubmit(onSave)();
    }
    onClose();
  };

  return (
    <Fragment>
      <Button variant="contained" size="small" color="error" onClick={onOpen}>
        {t("custom.button.save")}
      </Button>
      <Dialog open={open} onClose={onClose}>
        <DialogConfirmContent label={t("custom.button.save")} />
        <DialogActions>
          <DialogConfirmActions
            label={t("custom.button.save")}
            onConfirm={onConfirm}
          />
        </DialogActions>
      </Dialog>
    </Fragment>
  );
};

const AppSettingsLoading = () => {
  return "App Settings Loading...";
};

const AppSummarySettings = () => {
  const t = useTranslate();

  const api = useAPI();
  const { data, isFetching, isError } = useSuspenseQuery({
    queryKey: APP_SETTINGS_QUERY_KEY,
    queryFn: api?.selectAllAppSettings,
  });
  const defaultValues = useMemo(
    () => ({
      name: data?.name,
      guardEnabled: data?.guardEnabled,
      guardAccess: data?.guardAccess,
    }),
    [data]
  );
  const methods = useForm({
    defaultValues,
  });
  const { reset, control, watch } = methods;

  useEffect(() => {
    if (isFetching) return;
    if (isError) return;
    reset(defaultValues);
  }, [isFetching, isError]);

  const guardEnabled = watch("guardEnabled");
  const guardAccess = watch("guardAccess");

  return (
    <Fragment>
      <FormProvider {...methods}>
        <Stack
          padding={3}
          direction="row"
          spacing={2}
          alignItems="center"
          justifyContent="flex-end"
        >
          <SaveButton />
          <RefreshButton />
        </Stack>
      </FormProvider>
      <Stack
        padding={3}
        direction="row"
        spacing={5}
        alignItems="flex-start"
        justifyContent="flex-start"
        useFlexGap
        flexWrap="wrap"
      >
        <NodeQRCode name={data?.name} nodeId={data?.nodeId}>
          <Stack
            direction="row"
            spacing={1}
            alignItems="center"
            justifyContent={"space-between"}
          >
            <QRCodeDownloadButton />
            <Button variant="contained" color="error" size="small">
              {t("custom.button.renew")}
            </Button>
          </Stack>
        </NodeQRCode>
        <Stack spacing={2} minWidth={200} maxWidth={760} width={"70%"}>
          <TextField
            label={t("custom.app/settings.fields.rootPath")}
            fullWidth
            variant="filled"
            value={data?.rootPath}
          />
          <Controller
            control={control}
            name="name"
            render={({ field }) => (
              <TextField
                label={t("custom.app/settings.fields.name")}
                fullWidth
                variant="filled"
                {...field}
                disabled={isFetching}
              />
            )}
          />

          <TextField
            label={t("custom.app/settings.fields.nodeId")}
            fullWidth
            multiline
            rows={3}
            variant="filled"
            InputProps={{
              readOnly: true,
            }}
            required
            value={data?.nodeId}
          />
          <FormControl component="fieldset">
            <FormLabel component="legend">
              {t("custom.app/settings.fields.guardEnabled")}
            </FormLabel>
            <Controller
              control={control}
              name="guardEnabled"
              render={({ field }) => (
                <FormControlLabel
                  label={t(
                    guardEnabled === false
                      ? "custom.label.disabled"
                      : "custom.label.enabled"
                  )}
                  labelPlacement="end"
                  control={
                    <Switch
                      checked={guardEnabled}
                      {...field}
                      disabled={isFetching}
                    />
                  }
                />
              )}
            />
          </FormControl>
          <FormControl component="fieldset">
            <FormLabel component="legend">
              {t("custom.app/settings.fields.guardAccess")}
            </FormLabel>
            <Controller
              control={control}
              name="guardAccess"
              render={({ field }) => (
                <FormControlLabel
                  label={t(
                    guardAccess === false
                      ? "custom.label.refused"
                      : "custom.label.allowed"
                  )}
                  labelPlacement="end"
                  control={
                    <Switch
                      checked={guardAccess}
                      {...field}
                      disabled={isFetching || guardEnabled === false}
                    />
                  }
                />
              )}
            />
          </FormControl>
        </Stack>
      </Stack>
    </Fragment>
  );
};

type AddressEditDialogProps = {
  open: boolean;
  onClose?: () => void;
  onBlur?: () => void;
  value?: string;
  onChange?: (value: string) => void;
  isNew?: boolean;
};
const AddressEditDialog = ({
  open,
  onClose,
  value,
  onChange,
  onBlur,
  isNew,
}: AddressEditDialogProps) => {
  const t = useTranslate();
  const initValue = useCallback(() => {
    if (!value) return ["", 0];
    const parts = value.split(":");
    const port = Number(parts.at(-1));
    const ip = parts.slice(0, -1).join(":");
    return [ip, port];
  }, [value]);

  const [fields, setFields] = useState(initValue);

  const onSubmit = () => {
    if (onChange) onChange(fields.join(":"));
    if (onClose) onClose();
    setFields(initValue());
  };

  const onCancel = () => {
    if (onBlur) onBlur();
    if (onClose) onClose();
    setFields(initValue());
  };

  const onIPChange = (e: ChangeEvent<HTMLInputElement>) => {
    setFields([e.target.value, fields[1]]);
  };
  const onPortChange = (e: ChangeEvent<HTMLInputElement>) => {
    setFields([fields[0], e.target.value]);
  };

  return (
    <Dialog open={open} onClose={onCancel}>
      <DialogTitle>
        <Stack direction="row" spacing={0.2} alignItems="flex-start">
          <Typography variant="h6">
            {t("custom.app/settings.network-address")}
          </Typography>
          <Typography
            variant="caption"
            bgcolor="green"
            height={0.5}
            paddingX={0.5}
            borderRadius={0.5}
            hidden={!isNew}
          >
            {t("custom.button.new")}
          </Typography>
        </Stack>
      </DialogTitle>
      <DialogContent>
        <TextField
          label={t("custom.app/settings.ip")}
          value={fields[0]}
          onChange={onIPChange}
          variant="outlined"
          focused
        />
        <TextField
          label={t("custom.app/settings.port")}
          value={fields[1]}
          onChange={onPortChange}
          variant="outlined"
          focused
        />
      </DialogContent>
      <DialogActions>
        <DialogSubmitActions
          onSubmit={onSubmit}
          onCancel={onCancel}
          label={isNew ? t("custom.button.new") : t("custom.button.edit")}
        />
      </DialogActions>
    </Dialog>
  );
};

type AddressSource =
  | "webAddress"
  | "broadcastAddress"
  | "nodeAddress"
  | "publicAddress";

const NewAddressButton = ({
  source,
  defaultValue,
}: {
  source: AddressSource;
  defaultValue?: string;
}) => {
  const t = useTranslate();
  const [open, setOpen] = useState(false);
  const onOpen = () => setOpen(true);
  const onClose = () => setOpen(false);

  const { control } = useFormContext();
  const { append } = useFieldArray({
    control,
    name: source,
  });
  return (
    <Fragment>
      <Button variant="contained" color="success" size="small" onClick={onOpen}>
        {t("custom.button.new")}
      </Button>
      <Controller
        control={control}
        name={source}
        render={({ field: { ref, onChange, value, ...field_ } }) => (
          <AddressEditDialog
            isNew
            open={open}
            onClose={onClose}
            value={defaultValue}
            {...field_}
            onChange={(value) => append(value)}
          />
        )}
      />
    </Fragment>
  );
};

const EditAddressButton = ({ source }: { source: string }) => {
  const t = useTranslate();
  const [open, setOpen] = useState(false);
  const onOpen = () => setOpen(true);
  const onClose = () => setOpen(false);

  const { control } = useFormContext();
  return (
    <Fragment>
      <Button variant="outlined" size="small" onClick={onOpen}>
        {t("custom.button.edit")}
      </Button>
      <Controller
        control={control}
        name={source}
        render={({ field: { ref, ...field } }) => (
          <AddressEditDialog open={open} onClose={onClose} {...field} />
        )}
      />
    </Fragment>
  );
  3;
};

const RemoveAddressButton = ({
  source,
  offset,
}: {
  source: AddressSource;
  offset: number;
}) => {
  const t = useTranslate();
  const { control, getValues } = useFormContext();
  const { remove } = useFieldArray({
    control,
    name: source,
  });

  const value: string = useMemo(
    () => getValues(source).at(offset),
    [getValues, source, offset]
  );

  const [open, setOpen] = useState(false);
  const onOpen = () => setOpen(true);
  const onClose = () => setOpen(false);

  const onConfirm = (confirm: boolean) => {
    if (confirm) remove(offset);
    onClose();
  };

  return (
    <Fragment>
      <Button variant="contained" color="error" size="small" onClick={onOpen}>
        {t("custom.button.remove")}
      </Button>
      <Dialog open={open} onClose={onClose}>
        <DialogConfirmContent
          label={t("custom.button.remove")}
          contentLabel={`${t("custom.button.remove").toLowerCase()} ${value}`}
        />
        <DialogActions>
          <DialogConfirmActions
            label={t("custom.button.remove")}
            onConfirm={onConfirm}
          />
        </DialogActions>
      </Dialog>
    </Fragment>
  );
};

type AppAddressSettingsProps = {
  source: AddressSource;
  children: ReactNode;
};
const AppAddressSettings = ({ source, children }: AppAddressSettingsProps) => {
  const avatar = useMemo(() => source[0].toUpperCase(), [source]);

  const api = useAPI();
  const { data, isFetching, isError } = useSuspenseQuery({
    queryKey: APP_SETTINGS_QUERY_KEY,
    queryFn: api?.selectAllAppSettings,
  });
  const defaultValues = useMemo(() => ({ [source]: data[source] }), [data]);
  const methods = useForm({
    defaultValues,
  });
  const { reset, watch } = methods;

  useEffect(() => {
    if (isFetching) return;
    if (isError) return;
    reset(defaultValues);
  }, [isFetching, isError]);

  const addresses = watch(source);

  return (
    <FormProvider {...methods}>
      <Stack
        padding={3}
        direction="row"
        spacing={2}
        alignItems="center"
        justifyContent="flex-end"
      >
        {children}
      </Stack>

      <Stack
        direction="row"
        spacing={3}
        alignItems="center"
        useFlexGap
        flexWrap="wrap"
      >
        {addresses.map((address, i) => (
          <AppAddressSettingsItem
            key={`${source}-${i}`}
            avatar={avatar}
            address={address}
          >
            <EditAddressButton source={`${source}.${i}`} />
            <RemoveAddressButton source={source} offset={i} />
          </AppAddressSettingsItem>
        ))}
      </Stack>
    </FormProvider>
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
              label={t("custom.app/settings.ip")}
              value={ip}
              variant="outlined"
              inputProps={{ readOnly: true }}
            />
            <TextField
              label={t("custom.app/settings.port")}
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

const WebAddressSettings = () => {
  const source = "webAddress";
  return (
    <AppAddressSettings source={source}>
      <NewAddressButton source={source} defaultValue="127.0.0.1:9002" />
      <SaveButton />
      <RefreshButton />
    </AppAddressSettings>
  );
};

const NodeAddressSettings = () => {
  const source = "nodeAddress";
  return (
    <AppAddressSettings source={source}>
      <NewAddressButton source={source} defaultValue="0.0.0.0:9000" />
      <SaveButton />
      <RefreshButton />
    </AppAddressSettings>
  );
};

const BroadcastAddressSettings = () => {
  const source = "broadcastAddress";
  return (
    <AppAddressSettings source={source}>
      <NewAddressButton source={source} defaultValue="224.0.0.120:9100" />
      <SaveButton />
      <RefreshButton />
    </AppAddressSettings>
  );
};

const PublicAddressSettings = () => {
  const source = "publicAddress";
  return (
    <AppAddressSettings source={source}>
      <NewAddressButton source={source} defaultValue="0.0.0.0:9000" />
      <SaveButton />
      <RefreshButton />
    </AppAddressSettings>
  );
};
