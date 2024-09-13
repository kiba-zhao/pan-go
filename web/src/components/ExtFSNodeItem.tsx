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
}: {
  onSave: () => void;
  needConfirm?: boolean;
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
        color="error"
        onClick={handleClick}
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

const ExtFSNodeItemForm = ({ id }: { id?: ExtFSNodeItem["id"] }) => {
  const t = useTranslate();
  const navigate = useNavigate();
  const api = useAPI();
  const { data, refetch } = useQuery({
    queryKey: ["extfs-node-item", id],
    queryFn: () => api?.getExtFSNodeItem(id as ExtFSNodeItem["id"]),
    enabled: !!id,
  });

  const defaultValues = useMemo(
    () => data || { name: "", filepath: "", enabled: true },
    [data]
  );

  const { handleSubmit, control, watch, reset } = useForm<ExtFSNodeItemFields>({
    defaultValues,
  });

  useEffect(() => {
    reset(defaultValues);
  }, [defaultValues]);

  const { mutate, isPending } = useMutation({
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
      await mutate(data);
    })();
  };

  return (
    <Paper component={"form"} sx={{ padding: 3 }} onSubmit={handleSave}>
      <Stack
        direction="row"
        spacing={1}
        alignItems={"flex-start"}
        justifyContent={"space-between"}
      >
        <Typography variant="h6">
          {t(`custom.extfs/local-node-items.${id ? "editName" : "createName"}`)}
        </Typography>
        <SaveButton onSave={handleSave} needConfirm={!!id} />
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
          name="filepath"
          render={({ field }) => (
            <FilePathInput
              title={t("custom.extfs/local-node-items.input.filepath", {})}
              label={t("custom.extfs/local-node-items.fields.filepath")}
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
                  <Switch {...field} disabled={isPending} checked={enabled} />
                }
              />
            )}
          />
        </FormControl>
        <FormControl component="fieldset">
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
              <Switch checked={available} inputProps={{ disabled: true }} />
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
