import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import LinearProgress from "@mui/material/LinearProgress";
import Link from "@mui/material/Link";
import type { LinkProps } from "@mui/material/Link";
import Alert from "@mui/material/Alert";
import AlertTitle from "@mui/material/AlertTitle";

import { Fragment, useMemo, useState, useEffect } from "react";
import { Trans } from "react-i18next";

import { Title, useTranslate } from "react-admin";
import { useSearchParams, useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";

import NotFound from "./NotFound";
import { ExtFSNodeQueryKey } from "./ExtFS/NodeItem";
import { ExtFSNodeFileQueryKey } from "./ExtFS/NodeFile";
import { ExtFSRemoteQueryKey } from "./ExtFS/RemoteItem";
import { ExtFSRemoteFileQueryKey } from "./ExtFS/RemoteFile";
import { useAPI } from "./API";

export const RoutePath = "/extfs/browse-files";

const ExtFSBrowseFile = () => {
  const [query, _] = useSearchParams();
  const { peerId, itemId, filePath } = useMemo(
    () => parseQueryParams(query),
    [query]
  );

  if (itemId < 1) {
    return <NotFound />;
  }

  if (filePath.length <= 0) {
    return peerId.length > 0 ? (
      <ExtFSBrowseFileWithRemoteItem peerId={peerId} itemId={itemId} />
    ) : (
      <ExtFSBrowseFileWithNodeItem itemId={itemId} />
    );
  }

  return peerId.length > 0 ? (
    <ExtFSBrowseFileWithRemoteFile
      peerId={peerId}
      itemId={itemId}
      filePath={filePath}
    />
  ) : (
    <ExtFSBrowseFileWithNodeFile itemId={itemId} filePath={filePath} />
  );
};

type ExtFSBrowseFileQueryParams = {
  peerId: string;
  itemId: number;
  filePath: string;
};

function parseQueryParams(query: URLSearchParams): ExtFSBrowseFileQueryParams {
  const params = {} as ExtFSBrowseFileQueryParams;
  params.peerId = query.get("peerId") || "";
  params.filePath = query.get("filePath") || "";

  const itemId = parseInt(query.get("itemId") || "");
  params.itemId = isNaN(itemId) ? -1 : itemId;
  return params;
}

export default ExtFSBrowseFile;

function pow1024(num: number) {
  return Math.pow(1024, num);
}

function convertFileSize(size: number) {
  if (!size) return "";
  if (size < pow1024(1)) return size + " B";
  if (size < pow1024(2)) return (size / pow1024(1)).toFixed(2) + " KB";
  if (size < pow1024(3)) return (size / pow1024(2)).toFixed(2) + " MB";
  if (size < pow1024(4)) return (size / pow1024(3)).toFixed(2) + " GB";
  return (size / pow1024(4)).toFixed(2) + " TB";
}

type ExtFSBrowseFileViewProps = {
  fileName: string;
  fileSize: number;
  disabled: boolean;
  error: Error | null;
  isLoading: boolean;
} & Pick<LinkProps, "href">;
const ExtFSBrowseFileView = ({
  fileName,
  fileSize,
  disabled,
  error,
  isLoading,
  ...linkProps
}: ExtFSBrowseFileViewProps) => {
  const t = useTranslate();

  const fileSizeStr = useMemo(() => convertFileSize(fileSize), [fileSize]);

  const errorProps = useMemo(() => {
    if (error) {
      return {
        title: t("custom.extfs/browse-files.error"),
        desc: error.message,
      };
    }
    if (disabled) {
      return {
        title: t("custom.extfs/browse-files.disabled"),
        desc: t("custom.extfs/browse-files.disabled_desc", {
          fileName,
          fileSize: fileSizeStr,
        }),
      };
    }
    return void 0;
  }, [error, disabled]);

  return (
    <Fragment>
      <Title title={t("custom.extfs/browse-files.name")} />
      <Card
        sx={{
          height: "100%",
          paddingBottom: "76px",
          paddingTop: "56px",
          marginTop: "-56px",
        }}
      >
        {isLoading ? <LinearProgress /> : void 0}
        <CardContent sx={{ paddingTop: 2 }}>
          {!isLoading && errorProps ? (
            <ExtFSBrowseFileViewWithError {...errorProps} />
          ) : (
            void 0
          )}
          {!isLoading && !errorProps ? (
            <ExtFSBrowseFileViewWithSuccess
              fileName={fileName}
              fileSize={fileSizeStr}
              {...linkProps}
            />
          ) : (
            void 0
          )}
        </CardContent>
      </Card>
    </Fragment>
  );
};

const ExtFSBrowseFileViewWithError = ({
  title,
  desc,
}: {
  title: string;
  desc: string;
}) => {
  return (
    <Alert severity="error">
      <AlertTitle>{title}</AlertTitle>
      {desc}
    </Alert>
  );
};

const ExtFSBrowseFileViewWithSuccess = ({
  fileName,
  fileSize,
  ...linkProps
}: Omit<
  ExtFSBrowseFileViewProps,
  "fileSize" | "disabled" | "error" | "isLoading"
> & {
  fileSize: string;
}) => {
  const t = useTranslate();

  const navigate = useNavigate();

  const [countDown, setCountDown] = useState(6);
  useEffect(() => {
    if (countDown > 0) {
      const ret = setTimeout(() => setCountDown(countDown - 1), 1000);
      return () => clearTimeout(ret);
    }
    navigate(linkProps.href || "");
  });

  return (
    <Alert severity={countDown > 0 ? "success" : "info"}>
      <AlertTitle>{t("custom.extfs/browse-files.success")}</AlertTitle>
      {countDown > 0 ? (
        <Fragment>
          <Trans
            i18nKey={"custom.extfs/browse-files.success_countdown"}
            values={{ countDown }}
            components={{
              Typography: <span style={{ color: "red" }} />,
            }}
          ></Trans>
          <Trans
            i18nKey={"custom.extfs/browse-files.success_desc"}
            values={{ fileName, fileSize: fileSize }}
            components={{
              Link: <Link {...linkProps} />,
            }}
          ></Trans>
        </Fragment>
      ) : (
        t("custom.extfs/browse-files.success_delay_end", {
          fileName,
          fileSize: fileSize,
        })
      )}
    </Alert>
  );
};

type ExtFSBrowseFileWithRemoteItemProps = Omit<
  ExtFSBrowseFileQueryParams,
  "filePath"
>;
const ExtFSBrowseFileWithRemoteItem = ({
  peerId,
  itemId,
}: ExtFSBrowseFileWithRemoteItemProps) => {
  const api = useAPI();
  const { data, isFetching, error } = useQuery({
    queryKey: [...ExtFSRemoteQueryKey, peerId, itemId],
    queryFn: async () => await api?.selectExtFSRemoteItem(peerId, itemId),
    enabled: itemId > 0,
  });

  return (
    <ExtFSBrowseFileView
      fileName={data?.name || ""}
      fileSize={data?.size || 0}
      isLoading={isFetching}
      disabled={!data?.available}
      error={error}
      href={`${api.RootPath()}/extfs/remotes/${peerId}/remote-items/${itemId}/_stream/`}
    />
  );
};

type ExtFSBrowseFileWithNodeItemProps = Omit<
  ExtFSBrowseFileWithRemoteItemProps,
  "peerId"
>;
const ExtFSBrowseFileWithNodeItem = ({
  itemId,
}: ExtFSBrowseFileWithNodeItemProps) => {
  const api = useAPI();
  const { data, isFetching, error } = useQuery({
    queryKey: [...ExtFSNodeQueryKey, itemId],
    queryFn: async () => await api?.selectExtFSNodeItem(itemId),
    enabled: itemId > 0,
  });

  return (
    <ExtFSBrowseFileView
      fileName={data?.name || ""}
      fileSize={data?.size || 0}
      isLoading={isFetching}
      disabled={!data?.available}
      error={error}
      href={`${api.RootPath()}/extfs/node-items/${itemId}/_stream/`}
    />
  );
};

type ExtFSBrowseFileWithRemoteFileProps = ExtFSBrowseFileQueryParams;
const ExtFSBrowseFileWithRemoteFile = ({
  peerId,
  itemId,
  filePath,
}: ExtFSBrowseFileWithRemoteFileProps) => {
  const api = useAPI();
  const { data, isFetching, error } = useQuery({
    queryKey: [...ExtFSRemoteFileQueryKey, peerId, itemId],
    queryFn: async () =>
      await api?.selectExtFSRemoteFile(peerId, itemId, filePath),
    enabled: itemId > 0,
  });
  return (
    <ExtFSBrowseFileView
      fileName={data?.name || ""}
      fileSize={data?.size || 0}
      isLoading={isFetching}
      disabled={!data?.available}
      error={error}
      href={`${api.RootPath()}/extfs/remotes/${peerId}/remote-items/${itemId}/_stream/${filePath}`}
    />
  );
};

type ExtFSBrowseFileWithNodeFileProps = Omit<
  ExtFSBrowseFileWithRemoteFileProps,
  "peerId"
>;
const ExtFSBrowseFileWithNodeFile = ({
  itemId,
  filePath,
}: ExtFSBrowseFileWithNodeFileProps) => {
  const api = useAPI();
  const { data, isFetching, error } = useQuery({
    queryKey: [...ExtFSNodeFileQueryKey, { itemId, filePath }],
    queryFn: async () => await api?.selectExtFSNodeFile(itemId, filePath),
    enabled: itemId > 0,
  });

  return (
    <ExtFSBrowseFileView
      fileName={data?.name || ""}
      fileSize={data?.size || 0}
      isLoading={isFetching}
      disabled={!data?.available}
      error={error}
      href={`${api.RootPath()}/extfs/node-items/${itemId}/_stream/${filePath}`}
    />
  );
};
