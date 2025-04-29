/**
 * ExtFSBrowseFile Page Definition File
 *
 * Components used for browsing files
 */
import Alert from "@mui/material/Alert";
import AlertTitle from "@mui/material/AlertTitle";
import LinearProgress from "@mui/material/LinearProgress";
import type { LinkProps } from "@mui/material/Link";
import Link from "@mui/material/Link";
import Typography from "@mui/material/Typography";

import { Fragment, useEffect, useMemo, useState } from "react";
import { Trans } from "react-i18next";

import { useQuery } from "@tanstack/react-query";

import { useSearchParams } from "../Route/Router";
import {
  selectExtFSRemoteItem,
  ROOT_PATH,
  selectExtFSNodeItem,
  selectExtFSRemoteFile,
  selectExtFSNodeFile,
} from "./api";

import { ExtFSNodeFileQueryKey } from "../ExtFS/NodeFile";
import { ExtFSNodeQueryKey } from "../ExtFS/NodeItem";
import { ExtFSRemoteFileQueryKey } from "../ExtFS/RemoteFile";
import { ExtFSRemoteQueryKey } from "../ExtFS/RemoteItem";
import { useBrowser } from "../Browser";
import { useTranslation } from "../i18n/Context";
import { PageHeader } from "../Master/Header";
import { PageLayout, PageNotFound } from "../Master/Page";
import { ExtFSBrowseFileI18nKey } from "./Route";

const ExtFSBrowseFileTitle = () => {
  const { t } = useTranslation();
  return <Typography variant="h6">{t(ExtFSBrowseFileI18nKey)}</Typography>;
};

/**
 * ExtFSBrowseFile Component
 *
 * Component that renders a file or a directory,
 * depending on the presence of the filePath parameter
 *
 * @param {number} itemId id of the item in url query
 * @param {string} peerId id of the peer in url query (if any)
 * @param {string} filePath path to the file or directory in url query
 *
 * @returns {JSX.Element} a JSX element representing the file or directory
 */
const ExtFSBrowseFile = () => {
  const [query, _] = useSearchParams();
  const { peerId, itemId, filePath } = useMemo(
    () => parseQueryParams(query),
    [query]
  );

  if (itemId < 1) {
    return <PageNotFound />;
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

/**
 * Parse URL query parameters into a ExtFSBrowseFileQueryParams object
 *
 * @param {URLSearchParams} query the URL query parameters to parse
 *
 * @returns {ExtFSBrowseFileQueryParams} an object containing the parsed query parameters
 */
function parseQueryParams(query: URLSearchParams): ExtFSBrowseFileQueryParams {
  const params = {} as ExtFSBrowseFileQueryParams;
  params.peerId = query.get("peerId") || "";
  params.filePath = query.get("filePath") || "";

  const itemId = parseInt(query.get("itemId") || "");
  params.itemId = isNaN(itemId) ? -1 : itemId;
  return params;
}

export default ExtFSBrowseFile;

/**
 * Calculate 1024 to the power of a given number
 *
 * @param {number} num the power to which 1024 should be raised
 *
 * @returns {number} 1024 to the power of num
 */
function pow1024(num: number) {
  return Math.pow(1024, num);
}

/**
 * Converts a file size from bytes to a human-readable string format.
 *
 * @param {number} size - The file size in bytes.
 * @returns {string} The file size formatted as a string with appropriate units (B, KB, MB, GB, or TB).
 *                   Returns an empty string if the size is not provided or invalid.
 */

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
  fileType: string;
  disabled: boolean;
  error: Error | null;
  isLoading: boolean;
} & Pick<LinkProps, "href">;
/**
 * ExtFSBrowseFileView Component
 *
 * Browse page view for displaying file of extfs
 *
 * @param {string} fileName name of the file
 * @param {number} fileSize size of the file in bytes
 * @param {string} fileType type of the file ("F" for files and "D" for directories)
 * @param {boolean} disabled if true, the file is disabled
 * @param {Error | null} error error that occurred while fetching the file
 * @param {boolean} isLoading if true, the file is loading
 * @param {LinkProps} linkProps props to pass to the Link component
 *
 * @returns {JSX.Element} a JSX element representing the file or directory
 */
const ExtFSBrowseFileView = ({
  fileName,
  fileSize,
  fileType,
  disabled,
  error,
  isLoading,
  ...linkProps
}: ExtFSBrowseFileViewProps) => {
  const { t } = useTranslation();

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

    if (fileType !== "F") {
      return {
        title: t("custom.extfs/browse-files.invalid"),
        desc: t("custom.extfs/browse-files.invalid_desc", {
          fileName,
          fileSize: fileSizeStr,
        }),
      };
    }
    return void 0;
  }, [error, disabled]);

  if (error?.name === "FetchError.404") {
    return <PageNotFound />;
  }

  return (
    <PageLayout>
      <PageHeader title={<ExtFSBrowseFileTitle />} />
      {isLoading ? <LinearProgress /> : void 0}
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
    </PageLayout>
  );
};

/**
 * Component to render an error message for ExtFSBrowseFile
 *
 * @param {{ title: string, desc: string }} props
 * @prop {string} title - error title
 * @prop {string} desc - error description
 * @returns {JSX.Element} a JSX element representing the error message
 */
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

/**
 * ExtFSBrowseFile Component with fetch file info success
 *
 * @param {{ fileName: string, fileSize: string }} props
 * @prop {string} fileName - name of the file
 * @prop {string} fileSize - size of the file
 * @returns {JSX.Element} a JSX element representing the success message
 */
