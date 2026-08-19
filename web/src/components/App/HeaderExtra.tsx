import { Link } from "./Router";
import { useTranslation, useAppI18n, I18nVariant } from "./I18Next";
import type { ComponentProps } from "react";
import { PartialType } from "@/lib/utility_types";

type GoBackProps = PartialType<ComponentProps<typeof Link>, "to">;
export const GoBack = ({
  className = "text-sm",
  to,
  children,
  ...props
}: GoBackProps) => {
  const { namespace } = useAppI18n();
  const { t } = useTranslation(namespace);

  return (
    <Link className={className} {...props} to={to || "/"}>
      {children || t(`${I18nVariant.HeaderExtra}.back`)}
    </Link>
  );
};
