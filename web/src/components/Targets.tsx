import {
  BooleanField,
  BooleanInput,
  BulkExportButton,
  CloneButton,
  Create,
  CreateButton,
  DatagridConfigurable,
  DateField,
  DateTimeInput,
  Edit,
  EditButton,
  ExportButton,
  FilterList,
  FilterListItem,
  FilterLiveSearch,
  List,
  ListButton,
  SavedQueriesList,
  SelectColumnsButton,
  Show,
  ShowButton,
  SimpleForm,
  SimpleShowLayout,
  TextField,
  TextInput,
  TextInputProps,
  TopToolbar,
  WrapperField,
  useGetList,
  useNotify,
  useTranslate,
} from "react-admin";

import { useFormContext } from "react-hook-form";

import Block from "@mui/icons-material/Block";
import CloseIcon from "@mui/icons-material/Close";
import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";
import NavigateNextIcon from "@mui/icons-material/NavigateNext";
import RefreshIcon from "@mui/icons-material/Refresh";
import ToggleOn from "@mui/icons-material/ToggleOn";
import AppBar from "@mui/material/AppBar";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import Breadcrumbs from "@mui/material/Breadcrumbs";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CircularProgress from "@mui/material/CircularProgress";
import Dialog from "@mui/material/Dialog";
import DialogContent from "@mui/material/DialogContent";
import IconButton from "@mui/material/IconButton";
import InputAdornment from "@mui/material/InputAdornment";
import Link from "@mui/material/Link";
import ListItem from "@mui/material/ListItem";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemText from "@mui/material/ListItemText";
import Slide from "@mui/material/Slide";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import { useTheme } from "@mui/material/styles";
import { TransitionProps } from "@mui/material/transitions";
import useMediaQuery from "@mui/material/useMediaQuery";

import { Fragment, forwardRef, useMemo, useState } from "react";
import AutoSizer from "react-virtualized-auto-sizer";
import { FixedSizeList } from "react-window";

import { basename, dirname, generateParents } from "../lib/path";

const TargetBulkActions = () => (
  <>
    <BulkExportButton />
  </>
);

const TargetListActions = () => (
  <TopToolbar>
    <SelectColumnsButton preferenceKey="targets.datagrid" />
    <CreateButton />
    <ExportButton />
  </TopToolbar>
);

const TargetEnabledFilter = () => {
  return (
    <FilterList
      label="resources.extfs/targets.filters.has_enabled"
      icon={<ToggleOn />}
    >
      <FilterListItem
        label="resources.extfs/targets.filters.enabled"
        value={{ enabled: true }}
      />
      <FilterListItem
        label="resources.extfs/targets.filters.disabled"
        value={{ enabled: false }}
      />
    </FilterList>
  );
};

const TargetInvalidFilter = () => {
  return (
    <FilterList
      label="resources.extfs/targets.filters.has_available"
      icon={<Block />}
    >
      <FilterListItem
        label="resources.extfs/targets.filters.available"
        value={{ available: true }}
      />
      <FilterListItem
        label="resources.extfs/targets.filters.not_available"
        value={{ available: false }}
      />
    </FilterList>
  );
};

const TargetFilters = () => {
  return (
    <Box
      sx={{
        display: {
          xs: "none",
          sm: "block",
        },
        order: -1, // display on the left rather than on the right of the list
      }}
    >
      <Card sx={{ mr: 2, mt: 6, width: 200 }}>
        <CardContent>
          <SavedQueriesList />
          <FilterLiveSearch />
          <TargetEnabledFilter />
          <TargetInvalidFilter />
        </CardContent>
      </Card>
    </Box>
  );
};

export const Targets = () => {
  return (
    <List actions={<TargetListActions />} aside={<TargetFilters />}>
      <DatagridConfigurable
        bulkActionButtons={<TargetBulkActions />}
        preferenceKey="targets.datagrid"
      >
        <TextField source="name" />
        <TextField source="filepath" />
        <BooleanField source="enabled" />
        <BooleanField source="available" />
        <DateField source="createAt" showTime />
        <WrapperField label="others.table.actions">
          <EditButton />
          <ShowButton />
        </WrapperField>
      </DatagridConfigurable>
    </List>
  );
};

