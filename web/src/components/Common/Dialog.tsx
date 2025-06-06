import { useTranslation } from "../I18Next/Context";

import Button from "@mui/material/Button";
import type { DialogProps as MuiDialogProps } from "@mui/material/Dialog";
import MuiDialog from "@mui/material/Dialog";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";

import { Fragment } from "react";

export type DialogProps = MuiDialogProps;
/**
 * A MUI Dialog component wrapped with a TranslateProvider.
 *
 * @example
 * <Dialog open={true} onClose={() => {}}>
 *   <DialogTitle>Dialog Title</DialogTitle>
 *   <DialogContent>
 *     <DialogContentText>
 *       To use Google Maps, enable JavaScript by changing your browser options, and then try again.
 *     </DialogContentText>
 *   </DialogContent>
 *   <DialogActions>
 *     <Button onClick={() => {}}>Cancel</Button>
 *     <Button onClick={() => {}}>Subscribe</Button>
 *   </DialogActions>
 * </Dialog>
 *
 * @param {DialogProps} props - The properties to pass to the MUI Dialog component.
 * @param {(node: ReactNode) => ReactNode} [children] - The content of the Dialog.
 * @returns {ReactElement} A MUI Dialog component with a TranslateProvider.
 */
export const Dialog = MuiDialog;

/**
 * A pair of buttons to save and cancel a dialog.
 *
 * @example
 * <DialogSaveActions
 *   onSave={() => {}}
 *   onCancel={() => {}}
 *   label="Save Changes"
 * />
 *
 * @param {Object} props - The properties to pass to the DialogSaveActions component.
 * @param {Function} props.onSave - The function to call when the save button is clicked.
 * @param {Function} props.onCancel - The function to call when the cancel button is clicked.
 * @param {string} [props.label] - The label to use for the save button. Defaults to "Save".
 * @returns {ReactElement} A pair of buttons to save and cancel a dialog.
 */
export const DialogSaveActions = ({
  onSave,
  onCancel,
  label,
}: {
  onSave: () => void;
  onCancel: () => void;
  label?: string;
}) => {
  const { t } = useTranslation();
  return (
    <Fragment>
      <Button size="small" onClick={onCancel}>
        {t("buttons.cancel")}
      </Button>
      <Button size="small" onClick={onSave} autoFocus>
        {label ? label : t("buttons.save")}
      </Button>
    </Fragment>
  );
};

/**
 * A pair of buttons to submit and cancel a dialog.
 *
 * @example
 * <DialogSubmitActions
 *   onSubmit={() => {}}
 *   onCancel={() => {}}
 *   label="Submit"
 * />
 *
 * @param {Object} props - The properties to pass to the DialogSubmitActions component.
 * @param {Function} props.onSubmit - The function to call when the submit button is clicked.
 * @param {Function} props.onCancel - The function to call when the cancel button is clicked.
 * @param {string} [props.label] - The label to use for the submit button. Defaults to "Submit".
 * @returns {ReactElement} A pair of buttons to submit and cancel a dialog.
 */
export const DialogSubmitActions = ({
  onSubmit,
  onCancel,
  label,
}: {
  onSubmit: () => void;
  onCancel: () => void;
  label?: string;
}) => {
  const { t } = useTranslation();
  return (
    <Fragment>
      <Button size="small" onClick={onCancel}>
        {t("buttons.cancel")}
      </Button>
      <Button size="small" onClick={onSubmit} autoFocus>
        {label ? label : t("buttons.submit")}
      </Button>
    </Fragment>
  );
};

/**
 * A pair of buttons to confirm or cancel a dialog.
 *
 * @example
 * <DialogConfirmActions
 *   onConfirm={(confirm) => {}}
 *   label="Confirm"
 * />
 *
 * @param {Object} props - The properties to pass to the DialogConfirmActions component.
 * @param {Function} props.onConfirm - The function to call when the cancel or confirm button is clicked. The function receives a boolean argument indicating whether the dialog was confirmed or not.
 * @param {string} [props.label] - The label to use for the confirm button. Defaults to "Confirm".
 * @returns {ReactElement} A pair of buttons to confirm or cancel a dialog.
 */
export const DialogConfirmActions = ({
  onConfirm,
  label,
}: {
  onConfirm: (confirm: boolean) => void;
  label?: string;
}) => {
  const { t } = useTranslation();
  return (
    <Fragment>
      <Button size="small" onClick={() => onConfirm(false)}>
        {t("buttons.cancel")}
      </Button>
      <Button
        size="small"
        onClick={() => onConfirm(true)}
        autoFocus
        color="error"
      >
        {label ? label : t("buttons.confirm")}
      </Button>
    </Fragment>
  );
};

type DialogConfirmContentProps = {
  label?: string;
  contentLabel?: string;
};
/**
 * A component that renders a dialog with a confirm title and content text.
 *
 * @param {Object} props - The properties to pass to the DialogConfirmContent component.
 * @param {string} [props.label] - The label used in the dialog title. Defaults to an empty string.
 * @param {string} [props.contentLabel] - The label used in the dialog content text. If not provided, it defaults to the lowercased `label`.
 * @returns {ReactElement} A dialog component with a title and content text.
 */

export const DialogConfirmContent = ({
  label = "",
  contentLabel,
}: DialogConfirmContentProps) => {
  const { t } = useTranslation();
  return (
    <Fragment>
      <DialogTitle>{t("dialogs.confirm.title", { label })}</DialogTitle>
      <DialogContent>
        <DialogContentText>
          {t("dialogs.confirm.content", {
            label: contentLabel || label.toLowerCase(),
          })}
        </DialogContentText>
      </DialogContent>
    </Fragment>
  );
};
