/**
 * App Node Page Definition
 */
import {
  BooleanField,
  BooleanInput,
  Create,
  DateTimeInput,
  Edit,
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
} from "react-admin";

import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Stack from "@mui/material/Stack";
import type { NodeQRCodeValue } from "./QRCode/Node";
import { NodeFileQRScan, NodeQRCode, NodeQRScan } from "./QRCode/Node";

import { InfinitePagination } from "./List/Infinite";

import CastIcon from "@mui/icons-material/Cast";
import DesktopAccessDisabledIcon from "@mui/icons-material/DesktopAccessDisabled";

import LanIcon from "@mui/icons-material/Lan";

import { useFormContext, useWatch } from "react-hook-form";

export const AppNodeIcon = LanIcon;
export const AppNodeRoutePath = "/app/nodes";

type AppNode = {
  id: number;
  peerId: string;
  name: string;
  blocked: boolean;
  online: boolean;
  createdAt: Date;
  updatedAt: Date;
};

/**
 * A filter component for app nodes that filters by their online and blocked
 * status.
 *
 * @returns A React element representing the filter component.
 */
const AppNodeStatFilter = () => {
  // const t = useTranslate();
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

/**
 * A component that renders the filters for the app nodes list.
 *
 * The component renders a Card containing the following components:
 * - SavedQueriesList: a list of saved queries
 * - FilterLiveSearch: a live search filter
 * - AppNodeStatFilter: a filter for the online and blocked status of the app
 *   nodes
 *
 * The component is displayed on the left of the list when the screen size is
 * sm (small) or larger, and is not displayed otherwise.
 *
 * @returns A React element representing the filter component.
 */
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

/**
 * Page Component of AppNodes
 *
 * The AppNodes component renders a list of application nodes with pagination
 * and filtering capabilities. It uses an InfiniteList to display the nodes
 * and includes filters for searching and filtering by node status.
 *
 * Each node is displayed using a SimpleList, showing the node's name, online
 * status, and last updated timestamp. The online status is represented with
 * icons indicating whether the node is online or offline.
 *
 * The component also provides an aside section for additional filters to refine
 * the list of displayed nodes, including saved queries and live search.
 *
 * @returns A React element representing the app nodes list.
 */

export const AppNodes = () => {
  // const t = useTranslate();
  return (
    <InfiniteList
      pagination={<InfinitePagination />}
      filters={AppNodeSimpleFilters}
      aside={<AppNodeFilters />}
    >
      <SimpleList<AppNode>
        // linkType="show"
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
        tertiaryText={(record) => new Date(record.updatedAt).toLocaleString()}
      />
    </InfiniteList>
  );
};

/**
 * AppNodeCreateActions
 *
 * A custom actions component for the app nodes create page that
 * renders a list button to return to the list of app nodes.
 *
 * @returns {ReactElement} A React element representing the actions
 */
const AppNodeCreateActions = () => {
  return (
    <TopToolbar>
      <ListButton />
    </TopToolbar>
  );
};

/**
 * Page components for creating new AppNodes
 *
 * A Create component for the app nodes resource that renders a form with
 * input fields for the name and peerId fields of the app node, and a boolean
 * input for the blocked field.
 *
 * The component also renders an AppNodeQRScan component, which is a QR code
 * scanner for the app node's peerId. The AppNodeQRScan component is displayed
 * on the left side of the form, and the input fields are displayed on the right
 * side.
 *
 * The component is wrapped in a Create component, which is a special type of
 * component provided by react-admin that renders a form with a submit button
 * to create a record in the app nodes resource.
 *
 * The component also renders a ListButton component, which is a button that
 * when clicked, redirects the user to the list of app nodes.
 *
 * @returns {ReactElement} A React element representing the create form for
 * the app nodes resource.
 */
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
        <AppNodeQRScan />
        <Stack spacing={1} minWidth={200} maxWidth={760} width={"70%"}>
          <TextInput source="name" fullWidth />
          <TextInput source="peerId" fullWidth rows={3} multiline />
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

/**
 * QRCode Component of AppNode
 *
 * A component that renders a QR code for the app node's peerId. It uses
 * the `useWatch` hook from `react-hook-form` to get the values of the
 * `name` and `peerId` fields from the form context.
 *
 * @returns {ReactElement} A React element representing the QR code.
 */
const AppNodeQRCode = () => {
  const name = useWatch({ name: "name" });
  const peerId = useWatch({ name: "peerId" });
  return <NodeQRCode name={name} peerId={peerId} />;
};

/**
 * QR Scan Component of AppNode
 *
 * A component that integrates QR code scanning functionality for app nodes.
 * It uses the `useFormContext` hook from `react-hook-form` to set the values
 * of the `peerId` and `name` fields in the form context when a QR code is
 * scanned.
 *
 * The component renders a stack layout containing the AppNodeQRCode component
 * and two QR scanning components: NodeQRScan and NodeFileQRScan. Both scanners
 * share the same `onQRScan` callback function, which updates the form fields
 * upon successful QR code scanning.
 *
 * @returns {ReactElement} A React element representing the QR scan interface
 * for app nodes.
 */

const AppNodeQRScan = () => {
  const { setValue } = useFormContext();

  const onQRScan = ({ peerId, name }: NodeQRCodeValue) => {
    setValue("peerId", peerId, { shouldValidate: true, shouldDirty: true });
    setValue("name", name, { shouldValidate: true, shouldDirty: true });
  };

  return (
    <Stack
      padding={1}
      spacing={2}
      alignItems="center"
      justifyContent={"space-between"}
    >
      <AppNodeQRCode />
      <Stack
        direction="row"
        spacing={1}
        alignItems="center"
        justifyContent={"space-between"}
      >
        <NodeQRScan onQRScan={onQRScan} />
        <NodeFileQRScan onQRScan={onQRScan} />
      </Stack>
    </Stack>
  );
};

/**
 * Page components for editing AppNodes
 *
 * The AppNodeEdit component renders a form for editing an application node.
 * It uses an Edit component from react-admin to wrap a SimpleForm that
 * displays the fields for the app node. The form is configured with
 * mutationMode="pessimistic" to prevent data loss when the user navigates
 * away from the form before submitting the changes.
 *
 * The component renders a stack layout containing the AppNodeQRCode
 * component and another stack layout with input fields for the app node.
 * The input fields include the name, peerId, blocked, online, createdAt,
 * and updatedAt fields. The peerId field is marked as read-only and the
 * online field is marked as read-only.
 *
 * @returns {ReactElement} A React element representing the edit form for
 * the app nodes resource.
 */
export const APPNodeEdit = () => (
  <Edit mutationMode="pessimistic">
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
          <TextInput source="peerId" fullWidth rows={3} multiline readOnly />
          <BooleanInput
            source="blocked"
            defaultValue={false}
            fullWidth
            margin="dense"
          />
          <BooleanInput source="online" readOnly />
          <DateTimeInput source="createdAt" readOnly />
          <DateTimeInput source="updatedAt" readOnly />
        </Stack>
      </Stack>
    </SimpleForm>
  </Edit>
);
