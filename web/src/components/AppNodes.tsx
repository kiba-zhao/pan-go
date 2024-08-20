import {
  BooleanField,
  BooleanInput,
  Create,
  FilterList,
  FilterListItem,
  FilterLiveSearch,
  InfiniteList,
  ListButton,
  SavedQueriesList,
  SearchInput,
  SimpleForm,
  SimpleList,
  TextField,
  TextInput,
  TopToolbar,
  useTranslate,
} from "react-admin";

import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Stack from "@mui/material/Stack";
import { NodeQRCode, NodeQRScan } from "./QRCode/Node";

import { InfinitePagination } from "./List/Infinite";

import CastIcon from "@mui/icons-material/Cast";
import DesktopAccessDisabledIcon from "@mui/icons-material/DesktopAccessDisabled";

import LanIcon from "@mui/icons-material/Lan";

import { useWatch } from "react-hook-form";

export const AppNodeIcon = LanIcon;

type AppNode = {
  id: number;
  nodeId: string;
  name: string;
  blocked: boolean;
  online: boolean;
  createAt: Date;
  updateAt: Date;
};

const AppNodeStatFilter = () => {
  const t = useTranslate();
  return (
    <FilterList
      label="resources.app/nodes.filters.stat"
      icon={<InfoOutlinedIcon />}
    >
      <FilterListItem
        label="resources.app/nodes.fields.online"
        value={{ online: true }}
      />
      <FilterListItem
        label="resources.app/nodes.fields.blocked"
        value={{ blocked: true }}
      />
    </FilterList>
  );
};

const AppNodeFilters = () => {
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
          <AppNodeStatFilter />
        </CardContent>
      </Card>
    </Box>
  );
};

const AppNodeSimpleFilters = [
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

export const AppNodes = () => {
  const t = useTranslate();
  return (
    <InfiniteList
      pagination={<InfinitePagination />}
      filters={AppNodeSimpleFilters}
      aside={<AppNodeFilters />}
    >
      <SimpleList<AppNode>
        linkType="show"
        primaryText={<TextField source="name" />}
        secondaryText={(record) => (
          <BooleanField
            source="online"
            TrueIcon={CastIcon}
            FalseIcon={DesktopAccessDisabledIcon}
            color={record.online ? "green" : void 0}
            valueLabelTrue="resources.app/nodes.fields.online"
            valueLabelFalse="resources.app/nodes.fields.offline"
          />
        )}
        tertiaryText={(record) => new Date(record.updateAt).toLocaleString()}
      />
    </InfiniteList>
  );
};

const AppNodeCreateActions = () => {
  return (
    <TopToolbar>
      <ListButton />
    </TopToolbar>
  );
};

export const AppNodeCreate = () => (
  <Create actions={<AppNodeCreateActions />}>
    <SimpleForm>
      <Stack
        direction="row"
        spacing={5}
        alignItems="flex-start"
        justifyContent="flex-start"
        useFlexGap
        flexWrap="wrap"
        width={"100%"}
      >
        <AppNodeQRCode />
        <Stack spacing={1} minWidth={200} maxWidth={760} width={"70%"}>
          <TextInput source="name" fullWidth />
          <TextInput source="nodeId" fullWidth rows={3} multiline />
          <BooleanInput
            source="blocked"
            defaultValue={false}
            fullWidth
            margin="dense"
          />
        </Stack>
      </Stack>
    </SimpleForm>
  </Create>
);

const AppNodeQRCode = () => {
  const t = useTranslate();

  const name = useWatch({ name: "name" });
  const nodeId = useWatch({ name: "nodeId" });

  return (
    <Stack
      padding={1}
      spacing={2}
      alignItems="center"
      justifyContent={"space-between"}
    >
      <NodeQRCode name={name} nodeId={nodeId} />
      <Stack
        direction="row"
        spacing={1}
        alignItems="center"
        justifyContent={"space-between"}
      >
        <NodeQRScan />
      </Stack>
    </Stack>
  );
};
