import { useAppExtra } from "./Context";
import { type UseFormStateProps, type FieldValues, useFormState } from "./Form";
import { DialogExtraAction } from "./Dialog";
import type { ComponentType, ComponentProps } from "react";
import type { RequiredType } from "@/lib/utility_types";

export type AppExtraState<Type> = {
  type?: Type;
};

export type AppExtraBaseProps<Type, State extends AppExtraState<Type>> = {
  extraState: State;
};

type AppExtraProps<Type, Props> = {
  asProps?: Props;
  as: ComponentType<Props>;
  extraType: Type;
};
export const AppExtra = <
  Type,
  State extends AppExtraState<Type>,
  Props extends AppExtraBaseProps<Type, State>,
>({
  as,
  asProps: props,
  extraType,
}: AppExtraProps<Type, Props>) => {
  const extraState = useAppExtra<State>();
  const Component = as;
  if (extraState.type !== extraType) {
    return null;
  }
  return <Component {...(props || ({} as Props))} extraState={extraState} />;
};

type DialogExtraFormActionProps<T extends FieldValues> = {
  control: NonNullable<UseFormStateProps<T>["control"]>;
} & RequiredType<ComponentProps<typeof DialogExtraAction>, "form", "type">;
export const DialogExtraFormAction = <T extends FieldValues>({
  control,
  children,
  disabled,
  progress,
  ...props
}: DialogExtraFormActionProps<T>) => {
  const {
    disabled: formDisabled,
    isValid,
    isDirty,
    isSubmitting,
  } = useFormState<T>({ control });
  return (
    <DialogExtraAction
      {...props}
      type="submit"
      disabled={disabled || formDisabled || !isValid || !isDirty}
      progress={progress || isSubmitting}
    >
      {children}
    </DialogExtraAction>
  );
};
