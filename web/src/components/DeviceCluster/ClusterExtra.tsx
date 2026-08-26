import { ClusterName, ClusterRoutePath } from "./meta";
import { type ExtraProps, withExtraState, ExtraType } from "./ExtraBase";
import {
  ClusterSearchFilter,
  ClusterList,
  ClusterInfo,
  ClusterInfoInputFields,
  ClusterPassphraseInputField,
  PassphraseVariant,
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
import { Card, CardAction } from "@/components/ui/card";
import { Marker, MarkerContent } from "@/components/ui/marker";
import { FieldGroup, Field } from "@/components/ui/field";
import { Button } from "@/components/ui/button";

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
      <form id="cluster-info-form" onSubmit={handleSubmit}>
        <FieldGroup>
          <ClusterInfoInputFields />
        </FieldGroup>
      </form>
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
      <form id="cluster-passphrase-form" onSubmit={handleSubmit}>
        <FieldGroup>
          <ClusterPassphraseInputField variant={PassphraseVariant.Old} />
          <ClusterPassphraseInputField variant={PassphraseVariant.New} />
          <ClusterPassphraseInputField variant={PassphraseVariant.NewConfirm} />
        </FieldGroup>
      </form>
    </CardDialogExtra>
  );
};

export const ClusterAddExtra = ({ extraState }: ExtraProps) => {
  const { t } = useTranslation(ClusterName);

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
    <DialogExtra className="md:max-w-3xl">
      <Card className="overflow-hidden p-0 gap-0 flex flex-row">
        <div className="bg-muted relative max-md:hidden">
          <div className="flex flex-col items-center gap-2 text-center p-6 md:p-8">
            <h3 className="text-2xl font-bold ">加入设备组</h3>
            <p className="text-balance text-muted-foreground">
              使用其他设备扫描二维码，加入现有设备组
            </p>
          </div>
          <div className="flex flex-col items-center gap-14 px-6">
            <PassportQRCode className="px-10" />
            <div className="w-full grid gap-2 md:grid-rows-2">
              <Button className="w-full">更新二维码</Button>
              <Marker variant="separator">
                <MarkerContent>或者</MarkerContent>
              </Marker>
              <Button variant="outline" className="w-full">
                保存二维码
              </Button>
            </div>
          </div>
        </div>
        <div className="grow p-6 md:p-8">
          <FieldGroup>
            <div className="flex flex-col items-center gap-2 text-center">
              <h3 className="text-2xl font-bold">创建新设备组</h3>
              <p className="text-balance text-muted-foreground">
                创建全新的设备组。
              </p>
            </div>
            <form id="cluster-add-form" onSubmit={handleSubmit}>
              <ClusterInfoInputFields />
              <ClusterPassphraseInputField />
              <ClusterPassphraseInputField
                variant={PassphraseVariant.Confirm}
              />
            </form>
            <Field>
              <Button type="submit">创建设备组</Button>
            </Field>
          </FieldGroup>
        </div>
        <DialogExtraClose variant="ghost" className="absolute top-2 right-2" />
      </Card>
    </DialogExtra>
  );
};
