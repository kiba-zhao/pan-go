import { Fragment, useEffect, useMemo, useState } from "react";
import { useTranslate } from "react-admin";

import Button from "@mui/material/Button";
import DialogActions from "@mui/material/DialogActions";
import FormControl from "@mui/material/FormControl";
import FormControlLabel from "@mui/material/FormControlLabel";
import FormLabel from "@mui/material/FormLabel";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Switch from "@mui/material/Switch";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";

import { useMutation, useQuery } from "@tanstack/react-query";
import { Controller, useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router-dom";
import type { ExtFSNodeItem, ExtFSNodeItemFields } from "../API";
import { useAPI } from "../API";
import {
  Dialog,
  DialogConfirmActions,
  DialogConfirmContent,
} from "./Feedback/Dialog";
import { FilePathInput } from "./FilePath/Input";

export const ExtFSNodeItemRoutePath = "/extfs/local-node-items";
export const ExtFSNodeItemCreate = () => <ExtFSNodeItemForm />;

const SaveButton = ({
  onSave,
  needConfirm,
  disabled,
}: {
  onSave: () => void;
  needConfirm?: boolean;
  disabled?: boolean;
}) => {
  const t = useTranslate();

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
        {t("custom.button.save")}
      </Button>
      <Dialog open={open} onClose={handleClose}>
        <DialogConfirmContent label={t("custom.button.save")} />
        <DialogActions>
          <DialogConfirmActions
            label={t("custom.button.save")}
            onConfirm={handleConfirm}
          />
        </DialogActions>
      </Dialog>
    </Fragment>
  );
};

const DeleteButton = ({
  onDelete,
  disabled,
  hidden,
}: {
  onDelete: () => void;
  disabled?: boolean;
  hidden?: boolean;
}) => {
  const t = useTranslate();
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
  return (
    <Fragment>
      <Button
        variant="contained"
        size="small"
        color="error"
        onClick={handleClick}
        sx={{ display: hidden ? "none" : "block" }}
        disabled={disabled}
      >
        {t("custom.button.remove")}
      </Button>
      <Dialog open={open} onClose={handleClose}>
        <DialogConfirmContent label={t("custom.button.remove")} />
        <DialogActions>
          <DialogConfirmActions
            label={t("custom.button.remove")}
            onConfirm={handleConfirm}
          />
        </DialogActions>
      </Dialog>
    </Fragment>
  );
};

const ExtFSNodeItemForm = ({ id }: { id?: ExtFSNodeItem["id"] }) => {
  const t = useTranslate();
  const navigate = useNavigate();
  const api = useAPI();
  const { data, refetch } = useQuery({
    queryKey: ["extfs-node-item", id],
    queryFn: () => api?.selectExtFSNodeItem(id as ExtFSNodeItem["id"]),
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
  }, [defaultValues]);

  const { mutate: saveMutate, isPending: isSavePending } = useMutation({
    mutationFn: async (fields: ExtFSNodeItemFields) =>
      await api?.saveExtFSNodeItem(fields, id),
    onSuccess: (entity) => {
      if (entity && id === void 0) {
        navigate(`${ExtFSNodeItemRoutePath}/${entity.id}`);
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
    mutationFn: api?.deleteExtFSNodeItem,
    onSuccess: () => {
      navigate("/extfs");
    },
  });
  const handleDelete = async () => {
    deleteMutate(id as ExtFSNodeItem["id"]);
  };

  const isPending = useMemo(() => {
    return isSavePending || isDeletePending;
  }, [isSavePending, isDeletePending]);

  return (
    <Paper component={"form"} sx={{ padding: 3 }} onSubmit={handleSave}>
      <Stack
        direction="row"
        spacing={1}
        alignItems={"flex-start"}
        justifyContent={"flex-end"}
      >
        <Typography variant="h6" sx={{ flexGrow: 1 }}>
          {t(`custom.extfs/local-node-items.${id ? "editName" : "createName"}`)}
        </Typography>
        <Button
          variant="contained"
          size="small"
          onClick={handleReset}
          sx={{ display: id === void 0 ? "none" : "block" }}
          disabled={isPending}
        >
          {t("custom.button.reset")}
        </Button>
        <SaveButton
          onSave={handleSave}
          needConfirm={!!id}
          disabled={isPending}
        />
        <DeleteButton
          onDelete={handleDelete}
          disabled={isPending}
          hidden={id === void 0}
        />
      </Stack>
      <Stack spacing={2} marginTop={2}>
        <Controller
          control={control}
          name="name"
          render={({ field }) => (
            <TextField
              label={t("custom.extfs/local-node-items.fields.name")}
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
              title={t("custom.extfs/local-node-items.input.filePath", {})}
              label={t("custom.extfs/local-node-items.fields.filePath")}
              {...field}
            />
          )}
        />
        <FormControl component="fieldset">
          <FormLabel component="legend">
            {t("custom.extfs/local-node-items.fields.enabled")}
          </FormLabel>
          <Controller
            control={control}
            name="enabled"
            render={({ field }) => (
              <FormControlLabel
                label={t(
                  enabled === false
                    ? "custom.label.disabled"
                    : "custom.label.enabled"
                )}
                labelPlacement="end"
                control={
                  <Switch {...field} disabled={isPending} checked={!!enabled} />
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
            {t("custom.extfs/local-node-items.fields.available")}
          </FormLabel>
          <FormControlLabel
            label={t(
              available === false
                ? "custom.label.not_available"
                : "custom.label.available"
            )}
            labelPlacement="end"
            control={
              <Switch checked={!!available} inputProps={{ disabled: true }} />
            }
          />
        </FormControl>
      </Stack>
    </Paper>
  );
};

export const ExtFSNodeItemEdit = () => {
  const { id } = useParams();
  const id_ = id ? parseInt(id) : void 0;
  return <ExtFSNodeItemForm id={id_ === void 0 || isNaN(id_) ? void 0 : id_} />;
};
