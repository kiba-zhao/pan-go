import {Device} from '@pango/data';
import type {Dispatch, PropsWithChildren} from 'react';
import {
  createContext,
  Fragment,
  useContext,
  useEffect,
  useMemo,
  useReducer,
  useState,
} from 'react';
import {useNavigation, useRoute} from '../App/Navigation';
import Box from '../Common/Box';
import Icon from '../Common/Icon';
import {DeviceQRCode, DeviceQRCodeModal} from './QRCode';

import {
  useIsFetching,
  useMutation,
  useQuery,
  useQueryClient,
} from '../Common/ReactQuery';
import Text from '../Common/Text';
import {destroyDevice, selectDevice, updateDevice} from '../Spec/Device';
import {QueryKey} from './ReactQuery';
import {EditorParam} from './ScreenRoute';

import Clipboard from '@react-native-clipboard/clipboard';
import {HeaderActions, ScreenLayout} from '../App/ScreenBase';
import Alert from '../Common/Alert';
import {
  IconFieldItem,
  SwitchFieldItem,
  TextFieldItem,
  TextFieldModal,
} from '../Common/Field';
import Paper from '../Common/Paper';
import Pressable from '../Common/Pressable';
import {
  ScreenLoading,
  ScreenRefreshControl,
  ScreenSafetyFooter,
} from '../Common/ScreenBase';
import {
  default as Section,
  SectionHeader,
  SectionItem,
  SectionList,
} from '../Common/Section';
import {invalidateListQueryCache} from './ReactQuery';
import {DeviceSearchHeaderAction} from './ScreenBase';

/**
 *  define context for editor screen
 */
