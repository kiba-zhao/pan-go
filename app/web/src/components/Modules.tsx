import { ChangeEvent, Fragment, useEffect, useReducer, useState } from "react";

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

import {
  ModuleItem as ModuleItemType,
  ModuleSearchResult,
  useAPI,
} from "../api.tsx";

type ModuleResultState = ModuleSearchResult & { Version: number };
const initialModuleResultState: ModuleResultState = {
  Total: 0,
  Items: new Array<ModuleItemType>(),
  Version: 0,
};

type ModuleResultSetAction = { type: "set"; value: ModuleResultState };
type ModuleItemSetAction = {
  type: "setItem";
  item: ModuleItemType;
  version: number;
};

function modulesReducer(
  state: ModuleResultState,
  action: ModuleResultSetAction | ModuleItemSetAction
): ModuleResultState {
  switch (action.type) {
    case "set":
      return { ...action.value, Items: action.value.Items || [] };
    case "setItem":
      if (action.version != state.Version) return state;
      const idx = state.Items.findIndex(
        (item) => item.Name == action.item.Name
      );
      state.Items[idx] = action.item;
      return state;
    default:
      throw new Error("Unknown action");
  }
}

function Modules() {
  const { t } = useTranslation();

  const [moduleResults, dispatch] = useReducer(
    modulesReducer,
    initialModuleResultState
  );
  const [loading, setLoading] = useState(false);

  const [showFilter, setShowFilter] = useState(false);
  const onFilterClose = () => setShowFilter(false);
  const onFilterOpen = () => !loading && setShowFilter(true);

  const [keyword, setKeyword] = useState("");
  const [submitVersion, setSubmitVersion] = useState(0);
  const onFilterChange = async (value: string) => {
    if (loading) return;
    setKeyword(value);
    setSubmitVersion(submitVersion + 1);
  };

  const api = useAPI();
  const onSearch = async (keyword: string, version: number) => {
    setLoading(true);
    const results = await api.SearchModules(keyword);
    if (version != submitVersion) return;
    dispatch({ type: "set", value: { ...results, Version: version } });
    setLoading(false);
  };

  useEffect(() => {
    onSearch(keyword, submitVersion);
  }, [submitVersion]);

  return (
    <Fragment>
      {loading ? <LinearProgress /> : <></>}
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
          {moduleResults.Total > 0
            ? t("Modules.Total", { Total: moduleResults.Total })
            : ""}
        </Divider>
        <Grid container spacing={2} paddingTop={2}>
          {moduleResults.Items.map((v: ModuleItemType) => (
            <Grid item key={v.Name} xs={12} sm={6} md={4} lg={3}>
              <ModuleItem
                module={v}
                disabled={loading}
                scopeVersion={submitVersion}
              ></ModuleItem>
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
  scopeVersion: number;
}

function ModuleItem({ module, disabled, scopeVersion }: ModuleItemProps) {
  const { t } = useTranslation();

  const [loading, setLoading] = useState(false);
  const [enabledVersion, setEnabledVersion] = useState(0);
  const [enabled, setEnabled] = useState(module.Enabled);

  const onEnableChange = async (value: boolean) => {
    setEnabled(value);
    setEnabledVersion(enabledVersion + 1);
  };

  const [_, dispatch] = useReducer(modulesReducer, initialModuleResultState);

  const api = useAPI();
  const onEnableSync = async (enabled: boolean, version: number) => {
    if (version != enabledVersion) return;
    setLoading(true);
    try {
      const item = await api.SetModuleEnabled(module.Name, enabled);
      dispatch({ type: "setItem", item, version: scopeVersion });
    } catch {
      setEnabled(!enabled);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (disabled || module.ReadOnly || enabledVersion == 0) {
      return;
    }
    onEnableSync(enabled, enabledVersion);
  }, [enabledVersion]);

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
            disabled={disabled || module.ReadOnly || loading}
          ></Switch>
        </Box>
      </CardActions>
      {loading ? <LinearProgress /> : <></>}
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
