import { SearchIcon } from "./Icon";

import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";
import { useTranslation } from "@/components/App/I18Next";
import { SettingsName, MainI18nPrefix } from "./meta";

export const SearchFilter = () => {
  const { t } = useTranslation(SettingsName);
  return (
    <InputGroup>
      <InputGroupInput
        placeholder={t("search.placeholder", { keyPrefix: MainI18nPrefix })}
      />
      <InputGroupAddon>
        <SearchIcon />
      </InputGroupAddon>
    </InputGroup>
  );
};
