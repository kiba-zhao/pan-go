import LanIcon from "@mui/icons-material/Lan";
import { generateCreatePath, generateEditPath } from "../Route/utils";

export const AppNodeCreateI18nKey = "resources.app/nodes.create-name";
export const AppNodeEditI18nKey = "resources.app/nodes.edit-name";
export const AppNodeIcon = LanIcon;
export const AppNodePath = "/app/nodes";
export const AppNodeCreatePath = generateCreatePath(AppNodePath);
export const AppNodeEditPath = generateEditPath(AppNodePath);
