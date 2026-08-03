import { useAppExtra, useAppDispatch } from "@/components/App/Context";
import {
  Dialog,
  DialogTitle,
  DialogDescription,
  DialogAction,
  DialogExtra,
} from "@/components/App/Dialog";
import { X as CloseIcon, CircleX, CircleCheckBig } from "./Icon";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
  InputGroupText,
  InputGroupButton,
  InputGroupTextarea,
} from "@/components/ui/input-group";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ExtraState as AppExtraState } from "@/components/App/Extra";
import { useRef, useEffect } from "react";
import { useIMask } from "react-imask";

const SettingsExtra = () => (
  <>
    <MemoEditExtra />
    <NameEditExtra />
    <PeerPortEditExtra />
    <BroadcastAddrsEditExtra />
    <PublicAddrsEditExtra />
    <WebPortEditExtra />
  </>
);

export default SettingsExtra;

export enum ExtraType {
  NameEdit,
  MemoEdit,
  PeerPortEdit,
  BroadcastAddrsEdit,
  PublicAddrsEdit,
  WebPortEdit,
}
export type ExtraState<Type extends ExtraType> = AppExtraState<
  Type,
  { open?: boolean }
>;

export const switchExtra = (
  type?: ExtraType,
  open?: boolean,
): undefined | ExtraState<ExtraType> => {
  return type === void 0 ? type : { type, open: open !== false };
};

