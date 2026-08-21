import { SettingsName } from "./meta";
import { useTranslation, I18nVariant } from "@/components/App/I18Next";
import {
  DialogExtraState,
  withDialogExtraState,
} from "@/components/App/Dialog";
import { AppExtraState, DialogExtraFormAction } from "@/components/App/Extra";
import { DialogExtraTitle as AppDialogExtraTitle } from "@/components/App/Dialog";

export enum ExtraType {
  DeviceNameEdit = "deviceNameEdit",
  DeviceMemoEdit = "deviceMemoEdit",
  PeerPortEdit = "peerPortEdit",
  BroadcastAddrsEdit = "broadcastAddrsEdit",
  PublicAddrsEdit = "publicAddrsEdit",
  WebPortEdit = "webPortEdit",
  ClusterSelect = "clusterSelect",
}

export type ExtraState = {} & AppExtraState<ExtraType> & DialogExtraState;
export function withExtraState(extraState: ExtraState) {
  return withDialogExtraState(extraState);
}

export type ExtraProps = { extraState: ExtraState };

export const DialogExtraTitle = ({ extraState }: ExtraProps) => {
  const { type } = extraState;
  const { t } = useTranslation(SettingsName);

  return (
    <AppDialogExtraTitle
      text={t(`${I18nVariant.Extra}.${type}.title`)}
      smallText={t(`${I18nVariant.Badge}.edit`)}
    />
  );
};

export { DialogExtraFormAction };

// export const FormExtra = <T extends FieldValues = FieldValues>({
//   extraState,
//   control,
//   title,
//   badge,
//   desc,
//   error,
//   blank,
//   id,
//   children,
//   ...props
// }: FormDialogExtraProps<T>) => {
//   const { type, open } = extraState;
//   const { t } = useTranslation(SettingsName);
//   const dispatch = useAppDispatch();

//   const handleClose = () => {
//     dispatch?.(withDialogExtraState({ type, open: false }));
//   };

//   const {
//     disabled: formDisabled,
//     isValid,
//     isDirty,
//     isSubmitting,
//   } = useFormState<T>({ control });

//   return (
//     <Dialog open={open} onClose={handleClose}>
//       <DialogExtraClose onClose={handleClose} />
//       <DialogExtraTitle text={title || t(`${I18nVariant.Extra}.${type}.title`)}>
//         {badge || t(`${I18nVariant.Extra}.${type}.subtitle`)}
//       </DialogExtraTitle>
//       <DialogDescription className={error || isSubmitting ? "hidden" : ""}>
//         {desc}
//       </DialogDescription>
//       <DialogErrorAlert
//         error={error}
//         className={!error || isSubmitting ? "hidden" : ""}
//       >
//         {blank}
//       </DialogErrorAlert>
//       <DialogProgressAlert className={!isSubmitting ? "hidden" : ""}>
//         {t(`${type}.saving`)}
//         {blank}
//       </DialogProgressAlert>
//       <form id={id} {...props}>
//         {children}
//       </form>
//       <DialogAction className="bg-muted">
//         <DialogExtraAction
//           type="submit"
//           form={id}
//           disabled={formDisabled || !isValid || !isDirty}
//           progress={isSubmitting}
//           onClose={handleClose}
//         >
//           {t(`${I18nVariant.Action}.save`)}
//         </DialogExtraAction>
//       </DialogAction>
//     </Dialog>
//   );
// };