type FilePathBreadcrumbsProps = {
  value: string;
  disabled: boolean;
  onSelect: (value: string) => void;
};
const FilePathBreadcrumbs = ({
  disabled,
  value,
  onSelect,
}: FilePathBreadcrumbsProps) => {
  const [root, ...paths] = useMemo(
    () => (value == "" ? [value] : generateParents(value, true)),
    [value]
  );
  const rootOnly = useMemo(() => value === root, [value, root]);
  const dirPath = useMemo(
    () => (paths.length > 0 ? paths[paths.length - 1] : undefined),
    [paths]
  );
  const filename = useMemo(
    () => (rootOnly ? "" : basename(value, { dirname: dirPath, root })),
    [value, dirPath, rootOnly]
  );
  const onSelectPathName = (path: string) => {
    if (disabled) return;
    onSelect(path);
  };
  const elements = paths.map((path, idx) => {
    const dirPath = idx > 0 ? paths[idx - 1] : root;
    const filename = basename(path, { dirname: dirPath, root });
    return (
      <Link
        underline="hover"
        sx={{
          display: "flex",
          alignItems: "center",
          cursor: disabled ? "default" : "pointer",
        }}
        color="inherit"
        key={idx}
        onClick={() => onSelectPathName(path)}
      >
        {filename}
      </Link>
    );
  });
  return (
    <Breadcrumbs aria-label="breadcrumb" sx={{ pt: 2, pl: 2 }}>
      <Link
        key="root"
        underline="hover"
        sx={{
          display: "flex",
          alignItems: "center",
          cursor: disabled ? "default" : "pointer",
        }}
        color="inherit"
        onClick={() => onSelectPathName(root)}
      >
        <FolderIcon sx={{ mr: 0.5 }} fontSize="inherit" />
      </Link>
      {elements}
      <Typography color="text.primary" hidden={rootOnly}>
        {filename}
      </Typography>
    </Breadcrumbs>
  );
};

interface FilePathItem {
  id: string;
  name: string;
  filepath: string;
  parent: string;
  fileType: string;
  updatedAt: string;
}

type FilePathSelectorProps = {
  open: boolean;
  value: string;
  resource: string;
  onClose: () => void;
  onChange: (value: string) => void;
};

const FilePathSelectorTransition = forwardRef(
  (
    props: TransitionProps & {
      children: React.ReactElement<any, any>;
    },
    ref: React.Ref<unknown>
  ) => {
    return <Slide direction="up" ref={ref} {...props} />;
  }
);

