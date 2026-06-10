import { useTheme } from "@mui/material/styles";
import useMediaQuery from "@mui/material/useMediaQuery";

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
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";
import Alert from "@mui/material/Alert";

import _ from "lodash";
import {
  createContext,
  Fragment,
  memo,
  useCallback,
  useContext,
  useDeferredValue,
  useMemo,
  useState,
} from "react";

import type { QueryKey } from "@tanstack/react-query";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Controller, useForm } from "react-hook-form";

import { useTranslation } from "../I18Next/Context";
import type { ExtFSSearchItem, ExtFSSearchItemFields } from "./api";
import {
  saveExtFSSearchItem,
  searchExtFSSearchItems,
  deleteExtFSSearchItem,
} from "./api";
import type { ListItemData } from "../Common/Item";
import { ListItems, useListItems } from "../Common/Item";
import { newExtFSState } from "./SearchFile";
import { useExtFS } from "./State";

type SearchQuery = {
  queryKey: QueryKey;
};

const SearchContext = createContext<SearchQuery>(null!);

type SearchItemsProps = {
  onEsc: () => void;
  enabled: boolean;
};
/**
 * SearchItems Component
 *
 * This component provides a search input field that allows users to perform
 * searches within the application. It utilizes debouncing to optimize the
 * searching process, updating the query state with a delay after user input.
 * The component also handles form submission for saving search queries and
 * provides a button to clear the search input.
 *
 * @param {Function} onEsc - Callback function to execute when the escape button is clicked.
 * @param {boolean} enabled - Determines whether the search results should be displayed.
 *
 * @returns {JSX.Element} A JSX element containing the search input and related actions.
 */

export const SearchItems = ({ onEsc, enabled }: SearchItemsProps) => {
  const { t } = useTranslation();

  const [query, setQuery] = useState("");

  const setQueryDelay = useCallback(
    _.debounce(setQuery, 500, { maxWait: 1000 }),
    [],
  );

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    if (e === void 0 || e.target === void 0) return;
    const value = e.target.value;
    setQueryDelay(value);
  };

  const { handleSubmit, control } = useForm<ExtFSSearchItemFields>({
    defaultValues: { query },
  });

  const [extfs, setExtFS] = useExtFS();

  const { mutate: saveMutate, isPending: isSavePending } = useMutation({
    mutationFn: async (data: ExtFSSearchItemFields) => {
      return await saveExtFSSearchItem(data);
    },

    onSuccess: (data) => {
      onEsc();
      const state = newExtFSState(extfs, { query: data.query });
      setExtFS(state);
    },
  });

  const handleSave = async (event: React.FormEvent) => {
    event.preventDefault();
    await handleSubmit(async (data) => {
      await saveMutate(data);
    })();
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
        <Controller
          control={control}
          rules={{
            required: true,
          }}
          name="query"
          render={({ field: { onChange, ...field_ } }) => (
            <InputBase
              placeholder={t("placeholders.search-input")}
              fullWidth
              size="medium"
              autoFocus
              {...field_}
              onChange={(event) => {
                onChange(event);
                handleChange(event);
              }}
              disabled={isSavePending}
            />
          )}
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

/**
 * SearchItemsResults Component
 *
 * List of items that match the search query.
 */
const SearchItemsResults = memo(
  ({ query, enabled, onEsc }: SearchItemsResultsProps) => {
    const { t } = useTranslation();
    const theme = useTheme();
    const fullScreen = useMediaQuery(theme.breakpoints.down("sm"));

    const queryKey = useMemo(() => ["extfs-search-items", query], [query]);

    const { data, isFetching, error } = useQuery({
      queryKey,
      queryFn: async () =>
        await searchExtFSSearchItems({ q: query, limit: 10 }),
      enabled,
    });

    return (
      <Fragment>
        <LinearProgress
          sx={{ visibility: isFetching ? "visible" : "hidden" }}
        />
        {!isFetching && error ? (
          <Alert severity="error">
            {t(`errors.${error.name}`, { _: error.message })}
          </Alert>
        ) : (
          void 0
        )}
        <div style={{ height: fullScreen ? "100%" : 680 }}>
          <SearchContext.Provider value={{ queryKey }}>
            <ListItems items={data || []} isFetching={isFetching} itemSize={68}>
              <SearchItem onClick={onEsc} />
            </ListItems>
          </SearchContext.Provider>
        </div>
      </Fragment>
    );
  },
);

type SearchItemProps = {
  onClick: () => void;
};
/**
 * A component that renders a single search item in a list.
 *
 * @param {SearchItemProps} props - The props for the component.
 * @param {() => void} props.onClick - The function to be called when the list item is clicked.
 *
 * @returns {ReactElement} A JSX element representing the list item.
 */
export const SearchItem = ({ onClick }: SearchItemProps) => {
  const { style, item }: ListItemData<ExtFSSearchItem> = useListItems();

  const [extfs, setExtFS] = useExtFS();
  const handleClick = () => {
    const state = newExtFSState(extfs, { query: item.query });
    setExtFS(state);
    onClick();
  };

  return (
    <ListItem
      style={style}
      disableGutters
      secondaryAction={<SearchItemRemoveAction />}
    >
      <ListItemButton onClick={handleClick}>
        <ListItemText
          primary={item.query}
          sx={{ paddingRight: 5 }}
        ></ListItemText>
      </ListItemButton>
    </ListItem>
  );
};

/**
 * A component that renders a button to remove a search item.
 *
 * @returns {ReactElement} A JSX element representing the button.
 *
 * @example
 * <SearchItemRemoveAction />
 */
const SearchItemRemoveAction = () => {
  const { item }: ListItemData<ExtFSSearchItem> = useListItems();

  const { queryKey } = useContext(SearchContext);
  const queryClient = useQueryClient();

  const { mutate, isPending } = useMutation({
    mutationFn: deleteExtFSSearchItem,
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
