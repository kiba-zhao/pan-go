import {Device} from '@pango/datatype';
import {ComponentProps, Fragment} from 'react';
import {Card} from '../Common/Card';
import Icon from '../Common/Icon';
import {Layout} from '../Common/Layout';
import {ListItem, withList} from '../Common/List';
import {useQuery} from '../Common/ReactQuery';
import Text from '../Common/Text';
import {searchDevices} from '../Native/Device.spec';
import {HeaderAction, HeaderActions} from './Home.header';
import {Screen} from './Home.navigation';
import {newTabBarIcon} from './Home.tab';

export const DeviceName = `device`;

const DeviceScreen = () => {
  return (
    <Layout>
      <DeviceSection />
    </Layout>
  );
};

export default DeviceScreen;

const DeviceIcon = newTabBarIcon({
  focusedName: 'radio-sharp',
  defaultName: 'radio-outline',
});

const DeviceHeaderActions = () => {
  return (
    <HeaderActions>
      <DeviceScreenSearchButton />
      <DeviceScreenMenuButton />
    </HeaderActions>
  );
};

export const DeviceHomeScreen = () => (
  <Screen
    name={DeviceName}
    component={DeviceScreen}
    options={{
      title: DeviceName.toUpperCase(),
      tabBarIcon: DeviceIcon,
      headerRight: DeviceHeaderActions,
    }}
  />
);

const DeviceScreenSearchButton = () => {
  const handlePress = () => {
    console.log('search');
  };
  return (
    <HeaderAction onPress={handlePress}>
      <Icon name="search-sharp" color="textPrimary" />
    </HeaderAction>
  );
};

const DeviceScreenMenuButton = () => {
  const handlePress = () => {
    console.log('menu');
  };
  return (
    <HeaderAction onPress={handlePress}>
      <Icon name="menu-sharp" color="textPrimary" />
    </HeaderAction>
  );
};

const DeviceSection = () => {
  const condition = {
    _sort: 'updatedAt',
    _order: 'desc',
    _start: 0,
    _end: 10,
  } as Parameters<typeof searchDevices>[0];
  const {data, isFetching} = useQuery({
    queryKey: ['devices', condition],
    queryFn: async () => await searchDevices(condition),
  });

  const entities = data && data[1];
  return (
    <Fragment>
      <Text font="bold" size="small" color="textSecondary">
        Active Devices
      </Text>
      <DeviceBox data={entities || []}>
        <Text size="h6" font="medium" color="textPrimary">
          Empty Device
        </Text>
      </DeviceBox>
    </Fragment>
  );
};

const DeviceBox = withList<Device, ComponentProps<typeof Card>>(
  Card,
  (value, index, items) => {
    return (
      <ListItem
        key={index}
        icon={<Icon name="desktop-outline" color="textDisabled" />}
        text={<Text color="textPrimary">{value.name}</Text>}>
        <Icon name="chevron-forward" size="small" color="textPrimary" />
      </ListItem>
    );
  },
);
