import {
  generateCreatePath,
  generateEditPath,
  generateShowPath,
} from "../Route/utils";
import CreateNewFolderIcon from "@mui/icons-material/CreateNewFolder";

export const ExtFSNodeItemIcon = CreateNewFolderIcon;
export const ExtFSNodeItemI18nKey = "custom.extfs/local-node-items.name";
export const ExtFSNodeItemPath = "/extfs/local-node-items";
export const ExtFSNodeItemCreateI18nKey =
  "custom.extfs/local-node-items.create-name";
export const ExtFSNodeItemCreatePath = generateCreatePath(ExtFSNodeItemPath);
export const ExtFSNodeItemEditI18nKey =
  "custom.extfs/local-node-items.edit-name";
export const ExtFSNodeItemEditPath = generateEditPath(ExtFSNodeItemPath);
export const ExtFSNodeItemShowPath = generateShowPath(ExtFSNodeItemPath);
