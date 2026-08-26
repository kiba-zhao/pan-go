import { type ExtraProps, withExtraState } from "./ExtraBase";
import { DeviceList, DeviceFormInput } from "./Device";
import { PassportQRScan, PassportQRFileInputAction } from "./Passport";
import { ClusterInfo } from "./Cluster";
import type { FormEvent } from "react";

import { useAppDispatch } from "@/components/App/Context";

import {
  DialogExtra,
  DialogExtraClose,
  DialogExtraCloseAction,
  DialogExtraAction,
  DialogExtraTitle,
  CardDialogExtra,
} from "@/components/App/Dialog";
import { Marker, MarkerContent } from "@/components/ui/marker";
import { Card } from "@/components/ui/card";

export const DevicesRemoveExtra = ({ extraState }: ExtraProps) => {
  const { type } = extraState;
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };
  return (
    <CardDialogExtra
      title={
        <DialogExtraTitle text="相关设备" className="text-destructive">
          移除
        </DialogExtraTitle>
      }
      description={
        <p>
          确认将下列<span className="text-destructive">32</span>
          个设备从上面的设备组中移除吗？
        </p>
      }
      action={<DialogExtraClose />}
      footer={
        <DialogExtraAction variant="destructive" onClick={handleClose}>
          确认移除
        </DialogExtraAction>
      }
      footerClassName="justify-end gap-2"
    >
      <ClusterInfo />
      <Marker variant="separator">
        <MarkerContent>设备列表</MarkerContent>
      </Marker>
      <DeviceList />
    </CardDialogExtra>
  );
};

export const DeviceEditExtra = ({ extraState }: ExtraProps) => {
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
      title={<DialogExtraTitle text="相关设备">设置</DialogExtraTitle>}
      description={<p>编辑设备信息</p>}
      action={<DialogExtraClose />}
      footer={
        <DialogExtraAction
          type="submit"
          form="device-info-form"
          onClick={handleClose}
        >
          保存
        </DialogExtraAction>
      }
      footerClassName="justify-end gap-2"
    >
      <form id="device-info-form" onSubmit={handleSubmit}>
        <DeviceFormInput />
      </form>
    </CardDialogExtra>
  );
};

export const DeviceAddExtra = ({ extraState }: ExtraProps) => {
  const { type } = extraState;
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withExtraState({ type, open: false }));
  };

  return (
    <CardDialogExtra
      title={<DialogExtraTitle text="新设备">添加</DialogExtraTitle>}
      description={<p>扫码将新设备添加到设备组</p>}
      action={<DialogExtraClose />}
      footer={
        <>
          <PassportQRFileInputAction />
          <DialogExtraCloseAction variant="outline">
            关闭
          </DialogExtraCloseAction>
        </>
      }
      footerClassName="justify-between gap-2"
    >
      <PassportQRScan />
    </CardDialogExtra>
  );
};
