/**
 * App Node Page Definition
 */

import Button from "@mui/material/Button";
import DialogActions from "@mui/material/DialogActions";
import TextField from "@mui/material/TextField";
import Stack from "@mui/material/Stack";
import FormControl from "@mui/material/FormControl";
import FormLabel from "@mui/material/FormLabel";
import FormControlLabel from "@mui/material/FormControlLabel";
import Switch from "@mui/material/Switch";

import {
  useMemo,
  useRef,
  useState,
  Fragment,
  createContext,
  useContext,
  useEffect,
} from "react";
import type { FormEvent } from "react";
import {
  useQuery,
  useMutation,
  useIsMutating,
  useIsFetching,
  useQueryClient,
} from "@tanstack/react-query";
import {
  useForm,
  useFormContext,
  FormProvider,
  Controller,
} from "react-hook-form";

import type { AppNode, AppNodeFields } from "./api";
import { selectAppNode, saveAppNode, deleteAppNode } from "./api";
import { useNavigate, useParams } from "../Route/Router";
import {
  AppNodePath,
  AppNodeCreateI18nKey,
  AppNodeEditI18nKey,
  AppNodeCreatePath,
  AppNodeNS,
} from "./Route";
import { generateEditPath } from "../Route/utils";
import { useTranslation } from "../I18Next/Context";
import { PageI18Next } from "../I18Next/Page";
import type { NodeQRCodeValue } from "../AppSettings/QRCode";
import { NodeFileQRScan, NodeQRCode, NodeQRScan } from "../AppSettings/QRCode";
import { PageHeader } from "../Master/Header";
import {
  PageLayout,
  PageTopBar,
  PageNotFound,
  PageHeaderTitle,
} from "../Master/Page";
import { ExtFSPath } from "../ExtFS/Route";
import { RouteMore } from "../Route/More";

import {
  Dialog,
  DialogConfirmActions,
  DialogConfirmContent,
} from "../Common/Dialog";
import { PageProvider } from "../Master/Context";

const APP_NODE_QUERY_KEY = ["app-node"];
const APP_NODE_DELETE_MUTATION_KEY = ["app-node-delete"];
const APP_NODE_SAVE_MUTATION_KEY = ["app-node-save"];
const APP_NODE_FORM_KEY = "app-node-form";

type AppNodeContextType = { id: AppNode["id"]; data?: AppNode | null };
const AppNodeContext = createContext<AppNodeContextType | null>(null);

const AppNodeDataProvider = ({
  id,
  children,
}: {
  id: AppNode["id"];
  children: React.ReactNode;
}) => {
  const { data, isFetching } = useQuery({
    queryKey: [...APP_NODE_QUERY_KEY, id],
    queryFn: () => selectAppNode(id as AppNode["id"]),
  });

  const value = id ? { id, data: data || isFetching ? data : null } : null;
  return (
    <AppNodeContext.Provider value={value}>{children}</AppNodeContext.Provider>
  );
};

const AppNodeNetworkAddrSeq = "\n";
type AppNodeFormFields = { networkAddrContent: string } & Omit<
  AppNodeFields,
  "networkAddrs"
>;
const AppNodeFormProvider = ({ children }: { children: React.ReactNode }) => {
  const { data } = useContext(AppNodeContext) || {};

  const defaultValues = useMemo(
    () =>
      data
        ? {
            ...data,
            networkAddrContent: data.networkAddrs
              ? data.networkAddrs.join(AppNodeNetworkAddrSeq)
              : "",
            networkAddrs: void 0,
          }
        : { peerId: "", name: "", blocked: false, networkAddrContent: "" },
    [data]
  );

  const methods = useForm<AppNodeFormFields>({
    defaultValues,
  });

  const { reset } = methods;

  useEffect(() => {
    if (data && reset) {
      reset(defaultValues);
    }
  }, [
    defaultValues.name,
    defaultValues.peerId,
    defaultValues.blocked,
    defaultValues.networkAddrContent,
  ]);

  return <FormProvider {...methods}>{children}</FormProvider>;
};

