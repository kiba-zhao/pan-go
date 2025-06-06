/**
 * ExtFSNodeItem Page Definition File
 */
import { Fragment, useEffect, useMemo, useState } from "react";

import Button from "@mui/material/Button";
import DialogActions from "@mui/material/DialogActions";
import FormControl from "@mui/material/FormControl";
import FormControlLabel from "@mui/material/FormControlLabel";
import FormLabel from "@mui/material/FormLabel";
import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import Switch from "@mui/material/Switch";
import TextField from "@mui/material/TextField";
import CircleIcon from "@mui/icons-material/Circle";

import { useMutation, useQuery } from "@tanstack/react-query";
import { Controller, useForm } from "react-hook-form";

import { useNavigate, useParams } from "../Route/Router";
import { useTranslation } from "../I18Next/Context";
import { PageI18Next } from "../I18Next/Page";
import type { ExtFSNodeItem, ExtFSNodeItemFields } from "./api";
import {
  selectExtFSNodeItem,
  saveExtFSNodeItem,
  deleteExtFSNodeItem,
} from "./api";
import {
  Dialog,
  DialogConfirmActions,
  DialogConfirmContent,
} from "../Common/Dialog";
import { FilePathInput } from "../ServerFS/Input";
import {
  ExtFSNodeItemEditI18nKey,
  ExtFSNodeItemPath,
  ExtFSNodeItemCreateI18nKey,
  ExtFSNodeItemCreatePath,
  ExtFSNodeItemI18nKey,
  ExtFSNodeItemNS,
} from "./Route";
import { ExtFSPath } from "../ExtFS/Route";
import { generateEditPath } from "../Route/utils";
import { PageHeader } from "../Master/Header";
import { PageLayout, PageTopBar, PageHeaderTitle } from "../Master/Page";
import { RouteMore } from "../Route/More";

/**
 * Component for creating a new ExtFSNodeItem Panel
 *
 * @returns {JSX.Element} The JSX element representing the ExtFSNodeItem Panel
 */
export const ExtFSNodeItemCreatePage = () => (
  <Fragment>
    <PageI18Next ns={ExtFSNodeItemNS} defaultNS={ExtFSNodeItemNS} />
    <ExtFSNodeItemForm />
  </Fragment>
);

/**
 * SaveButton Component
 *
 * A button component that triggers a save action. It optionally
 * displays a confirmation dialog before executing the save action.
 *
 * @param {Object} props - The properties for the SaveButton component.
 * @param {Function} props.onSave - The function to call when the save action is triggered.
 * @param {boolean} [props.needConfirm] - Whether a confirmation dialog is needed before saving.
 * @param {boolean} [props.disabled] - Whether the button is disabled.
 *
 * @returns {JSX.Element} A JSX element representing the save button with optional confirmation dialog.
 */

