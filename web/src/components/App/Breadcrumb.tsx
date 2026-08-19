import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { useAppHeader } from "./Context";
import { Link } from "./Router";
import { Fragment, type ComponentProps } from "react";

const AppBreadcrumb = ({ ...props }: ComponentProps<typeof Breadcrumb>) => {
  const { title, breadcrumbs } = useAppHeader();
  return (
    <Breadcrumb {...props}>
      <BreadcrumbList>
        {breadcrumbs?.map(([text, href], index) => (
          <Fragment key={index}>
            <BreadcrumbItem className="hidden md:block">
              <BreadcrumbLink render={<Link to={href}>{text}</Link>} />
            </BreadcrumbItem>
            <BreadcrumbSeparator className="hidden md:block" />
          </Fragment>
        ))}
        <BreadcrumbItem>
          <BreadcrumbPage>{title}</BreadcrumbPage>
        </BreadcrumbItem>
      </BreadcrumbList>
    </Breadcrumb>
  );
};

export default AppBreadcrumb;
