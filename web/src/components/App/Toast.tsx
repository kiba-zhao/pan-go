import { Toaster, toast } from "@/components/ui/toast";
export { Toaster, toast };

export enum ToastType {
  Default = "default",
  Success = "success",
  Error = "error",
  Warning = "warning",
  Info = "info",
  Loading = "loading",
}

type ToastManagerAddOptions<Data> = Parameters<typeof toast.add<Data>>[0];
type ToastManagerUpdateOptions<Data> = Parameters<typeof toast.update<Data>>[1];
export type ToastID = ReturnType<typeof toast.add>;
export function withToast<
  Data extends object = any,
  ToastOptions = ToastManagerAddOptions<Data> | ToastManagerUpdateOptions<Data>,
>(type: ToastType = ToastType.Default, opts?: Omit<ToastOptions, "type">) {
  return {
    type: type === ToastType.Default ? void 0 : type,
    priority: type === ToastType.Error ? "high" : "low",
    timeout: type === ToastType.Loading ? 0 : void 0,
    ...opts,
  } as ToastOptions;
}
