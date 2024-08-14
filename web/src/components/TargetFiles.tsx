import CheckCircleOutlineOutlinedIcon from "@mui/icons-material/CheckCircleOutlineOutlined";
import DoDisturbOnOutlinedIcon from "@mui/icons-material/DoDisturbOnOutlined";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Tooltip from "@mui/material/Tooltip";

import {
  BooleanField,
  DateField,
  FilterList,
  FilterListItem,
  FilterLiveSearch,
  FunctionField,
  InfiniteList,
  ReferenceField,
  SavedQueriesList,
  SearchInput,
  Show,
  SimpleList,
  TabbedShowLayout,
  TextField,
  useTranslate,
} from "react-admin";

import { formatBytes } from "../lib/byte";
import { basename } from "../lib/path";
import { InfinitePagination } from "./custom/List";

export const TargetFileIcon = InsertDriveFileIcon;

type TargetFile = {
  id: number;
  targetId: number;
  filepath: string;
  mimeType: string;
  size: number;
  checkSum: string;
  modTime: Date;
  available: boolean;
  createAt: Date;
  updateAt: Date;
};

const TargetFileInvalidFilter = () => {
  const t = useTranslate();
  return (
    <FilterList
      label="resources.extfs/target-files.filters.has_available"
      icon={
        <Tooltip
          title={t("resources.extfs/target-files.filters.help_available")}
        >
          <InfoOutlinedIcon />
        </Tooltip>
      }
    >
      <FilterListItem
        label="resources.extfs/target-files.filters.available"
        value={{ available: true }}
      />
      <FilterListItem
        label="resources.extfs/target-files.filters.not_available"
        value={{ available: false }}
      />
    </FilterList>
  );
};

const TargetFileFilters = () => {
  return (
    <Box
      sx={{
        display: {
          xs: "none",
          sm: "block",
        },
        order: -1, // display on the left rather than on the right of the list
      }}
    >
      <Card sx={{ mr: 2, mt: 8, width: 200 }}>
        <CardContent>
          <SavedQueriesList />
          <FilterLiveSearch />
          <TargetFileInvalidFilter />
        </CardContent>
      </Card>
    </Box>
  );
};

const TargetFileSimpleFilters = [
  <SearchInput
    sx={{
      display: {
        xs: "block",
        sm: "none",
      },
    }}
    source="q"
    alwaysOn
  />,
];

export const TargetFiles = () => {
  return (
    <InfiniteList
      aside={<TargetFileFilters />}
      filters={TargetFileSimpleFilters}
      pagination={<InfinitePagination />}
    >
      <SimpleList<TargetFile>
        linkType="show"
        primaryText={(record) => basename(record.filepath)}
        secondaryText={(record) => formatBytes(record.size)}
        tertiaryText={(record) => new Date(record.modTime).toLocaleString()}
        leftIcon={(record) =>
          record.available ? (
            <CheckCircleOutlineOutlinedIcon color="success" />
          ) : (
            <DoDisturbOnOutlinedIcon color="error" />
          )
        }
      />
    </InfiniteList>
  );
};

export const TargetFileShow = () => {
  return (
    <Show>
      <TabbedShowLayout>
        <TabbedShowLayout.Tab label="resources.extfs/target-files.show.summary">
          <ReferenceField
            source="targetId"
            reference="extfs/targets"
            label="resources.extfs/target-files.fields.target"
          >
            <TextField source="name" />
          </ReferenceField>
          <TextField source="filepath" />
          <TextField source="mimeType" />
          <FunctionField<TargetFile>
            source="size"
            render={(record) => formatBytes(record.size)}
          />
          <TextField source="checkSum" />
          <DateField source="modTime" showTime />
          <BooleanField source="available" />
          <DateField source="createAt" showTime />
          <DateField source="updateAt" showTime />
        </TabbedShowLayout.Tab>
      </TabbedShowLayout>
    </Show>
  );
};
