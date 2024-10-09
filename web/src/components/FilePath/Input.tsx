import { useQuery } from "@tanstack/react-query";
import type { ReactElement, Ref } from "react";
import {
  forwardRef,
  Fragment,
  useEffect,
  useId,
  useMemo,
  useState,
} from "react";

import CloseIcon from "@mui/icons-material/Close";
import FolderIcon from "@mui/icons-material/Folder";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";
import NavigateNextIcon from "@mui/icons-material/NavigateNext";
import RefreshIcon from "@mui/icons-material/Refresh";
import AppBar from "@mui/material/AppBar";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import Breadcrumbs from "@mui/material/Breadcrumbs";
import CircularProgress from "@mui/material/CircularProgress";
import Dialog from "@mui/material/Dialog";
import DialogContent from "@mui/material/DialogContent";
import IconButton from "@mui/material/IconButton";
import InputAdornment from "@mui/material/InputAdornment";
import Link from "@mui/material/Link";
import ListItem from "@mui/material/ListItem";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemSecondaryAction from "@mui/material/ListItemSecondaryAction";
import ListItemText from "@mui/material/ListItemText";
import Slide from "@mui/material/Slide";
import { useTheme } from "@mui/material/styles";
import TextField from "@mui/material/TextField";
import Toolbar from "@mui/material/Toolbar";
import { TransitionProps } from "@mui/material/transitions";
import Typography from "@mui/material/Typography";
import useMediaQuery from "@mui/material/useMediaQuery";

import AutoSizer from "react-virtualized-auto-sizer";
import { FixedSizeList } from "react-window";

import { DiskFile, useAPI } from "../../API";
import { basename, generateParents } from "./path";

const SlideUPTransition = forwardRef(
  (
    props: TransitionProps & {
      children: ReactElement<any, any>;
    },
    ref: Ref<unknown>
  ) => {
    return <Slide direction="up" ref={ref} {...props} />;
  }
);

type InputProps<T = any> = {
  label: string;
} & T;

export type FilePathInputProps = InputProps<{
  value: string;
  onChange: (value: string) => void;
  onBlur?: () => void;
  fileType?: string;
  title?: string;
}>;
export const FilePathInput = forwardRef(
  (
    { label, value, onBlur, onChange, fileType, title }: FilePathInputProps,
    ref: Ref<HTMLInputElement>
  ) => {
    const [open, setOpen] = useState(false);
    const onOpen = () => setOpen(true);
    const onClose = () => setOpen(false);
    return (
      <Fragment>
        <TextField
          variant="filled"
          label={label}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onBlur={onBlur}
          InputProps={{
            endAdornment: (
              <InputAdornment position="end">
                <IconButton aria-label="directions" onClick={onOpen}>
                  <FolderIcon />
                </IconButton>
              </InputAdornment>
            ),
          }}
          inputRef={ref}
        />
        {open && (
          <FilePathSelect
            open={open}
            onClose={onClose}
            value={value}
            onChange={onChange}
            fileType={fileType}
            title={title}
          />
        )}
      </Fragment>
    );
  }
);

type FilePathSelectProps = {
  value: string;
  onChange: (value: string) => void;
  open: boolean;
  onClose: () => void;
  fileType?: string;
  title?: string;
};
const FilePathSelect = ({
  open,
  onClose,
  value,
  onChange,
  fileType,
  title,
}: FilePathSelectProps) => {
  const theme = useTheme();
  const fullScreen = useMediaQuery(theme.breakpoints.down("sm"));

  const api = useAPI();

  const { isFetching: isFileFetching, data: results } = useQuery({
    queryKey: ["app-disk-files", { filePath: value }],
    queryFn: async () => await api?.searchDiskFiles({ filePath: value }),
    select: (data) => (data && data[1] ? data[1][0] : void 0),
    enabled: open && !!api,
  });

  const [parentPath, setParentPath] = useState<string | null>(null);

  useEffect(() => {
    if (isFileFetching) return;
    setParentPath(results ? results.parentPath : "");
  }, [isFileFetching, results]);

  const { isFetching, data, refetch } = useQuery({
    queryKey: ["app-disk-files", { parentPath }],

    queryFn: async () =>
      await api?.searchDiskFiles(
        fileType
          ? { parentPath: parentPath as string, fileType }
          : { parentPath: parentPath as string }
      ),
    enabled: parentPath !== null && !isFileFetching,
  });

  const onRefresh = () => {
    if (!isFetching) refetch();
  };

  const onSelected = (item: DiskFile) => {
    onChange(item.filePath);
    onClose();
  };

  const onEnter = (item: DiskFile) => {
    setParentPath(item.filePath);
  };

  const [_, rows] = useMemo(() => data || [0, []], [data]);

  const id = useId();

  return (
    <Dialog
      fullScreen={fullScreen}
      open={open}
      onClose={onClose}
      TransitionComponent={SlideUPTransition}
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
            {title}
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
      {parentPath !== null && (
        <FilePathBreadcrumbs
          value={parentPath}
          onSelect={setParentPath}
          disabled={isFetching}
        />
      )}

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
              itemSize={72}
              itemCount={rows.length}
              overscanCount={20}
            >
              {({ index, style }) => (
                <ListItem style={style} key={`${id}-${index}`} disablePadding>
                  <ListItemButton
                    selected={rows[index].filePath === value}
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
                    <ListItemSecondaryAction>
                      <IconButton
                        edge="end"
                        aria-label="enter"
                        onClick={(e) => {
                          e.stopPropagation();
                          onEnter(rows[index]);
                        }}
                        sx={{
                          visibility:
                            rows[index].fileType === "D" ? "visible" : "hidden",
                        }}
                      >
                        <NavigateNextIcon />
                      </IconButton>
                    </ListItemSecondaryAction>
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
