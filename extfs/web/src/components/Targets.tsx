import {
  BooleanInput,
  BulkDeleteButton,
  Datagrid,
  Edit,
  EditButton,
  List,
  SimpleForm,
  TextField,
  TextInput,
} from "react-admin";

const TargetBulkActionButtons = () => (
  <>
    <BulkDeleteButton />
  </>
);

export const Targets = () => (
  <List>
    <Datagrid bulkActionButtons={false}>
      <TextField source="id" />
      <TextField source="name" />
      <TextField source="filepath" />
      <TextField source="enabled" />
      <TextField source="createAt" />
      <TextField source="updateAt" />
      <EditButton />
    </Datagrid>
  </List>
);

export const TargetEdit = () => (
  <Edit>
    <SimpleForm>
      <TextInput source="id" readOnly={true} />
      <TextInput source="name" />
      <TextInput source="filepath" />
      <BooleanInput source="enabled" />
    </SimpleForm>
  </Edit>
);
