import { ClusterName } from "./meta";
import { type ExtraProps, withExtraState, ExtraType } from "./ExtraBase";
import {
  ClusterSearchFilter,
  ClusterList,
  ClusterInfo,
  ClusterInfoForm,
  ClusterPassphraseForm,
} from "./Cluster";

import { useAppDispatch } from "@/components/App/Context";
import { useLocation } from "@/components/App/Router";
import {
  I18nVariant,
  useTranslation,
  useExternalNamespace,
} from "@/components/App/I18Next";
import {
  DialogExtraTitle,
  DialogExtraClose,
  DialogExtraAction,
} from "@/components/App/Extra";
import {
  Dialog,
  DialogVariant,
  DialogDescription,
  DialogAction,
} from "@/components/App/Dialog";
import { useRef } from "react";
import type { FormEvent, ComponentProps } from "react";

import { Separator } from "@/components/ui/separator";

export const ClusterSwitchExtra = ({ extraState }: ExtraProps) => {
  const { type, open } = extraState;
  const dispatch = useAppDispatch();
  const location = useLocation();

  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };
  return (
    <ClusterSwitchExtraBase
      type={type}
      open={open}
      onClose={handleClose}
      locationState={location.state}
    />
  );
};

type ClusterSwitchExtraProps = Pick<
  ComponentProps<typeof Dialog>,
  "open" | "onClose"
> & { type?: string; locationState?: any };
export const ClusterSwitchExtraBase = ({
  open,
  onClose,
  type = ExtraType.ClusterSwitch,
  locationState,
}: ClusterSwitchExtraProps) => {
  const namespace = useExternalNamespace(ClusterName);
  const { t } = useTranslation(namespace);

  const ref = useRef<HTMLDialogElement>(null);
  const handleEsc = () => {
    ref.current?.close();
  };

  return (
    <Dialog
      variant={DialogVariant.Modal}
      open={open}
      onClose={onClose}
      className="gap-0 w-full max-w-xl p-0 border-none bg-muted text-sm"
      ref={ref}
    >
      <div className="px-2 py-3">
        <ClusterSearchFilter onEsc={handleEsc} />
      </div>
      <Separator />
      <h3 className="px-4 h-12 leading-12 font-bold text-muted-foreground text-xs">
        {t(`${I18nVariant.Extra}.${type}.recentlyUsed`)}
      </h3>
      <ClusterList locationState={locationState} />
      <div className="px-4 h-12 leading-12 text-right text-muted-foreground">
        {t(`${I18nVariant.Extra}.${type}.footer`, { count: 12 })}
        {/* Search by xxxx */}
      </div>
    </Dialog>
  );
};

export const ClusterRemoveExtra = ({ extraState }: ExtraProps) => {
  const { type, open } = extraState;
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };
  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtraClose onClose={handleClose} />
      <DialogExtraTitle text="设备组" className="text-destructive">
        移除
      </DialogExtraTitle>
      <DialogDescription>
        <p>确认将当前设备将要从下面的设备组中移除吗？</p>
      </DialogDescription>
      <ClusterInfo variant="disused" />
      <DialogAction>
        <DialogExtraAction
          variant="destructive"
          onClose={handleClose}
          onClick={handleClose}
        >
          确认移除
        </DialogExtraAction>
      </DialogAction>
    </Dialog>
  );
};

export const ClusterEditExtra = ({ extraState }: ExtraProps) => {
  const { type, open } = extraState;
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    handleClose();
  };

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtraClose onClose={handleClose} />
      <DialogExtraTitle text="设备组">设置</DialogExtraTitle>
      <DialogDescription>编辑设备组信息</DialogDescription>
      <ClusterInfoForm id="cluster-info-form" onSubmit={handleSubmit} />
      <DialogAction>
        <DialogExtraAction
          onClose={handleClose}
          type="submit"
          form="cluster-info-form"
        >
          保存
        </DialogExtraAction>
      </DialogAction>
    </Dialog>
  );
};

export const ClusterPassphraseEditExtra = ({ extraState }: ExtraProps) => {
  const { type, open } = extraState;
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    handleClose();
  };
  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtraClose onClose={handleClose} />
      <DialogExtraTitle text="管理口令">设置</DialogExtraTitle>
      <ClusterInfo />
      <DialogDescription>
        <p>设置新的管理口令,请确保口令强度足够.</p>
      </DialogDescription>

      <ClusterPassphraseForm
        id="cluster-passphrase-form"
        onSubmit={handleSubmit}
      />
      <DialogAction>
        <DialogExtraAction
          onClose={handleClose}
          type="submit"
          form="cluster-passphrase-form"
        >
          保存新密码
        </DialogExtraAction>
      </DialogAction>
    </Dialog>
  );
};
