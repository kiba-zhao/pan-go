import {DeviceFields} from '@pango/data';
import Clipboard from '@react-native-clipboard/clipboard';
import type {Dispatch, PropsWithChildren} from 'react';
import {createContext, Fragment, useContext, useReducer} from 'react';
import {Button} from 'react-native';
import {CommonActions, useNavigation, useRoute} from '../App/Navigation';
import Alert from '../Common/Alert';
import {SwitchFieldItem, TextFieldItem, TextFieldModal} from '../Common/Field';
import Paper from '../Common/Paper';
import {useMutation, useQueryClient} from '../Common/ReactQuery';
import {ScreenLayout} from '../Common/ScreenBase';
import {
  default as Section,
  SectionHeader,
  SectionItem,
} from '../Common/Section';
import Text from '../Common/Text';
import {createDevice} from '../Spec/Device';
import {invalidateListQueryCache} from './ReactQuery';
import {CreatorParam, DeviceEditorScreenName} from './ScreenRoute';

type DeviceScreenState = {
  fields?: DeviceFields;
  nameModalVisible?: boolean;
};
type DeviceScreenStateAction = DeviceScreenState;

const reducer = (
  state: DeviceScreenState,
  action: DeviceScreenStateAction,
) => ({
  ...state,
  ...action,
});

const Context = createContext<DeviceScreenState>({});
const DispatchContext = createContext<Dispatch<DeviceScreenStateAction> | null>(
  null,
);

const DeviceScreenStateProvider = ({
  fields,
  children,
}: PropsWithChildren<Pick<DeviceScreenState, 'fields'>>) => {
  const [state, dispatch] = useReducer(reducer, {fields});
  return (
    <Context.Provider value={state}>
      <DispatchContext.Provider value={dispatch}>
        {children}
      </DispatchContext.Provider>
    </Context.Provider>
  );
};

const DeviceCreatorScreen = () => {
  const route = useRoute();
  const {name, peerId} = route.params as CreatorParam;

  return (
    <ScreenLayout>
      <DeviceScreenStateProvider fields={{name, peerId, enabled: true}}>
        <DeviceFieldsSection />
        <DeviceActionSection />
        <DeviceNameModal />
      </DeviceScreenStateProvider>
    </ScreenLayout>
  );
};

export default DeviceCreatorScreen;

const DeviceFieldsSection = () => {
  const {fields} = useContext(Context);
  const dispatch = useContext(DispatchContext);
  const handlePressName = () => {
    dispatch?.({nameModalVisible: true});
  };

  const handlePressPeerID = () => {
    if (fields === void 0) return;
    Clipboard.setString(fields.peerId);
  };

  const handleEnabledChange = (enabled: boolean) => {
    dispatch?.({fields: {...fields, enabled} as DeviceFields});
  };
  return (
    <Fragment>
      <Section>
        <SectionItem variant="row-start" onPress={handlePressName}>
          <TextFieldItem label="Name" text={fields?.name} editable />
        </SectionItem>
        <SectionItem variant="row-end">
          <SwitchFieldItem
            label="Enabled"
            value={fields?.enabled}
            onValueChange={handleEnabledChange}
          />
        </SectionItem>
      </Section>
      <Section>
        <SectionHeader>
          <Text size="small">Peer ID</Text>
        </SectionHeader>
        <Paper padding={1} borderColor="divider">
          <Text
            size="small"
            color="textSecondary"
            lines={5}
            onPress={handlePressPeerID}>
            {fields?.peerId || ''}
          </Text>
        </Paper>
      </Section>
    </Fragment>
  );
};

const DeviceActionSection = () => {
  const {fields} = useContext(Context);

  const navigation = useNavigation();
  const queryClient = useQueryClient();
  const {mutate} = useMutation({
    mutationFn: createDevice,
    onSuccess: data => {
      invalidateListQueryCache(queryClient, data);

      navigation.dispatch(state => {
        const routes = state.routes.slice(0, -2);
        return CommonActions.reset({
          ...state,
          routes: [
            ...routes,
            {name: DeviceEditorScreenName, params: {id: data.id}},
          ],
          index: routes.length,
        });
      });
    },
    onError: error => {
      Alert.alert(error.name, error.message);
    },
  });
  const handleSubmit = () => {
    if (fields === void 0) return;
    mutate(fields);
  };

  return (
    <Paper bgColor="transparent" padding={[3, 0, 1, 0]}>
      <Button onPress={handleSubmit} title="Save" />
    </Paper>
  );
};

const DeviceNameModal = () => {
  const {fields, nameModalVisible} = useContext(Context);
  const dispatch = useContext(DispatchContext);

  const handleClose = () => {
    dispatch?.({nameModalVisible: false});
  };

  const handleSubmit = (name?: string) => {
    dispatch?.({
      fields: {...fields, name} as DeviceFields,
      nameModalVisible: false,
    });
  };

  return (
    <TextFieldModal
      title="Device Name"
      submitLabel="修改"
      visible={!!nameModalVisible}
      value={fields?.name}
      onClose={handleClose}
      onSubmit={handleSubmit}
    />
  );
};
