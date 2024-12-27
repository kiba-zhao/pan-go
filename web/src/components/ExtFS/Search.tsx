import { useTheme } from "@mui/material/styles";
import useMediaQuery from "@mui/material/useMediaQuery";
import { useTranslate } from "react-admin";

import ClearIcon from "@mui/icons-material/Clear";
import SearchIcon from "@mui/icons-material/Search";
import ButtonBase from "@mui/material/ButtonBase";
import Chip from "@mui/material/Chip";
import Divider from "@mui/material/Divider";
import IconButton from "@mui/material/IconButton";
import InputBase from "@mui/material/InputBase";
import LinearProgress from "@mui/material/LinearProgress";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemSecondaryAction from "@mui/material/ListItemSecondaryAction";
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";

import _ from "lodash";
import {
  createContext,
  Fragment,
  useCallback,
  useContext,
  useMemo,
  useState,
  memo,
  useDeferredValue,
} from "react";

import type { QueryKey } from "@tanstack/react-query";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import type { ExtFSSearchItem } from "../../API";
import { useAPI } from "../../API";
import type { ListItemData } from "../List/Item";
import { ListItems, useListItems } from "../List/Item";
import { useExtFS } from "./State";
import { newExtFSState } from "./SearchFile";

type SearchQuery = {
  queryKey: QueryKey;
};

const SearchContext = createContext<SearchQuery>(null!);

type SearchItemsProps = {
  onEsc: () => void;
  enabled: boolean;
};
export const SearchItems = ({ onEsc, enabled }: SearchItemsProps) => {
  const t = useTranslate();

  const [query, setQuery] = useState("");

  const setQueryDelay = useCallback(
    _.debounce(setQuery, 500, { maxWait: 1000 }),
    []
  );

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    if (e === void 0 || e.target === void 0) return;
    const value = e.target.value;
    setQueryDelay(value);
  };

  const [extfs, setExtFS] = useExtFS();

  const handleSave = async (event: React.FormEvent) => {
    event.preventDefault();
    const state = newExtFSState(extfs, query);
    setExtFS(state);
    event.stopPropagation();
  };

  const deferredQuery = useDeferredValue(query);

  return (
    <Fragment>
      <Stack
        padding={1}
        component="form"
        spacing={0.5}
        direction="row"
        alignItems="center"
        justifyContent="space-between"
        onSubmit={handleSave}
      >
        <SearchIcon />
        <InputBase
          placeholder={t("custom.placeholder.search-input")}
          fullWidth
          size="medium"
          autoFocus
          onChange={handleChange}
        />
        <ButtonBase onClick={onEsc}>
          <Chip
            label="esc"
            variant="outlined"
            sx={{ borderRadius: 1 }}
            size="small"
          />
        </ButtonBase>
      </Stack>
      <Divider />
      <SearchItemsResults
        query={deferredQuery}
        enabled={enabled}
        onEsc={onEsc}
      />
    </Fragment>
  );
};

type SearchItemsResultsProps = {
  onEsc: () => void;
  enabled: boolean;
  query: string;
};
const SearchItemsResults = memo(
  ({ query, enabled, onEsc }: SearchItemsResultsProps) => {
    const theme = useTheme();
    const fullScreen = useMediaQuery(theme.breakpoints.down("sm"));

    const queryKey = useMemo(() => ["extfs-search-items", query], [query]);

    const api = useAPI();
    const { data, isFetching } = useQuery({
      queryKey,
      queryFn: async () =>
        await api?.searchExtFSSearchItems({ q: query, limit: 10 }),
      enabled,
    });
    return (
      <Fragment>
        <LinearProgress
          sx={{ visibility: isFetching ? "visible" : "hidden" }}
        />
        <div style={{ height: fullScreen ? "100%" : 680 }}>
          <SearchContext.Provider value={{ queryKey }}>
            <ListItems items={data || []} isFetching={isFetching} itemSize={68}>
              <SearchItem onClick={onEsc} />
            </ListItems>
          </SearchContext.Provider>
        </div>
      </Fragment>
    );
  }
);

type SearchItemProps = {
  onClick: () => void;
};
export const SearchItem = ({ onClick }: SearchItemProps) => {
  const { style, item }: ListItemData<ExtFSSearchItem> = useListItems();

  const [extfs, setExtFS] = useExtFS();
  const handleClick = () => {
    const state = newExtFSState(extfs, item.query);
    setExtFS(state);
    onClick();
  };

  return (
    <ListItem style={style} disableGutters>
      <ListItemButton onClick={handleClick}>
        <ListItemText
          primary={item.query}
          sx={{ paddingRight: 5 }}
        ></ListItemText>
        <ListItemSecondaryAction>
          <SearchItemRemoveAction />
        </ListItemSecondaryAction>
      </ListItemButton>
    </ListItem>
  );
};

const SearchItemRemoveAction = () => {
  const { item }: ListItemData<ExtFSSearchItem> = useListItems();

  const { queryKey } = useContext(SearchContext);
  const queryClient = useQueryClient();

  const api = useAPI();
  const { mutate, isPending } = useMutation({
    mutationFn: api?.deleteExtFSSearchItem,
    onSuccess: () => {
      queryClient.refetchQueries({ queryKey });
    },
  });

  const handleDelete = async (event: React.MouseEvent) => {
    event.preventDefault();
    mutate(item.id);
    event.stopPropagation();
  };

  return (
    <IconButton size="small" disabled={isPending} onClick={handleDelete}>
      <ClearIcon fontSize="small" />
    </IconButton>
  );
};
