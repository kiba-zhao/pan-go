import {
  generateCreatePath,
  generateEditPath,
  generateShowPath,
} from "../Route/utils";
import CreateNewFolderIcon from "@mui/icons-material/CreateNewFolder";

export const ExtFSNodeItemNS = "extfs_node_item";
export const ExtFSNodeItemIcon = CreateNewFolderIcon;
export const ExtFSNodeItemI18nKey = "resources.extfs/local-node-items.name";
export const ExtFSNodeItemPath = "/extfs/local-node-items";
export const ExtFSNodeItemCreateI18nKey =
  "resources.extfs/local-node-items.create-name";
export const ExtFSNodeItemCreatePath = generateCreatePath(ExtFSNodeItemPath);
export const ExtFSNodeItemEditI18nKey =
  "resources.extfs/local-node-items.edit-name";
export const ExtFSNodeItemEditPath = generateEditPath(ExtFSNodeItemPath);
export const ExtFSNodeItemShowPath = generateShowPath(ExtFSNodeItemPath);
