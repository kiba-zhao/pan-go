import {Device} from '@pango/datatype';
import {Fragment} from 'react';
import {View} from 'react-native';
import {CommonActions, useNavigation} from '../App/App.navigation';
import Icon from '../Common/Icon';
import {FlexRowLayout, RowLayout, ScreenLayout} from '../Common/Layout';
import {useQuery} from '../Common/ReactQuery';
import {
  default as Section,
  SectionHeader,
  SectionList,
} from '../Common/Section';
import Text from '../Common/Text';
import {DeviceScreenName} from '../Device/Device.screen';
import {searchDevices} from '../Native/Device.spec';
import {HeaderAction, HeaderActions} from './Home.header';
import {Screen} from './Home.navigation';
import {newTabBarIcon} from './Home.tab';

export const DeviceName = `home.device`;

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
      title: 'Device',
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

  const navigation = useNavigation();
  const handleItemPress = (entity: Device) => {
    navigation.dispatch(
      CommonActions.navigate(DeviceScreenName, {
        id: entity.id,
      }),
    );
  };
  return (
    <Section>
      <SectionHeader>
        <Text font="bold" size="small" color="textSecondary">
          Recently Devices
        </Text>
      </SectionHeader>
      <SectionList<Device>
        itemProps={{
          gap: 1,
          onPress: handleItemPress,
        }}
        entities={entities || []}
        extractKey={device => device.id}
        renderItem={(entity, index, entities) => (
          <DeviceItem entity={entity} index={index} entities={entities} />
        )}></SectionList>
    </Section>
  );
};

type DeviceItemProps = {
  entity: Device;
  index: number;
  entities: Device[];
};
const DeviceItem = ({entity, index, entities}: DeviceItemProps) => {
  return (
    <Fragment>
      <Icon name="desktop-outline" color="textPrimary" />
      <FlexRowLayout style={{alignItems: 'stretch'}}>
        <View style={{flex: 1}}>
          <Text color="textPrimary">{entity.name}</Text>
        </View>
        <RowLayout style={{alignItems: 'center'}}>
          <Icon name="radio-button-off" size="small" color="textPrimary" />
        </RowLayout>
      </FlexRowLayout>
    </Fragment>
  );
};
