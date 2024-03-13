import { ChangeEvent, Fragment, useMemo, useState } from "react";

import Autocomplete from "@mui/material/Autocomplete";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardContent from "@mui/material/CardContent";
import Divider from "@mui/material/Divider";
import Drawer from "@mui/material/Drawer";
import Grid from "@mui/material/Grid";
import LinearProgress from "@mui/material/LinearProgress";
import Skeleton from "@mui/material/Skeleton";
import Stack from "@mui/material/Stack";
import Switch from "@mui/material/Switch";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";

import { useTranslation } from "react-i18next";

import { useMutation, useQuery } from "@tanstack/react-query";
import { ModuleItem as ModuleItemType, useAPI } from "../api.tsx";

const MODULES_QUERY_KEY = "modules";

function Modules() {
  const { t } = useTranslation();

  const [showFilter, setShowFilter] = useState(false);
  const onFilterClose = () => setShowFilter(false);

  const [keyword, setKeyword] = useState("");

  const api = useAPI();

  const {
    isFetching,
    data: results,
    refetch,
  } = useQuery({
    queryKey: [MODULES_QUERY_KEY, keyword],
    queryFn: () => api.SearchModules(keyword),
  });
  const [total, items] = results || [0, []];

  const onFilterOpen = () => !isFetching && setShowFilter(true);

  const onFilterChange = async (value: string) => {
    if (isFetching) return;
    if (value == keyword) return await refetch();
    setKeyword(value);
  };

  return (
    <Fragment>
      {isFetching ? <LinearProgress /> : <></>}
      <Box padding={2}>
        <ModuleFilter
          value={keyword}
          show={showFilter}
          onClose={onFilterClose}
          onChange={onFilterChange}
        ></ModuleFilter>
        <Typography variant="h6">
          {t("Header.Modules")} &nbsp;
          <Typography
            component={"small"}
            color={"primary"}
            sx={{ cursor: "pointer" }}
            onClick={onFilterOpen}
            tabIndex={0}
          >
            {t("Modules.Settings")}
          </Typography>
        </Typography>
        <Divider>
          {total > 0 ? t("Modules.Total", { Total: total }) : ""}
        </Divider>
        <Grid container spacing={2} paddingTop={2}>
          {items.map((v: ModuleItemType) => (
            <Grid item key={v.Name} xs={12} sm={6} md={4} lg={3}>
              <ModuleItem module={v} disabled={isFetching}></ModuleItem>
            </Grid>
          ))}
        </Grid>
      </Box>
    </Fragment>
  );
}

interface ModuleItemProps {
  module: ModuleItemType;
  disabled: boolean;
}

function ModuleItem({ module, disabled }: ModuleItemProps) {
  const { t } = useTranslation();

  const [enabled, setEnabled] = useState(module.Enabled);

  const api = useAPI();
  const { isPending, mutate } = useMutation({
    mutationFn: (value: boolean) => api.SetModuleEnabled(module.Name, value),
    onMutate: (value) => {
      const ctx = { enabled };
      setEnabled(value);
      return ctx;
    },
    onError: (_, v, ctx) => ctx && setEnabled(ctx.enabled),
  });

  const onEnableChange = mutate;

  const disabledSwitch = useMemo(() => {
    return disabled || module.ReadOnly || isPending;
  }, [disabled, module.ReadOnly, isPending]);

  return (
    <Card raised={true}>
      <CardContent>
        <Stack spacing={2} direction="row" alignItems="center">
          {disabled ? (
            <Skeleton variant="circular">
              <ModuleItemAvatar url={module.Avatar} name={module.Name} />
            </Skeleton>
          ) : (
            <ModuleItemAvatar url={module.Avatar} name={module.Name} />
          )}
          {disabled ? (
            <Skeleton variant="text" width="100%">
              <Typography padding={1} noWrap>
                {module.Desc}
              </Typography>
            </Skeleton>
          ) : (
            <Typography padding={1} noWrap>
              {module.Desc}
            </Typography>
          )}
        </Stack>
      </CardContent>
      <CardActions>
        {module.HasWeb ? (
          <Button
            variant="text"
            size="small"
            href={`/${module.Name}`}
            disabled={disabled}
          >
            {t("Modules.Details")}
          </Button>
        ) : (
          <></>
        )}
        <Button variant="text" size="small" disabled={disabled}>
          {t("Modules.Remove")}
        </Button>
        <Box
          width={"100%"}
          sx={{ display: "flex", justifyContent: "flex-end" }}
        >
          <Switch
            checked={enabled}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              onEnableChange(event.target.checked)
            }
            disabled={disabledSwitch}
          ></Switch>
        </Box>
      </CardActions>
      {isPending ? <LinearProgress /> : <></>}
    </Card>
  );
}

interface ModuleItemAvatarProps {
  url: string;
  name: string;
}
function ModuleItemAvatar({ url, name }: ModuleItemAvatarProps) {
  return url.trim().length > 0 ? (
    <Avatar alt={name} src={url} />
  ) : (
    <Avatar>{name[0].toUpperCase()}</Avatar>
  );
}

interface ModuleFilterProps {
  value: string | undefined;
  show: boolean;
  onChange: (keyword: string) => void;
  onClose: () => void;
}

function ModuleFilter({ show, value, onClose, onChange }: ModuleFilterProps) {
  const { t } = useTranslation();
  const searchOpts = [{ label: "app" }];

  const [keyword, setKeyword] = useState(value);
  const onKeywordChange = (_: React.SyntheticEvent, newValue: string) =>
    setKeyword(newValue);
  const onCloseAndReset = () => {
    onClose();
    if (value != keyword) setKeyword(value);
  };
  const onSave = () => {
    onChange(keyword || "");
    onClose();
  };

  return (
    <Drawer anchor="right" open={show} onClose={onClose}>
      <Box component="form" padding={1} noValidate autoComplete="off">
        <Typography noWrap>{t("Modules.Filter")}</Typography>
        <Autocomplete
          id="keyword"
          freeSolo
          inputValue={keyword}
          onInputChange={onKeywordChange}
          options={searchOpts}
          renderInput={(params) => (
            <TextField
              {...params}
              variant="standard"
              label={t("Modules.Search")}
            />
          )}
        ></Autocomplete>
        <Stack direction="row" spacing={1} padding={3} width={"100%"}>
          <Button variant="outlined" onClick={onCloseAndReset}>
            {t("button.Close")}
          </Button>
          <Button variant="contained" onClick={onSave}>
            {t("button.Submit")}
          </Button>
        </Stack>
      </Box>
    </Drawer>
  );
}

export default Modules;
