import type { ComponentProps, ElementType, FormEvent } from "react";
import { Fragment, useState } from "react";
import type { FieldValues, UseFormHandleSubmit } from "react-hook-form";
import { useFormContext } from "react-hook-form";

import { default as MuiBox } from "@mui/material/Box";
import { default as MuiButton } from "@mui/material/Button";

type MuiFormProps = ComponentProps<typeof MuiBox>;
export type BaseFormProps = Pick<
  ComponentProps<"form">,
  "id" | "children" | "onSubmit"
>;
type FormElementType = BaseFormProps | MuiFormProps;
type FormElementProps<T extends FormElementType> = T extends BaseFormProps
  ? Omit<T, "onSubmit"> & { BaseForm: ElementType<T> }
  : Omit<MuiFormProps, "onSubmit" | "component"> & { BaseForm: never };

export const Form = <T extends FieldValues, E extends FormElementType>({
  onValid,
  onInvalid,
  children,
  BaseForm,
  ...props
}: {
  onValid: Parameters<UseFormHandleSubmit<T>>[0];
  onInvalid?: Parameters<UseFormHandleSubmit<T>>[1];
} & FormElementProps<E>) => {
  const { handleSubmit } = useFormContext<T>() || {};

  const handleFormSubmit = async (event: FormEvent) => {
    event.preventDefault();
    await handleSubmit(onValid, onInvalid);
  };

  if (BaseForm === void 0)
    return (
      <MuiBox component="form" onSubmit={handleFormSubmit} {...props}>
        {children}
      </MuiBox>
    );

  return (
    <BaseForm {...props} onSubmit={handleFormSubmit}>
      {children}
    </BaseForm>
  );
};

type MuiButtonProps = ComponentProps<typeof MuiButton>;
export type BaseButtonProps = Pick<
  ComponentProps<"button">,
  "onClick" | "disabled" | "children"
>;
type ButtonElementType = BaseButtonProps | MuiButtonProps;
type ButtonProps<T extends ButtonElementType> = T extends BaseButtonProps
  ? T & { BaseButton: ElementType<T> }
  : Omit<MuiButtonProps, "type" | "form"> & { BaseButton: never };

type BaseButtonEvent<T extends ButtonElementType> = Parameters<
  NonNullable<T["onClick"]>
>[0];

export const Button = <T extends ButtonElementType>({
  onClick,
  children,
  BaseButton = MuiButton,
  ...props
}: ButtonProps<T>) => {
  const [event, setEvent] = useState<BaseButtonEvent<T> | null>(null);

  const handleClick = (event: BaseButtonEvent<T>) => {
    setEvent(event);
  };

  const handleClose = () => {
    setEvent(null);
  };

  const handleConfirm = () => {
    if (!event) return;
    onClick?.(event);
    handleClose();
  };

  return (
    <Fragment>
      <BaseButton onClick={handleClick} {...props}>
        {children}
      </BaseButton>
    </Fragment>
  );
};

// type SubmitFormType = ComponentProps<"button">["form"];

// export const Submit = <T extends ButtonElementType>({
//   form,
//   BaseButton,
//   ...props
// }: {
//   form: SubmitFormType;
// } & ButtonProps<T>) => {
//   const ref = useRef<HTMLButtonElement>(null);
//   const handleClick = () => {
//     ref.current?.click();
//   };

//   const SubmitButton = BaseButton === void 0 ? Button<T> : BaseButton;

//   return (
//     <Fragment>
//       <SubmitButton onClick={(_) => handleClick} {...props} />
//       <button type="submit" form={form} ref={ref} hidden={true} />
//     </Fragment>
//   );
// };
