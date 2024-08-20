import {
  BooleanField,
  BooleanInput,
  BulkExportButton,
  CloneButton,
  Create,
  CreateButton,
  DatagridConfigurable,
  DateField,
  DateTimeInput,
  Edit,
  EditButton,
  ExportButton,
  FilterList,
  FilterListItem,
  FilterLiveSearch,
  InfiniteList,
  ListButton,
  SavedQueriesList,
  SelectColumnsButton,
  Show,
  ShowButton,
  SimpleForm,
  SimpleShowLayout,
  TextField,
  TextInput,
  TopToolbar,
  WrapperField,
  useResourceContext,
  useTranslate,
} from "react-admin";

import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";
import ToggleOn from "@mui/icons-material/ToggleOn";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Tooltip from "@mui/material/Tooltip";

import { Controller, useFormContext } from "react-hook-form";
import { FilePathInput } from "./FilePath/Input";
import { InfinitePagination } from "./List/Infinite";

import type { TextInputProps } from "react-admin";

const TargetBulkActions = () => (
  <>
    <BulkExportButton />
  </>
);

const TargetListActions = () => (
  <TopToolbar>
    <SelectColumnsButton preferenceKey="targets.datagrid" />
    <CreateButton />
    <ExportButton />
  </TopToolbar>
);

const TargetEnabledFilter = () => {
  return (
    <FilterList
      label="resources.extfs/targets.filters.has_enabled"
      icon={<ToggleOn />}
    >
      <FilterListItem
        label="resources.extfs/targets.filters.enabled"
        value={{ enabled: true }}
      />
      <FilterListItem
        label="resources.extfs/targets.filters.disabled"
        value={{ enabled: false }}
      />
    </FilterList>
  );
};

const TargetInvalidFilter = () => {
  const t = useTranslate();
  return (
    <FilterList
      label="resources.extfs/targets.filters.has_available"
      icon={
        <Tooltip title={t("resources.extfs/targets.filters.help_available")}>
          <InfoOutlinedIcon />
        </Tooltip>
      }
    >
      <FilterListItem
        label="resources.extfs/targets.filters.available"
        value={{ available: true }}
      />
      <FilterListItem
        label="resources.extfs/targets.filters.not_available"
        value={{ available: false }}
      />
    </FilterList>
  );
};

const TargetFilters = () => {
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
      <Card sx={{ mr: 2, mt: 6, width: 200 }}>
        <CardContent>
          <SavedQueriesList />
          <FilterLiveSearch />
          <TargetEnabledFilter />
          <TargetInvalidFilter />
        </CardContent>
      </Card>
    </Box>
  );
};

export const Targets = () => {
  return (
    <InfiniteList
      actions={<TargetListActions />}
      aside={<TargetFilters />}
      pagination={<InfinitePagination />}
    >
      <DatagridConfigurable
        bulkActionButtons={<TargetBulkActions />}
        preferenceKey="targets.datagrid"
      >
        <TextField source="name" />
        <TextField source="filepath" />
        <BooleanField source="enabled" />
        <BooleanField source="available" />
        <DateField source="createAt" showTime />
        <WrapperField label="custom.table.actions">
          <EditButton />
          <ShowButton />
        </WrapperField>
      </DatagridConfigurable>
    </InfiniteList>
  );
};

export const TargetCreate = () => (
  <Create>
    <SimpleForm>
      <TextInput source="name" />
      <RAFilePathInput source="filepath" />
      <BooleanInput source="enabled" defaultValue={true} />
    </SimpleForm>
  </Create>
);

const RAFilePathInput = ({ source }: TextInputProps) => {
  const t = useTranslate();
  const resource = useResourceContext();
  const { control } = useFormContext();
  return (
    <Controller
      control={control}
      name={source}
      render={({ field }) => (
        <FilePathInput
          title={t(`resources.${resource}.input.${source}`)}
          label={t(`resources.${resource}.fields.${source}`, { _: source })}
          {...field}
        />
      )}
    />
  );
};

const TargetEditActions = () => (
  <TopToolbar>
    <CreateButton />
    <CloneButton />
    <ListButton />
    <ShowButton />
  </TopToolbar>
);

export const TargetEdit = () => (
  <Edit actions={<TargetEditActions />} mutationMode="pessimistic">
    <SimpleForm>
      <TextInput source="id" readOnly={true} />
      <TextInput source="name" />
      <RAFilePathInput source="filepath" />
      <BooleanInput source="enabled" />
      <BooleanInput source="available" disabled={true} />
      <DateTimeInput source="createAt" disabled={true} />
      <DateTimeInput source="updatedAt" disabled={true} />
    </SimpleForm>
  </Edit>
);

const TargetShowActions = () => (
  <TopToolbar>
    <EditButton />
    <ListButton />
  </TopToolbar>
);

export const TargetShow = () => (
  <Show actions={<TargetShowActions />}>
    <SimpleShowLayout>
      <TextField source="id" />
      <TextField source="name" />
      <TextField source="filepath" />
      <BooleanField source="enabled" />
      <BooleanField source="available" />
      <DateField source="createAt" showTime />
      <DateField source="updatedAt" showTime />
    </SimpleShowLayout>
  </Show>
);
