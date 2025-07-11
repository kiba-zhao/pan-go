import {Device} from '@pango/datatype';
import {ComponentProps} from 'react';
import Icon from '../Common/Icon';
import {ScreenLayout} from '../Common/Layout';
import {withList} from '../Common/List';
import {useQuery} from '../Common/ReactQuery';
import {
  Section,
  SectionBox,
  SectionHeader,
  SectionItem,
} from '../Common/Section';
import Text from '../Common/Text';
import {searchDevices} from '../Native/Device.spec';
import {HeaderAction, HeaderActions} from './Home.header';
import {Screen} from './Home.navigation';
import {newTabBarIcon} from './Home.tab';

export const DeviceName = `device`;

const DeviceScreen = () => {
  return (
    <ScreenLayout>
      <DeviceSection />
    </ScreenLayout>
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
    <Section>
      <SectionHeader>
        <Text font="bold" size="small" color="textSecondary">
          Active Devices
        </Text>
      </SectionHeader>
      <DeviceBox data={entities || []}>
        <Text size="title" font="medium" color="textPrimary">
          Empty Device
        </Text>
      </DeviceBox>
    </Section>
  );
};

const DeviceBox = withList<Device, ComponentProps<typeof SectionBox>>(
  SectionBox,
  (value, index, items) => {
    return (
      <SectionItem
        key={index}
        variant={index < items.length - 1 ? 'divider' : 'default'}
        icon={<Icon name="desktop-outline" color="textPrimary" />}
        extra={
          <Icon name="chevron-forward" size="small" color="textPrimary" />
        }>
        <Text color="textPrimary">{value.name}</Text>
      </SectionItem>
    );
  },
);
