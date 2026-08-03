import { SearchIcon } from "./Icon";
import { cn } from "@/lib/utils";

import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";

export const SearchFilter = () => {
  return (
    <form className="font-medium text-xl">
      <InputGroup>
        <InputGroupInput placeholder="Search..." />
        <InputGroupAddon>
          <SearchIcon />
        </InputGroupAddon>
      </InputGroup>
    </form>

    // <form className="flex flex-row w-3xl min-w-auto p-2 gap-1 rounded-4xl hover:bg-[#fff]/15">
    //   <SearchIcon />
    //   <input
    //     className="w-full outline-none font-medium"
    //     type="text"
    //     placeholder="Search Filters"
    //   />
    // </form>
  );
};