const ExtFSBrowseFileViewWithSuccess = ({
  fileName,
  fileSize,
  ...linkProps
}: Omit<
  ExtFSBrowseFileViewProps,
  "fileSize" | "disabled" | "error" | "isLoading" | "fileType"
> & {
  fileSize: string;
}) => {
  const { t } = useTranslation();
  const browser = useBrowser();

  const [countDown, setCountDown] = useState(6);
  useEffect(() => {
    if (countDown > 0) {
      const ret = setTimeout(() => setCountDown(countDown - 1), 1000);
      return () => clearTimeout(ret);
    }
    const { window } = browser || {};
    if (window) window.location.replace(linkProps.href || "");
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
/**
 * ExtFSBrowseFile Component with fetch remote item info
 *
 * @param {{ peerId: string, itemId: number }} props
 * @prop {string} peerId - id of the peer in url query
 * @prop {number} itemId - id of the item in url query
 *
 * @returns {JSX.Element} a JSX element representing the file or directory
 */
const ExtFSBrowseFileWithRemoteItem = ({
  peerId,
  itemId,
}: ExtFSBrowseFileWithRemoteItemProps) => {
  const { data, isFetching, error } = useQuery({
    queryKey: [...ExtFSRemoteQueryKey, peerId, itemId],
    queryFn: async () => await selectExtFSRemoteItem(peerId, itemId),
    enabled: itemId > 0,
  });

  return (
    <ExtFSBrowseFileView
      fileName={data?.name || ""}
      fileSize={data?.size || 0}
      fileType={data?.fileType || ""}
      isLoading={isFetching}
      disabled={!data?.available}
      error={error}
      href={`${ROOT_PATH}/extfs/remotes/${peerId}/remote-items/${itemId}/_stream/`}
    />
  );
};

type ExtFSBrowseFileWithNodeItemProps = Omit<
  ExtFSBrowseFileWithRemoteItemProps,
  "peerId"
>;
/**
 * ExtFSBrowseFile Component with fetch node item info
 *
 * @param {{ itemId: number }} props
 * @prop {number} itemId - id of the item in url query
 *
 * @returns {JSX.Element} a JSX element representing the file or directory
 */
const ExtFSBrowseFileWithNodeItem = ({
  itemId,
}: ExtFSBrowseFileWithNodeItemProps) => {
  const { data, isFetching, error } = useQuery({
    queryKey: [...ExtFSNodeQueryKey, itemId],
    queryFn: async () => await selectExtFSNodeItem(itemId),
    enabled: itemId > 0,
  });

  return (
    <ExtFSBrowseFileView
      fileName={data?.name || ""}
      fileSize={data?.size || 0}
      fileType={data?.fileType || ""}
      isLoading={isFetching}
      disabled={!data?.available}
      error={error}
      href={`${ROOT_PATH}/extfs/node-items/${itemId}/_stream/`}
    />
  );
};

type ExtFSBrowseFileWithRemoteFileProps = ExtFSBrowseFileQueryParams;
/**
 * ExtFSBrowseFile Component with fetch remote file info
 *
 * @param {{ peerId: string, itemId: number, filePath: string }} props
 * @prop {string} peerId - id of the peer in url query
 * @prop {number} itemId - id of the item in url query
 * @prop {string} filePath - path to the file in url query
 *
 * @returns {JSX.Element} a JSX element representing the file or directory
 */

const ExtFSBrowseFileWithRemoteFile = ({
  peerId,
  itemId,
  filePath,
}: ExtFSBrowseFileWithRemoteFileProps) => {
  const { data, isFetching, error } = useQuery({
    queryKey: [...ExtFSRemoteFileQueryKey, peerId, itemId],
    queryFn: async () => await selectExtFSRemoteFile(peerId, itemId, filePath),
    enabled: itemId > 0,
  });
  return (
    <ExtFSBrowseFileView
      fileName={data?.name || ""}
      fileSize={data?.size || 0}
      fileType={data?.fileType || ""}
      isLoading={isFetching}
      disabled={!data?.available}
      error={error}
      href={`${ROOT_PATH}/extfs/remotes/${peerId}/remote-items/${itemId}/_stream/${filePath}`}
    />
  );
};

type ExtFSBrowseFileWithNodeFileProps = Omit<
  ExtFSBrowseFileWithRemoteFileProps,
  "peerId"
>;
/**
 * ExtFSBrowseFile Component with fetch node file info
 *
 * @param {{ itemId: number, filePath: string }} props
 * @prop {number} itemId - id of the item in url query
 * @prop {string} filePath - path to the file in url query
 *
 * @returns {JSX.Element} a JSX element representing the file or directory
 */
const ExtFSBrowseFileWithNodeFile = ({
  itemId,
  filePath,
}: ExtFSBrowseFileWithNodeFileProps) => {
  const { data, isFetching, error } = useQuery({
    queryKey: [...ExtFSNodeFileQueryKey, { itemId, filePath }],
    queryFn: async () => await selectExtFSNodeFile(itemId, filePath),
    enabled: itemId > 0,
  });

  return (
    <ExtFSBrowseFileView
      fileName={data?.name || ""}
      fileSize={data?.size || 0}
      fileType={data?.fileType || ""}
      isLoading={isFetching}
      disabled={!data?.available}
      error={error}
      href={`${ROOT_PATH}/extfs/node-items/${itemId}/_stream/${filePath}`}
    />
  );
};
