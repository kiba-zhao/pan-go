import { useId, useMemo } from "react";

import FolderIcon from "@mui/icons-material/Folder";
import FilledInput from "@mui/material/FilledInput";
import FormControl from "@mui/material/FormControl";
import IconButton from "@mui/material/IconButton";
import InputAdornment from "@mui/material/InputAdornment";
import InputLabel from "@mui/material/InputLabel";

type InputProps<T = any> = {
  label: string;
} & T;

export type FilePathInputProps = InputProps<{}>;
export const FilePathInput = ({ label }: FilePathInputProps) => {
  const id = useId();
  const elementId = useMemo(() => `custom-filepath-input-${id}`, [id]);
  return (
    <FormControl variant="filled" fullWidth>
      <InputLabel htmlFor={elementId}>{label}</InputLabel>
      <FilledInput
        id={elementId}
        endAdornment={
          <InputAdornment position="end">
            <IconButton aria-label="directions">
              <FolderIcon />
            </IconButton>
          </InputAdornment>
        }
      />
    </FormControl>
  );
};