const AppNodeForm = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const { id } = useContext(AppNodeContext) || {};
  const { handleSubmit, control, watch } =
    useFormContext<AppNodeFormFields>() || {};

  const { mutate: saveMutate, isPending } = useMutation({
    mutationKey: [...APP_NODE_SAVE_MUTATION_KEY, id],
    mutationFn: async (fields: AppNodeFields) => await saveAppNode(fields, id),
    onSuccess: (entity) => {
      if (entity && id === void 0) {
        navigate(generateEditPath(AppNodePath, entity.id));
        return;
      }
    },
  });

  const isMutation = useIsMutating({
    mutationKey: [...APP_NODE_DELETE_MUTATION_KEY, id],
  });

  const blocked = watch && watch("blocked");

  const handleSave = async (event: FormEvent) => {
    event.preventDefault();
    await handleSubmit(async (data) => {
      const { networkAddrContent, ...variables } = data;
      const networkAddrs: string[] = [];

      if (networkAddrContent && networkAddrContent.length > 0) {
        const lines = networkAddrContent.trim().split(AppNodeNetworkAddrSeq);
        for (const line of lines) {
          const trimmedLine = line.trim();
          if (trimmedLine.length > 0 && !networkAddrs.includes(trimmedLine)) {
            networkAddrs.push(trimmedLine);
          }
        }
      }
      await saveMutate({ ...variables, networkAddrs });
    })();
  };

  const disabled = useMemo(() => {
    return !handleSubmit || isPending || isMutation > 0;
  }, [isPending, isMutation, handleSubmit]);

  if (!control) {
    return null;
  }

  return (
    <Stack
      component={"form"}
      id={`${APP_NODE_FORM_KEY}-${id || ""}`}
      onSubmit={handleSave}
      direction="row"
      spacing={5}
      alignItems="flex-start"
      justifyContent="flex-start"
      useFlexGap
      flexWrap="wrap"
      width={"100%"}
    >
      <AppNodeQRScan />
      <Stack
        spacing={1}
        minWidth={200}
        width={{ sm: "100%", md: "calc(100% - 240px)" }}
      >
        <Controller
          control={control}
          name="name"
          rules={{ required: true }}
          disabled={disabled}
          render={({ field }) => (
            <TextField
              label={t("resources.app/nodes.fields.name")}
              fullWidth
              variant="filled"
              {...field}
            />
          )}
        />
        <Controller
          control={control}
          name="peerId"
          rules={{ required: true }}
          disabled={disabled}
          render={({ field }) => (
            <TextField
              label={t("resources.app/nodes.fields.peerId")}
              fullWidth
              multiline
              rows={3}
              variant="filled"
              {...field}
            />
          )}
        />
        <FormControl component="fieldset">
          <FormLabel component="legend">
            {t("resources.app/nodes.fields.blocked")}
          </FormLabel>
          <Controller
            control={control}
            name="blocked"
            disabled={disabled}
            render={({ field }) => (
              <FormControlLabel
                label={t(
                  blocked === false ? "labels.disabled" : "labels.enabled"
                )}
                labelPlacement="end"
                control={<Switch {...field} checked={!!blocked} />}
              />
            )}
          />
        </FormControl>
        <Controller
          control={control}
          name="networkAddrContent"
          disabled={disabled}
          render={({ field }) => (
            <TextField
              label={t("resources.app/nodes.fields.networkAddr")}
              fullWidth
              multiline
              rows={8}
              variant="filled"
              {...field}
            />
          )}
        />
      </Stack>
    </Stack>
  );
};

const ResetButton = ({ children }: { children?: React.ReactNode }) => {
  const { t } = useTranslation();
  const { id } = useContext(AppNodeContext) || {};
  const queryClient = useQueryClient();

  const isDeleteMutation = useIsMutating({
    mutationKey: [...APP_NODE_DELETE_MUTATION_KEY, id],
  });
  const isSaveMutation = useIsMutating({
    mutationKey: [...APP_NODE_SAVE_MUTATION_KEY, id],
  });
  const isFetching = useIsFetching({ queryKey: [...APP_NODE_QUERY_KEY, id] });
  const disabled = useMemo(
    () => isFetching > 0 || isDeleteMutation > 0 || isSaveMutation > 0,
    [isDeleteMutation, isSaveMutation, isFetching]
  );

  const handleClick = async () => {
    await queryClient.refetchQueries({
      queryKey: [...APP_NODE_QUERY_KEY, id],
      type: "active",
    });
  };

  return (
    <Button
      variant="contained"
      size="small"
      color="warning"
      onClick={handleClick}
      disabled={disabled}
      hidden={!id}
    >
      {children || t("buttons.reset")}
    </Button>
  );
};