const FilePathSelector = ({
  open,
  value,
  onClose,
  onChange,
  resource,
}: FilePathSelectorProps) => {
  const t = useTranslate();

  const theme = useTheme();
  const fullScreen = useMediaQuery(theme.breakpoints.down("sm"));

  const notify = useNotify();

  const onError = (error: Error) => {
    notify(error.message, { type: "error" });
  };

  const [parent, setParent] = useState<string>(() => {
    return value !== "" ? dirname(value) : "";
  });

  const { isFetching, data, refetch } = useGetList<FilePathItem>(
    resource,
    {
      meta: { noPagination: true },
      filter: { parent },
    },
    { onError, enabled: open }
  );

  const onRefresh = () => {
    if (!isFetching) refetch();
  };

  const onSelected = (item: FilePathItem) => {
    onChange(item.filepath);
    onClose();
  };

  const onEnter = (item: FilePathItem) => {
    setParent(item.filepath);
  };

  const rows = useMemo(() => data || [], [data]);

  return (
    <Dialog
      fullScreen={fullScreen}
      open={open}
      onClose={onClose}
      TransitionComponent={FilePathSelectorTransition}
    >
      <AppBar sx={{ position: "relative" }}>
        <Toolbar>
          <IconButton
            edge="start"
            color="inherit"
            onClick={onClose}
            aria-label="close"
          >
            <CloseIcon />
          </IconButton>
          <Typography sx={{ ml: 2, flex: 1 }} variant="h6" component="div">
            {t("others.input.filepath.title")}
          </Typography>
          <IconButton
            color="inherit"
            onClick={onRefresh}
            aria-label="refresh"
            disabled={isFetching}
          >
            <RefreshIcon />
          </IconButton>
        </Toolbar>
      </AppBar>
      <FilePathBreadcrumbs
        value={parent}
        onSelect={setParent}
        disabled={isFetching}
      ></FilePathBreadcrumbs>

      <DialogContent sx={fullScreen ? void 0 : { height: 680, width: 552 }}>
        <Box
          height="100%"
          sx={{
            display: isFetching ? "flex" : "none",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          <CircularProgress />
        </Box>
        <AutoSizer disableWidth={true} hidden={isFetching}>
          {({ height }) => (
            <FixedSizeList
              height={height}
              width="100%"
              itemSize={68}
              itemCount={rows.length}
              overscanCount={20}
            >
              {({ index, style }) => (
                <ListItem
                  style={style}
                  key={index}
                  component="div"
                  disablePadding
                  secondaryAction={
                    <IconButton
                      edge="end"
                      aria-label="enter"
                      onClick={() => onEnter(rows[index])}
                      sx={{
                        visibility:
                          rows[index].fileType === "D" ? "visible" : "hidden",
                      }}
                    >
                      <NavigateNextIcon />
                    </IconButton>
                  }
                >
                  <ListItemButton
                    selected={rows[index].filepath === value}
                    onClick={() => onSelected(rows[index])}
                  >
                    <ListItemAvatar>
                      <Avatar>
                        {rows[index].fileType === "D" ? (
                          <FolderIcon />
                        ) : (
                          <InsertDriveFileIcon />
                        )}
                      </Avatar>
                    </ListItemAvatar>
                    <ListItemText
                      primary={rows[index].name}
                      secondary={rows[index].updatedAt}
                    />
                  </ListItemButton>
                </ListItem>
              )}
            </FixedSizeList>
          )}
        </AutoSizer>
      </DialogContent>
    </Dialog>
  );
};

type FilePathInputProps = TextInputProps;
const FilePathInput = ({
  name,
  source,
  InputProps,
  ...rest
}: FilePathInputProps) => {
  const [open, setOpen] = useState(false);

  const { setValue, getValues } = useFormContext();
  const finalName = name || source;
  const onChange = (value: string) => {
    setValue(finalName, value);
  };

  const [filepath, setFilePath] = useState<string>(
    typeof rest.value === "string" ? rest.value : ""
  );
  const onOpen = () => {
    setFilePath(getValues(source) || "");
    setOpen(true);
  };

  const onClose = () => {
    setOpen(false);
  };
  const endAdornment = (
    <InputAdornment position="end">
      <IconButton aria-label="directions" onClick={onOpen}>
        <FolderIcon />
      </IconButton>
    </InputAdornment>
  );

  return (
    <Fragment>
      <TextInput
        {...rest}
        value={filepath}
        name={name}
        source={source}
        resettable={false}
        InputProps={{ ...InputProps, endAdornment }}
      />
      {open && (
        <FilePathSelector
          open={open}
          value={filepath}
          resource="extfs/disk-files"
          onChange={onChange}
          onClose={onClose}
        />
      )}
    </Fragment>
  );
};

export const TargetCreate = () => (
  <Create>
    <SimpleForm>
      <TextInput source="name" />
      <FilePathInput source="filepath" />
      <BooleanInput source="enabled" defaultValue={true} />
    </SimpleForm>
  </Create>
);

const TargetEditActions = () => (
  <TopToolbar>
    <CreateButton />
    <CloneButton />
    <ListButton />
    <ShowButton />
  </TopToolbar>
);

export const TargetEdit = () => (
  <Edit actions={<TargetEditActions />} mutationMode="pessimistic">
    <SimpleForm>
      <TextInput source="id" readOnly={true} />
      <TextInput source="name" />
      <FilePathInput source="filepath" />
      <BooleanInput source="enabled" />
      <BooleanInput source="available" disabled={true} />
      <DateTimeInput source="createAt" disabled={true} />
      <DateTimeInput source="updatedAt" disabled={true} />
    </SimpleForm>
  </Edit>
);

const TargetShowActions = () => (
  <TopToolbar>
    <EditButton />
    <ListButton />
  </TopToolbar>
);

export const TargetShow = () => (
  <Show actions={<TargetShowActions />}>
    <SimpleShowLayout>
      <TextField source="id" />
      <TextField source="name" />
      <TextField source="filepath" />
      <BooleanField source="enabled" />
      <BooleanField source="available" />
      <DateField source="createAt" showTime />
      <DateField source="updatedAt" showTime />
    </SimpleShowLayout>
  </Show>
);