const NameEditExtra = () => {
  const { type, open } = useAppExtra<ExtraState<ExtraType.NameEdit>>({});
  const dispatch = useAppDispatch();
  const inputRef = useRef<HTMLInputElement>(null);

  const handleClose = () => {
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  useEffect(() => {
    if (type === ExtraType.NameEdit && open) {
      inputRef.current?.focus();
    }
  }, [type, open, inputRef]);

  if (type !== ExtraType.NameEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        备注名
        <small className="text-muted-foreground">&nbsp;设置</small>
      </DialogTitle>
      <DialogDescription>
        备注名用作标识当前设备,请输入易于辨别的名称。支持最多32个字符。
      </DialogDescription>
      <form>
        <InputGroup>
          <InputGroupInput
            type="text"
            name="name"
            placeholder="请输入备注名"
            maxLength={32}
            required={true}
            ref={inputRef}
          />
          <InputGroupAddon align="block-end">
            <InputGroupText className="text-xs text-muted-foreground">
              2 / 32
            </InputGroupText>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-sm:hidden"
            >
              重置
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-sm:hidden"
          onClick={handleClose}
          variant="outline"
        >
          取消
        </Button>
        <Button disabled>保存</Button>
      </DialogAction>
    </Dialog>
  );
};

const MemoEditExtra = () => {
  const { type, open } = useAppExtra<ExtraState<ExtraType.MemoEdit>>({});
  const dispatch = useAppDispatch();
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const handleClose = () => {
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  useEffect(() => {
    if (type === ExtraType.MemoEdit && open) {
      inputRef.current?.focus();
    }
  }, [type, open, inputRef]);

  if (type !== ExtraType.MemoEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        我的备忘
        <small className="text-muted-foreground">&nbsp;设置</small>
      </DialogTitle>
      <DialogDescription>
        我的备忘用于记录设备相关的信息,支持最多255个字符。
      </DialogDescription>
      <form>
        <InputGroup>
          <InputGroupTextarea
            placeholder="请输入备忘"
            name="memo"
            className="min-h-29 max-h-29 scrollbar-thin"
            maxLength={255}
            required={true}
            ref={inputRef}
          />
          <InputGroupAddon align="block-end">
            <InputGroupText className="text-xs text-muted-foreground">
              2 / 255
            </InputGroupText>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-sm:hidden"
            >
              重置
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-sm:hidden"
          onClick={handleClose}
          variant="outline"
        >
          取消
        </Button>
        <Button disabled>保存</Button>
      </DialogAction>
    </Dialog>
  );
};

const PeerPortEditExtra = () => {
  const { type, open } = useAppExtra<ExtraState<ExtraType.PeerPortEdit>>({});
  const dispatch = useAppDispatch();

  const { ref: inputRef } = useIMask<HTMLInputElement>({
    mask: Number,
    min: 0,
    max: 65535,
    scale: 0,
    autofix: true,
  });

  const handleClose = () => {
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  useEffect(() => {
    if (type === ExtraType.PeerPortEdit && open) {
      inputRef.current?.focus();
    }
  }, [type, open, inputRef]);

  if (type !== ExtraType.PeerPortEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        网络端口
        <small className="text-muted-foreground">&nbsp;设置</small>
      </DialogTitle>
      <DialogDescription>
        网络端口用于设备之间的通信,请输入一个可用的端口号。
      </DialogDescription>
      <form>
        <InputGroup>
          <InputGroupInput
            type="text"
            name="peerPort"
            placeholder="请输入网络端口"
            required={true}
            ref={inputRef}
            defaultValue="9000"
          />
          <InputGroupAddon align="block-end">
            <InputGroupText className="text-xs text-muted-foreground">
              0 ~ 65535
            </InputGroupText>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-sm:hidden"
            >
              重置
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-sm:hidden"
          onClick={handleClose}
          variant="outline"
        >
          取消
        </Button>
        <Button disabled>保存</Button>
      </DialogAction>
    </Dialog>
  );
};

const BroadcastAddrsEditExtra = () => {
  const { type, open } = useAppExtra<ExtraState<ExtraType.BroadcastAddrsEdit>>(
    {},
  );
  const dispatch = useAppDispatch();
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const handleClose = () => {
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  useEffect(() => {
    if (type === ExtraType.BroadcastAddrsEdit && open) {
      inputRef.current?.focus();
    }
  }, [type, open, inputRef]);

  if (type !== ExtraType.BroadcastAddrsEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        广播地址
        <small className="text-muted-foreground">&nbsp;设置</small>
      </DialogTitle>
      <div>
        <DialogDescription>广播地址用于访问设备的网络服务。</DialogDescription>
        <DialogDescription>
          支持域名/主机名/IP地址 + 端口，每行一个地址。
        </DialogDescription>
      </div>
      <form>
        <InputGroup>
          <InputGroupTextarea
            placeholder="224.0.0.1:9001"
            name="broadcastAddrs"
            className="min-h-29 max-h-29 scrollbar-thin"
            maxLength={255}
            required={true}
            ref={inputRef}
          />
          <InputGroupAddon align="block-end">
            <Badge variant="secondary">
              <CircleCheckBig data-icon="inline-start" /> 5
            </Badge>
            <Badge variant="destructive">
              <CircleX data-icon="inline-start" /> 3
            </Badge>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-sm:hidden"
            >
              重置
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-sm:hidden"
          onClick={handleClose}
          variant="outline"
        >
          取消
        </Button>
        <Button disabled>保存</Button>
      </DialogAction>
    </Dialog>
  );
};

const PublicAddrsEditExtra = () => {
  const { type, open } = useAppExtra<ExtraState<ExtraType.PublicAddrsEdit>>({});
  const dispatch = useAppDispatch();
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const handleClose = () => {
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  useEffect(() => {
    if (type === ExtraType.PublicAddrsEdit && open) {
      inputRef.current?.focus();
    }
  }, [type, open, inputRef]);

  if (type !== ExtraType.PublicAddrsEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        公共地址
        <small className="text-muted-foreground">&nbsp;设置</small>
      </DialogTitle>
      <div>
        <DialogDescription>公共地址用于访问设备的网络服务。</DialogDescription>
        <DialogDescription>
          支持域名/主机名/IP地址 + 端口，每行一个地址。
        </DialogDescription>
      </div>
      <form>
        <InputGroup>
          <InputGroupTextarea
            placeholder="www.example.com:9000"
            name="publicAddrs"
            className="min-h-29 max-h-29 scrollbar-thin"
            maxLength={255}
            required={true}
            ref={inputRef}
          />
          <InputGroupAddon align="block-end">
            <Badge variant="secondary">
              <CircleCheckBig data-icon="inline-start" /> 5
            </Badge>
            <Badge variant="destructive">
              <CircleX data-icon="inline-start" /> 3
            </Badge>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-sm:hidden"
            >
              重置
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-sm:hidden"
          onClick={handleClose}
          variant="outline"
        >
          取消
        </Button>
        <Button disabled>保存</Button>
      </DialogAction>
    </Dialog>
  );
};

const WebPortEditExtra = () => {
  const { type, open } = useAppExtra<ExtraState<ExtraType.WebPortEdit>>({});
  const dispatch = useAppDispatch();

  const { ref: inputRef } = useIMask<HTMLInputElement>({
    mask: Number,
    min: 0,
    max: 65535,
    scale: 0,
    autofix: true,
  });

  const handleClose = () => {
    dispatch?.({
      extra: switchExtra(type, false),
    });
  };

  useEffect(() => {
    if (type === ExtraType.WebPortEdit && open) {
      inputRef.current?.focus();
    }
  }, [type, open, inputRef]);

  if (type !== ExtraType.WebPortEdit) {
    return null;
  }

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogExtra>
        <Button size="icon-sm" variant="ghost" onClick={handleClose}>
          <CloseIcon />
        </Button>
      </DialogExtra>
      <DialogTitle>
        Web端口
        <small className="text-muted-foreground">&nbsp;设置</small>
      </DialogTitle>
      <DialogDescription>
        提供Web访问的端口号,请输入一个可用的端口。
      </DialogDescription>
      <form>
        <InputGroup>
          <InputGroupInput
            type="text"
            name="webPort"
            placeholder="请输入网页服务端口"
            required={true}
            ref={inputRef}
            defaultValue="9000"
          />
          <InputGroupAddon align="block-end">
            <InputGroupText className="text-xs text-muted-foreground">
              0 ~ 65535
            </InputGroupText>
            <InputGroupButton
              variant="default"
              size="sm"
              className="ml-auto max-sm:hidden"
            >
              重置
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
      <DialogAction className="bg-muted">
        <Button
          className="max-sm:hidden"
          onClick={handleClose}
          variant="outline"
        >
          取消
        </Button>
        <Button disabled>保存</Button>
      </DialogAction>
    </Dialog>
  );
};