const SaveButton = ({ children }: { children?: React.ReactNode }) => {
  const { t } = useTranslation();
  const { id } = useContext(AppNodeContext) || {};

  const { formState } = useFormContext<AppNodeFields>() || {};
  const isMutation = useIsMutating({
    mutationKey: [...APP_NODE_DELETE_MUTATION_KEY, id],
  });
  const isFetching = useIsFetching({ queryKey: [...APP_NODE_QUERY_KEY, id] });

  const disabled = useMemo(
    () =>
      isFetching > 0 ||
      isMutation > 0 ||
      formState?.isSubmitting ||
      !formState?.isValid,
    [formState, isFetching, isMutation]
  );

  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLButtonElement>(null);
  const handleClick = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleConfirm = (confirm: boolean) => {
    if (confirm) {
      ref.current?.click();
    }
    handleClose();
  };

  return (
    <Fragment>
      <Button
        variant="contained"
        size="small"
        onClick={handleClick}
        disabled={disabled}
      >
        {children ?? t("buttons.save")}
      </Button>
      <button
        type="submit"
        form={`${APP_NODE_FORM_KEY}-${id || ""}`}
        ref={ref}
        hidden={true}
        disabled={disabled}
      />
      <Dialog open={open} onClose={handleClose}>
        <DialogConfirmContent label={t("buttons.save")} />
        <DialogActions>
          <DialogConfirmActions
            label={t("buttons.save")}
            onConfirm={handleConfirm}
          />
        </DialogActions>
      </Dialog>
    </Fragment>
  );
};

const DeleteButton = ({ children }: { children?: React.ReactNode }) => {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const { id } = useContext(AppNodeContext) || {};

  const isMutation = useIsMutating({
    mutationKey: [...APP_NODE_SAVE_MUTATION_KEY, id],
  });
  const isFetching = useIsFetching({ queryKey: [...APP_NODE_QUERY_KEY, id] });

  const { mutate, isPending } = useMutation({
    mutationKey: [...APP_NODE_DELETE_MUTATION_KEY, id],
    mutationFn: deleteAppNode,
    onSuccess: () => {
      navigate(ExtFSPath);
    },
  });

  const disabled = useMemo(
    () => isFetching > 0 || isMutation > 0 || isPending,
    [isMutation, isFetching, isPending]
  );

  const handleClick = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleConfirm = (confirm: boolean) => {
    if (confirm) mutate(id as AppNode["id"]);
    else handleClose();
  };

  return (
    <Fragment>
      <Button
        variant="contained"
        size="small"
        color="error"
        onClick={handleClick}
        hidden={!id}
        disabled={disabled}
      >
        {children ?? t("buttons.remove")}
      </Button>
      <Dialog open={open} onClose={handleClose}>
        <DialogConfirmContent label={t("buttons.remove")} />
        <DialogActions>
          <DialogConfirmActions
            label={t("buttons.remove")}
            onConfirm={handleConfirm}
          />
        </DialogActions>
      </Dialog>
    </Fragment>
  );
};

/**
 * QRCode Component of AppNode
 *
 * A component that renders a QR code for the app node's peerId. It uses
 * the `useWatch` hook from `react-hook-form` to get the values of the
 * `name` and `peerId` fields from the form context.
 *
 * @returns {ReactElement} A React element representing the QR code.
 */
const AppNodeQRCode = () => {
  const { watch } = useFormContext() || {};
  const name = watch ? watch("name") : "";
  const peerId = watch ? watch("peerId") : "";
  return <NodeQRCode name={name} peerId={peerId} />;
};

/**
 * QR Scan Component of AppNode
 *
 * A component that integrates QR code scanning functionality for app nodes.
 * It uses the `useFormContext` hook from `react-hook-form` to set the values
 * of the `peerId` and `name` fields in the form context when a QR code is
 * scanned.
 *
 * The component renders a stack layout containing the AppNodeQRCode component
 * and two QR scanning components: NodeQRScan and NodeFileQRScan. Both scanners
 * share the same `onQRScan` callback function, which updates the form fields
 * upon successful QR code scanning.
 *
 * @returns {ReactElement} A React element representing the QR scan interface
 * for app nodes.
 */

