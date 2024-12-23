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
} from "react";

import type { QueryKey } from "@tanstack/react-query";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";

import type { ExtFSSearchItem, ExtFSSearchItemFields } from "../../API";
import { useAPI } from "../../API";
import type { ListItemData } from "../List/Item";
import { ListItems, useListItems } from "../List/Item";

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
  const theme = useTheme();
  const fullScreen = useMediaQuery(theme.breakpoints.down("sm"));

  const [query, setQuery] = useState("");
  const queryKey = useMemo(() => ["extfs-search-items", query], [query]);

  const api = useAPI();
  const { data, isFetching } = useQuery({
    queryKey,
    queryFn: async () =>
      await api?.searchExtFSSearchItems({ q: query, limit: 10 }),
    enabled,
  });

  const setQueryDelay = useCallback(
    _.debounce(setQuery, 500, { maxWait: 1000 }),
    []
  );

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setQueryDelay(value);
  };

  const { handleSubmit, register } = useForm<ExtFSSearchItemFields>();

  const { mutate: saveMutate, isPending: isSavePending } = useMutation({
    mutationFn: api?.saveExtFSSearchItem,
    onSuccess: () => {
      // TODO: redirect to search file mode
    },
  });

  const handleSave = async (event: React.FormEvent) => {
    event.preventDefault();
    await handleSubmit(async (data) => {
      await saveMutate(data);
    })();
    event.stopPropagation();
  };

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
          disabled={isSavePending}
          {...register("query", { onChange: handleChange })}
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
      <LinearProgress sx={{ visibility: isFetching ? "visible" : "hidden" }} />
      <div style={{ height: fullScreen ? "100%" : 680 }}>
        <SearchContext.Provider value={{ queryKey }}>
          <ListItems items={data || []} isFetching={isFetching} itemSize={68}>
            <SearchItem />
          </ListItems>
        </SearchContext.Provider>
      </div>
    </Fragment>
  );
};

export const SearchItem = () => {
  const { style, item }: ListItemData<ExtFSSearchItem> = useListItems();

  const handleClick = () => {
    // TODO: redirect to search file mode
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
      queryClient.prefetchQuery({ queryKey });
    },
  });

  const handleDelete = async () => {
    mutate(item.id);
  };

  return (
    <IconButton size="small" disabled={isPending} onClick={handleDelete}>
      <ClearIcon fontSize="small" />
    </IconButton>
  );
};
