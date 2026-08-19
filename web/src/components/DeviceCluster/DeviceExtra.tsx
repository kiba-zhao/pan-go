import { type ExtraProps } from "./ExtraBase";
import { DeviceList, DeviceListMarker, DeviceInfoForm } from "./Device";
import { ClusterInfo } from "./Cluster";
import type { FormEvent } from "react";

import { useAppDispatch } from "@/components/App/Context";
import {
  withDialogExtraState,
  DialogExtraTitle,
  DialogExtraClose,
  DialogExtraAction,
} from "@/components/App/Extra";
import {
  Dialog,
  DialogDescription,
  DialogAction,
} from "@/components/App/Dialog";

export const DevicesRemoveExtra = ({ extraState }: ExtraProps) => {
  const { type, open } = extraState;
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withDialogExtraState({ type, open: false }));
  };
  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtraClose onClose={handleClose} />
      <DialogExtraTitle text="相关设备" className="text-destructive">
        移除
      </DialogExtraTitle>
      <ClusterInfo />
      <DialogDescription>
        <p>
          确认将下列<span className="text-destructive">32</span>
          个设备从上面的设备组中移除吗？
        </p>
      </DialogDescription>
      <DeviceListMarker />
      <DeviceList />
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

export const DeviceEditExtra = ({ extraState }: ExtraProps) => {
  const { type, open } = extraState;
  const dispatch = useAppDispatch();

  const handleClose = () => {
    dispatch?.(withDialogExtraState({ type, open: false }));
  };

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    handleClose();
  };

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtraClose onClose={handleClose} />
      <DialogExtraTitle text="相关设备">设置</DialogExtraTitle>
      <DialogDescription>编辑设备信息</DialogDescription>
      <DeviceInfoForm id="device-info-form" onSubmit={handleSubmit} />
      <DialogAction>
        <DialogExtraAction
          onClose={handleClose}
          type="submit"
          form="device-info-form"
        >
          保存
        </DialogExtraAction>
      </DialogAction>
    </Dialog>
  );
};
