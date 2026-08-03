import {
  NavList,
  NavListItem,
  NavListItemLink,
  NavListItemText,
  NavListItemAvatar,
  type NavigationProps,
} from "@/components/App/Navigation";
import { useTranslation } from "@/components/App/I18Next";

import { ClustersName } from "./meta";
import {
  default as ClusterOutlinedIcon,
  HelpOutlineIcon,
  PassportOutlineIcon,
} from "./Icon";

const ClustersNavigation = ({ className }: NavigationProps) => {
  // TODO: Chat Navigation Sample
  const { t } = useTranslation();
  return (
    <NavList className={className}>
      <NavListItem>
        <NavListItemLink to="#">
          <NavListItemAvatar as={ClusterOutlinedIcon} />
          <NavListItemText>Clusters Overview</NavListItemText>
        </NavListItemLink>
      </NavListItem>
      <NavListItem>
        <NavListItemLink to="#">
          <NavListItemAvatar as={PassportOutlineIcon} />
          <NavListItemText>Passports</NavListItemText>
        </NavListItemLink>
      </NavListItem>
      <NavListItem>
        <NavListItemLink to="#">
          <NavListItemAvatar as={HelpOutlineIcon} />
          <NavListItemText>Help</NavListItemText>
        </NavListItemLink>
      </NavListItem>
      <NavListItem>
        <NavListItemLink to="#">
          <NavListItemText>Test Access</NavListItemText>
        </NavListItemLink>
      </NavListItem>
    </NavList>
  );
};

export default ClustersNavigation;
