import { ClusterName, ClusterRoutePath } from "./meta";
import { type ExtraProps, withExtraState, ExtraType } from "./ExtraBase";
import {
  ClusterSearchFilter,
  ClusterList,
  ClusterInfo,
  ClusterInfoForm,
  ClusterPassphraseForm,
} from "./Cluster";
import { PassportQRCode } from "./Passport";

import { useAppDispatch } from "@/components/App/Context";
import {
  useLocation,
  useParams,
  useNavigate,
  generatePath,
} from "@/components/App/Router";
import {
  I18nVariant,
  useTranslation,
  useExternalNamespace,
} from "@/components/App/I18Next";

import {
  Dialog,
  DialogExtra,
  DialogVariant,
  DialogExtraClose,
  DialogExtraAction,
  DialogExtraCloseAction,
  DialogExtraTitle,
  CardDialogExtra,
} from "@/components/App/Dialog";

import { useRef } from "react";
import type { FormEvent, ComponentProps } from "react";

import { Separator } from "@/components/ui/separator";
import { Card, CardFooter } from "@/components/ui/card";
import { Marker, MarkerContent } from "../ui/marker";

export const ClusterSwitchExtra = ({ extraState }: ExtraProps) => {
  const { type, open } = extraState;
  const dispatch = useAppDispatch();
  const location = useLocation();
  const { clusterId } = useParams();

  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };
  return (
    <ClusterSwitchExtraBase
      type={type}
      open={open}
      onClose={handleClose}
      locationState={location.state}
      clusterId={Number(clusterId)}
    />
  );
};

type ClusterSwitchExtraProps = Pick<
  ComponentProps<typeof Dialog>,
  "open" | "onClose"
> & { type?: string; locationState?: any; clusterId?: number };
export const ClusterSwitchExtraBase = ({
  open,
  onClose,
  type = ExtraType.ClusterSwitch,
  locationState,
  clusterId,
}: ClusterSwitchExtraProps) => {
  const namespace = useExternalNamespace(ClusterName);
  const { t } = useTranslation(namespace);

  const ref = useRef<HTMLDialogElement>(null);
  const handleEsc = () => {
    ref.current?.close();
  };

  const navigate = useNavigate();
  const handleSelect = (clusterId: number) => {
    navigate(
      generatePath(ClusterRoutePath, {
        clusterId: clusterId.toString(),
      }),
      { state: locationState },
    );
    handleEsc();
  };
  return (
    <Dialog
      variant={DialogVariant.Modal}
      open={open}
      onClose={onClose}
      className="md:max-w-xl bg-muted text-foreground text-sm"
      ref={ref}
    >
      <div className="px-2 py-3">
        <ClusterSearchFilter onEsc={handleEsc} />
      </div>
      <Separator />
      <h3 className="px-4 h-12 leading-12 font-bold text-muted-foreground text-xs">
        {t(`${I18nVariant.Extra}.${type}.recentlyUsed`)}
      </h3>
      <ClusterList selected={clusterId} onSelect={handleSelect} />
      <div className="px-4 h-12 leading-12 text-right text-muted-foreground">
        {t(`${I18nVariant.Extra}.${type}.footer`, { count: 12 })}
      </div>
    </Dialog>
  );
};

export const ClusterRemoveExtra = ({ extraState }: ExtraProps) => {
  const { type } = extraState;
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };
  return (
    <CardDialogExtra
      title={
        <DialogExtraTitle
          text="设备组"
          smallText="移除"
          className="text-destructive"
        />
      }
      description={<p>确认将当前设备将要从下面的设备组中移除吗？</p>}
      action={<DialogExtraClose />}
      footer={
        <DialogExtraAction variant="destructive" onClick={handleClose}>
          确认移除
        </DialogExtraAction>
      }
      footerClassName="justify-end gap-2"
    >
      <ClusterInfo variant="disused" />
    </CardDialogExtra>
  );
};

export const ClusterEditExtra = ({ extraState }: ExtraProps) => {
  const { type } = extraState;
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    handleClose();
  };

  return (
    <CardDialogExtra
      title={<DialogExtraTitle text="设备组" smallText="设置" />}
      description={<p>编辑设备组信息.</p>}
      action={<DialogExtraClose />}
      footer={
        <DialogExtraAction type="submit" form="cluster-info-form">
          保存
        </DialogExtraAction>
      }
      footerClassName="justify-end gap-2"
    >
      <ClusterInfoForm id="cluster-info-form" onSubmit={handleSubmit} />
    </CardDialogExtra>
  );
};

export const ClusterPassphraseEditExtra = ({ extraState }: ExtraProps) => {
  const { type } = extraState;
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    handleClose();
  };
  return (
    <CardDialogExtra
      title={<DialogExtraTitle text="管理口令" smallText="设置" />}
      description={<p>设置新的管理口令,请确保口令强度足够.</p>}
      action={<DialogExtraClose />}
      footer={
        <DialogExtraAction type="submit" form="cluster-passphrase-form">
          保存新密码
        </DialogExtraAction>
      }
      footerClassName="justify-end gap-2"
    >
      <ClusterInfo />
      <Marker variant="separator">
        <MarkerContent>设置新口令</MarkerContent>
      </Marker>
      <ClusterPassphraseForm
        id="cluster-passphrase-form"
        onSubmit={handleSubmit}
      />
    </CardDialogExtra>
  );
};

export const ClusterAddExtra = ({ extraState }: ExtraProps) => {
  const { t } = useTranslation(ClusterName);

  return (
    <DialogExtra>
      <Card className="pt-0">
        <PassportQRCode />
        <CardFooter>
          <DialogExtraCloseAction className="w-full" variant="outline">
            {t(`${I18nVariant.Action}.cancel`)}
          </DialogExtraCloseAction>
        </CardFooter>
      </Card>
    </DialogExtra>
  );
};