const SaveButton = ({
  onSave,
  needConfirm,
  disabled,
}: {
  onSave: () => void;
  needConfirm?: boolean;
  disabled?: boolean;
}) => {
  const { t } = useTranslation();

  const [open, setOpen] = useState(false);
  const handleClick = () => {
    if (!needConfirm) {
      onSave();
      return;
    }
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleConfirm = (confirm: boolean) => {
    if (confirm) {
      onSave();
    }
    setOpen(false);
  };
  return (
    <Fragment>
      <Button
        variant="contained"
        size="small"
        onClick={handleClick}
        disabled={disabled}
      >
        {t("buttons.save")}
      </Button>
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

/**
 * DeleteButton Component
 *
 * A button component that triggers a delete action. It optionally
 * displays a confirmation dialog before executing the delete action.
 *
 * @param {Object} props - The properties for the DeleteButton component.
 * @param {Function} props.onDelete - The function to call when the delete action is triggered.
 * @param {boolean} [props.disabled] - Whether the button is disabled.
 * @param {boolean} [props.hidden] - Whether the button is hidden.
 *
 * @returns {JSX.Element} A JSX element representing the delete button with optional confirmation dialog.
 */
const DeleteButton = ({
  onDelete,
  disabled,
  hidden,
}: {
  onDelete: () => void;
  disabled?: boolean;
  hidden?: boolean;
}) => {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const handleClick = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleConfirm = (confirm: boolean) => {
    if (confirm) onDelete();
    else handleClose();
  };
  if (hidden) return null;
  return (
    <Fragment>
      <Button
        variant="contained"
        size="small"
        color="error"
        onClick={handleClick}
        disabled={disabled}
      >
        {t("buttons.remove")}
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
 * Component for submitting a local node item fields.
 *
 * @param {{ id?: ExtFSNodeItem["id"] }} props
 * @prop {ExtFSNodeItem["id"]|undefined} id - id of the node item, if undefined, create a new one
 *
 * @returns {JSX.Element} a JSX element representing the form
 */
const ExtFSNodeItemForm = ({ id }: { id?: ExtFSNodeItem["id"] }) => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const { data, refetch } = useQuery({
    queryKey: ["extfs-node-item", id],
    queryFn: () => selectExtFSNodeItem(id as ExtFSNodeItem["id"]),
    enabled: !!id,
  });

  const defaultValues = useMemo(
    () => data || { name: "", filePath: "", enabled: true },
    [data]
  );

  const { handleSubmit, control, watch, reset } = useForm<ExtFSNodeItemFields>({
    defaultValues,
  });

  useEffect(() => {
    reset(defaultValues);
  }, [defaultValues.name, defaultValues.filePath, defaultValues.enabled]);

  const { mutate: saveMutate, isPending: isSavePending } = useMutation({
    mutationFn: async (fields: ExtFSNodeItemFields) =>
      await saveExtFSNodeItem(fields, id),
    onSuccess: (entity) => {
      if (entity && id === void 0) {
        navigate(generateEditPath(ExtFSNodeItemPath, entity.id));
        return;
      }
      refetch();
    },
  });

  const enabled = watch("enabled");
  const available = useMemo(() => data?.available, [data]);

  const handleSave = async () => {
    await handleSubmit(async (data) => {
      await saveMutate(data);
    })();
  };

  const handleReset = () => {
    refetch();
    reset(defaultValues);
  };

  const { mutate: deleteMutate, isPending: isDeletePending } = useMutation({
    mutationFn: deleteExtFSNodeItem,
    onSuccess: () => {
      navigate(ExtFSPath);
    },
  });
  const handleDelete = async () => {
    deleteMutate(id as ExtFSNodeItem["id"]);
  };

  const isPending = useMemo(() => {
    return isSavePending || isDeletePending;
  }, [isSavePending, isDeletePending]);

  return (
    <PageLayout>
      <PageHeader
        title={<ExtFSNodeItemEditTitle />}
        addons={<ExtFSNodeItemAddons />}
      />
      <Box component={"form"} onSubmit={handleSave}>
        <PageTopBar>
          <DeleteButton
            onDelete={handleDelete}
            disabled={isPending}
            hidden={id === void 0}
          />
          {id === void 0 ? null : (
            <Button
              variant="contained"
              size="small"
              color="warning"
              onClick={handleReset}
              disabled={isPending}
            >
              {t("buttons.reset")}
            </Button>
          )}
          <SaveButton
            onSave={handleSave}
            needConfirm={!!id}
            disabled={isPending}
          />
        </PageTopBar>
        <Stack spacing={2}>
          <Controller
            control={control}
            name="name"
            render={({ field }) => (
              <TextField
                label={t("resources.extfs/local-node-items.fields.name")}
                fullWidth
                variant="filled"
                {...field}
                disabled={isPending}
              />
            )}
          />
          <Controller
            control={control}
            name="filePath"
            render={({ field }) => (
              <FilePathInput
                title={t("resources.extfs/local-node-items.input.filePath", {})}
                label={t("resources.extfs/local-node-items.fields.filePath")}
                {...field}
              />
            )}
          />
          <FormControl component="fieldset">
            <FormLabel component="legend">
              {t("resources.extfs/local-node-items.fields.enabled")}
            </FormLabel>
            <Controller
              control={control}
              name="enabled"
              render={({ field }) => (
                <FormControlLabel
                  label={t(
                    enabled === false ? "labels.disabled" : "labels.enabled"
                  )}
                  labelPlacement="end"
                  control={
                    <Switch
                      {...field}
                      disabled={isPending}
                      checked={!!enabled}
                    />
                  }
                />
              )}
            />
          </FormControl>
          <FormControl
            component="fieldset"
            sx={{ display: id === void 0 ? "none" : "block" }}
          >
            <FormLabel component="legend">
              {t("resources.extfs/local-node-items.fields.available")}
            </FormLabel>
            <FormControlLabel
              label={t(
                available === false
                  ? "labels.not_available"
                  : "labels.available"
              )}
              labelPlacement="end"
              control={
                <Switch
                  checked={!!available}
                  slotProps={{ input: { disabled: true } }}
                />
              }
            />
          </FormControl>
        </Stack>
      </Box>
    </PageLayout>
  );
};

const ExtFSNodeItemEditTitle = () => {
  const { t } = useTranslation();

  const { id } = useParams();
  const id_ = id ? parseInt(id) : 0;
  return (
    <PageHeaderTitle>
      {t(
        id_ == 0 || isNaN(id_)
          ? ExtFSNodeItemCreateI18nKey
          : ExtFSNodeItemEditI18nKey
      )}
    </PageHeaderTitle>
  );
};

const ExtFSNodeItemAddons = () => {
  const { id } = useParams();
  const path = id === void 0 ? ExtFSNodeItemCreatePath : void 0;
  return <RouteMore path={path}></RouteMore>;
};

export const ExtFSNodeItemEditPage = () => (
  <Fragment>
    <PageI18Next ns={ExtFSNodeItemNS} defaultNS={ExtFSNodeItemNS} />
    <ExtFSNodeItemEdit />
  </Fragment>
);
/**
 * Component for editing an ExtFS node item.
 *
 * This component uses the URL parameter `id` to determine whether
 * an existing node item is being edited or a new one is being created.
 * It parses the `id` from the URL and passes it to the `ExtFSNodeItemForm`
 * component. If the `id` is not present or invalid, a new node item form
 * is rendered.
 *
 * @returns {JSX.Element} A form for editing or creating an ExtFS node item.
 */

const ExtFSNodeItemEdit = () => {
  const { id } = useParams();
  const id_ = id ? parseInt(id) : 0;
  if (id_ == 0 || isNaN(id_)) return <ExtFSNodeItemForm />;
  return <ExtFSNodeItemForm id={id_} />;
};

const ExtFSNodeItemViewTitle = () => {
  const { t } = useTranslation();
  return <PageHeaderTitle>{t(ExtFSNodeItemI18nKey)}</PageHeaderTitle>;
};

export const ExtFSNodeItemViewPage = () => (
  <Fragment>
    <PageI18Next ns={ExtFSNodeItemNS} defaultNS={ExtFSNodeItemNS} />
    <ExtFSNodeItemView />
  </Fragment>
);
/**
 * Component for viewing an ExtFS node item.
 *
 * This component renders a form displaying the `name` and `filePath` of an
 * ExtFS node item, along with a "Refresh" and "Open" button. The "Refresh"
 * button will retrieve the latest data for the node item from the API. The
 * "Open" button will open the node item in a new window.
 *
 * The component will only render if the `id` parameter is present in the URL
 * and is a valid number.
 *
 * @returns {JSX.Element} A form for viewing an ExtFS node item.
 */
const ExtFSNodeItemView = () => {
  const { t } = useTranslation();
  const { id: paramId } = useParams();
  const id = parseInt(paramId as string);

  const { data, isFetching, refetch } = useQuery({
    queryKey: ["extfs-node-item", id],
    queryFn: () => selectExtFSNodeItem(id as ExtFSNodeItem["id"]),
    enabled: !isNaN(id),
  });

  const handleRefresh = () => {
    refetch();
  };

  return (
    <PageLayout>
      <PageHeader
        title={<ExtFSNodeItemViewTitle />}
        addons={<ExtFSNodeItemAddons />}
      />
      <PageTopBar>
        <Button
          variant="contained"
          size="small"
          onClick={handleRefresh}
          sx={{ display: id === void 0 ? "none" : "block" }}
          disabled={isFetching}
        >
          {t("buttons.refresh")}
        </Button>
        <Button
          variant="contained"
          size="small"
          sx={{ display: id === void 0 ? "none" : "block" }}
          disabled={isFetching}
          color="success"
        >
          {t("buttons.open")}
        </Button>
      </PageTopBar>
      <Stack spacing={2} marginTop={2}>
        <TextField
          label={t("resources.extfs/local-node-items.fields.name")}
          fullWidth
          variant="filled"
          value={data?.name || ""}
        />
        <TextField
          label={t("resources.extfs/local-node-items.fields.filePath")}
          fullWidth
          variant="filled"
          value={data?.filePath || ""}
        />
        <CircleIcon
          fontSize="small"
          color={data?.available ? "success" : "disabled"}
        />
      </Stack>
    </PageLayout>
  );
};
