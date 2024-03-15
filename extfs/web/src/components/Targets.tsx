import {
  BooleanField,
  BooleanInput,
  BulkDeleteButton,
  BulkExportButton,
  CloneButton,
  Create,
  CreateButton,
  DatagridConfigurable,
  DateField,
  Edit,
  ExportButton,
  List,
  ListButton,
  SearchInput,
  SelectColumnsButton,
  Show,
  ShowButton,
  SimpleForm,
  SimpleShowLayout,
  TextField,
  TextInput,
  TopToolbar,
} from "react-admin";

const TargetBulkActions = () => (
  <>
    <BulkDeleteButton />
    <BulkExportButton />
  </>
);

const TargetListActions = () => (
  <TopToolbar>
    <SelectColumnsButton />
    <CreateButton />
    <ExportButton />
  </TopToolbar>
);

const TargetFilters = [<SearchInput source="q" alwaysOn />];

export const Targets = () => {
  return (
    <List actions={<TargetListActions />} filters={TargetFilters}>
      <DatagridConfigurable
        rowClick="edit"
        bulkActionButtons={<TargetBulkActions />}
      >
        <TextField source="id" />
        <TextField source="name" />
        <TextField source="filepath" />
        <BooleanField source="enabled" />
        <DateField source="createAt" showTime />
        <DateField source="updateAt" showTime />
      </DatagridConfigurable>
    </List>
  );
};

export const TargetCreate = () => (
  <Create>
    <SimpleForm>
      <TextInput source="name" />
      <TextInput source="filepath" />
      <BooleanInput source="enabled" />
    </SimpleForm>
  </Create>
);

const TargetEditActions = () => (
  <TopToolbar>
    <CreateButton />
    <CloneButton />
    <ShowButton />
    <ListButton />
  </TopToolbar>
);

export const TargetEdit = () => (
  <Edit actions={<TargetEditActions />}>
    <SimpleForm>
      <TextInput source="id" readOnly={true} />
      <TextInput source="name" />
      <TextInput source="filepath" />
      <BooleanInput source="enabled" />
      <DateField source="createAt" showTime />
      <DateField source="updateAt" showTime />
    </SimpleForm>
  </Edit>
);

export const TargetShow = () => (
  <Show>
    <SimpleShowLayout>
      <TextField source="id" />
      <TextField source="name" />
      <TextField source="filepath" />
      <BooleanField source="enabled" />
      <DateField source="createAt" showTime />
      <DateField source="updateAt" showTime />
    </SimpleShowLayout>
  </Show>
);