type DeviceScreenState = {
  nameModalVisible?: boolean;
  qrCodeModalVisible?: boolean;
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

const DeviceScreenStateProvider = ({children}: PropsWithChildren<{}>) => {
  const [state, dispatch] = useReducer(reducer, {});
  return (
    <Context.Provider value={state}>
      <DispatchContext.Provider value={dispatch}>
        {children}
      </DispatchContext.Provider>
    </Context.Provider>
  );
};
// end define context

const DeviceEditorScreen = () => {
  const route = useRoute();
  const {id} = route.params as EditorParam;
  const {data, isFetching, refetch, error} = useQuery({
    queryKey: [...QueryKey, id],
    queryFn: async () => await selectDevice(id),
    enabled: id > 0,
  });

  if (error) {
    Alert.alert(error.name, error.message);
  }

  const handleRefresh = () => {
    refetch();
  };

  return (
    <Fragment>
      <DeviceEditorHeaderActions />
      {isFetching && <ScreenLoading />}
      <ScreenLayout
        refreshControl={<ScreenRefreshControl onRefresh={handleRefresh} />}>
        <DeviceScreenStateProvider>
          <DeviceTopSection id={id} device={data} />
          <DeviceActionSection id={id} device={data} />
          <DeviceFieldsSection id={id} device={data} />
          <DeviceNetworkAddressSection id={id} device={data} />
          <DeviceNameModal id={id} device={data} />
          <DeviceQRModal id={id} device={data} />
        </DeviceScreenStateProvider>
        <ScreenSafetyFooter />
      </ScreenLayout>
    </Fragment>
  );
};

export default DeviceEditorScreen;

const DeviceEditorHeaderActions = () => {
  const route = useRoute();
  const {id} = route.params as EditorParam;

  const isFetching = useIsFetching({
    queryKey: [...QueryKey, id],
  });

  return (
    <HeaderActions>
      <DeviceSearchHeaderAction disabled={isFetching > 0} />
    </HeaderActions>
  );
};

type DeviceFieldItemProps = {device?: Device};
type DeviceSectionProps = {
  id: Device['id'];
} & DeviceFieldItemProps;
const DeviceTopSection = ({id, device}: DeviceSectionProps) => {
  const dispatch = useContext(DispatchContext);
  const handlePress = () => {
    dispatch?.({qrCodeModalVisible: true});
  };
  return (
    <Paper
      bgColor="transparent"
      padding={[0, 0, 2.5, 0]}
      style={{alignItems: 'center'}}>
      <Pressable onPress={handlePress}>
        <DeviceQRCode
          disabled={!device}
          name={device?.name || ''}
          peerId={device?.peerId || ''}
        />
      </Pressable>
    </Paper>
  );
};

const DeviceQRModal = ({id, device}: DeviceSectionProps) => {
  const {qrCodeModalVisible} = useContext(Context);
  const dispatch = useContext(DispatchContext);

  const handleClose = () => {
    if (dispatch) {
      dispatch({qrCodeModalVisible: false});
    }
  };
  return (
    <DeviceQRCodeModal
      visible={!!qrCodeModalVisible}
      onClose={handleClose}
      name={device?.name || ''}
      peerId={device?.peerId || ''}
    />
  );
};

const DeviceActionSection = (props: DeviceSectionProps) => {
  return (
    <Paper
      bgColor="transparent"
      padding={[0, 0, 2.5, 0]}
      gap={3}
      style={{
        flexDirection: 'row',
        justifyContent: 'center',
      }}>
      <DeviceSyncAction {...props} />
      <DeviceRemoveAction {...props} />
      <DeviceExportAction {...props} />
    </Paper>
  );
};

type DeviceActionProps = DeviceSectionProps;
const DeviceRemoveAction = ({id, device}: DeviceActionProps) => {
  const navigation = useNavigation();

  const queryClient = useQueryClient();
  const {mutate} = useMutation({
    mutationFn: async () => await destroyDevice(device?.id || 0),
    onSuccess: data => {
      invalidateListQueryCache(queryClient, data);
      navigation.goBack();
    },
    onError: error => {
      Alert.alert(error.name, error.message);
    },
  });

  const handleConfirm = () => {
    mutate();
  };
  const handlePress = () => {
    Alert.alert(`警告`, `确定将永久移除设备"${device?.name}"`, [
      {text: '取消'},
      {text: '确定', onPress: handleConfirm},
    ]);
  };

  return (
    <Box
      style={{alignItems: 'center'}}
      padding={[1.5, 3]}
      borderColor="transparent"
      gap={0.5}
      onPress={handlePress}>
      <Icon
        name="trash-outline"
        color={device ? 'textSecondary' : 'textDisabled'}
      />
      <Text size="small" color={device ? 'textSecondary' : 'textDisabled'}>
        移除
      </Text>
    </Box>
  );
};

const DeviceSyncAction = ({id, device}: DeviceActionProps) => {
  return (
    <Paper style={{alignItems: 'center'}} gap={0.5} padding={[1.5, 3]}>
      <Icon
        name="sync-outline"
        color={device ? 'textSecondary' : 'textDisabled'}
      />
      <Text size="small" color={device ? 'textSecondary' : 'textDisabled'}>
        同步
      </Text>
    </Paper>
  );
};

const DeviceExportAction = ({id, device}: DeviceActionProps) => {
  return (
    <Paper style={{alignItems: 'center'}} gap={0.5} padding={[1.5, 3]}>
      <Icon
        name="share-social-outline"
        color={device ? 'textSecondary' : 'textDisabled'}
      />
      <Text size="small" color={device ? 'textSecondary' : 'textDisabled'}>
        导出
      </Text>
    </Paper>
  );
};

const DeviceFieldsSection = ({device}: DeviceSectionProps) => {
  const dispatch = useContext(DispatchContext);
  const handlePressName = () => {
    dispatch?.({nameModalVisible: true});
  };

  const handlePressPeerID = () => {
    if (!device?.peerId) return;
    Clipboard.setString(device.peerId);
  };

  return (
    <Section>
      <SectionItem variant="row-start" onPress={handlePressName}>
        <TextFieldItem label="Name" text={device?.name} editable />
      </SectionItem>
      <SectionItem variant="row">
        <DeviceEnabledItem device={device} />
      </SectionItem>
      <SectionItem variant="row">
        <IconFieldItem
          label="Peer ID"
          disabled={!device}
          name="open-outline"
          onPress={handlePressPeerID}
        />
      </SectionItem>
      <SectionItem variant="row">
        <TextFieldItem label="Created At" text={device?.createdAt} />
      </SectionItem>
      <SectionItem variant="row-end">
        <TextFieldItem label="Updated At" text={device?.updatedAt} />
      </SectionItem>
    </Section>
  );
};

const DeviceEnabledItem = ({device}: DeviceFieldItemProps) => {
  const [enabled, setEnabled] = useState(device?.enabled);

  useEffect(() => {
    setEnabled(device?.enabled);
  }, [device?.enabled]);

  const queryClient = useQueryClient();
  const {mutate} = useMutation({
    mutationFn: async (enabled: boolean) =>
      await updateDevice({enabled}, device?.id || 0),
    onSuccess: data => {
      queryClient.setQueryData([...QueryKey, data.id], data);
      invalidateListQueryCache(queryClient, data);
    },
    onError: error => {
      Alert.alert(error.name, error.message);
      setEnabled(device?.enabled);
    },
  });

  const handleConfirm = () => {
    mutate(false);
  };

  const handleCancel = () => {
    setEnabled(device?.enabled);
  };

  const handleChange = (value: boolean) => {
    setEnabled(value);
    if (!value) {
      Alert.alert(`警告`, `确定将禁用设备"${device?.name}"`, [
        {text: '取消', onPress: handleCancel},
        {text: '确定', onPress: handleConfirm},
      ]);
      return;
    }
    mutate(value);
  };

  return (
    <SwitchFieldItem
      label="Enabled"
      value={enabled}
      onValueChange={handleChange}
    />
  );
};

const DeviceNetworkAddressSection = ({id, device}: DeviceSectionProps) => {
  const addrs = useMemo(
    () => device?.networkAddrs || [],
    [device?.networkAddrs],
  );
  return (
    <Section>
      <SectionHeader>
        <Text
          size="small"
          color={addrs.length > 0 ? 'textSecondary' : 'textDisabled'}>
          Network Address
        </Text>
      </SectionHeader>
      <SectionList<string>
        entities={addrs || []}
        extractKey={(_, index) => `device-${id}-${index}`}
        renderItem={addr => <Text size="small">{addr}</Text>}></SectionList>
    </Section>
  );
};

const DeviceNameModal = ({id, device}: DeviceSectionProps) => {
  const {nameModalVisible} = useContext(Context);
  const dispatch = useContext(DispatchContext);

  const handleClose = () => {
    dispatch?.({nameModalVisible: false});
  };

  const queryClient = useQueryClient();
  const {mutate} = useMutation({
    mutationFn: async (name: string) => await updateDevice({name}, id),
    onSuccess: data => {
      queryClient.setQueryData([...QueryKey, id], data);
      dispatch?.({nameModalVisible: false});
      invalidateListQueryCache(queryClient, data);
    },
    onError: error => {
      Alert.alert(error.name, error.message);
      handleClose();
    },
  });

  const handleSubmit = (name?: string) => {
    if (!name) return;
    mutate(name);
  };

  return (
    <TextFieldModal
      labelText="Device Name"
      visible={!!nameModalVisible}
      value={device?.name}
      onClose={handleClose}
      onSubmit={handleSubmit}
    />
  );
};
