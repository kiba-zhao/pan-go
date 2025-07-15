import {Device} from '@pango/datatype';
import {HeaderAction, HeaderActions} from '../App/App.header';
import type {
  ScreenComponent,
  ScreenComponentProps,
} from '../App/App.navigation';
import {CommonActions, Screen, useNavigation} from '../App/App.navigation';
import Icon from '../Common/Icon';
import {ScreenLayout} from '../Common/Layout';
import {useQuery} from '../Common/ReactQuery';
import Text from '../Common/Text';

import {Fragment, useEffect} from 'react';
import Paper from '../Common/Paper';
import {LoadingScreen} from '../Common/Screen';
import {
  default as Section,
  SectionHeader,
  SectionItem,
  SectionList,
} from '../Common/Section';
import Switch from '../Common/Switch';
import {selectDevice} from '../Native/Device.spec';
import {DeviceName} from './Device.constants';

export const DeviceScreenName = DeviceName;

type DeviceParam = Pick<Device, 'id'> & {enabled?: boolean};

export const DeviceScreen = ({route}: ScreenComponentProps<DeviceParam>) => {
  const {id, enabled} = route.params;

  const {data, isFetching} = useQuery({
    queryKey: ['devices', id],
    queryFn: async () => await selectDevice(id),
  });

  if (isFetching) {
    return <LoadingScreen />;
  }

  return (
    <ScreenLayout>
      <DeviceEditor enabled={enabled} id={id} device={data} />
    </ScreenLayout>
  );
};

export default DeviceScreen;

type DeviceHeaderActionProps = Pick<DeviceParam, 'id' | 'enabled'>;
const DeviceHeaderActions = ({id, enabled}: DeviceHeaderActionProps) => {
  return (
    <HeaderActions>
      <DeviceHeaderRefreshAction id={id} disabled={!enabled} />
    </HeaderActions>
  );
};

export const DeviceAppScreen = () => (
  <Screen
    name={DeviceScreenName}
    component={DeviceScreen as ScreenComponent}
    options={({route}) => ({
      title: 'Device Settings',
      headerRight: () =>
        DeviceHeaderActions(route.params as DeviceHeaderActionProps),
    })}
  />
);

const DeviceHeaderRefreshAction = ({
  id,
  disabled,
}: {
  disabled?: boolean;
  id: Device['id'];
}) => {
  const handlePress = () => {
    console.log('device remove');
  };

  return (
    <HeaderAction onPress={handlePress} disabled={disabled}>
      <Icon
        name="refresh-outline"
        size="medium"
        color={disabled ? 'textDisabled' : 'textPrimary'}
      />
    </HeaderAction>
  );
};

type DeviceEditorProps = {
  device?: Device;
} & DeviceParam;
const DeviceEditor = ({id, enabled, device}: DeviceEditorProps) => {
  const navigation = useNavigation();

  useEffect(() => {
    if (enabled === void 0)
      navigation.dispatch(CommonActions.setParams({enabled: true}));
  }, [device]);

  return (
    <Fragment>
      <DeviceTopSection id={id} device={device} />
      <DeviceActionSection id={id} device={device} />
      <DeviceSettingsSection id={id} device={device} />
      <DeviceNetworkAddressSection id={id} device={device} />
    </Fragment>
  );
};

type DeviceSectionProps = Pick<DeviceEditorProps, 'id' | 'device'>;
const DeviceTopSection = ({id, device}: DeviceSectionProps) => {
  return (
    <Paper
      bgColor="transparent"
      padding={[0, 0, 2.5, 0]}
      style={{alignItems: 'center'}}>
      <Paper bgColor="textDisabled" style={{width: 120, height: 120}} />
    </Paper>
  );
};

const DeviceActionSection = ({}: DeviceSectionProps) => {
  return (
    <Paper
      bgColor="transparent"
      padding={[0, 0, 2.5, 0]}
      gap={1}
      style={{
        flexDirection: 'row',
        justifyContent: 'center',
      }}>
      <DeviceSyncAction />
      <DeviceRemoveAction />
      <DeviceExportAction />
    </Paper>
  );
};

const DeviceRemoveAction = () => {
  return (
    <Paper style={{alignItems: 'center'}} padding={[1.5, 3]}>
      <Icon name="trash-outline" color="textSecondary" />
      <Text size="small" color="textSecondary">
        移除设备
      </Text>
    </Paper>
  );
};

const DeviceSyncAction = () => {
  return (
    <Paper style={{alignItems: 'center'}} padding={[1.5, 3]}>
      <Icon name="sync-outline" color="textSecondary" />
      <Text size="small" color="textSecondary">
        同步设备
      </Text>
    </Paper>
  );
};

const DeviceExportAction = () => {
  return (
    <Paper style={{alignItems: 'center'}} padding={[1.5, 3]}>
      <Icon name="share-social-outline" color="textSecondary" />
      <Text size="small" color="textSecondary">
        导出设备
      </Text>
    </Paper>
  );
};

const DeviceSettingsSection = ({device}: DeviceSectionProps) => {
  return (
    <Section>
      <SectionItem variant="row-start">
        <DeviceNameItem name={device?.name} />
      </SectionItem>
      <SectionItem variant="row">
        <DeviceBlockedItem blocked={device?.blocked} />
      </SectionItem>
      <SectionItem variant="row">
        <DevicePeerIDItem peerId={device?.peerId} />
      </SectionItem>
      <SectionItem variant="row">
        <DeviceCreatedAtItem createdAt={device?.createdAt} />
      </SectionItem>
      <SectionItem variant="row-end">
        <DeviceUpdatedAtItem updatedAt={device?.updatedAt} />
      </SectionItem>
    </Section>
  );
};

const DeviceNameItem = ({name}: {name?: Device['name']}) => (
  <Fragment>
    <Text size="small" style={{flex: 1}}>
      Name
    </Text>
    <Text size="small" color="textSecondary">
      {name}
    </Text>
  </Fragment>
);

const DeviceBlockedItem = ({blocked}: {blocked?: Device['blocked']}) => (
  <Fragment>
    <Text size="small" style={{flex: 1}}>
      Blocked
    </Text>
    <Switch value={blocked} />
  </Fragment>
);

const DevicePeerIDItem = ({peerId}: {peerId?: Device['peerId']}) => (
  <Fragment>
    <Text size="small" style={{flex: 1}}>
      Peer ID
    </Text>
    <Icon name="open-outline" />
  </Fragment>
);

const DeviceCreatedAtItem = ({
  createdAt,
}: {
  createdAt?: Device['createdAt'];
}) => (
  <Fragment>
    <Text size="small" style={{flex: 1}} color="textSecondary">
      Created At
    </Text>
    <Text size="small">{createdAt}</Text>
  </Fragment>
);

const DeviceUpdatedAtItem = ({
  updatedAt,
}: {
  updatedAt?: Device['updatedAt'];
}) => (
  <Fragment>
    <Text size="small" style={{flex: 1}} color="textSecondary">
      Updated At
    </Text>
    <Text size="small">{updatedAt}</Text>
  </Fragment>
);

const DeviceNetworkAddressSection = ({id, device}: DeviceSectionProps) => {
  return (
    <Section>
      <SectionHeader>
        <Text size="small" color="textSecondary">
          Network Address
        </Text>
      </SectionHeader>
      <SectionList<string>
        entities={device?.networkAddrs || []}
        extractKey={(_, index) => `device-${id}-${index}`}
        renderItem={addr => <Text size="small">{addr}</Text>}></SectionList>
    </Section>
  );
};
