import {
  BooleanField,
  BooleanInput,
  BulkDeleteButton,
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
import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";
import NavigateNextIcon from "@mui/icons-material/NavigateNext";
import ToggleOn from "@mui/icons-material/ToggleOn";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Dialog from "@mui/material/Dialog";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import IconButton from "@mui/material/IconButton";
import InputAdornment from "@mui/material/InputAdornment";
import ListItem from "@mui/material/ListItem";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemText from "@mui/material/ListItemText";
import Slide from "@mui/material/Slide";
import { useTheme } from "@mui/material/styles";
import { TransitionProps } from "@mui/material/transitions";
import useMediaQuery from "@mui/material/useMediaQuery";

import { Fragment, forwardRef, useMemo, useState } from "react";
import AutoSizer from "react-virtualized-auto-sizer";
import { FixedSizeList } from "react-window";

const TargetBulkActions = () => (
  <>
    <BulkDeleteButton mutationMode="pessimistic" />
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
  const t = useTranslate();
  return (
    <FilterList
      label={t("resources.targets.filters.has_enabled")}
      icon={<ToggleOn />}
    >
      <FilterListItem
        label={t("resources.targets.filters.enabled")}
        value={{ enabled: true }}
      />
      <FilterListItem
        label={t("resources.targets.filters.disabled")}
        value={{ enabled: false }}
      />
    </FilterList>
  );
};

const TargetInvalidFilter = () => {
  const t = useTranslate();
  return (
    <FilterList
      label={t("resources.targets.filters.has_invalid")}
      icon={<Block />}
    >
      <FilterListItem
        label={t("resources.targets.filters.invalid")}
        value={{ invalid: true }}
      />
      <FilterListItem
        label={t("resources.targets.filters.valid")}
        value={{ invalid: false }}
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
  const t = useTranslate();
  return (
    <List actions={<TargetListActions />} aside={<TargetFilters />}>
      <DatagridConfigurable
        bulkActionButtons={<TargetBulkActions />}
        preferenceKey="targets.datagrid"
      >
        <TextField source="name" />
        <TextField source="filepath" />
        <BooleanField source="enabled" />
        <BooleanField source="invalid" />
        <DateField source="createAt" showTime />
        <WrapperField label={t("others.table.actions")}>
          <EditButton />
          <ShowButton />
        </WrapperField>
      </DatagridConfigurable>
    </List>
  );
};

interface FilePathItem {
  id: string;
  filepath: string;
  parent: string;
  fileType: string;
  mimeType?: string;
}

type FilePathSelectorProps = {
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
  value,
  onClose,
  onChange,
  resource,
}: FilePathSelectorProps) => {
  const theme = useTheme();
  const fullScreen = useMediaQuery(theme.breakpoints.down("sm"));

  const notify = useNotify();

  const onError = (error: Error) => {
    notify(error.message, { type: "error" });
  };

  const [parent, setParent] = useState<string>(() => {
    if (!value || value.trim() === "") return "/";
    const parentValues = value.split("/").slice(0, -1);
    return parentValues.length <= 1 ? "/" : parentValues.join("/");
  });
  const { isFetching, data, refetch } = useGetList<FilePathItem>(
    resource,
    {
      meta: { noPagination: true },
      filter: { parent },
    },
    { onError }
  );

  const rows = useMemo(() => {
    if (data === void 0) {
      return [];
    }
    return data.map((item) => ({
      name:
        item.filepath === "/" ? item.filepath : item.filepath.split("/").pop(),
      ...item,
    }));
  }, [data]);

  const onSelected = (item: FilePathItem) => {
    onChange(item.filepath);
    onClose();
  };

  const onEnter = (item: FilePathItem) => {
    setParent(item.filepath);
  };
  return (
    <Dialog
      fullScreen={fullScreen}
      open={true}
      onClose={onClose}
      TransitionComponent={FilePathSelectorTransition}
    >
      <DialogTitle>Disk Path Selector</DialogTitle>
      <DialogContent sx={fullScreen ? void 0 : { height: 680, width: 552 }}>
        <AutoSizer disableWidth={true}>
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
                    <ListItemText primary={rows[index].name} />
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
    setFilePath(getValues(source));
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
          value={filepath}
          resource="disk-files"
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
    <ShowButton />
    <ListButton />
  </TopToolbar>
);

export const TargetEdit = () => (
  <Edit actions={<TargetEditActions />} mutationMode="pessimistic">
    <SimpleForm>
      <TextInput source="id" readOnly={true} />
      <TextInput source="name" />
      <FilePathInput source="filepath" />
      <BooleanInput source="enabled" />
      <BooleanInput source="invalid" disabled={true} />
      <DateTimeInput source="createAt" disabled={true} />
      <DateTimeInput source="updateAt" disabled={true} />
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
      <BooleanField source="invalid" />
      <DateField source="createAt" showTime />
      <DateField source="updateAt" showTime />
    </SimpleShowLayout>
  </Show>
);
