import {
  FieldValues,
  UseControllerProps,
  RegisterOptions,
} from "@/components/App/Form";
import { RequiredType } from "@/lib/utility_types";
export type FormInputProps<T extends FieldValues> = Pick<
  RegisterOptions<T>,
  "onChange"
> &
  RequiredType<
    UseControllerProps<T>,
    "name" | "control" | "defaultValue",
    "rules"
  >;
