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
import {Button, TextInput} from 'react-native';
import {useNavigation, useRoute} from '../App/Navigation';
import Box from '../Common/Box';
import Icon from '../Common/Icon';
import {DeviceQRCode, DeviceQRCodeModal} from './QRCode';

import {useMutation, useQuery, useQueryClient} from '../Common/ReactQuery';
import Text from '../Common/Text';
import {destroyDevice, selectDevice, updateDevice} from '../Spec/Device';
import {QueryKey} from './ReactQuery';
import {EditorParam} from './ScreenRoute';

import Clipboard from '@react-native-clipboard/clipboard';
import {HeaderActions} from '../App/ScreenBase';
import Alert from '../Common/Alert';
import {RowLayout} from '../Common/Layout';
import Paper from '../Common/Paper';
import Pressable from '../Common/Pressable';
import {
  ScreenLayout,
  ScreenLoading,
  ScreenModal,
  ScreenRefreshControl,
  ScreenSafetyFooter,
} from '../Common/ScreenBase';
import {
  default as Section,
  SectionHeader,
  SectionItem,
  SectionList,
} from '../Common/Section';
import Switch from '../Common/Switch';
import {RecentlyQueryKey} from './ReactQuery';
import {DeviceRefreshHeaderAction} from './ScreenBase';

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

export const DeviceEditorHeaderActions = () => {
  const route = useRoute();
  const {id} = route.params as EditorParam;
  return (
    <HeaderActions>
      <DeviceRefreshHeaderAction id={id} />
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
      gap={1}
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
      queryClient.invalidateQueries({
        queryKey: RecentlyQueryKey,
      });
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
      onPress={handlePress}>
      <Icon
        name="trash-outline"
        color={device ? 'textSecondary' : 'textDisabled'}
      />
      <Text size="small" color={device ? 'textSecondary' : 'textDisabled'}>
        移除设备
      </Text>
    </Box>
  );
};

const DeviceSyncAction = ({id, device}: DeviceActionProps) => {
  return (
    <Paper style={{alignItems: 'center'}} padding={[1.5, 3]}>
      <Icon
        name="sync-outline"
        color={device ? 'textSecondary' : 'textDisabled'}
      />
      <Text size="small" color={device ? 'textSecondary' : 'textDisabled'}>
        同步设备
      </Text>
    </Paper>
  );
};

const DeviceExportAction = ({id, device}: DeviceActionProps) => {
  return (
    <Paper style={{alignItems: 'center'}} padding={[1.5, 3]}>
      <Icon
        name="share-social-outline"
        color={device ? 'textSecondary' : 'textDisabled'}
      />
      <Text size="small" color={device ? 'textSecondary' : 'textDisabled'}>
        导出设备
      </Text>
    </Paper>
  );
};

const DeviceFieldsSection = ({device}: DeviceSectionProps) => {
  const dispatch = useContext(DispatchContext);
  const handlePressName = () => {
    dispatch?.({nameModalVisible: true});
  };
  return (
    <Section>
      <SectionItem variant="row-start" onPress={handlePressName}>
        <DeviceNameItem device={device} />
      </SectionItem>
      <SectionItem variant="row">
        <DeviceEnabledItem device={device} />
      </SectionItem>
      <SectionItem variant="row">
        <DevicePeerIDItem device={device} />
      </SectionItem>
      <SectionItem variant="row">
        <DeviceCreatedAtItem device={device} />
      </SectionItem>
      <SectionItem variant="row-end">
        <DeviceUpdatedAtItem device={device} />
      </SectionItem>
    </Section>
  );
};

const DeviceNameItem = ({device}: DeviceFieldItemProps) => {
  return (
    <Fragment>
      <Text
        size="small"
        style={{flex: 1}}
        color={device?.name ? 'textPrimary' : 'textDisabled'}>
        Name
      </Text>
      <Text size="small" color="textSecondary">
        {device?.name}
      </Text>
    </Fragment>
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
      queryClient.invalidateQueries({
        queryKey: RecentlyQueryKey,
      });
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
    <Fragment>
      <Text
        size="small"
        style={{flex: 1}}
        color={enabled === void 0 ? 'textDisabled' : 'textPrimary'}>
        Enabled
      </Text>
      <Switch
        value={!!enabled}
        disabled={enabled === void 0}
        onValueChange={handleChange}
      />
    </Fragment>
  );
};

const DevicePeerIDItem = ({device}: DeviceFieldItemProps) => {
  const handlePress = () => {
    if (!device?.peerId) return;
    Clipboard.setString(device.peerId);
  };

  return (
    <Fragment>
      <Text
        size="small"
        style={{flex: 1}}
        color={device?.peerId ? 'textPrimary' : 'textDisabled'}>
        Peer ID
      </Text>
      <Icon
        onPress={handlePress}
        name="open-outline"
        color={device?.peerId ? 'textPrimary' : 'textDisabled'}
      />
    </Fragment>
  );
};

const DeviceCreatedAtItem = ({device}: DeviceFieldItemProps) => (
  <Fragment>
    <Text
      size="small"
      style={{flex: 1}}
      color={device?.createdAt ? 'textSecondary' : 'textDisabled'}>
      Created At
    </Text>
    <Text size="small">{device?.createdAt}</Text>
  </Fragment>
);

const DeviceUpdatedAtItem = ({device}: DeviceFieldItemProps) => (
  <Fragment>
    <Text
      size="small"
      style={{flex: 1}}
      color={device?.updatedAt ? 'textSecondary' : 'textDisabled'}>
      Updated At
    </Text>
    <Text size="small">{device?.updatedAt}</Text>
  </Fragment>
);

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
  const [name, setName] = useState(device?.name);

  useEffect(() => {
    setName(device?.name);
  }, [device?.name]);

  const {nameModalVisible} = useContext(Context);
  const dispatch = useContext(DispatchContext);

  const handleClose = () => {
    dispatch?.({nameModalVisible: false});
    setName(device?.name);
  };

  const queryClient = useQueryClient();
  const {mutate} = useMutation({
    mutationFn: async (name: string) =>
      await updateDevice({name}, device?.id || 0),
    onSuccess: data => {
      queryClient.setQueryData([...QueryKey, data.id], data);
      dispatch?.({nameModalVisible: false});
      queryClient.invalidateQueries({
        queryKey: RecentlyQueryKey,
      });
    },
    onError: error => {
      Alert.alert(error.name, error.message);
      handleClose();
    },
  });

  const handleSubmit = () => {
    if (!name) return;
    mutate(name);
  };

  return (
    <ScreenModal visible={!!nameModalVisible} onClose={handleClose}>
      <RowLayout>
        <Text style={{flex: 1}}>Device Name</Text>
        <Icon name="close-outline" onPress={handleClose} />
      </RowLayout>
      <TextInput
        autoComplete="name"
        inputMode="text"
        placeholder={device?.name}
        value={name}
        onChangeText={setName}
        onSubmitEditing={handleSubmit}
        style={{width: '100%', borderColor: 'black', borderBottomWidth: 1}}
      />
      <Paper bgColor="transparent" padding={[3, 0, 1, 0]}>
        <Button onPress={handleSubmit} title="Submit" disabled={!name} />
      </Paper>
    </ScreenModal>
  );
};