const AppNodeQRScan = () => {
  const { setValue } = useFormContext() || {};

  const onQRScan = ({ peerId, name }: NodeQRCodeValue) => {
    if (!setValue) return;
    setValue("peerId", peerId, { shouldValidate: true, shouldDirty: true });
    setValue("name", name, { shouldValidate: true, shouldDirty: true });
  };

  return (
    <Stack spacing={2} alignItems="center" justifyContent={"space-between"}>
      <AppNodeQRCode />
      <Stack
        direction="row"
        spacing={1}
        alignItems="center"
        justifyContent={"space-between"}
      >
        <NodeQRScan onQRScan={onQRScan} />
        <NodeFileQRScan onQRScan={onQRScan} />
      </Stack>
    </Stack>
  );
};

const AppNodeAddons = () => {
  const { id } = useParams();
  const path = id === void 0 ? AppNodeCreatePath : void 0;
  return <RouteMore path={path}></RouteMore>;
};

const AppNodeCreateTitle = () => {
  const { t } = useTranslation();
  return <PageHeaderTitle>{t(AppNodeCreateI18nKey)}</PageHeaderTitle>;
};

const AppNodeCreateProvider = ({
  children,
}: {
  children?: React.ReactNode;
}) => (
  <AppNodeContext.Provider value={null}>
    <AppNodeFormProvider>{children}</AppNodeFormProvider>
  </AppNodeContext.Provider>
);

const AppNodeCreateForm = () => (
  <PageLayout>
    <PageTopBar>
      <SaveButton />
    </PageTopBar>
    <AppNodeForm />
  </PageLayout>
);

/**
 * Page components for creating new AppNodes
 *
 * A Create component for the app nodes resource that renders a form with
 * input fields for the name and peerId fields of the app node, and a boolean
 * input for the blocked field.
 *
 * The component also renders an AppNodeQRScan component, which is a QR code
 * scanner for the app node's peerId. The AppNodeQRScan component is displayed
 * on the left side of the form, and the input fields are displayed on the right
 * side.
 *
 * The component is wrapped in a Create component, which is a special type of
 * component provided by react-admin that renders a form with a submit button
 * to create a record in the app nodes resource.
 *
 * The component also renders a ListButton component, which is a button that
 * when clicked, redirects the user to the list of app nodes.
 *
 * @returns {ReactElement} A React element representing the create form for
 * the app nodes resource.
 */
export const AppNodeCreatePage = () => {
  return (
    <Fragment>
      <PageHeader title={<AppNodeCreateTitle />} addons={<AppNodeAddons />} />
      <PageI18Next ns={AppNodeNS} defaultNS={AppNodeNS} />
      <PageProvider Component={AppNodeCreateProvider} />
      <AppNodeCreateForm />
    </Fragment>
  );
};

const AppNodeEditTitle = () => {
  const { t } = useTranslation();
  return <PageHeaderTitle>{t(AppNodeEditI18nKey)}</PageHeaderTitle>;
};

const AppNodeEditProvider = ({ children }: { children?: React.ReactNode }) => {
  const { id } = useParams();
  const id_ = id ? parseInt(id) : 0;
  if (id_ == 0 || isNaN(id_))
    return (
      <AppNodeContext.Provider value={null}>
        <AppNodeFormProvider>{children}</AppNodeFormProvider>
      </AppNodeContext.Provider>
    );
  return (
    <AppNodeDataProvider id={id_}>
      <AppNodeFormProvider>{children}</AppNodeFormProvider>
    </AppNodeDataProvider>
  );
};

const AppNodeEditForm = () => {
  const { id, data } = useContext(AppNodeContext) || {};

  if (!id || data === null) return <PageNotFound />;
  return (
    <PageLayout>
      <PageTopBar>
        <DeleteButton />
        <ResetButton />
        <SaveButton />
      </PageTopBar>
      <AppNodeForm />
    </PageLayout>
  );
};

export const AppNodeEditPage = () => {
  return (
    <Fragment>
      <PageHeader title={<AppNodeEditTitle />} addons={<AppNodeAddons />} />
      <PageI18Next ns={AppNodeNS} defaultNS={AppNodeNS} />
      <PageProvider Component={AppNodeEditProvider} />
      <AppNodeEditForm />
    </Fragment>
  );
};
