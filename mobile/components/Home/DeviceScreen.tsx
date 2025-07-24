import type {Device, SearchResults} from '@pango/data';
import {Fragment} from 'react';
import {View} from 'react-native';
import {CommonActions, useNavigation} from '../App/Navigation';
import {QRScannerScreenName} from '../CameraScanner/Screen';
import {QRScannerScope} from '../CameraScanner/ScreenRoute';

import Icon from '../Common/Icon';
import {FlexRowLayout, RowLayout} from '../Common/Layout';
import {useIsFetching, useQuery, useQueryClient} from '../Common/ReactQuery';
import {ScreenLayout, ScreenRefreshControl} from '../Common/ScreenBase';
import {
  default as Section,
  SectionHeader,
  SectionList,
} from '../Common/Section';
import Text from '../Common/Text';
import {RecentlyQueryKey as HomeDeviceRecentlyQueryKey} from '../Device/ReactQuery';
import {DeviceEditorScreenName} from '../Device/Screen';
import {searchDevices} from '../Spec/Device';
import {HeaderAction, HeaderActions} from './ScreenBase';

const DeviceScreen = () => {
  const isFetching = useIsFetching({
    queryKey: HomeDeviceRecentlyQueryKey,
  });
  const queryClient = useQueryClient();
  const handleRefresh = () => {
    if (isFetching) return;
    queryClient.refetchQueries({
      queryKey: HomeDeviceRecentlyQueryKey,
      type: 'active',
    });
  };

  return (
    <ScreenLayout
      refreshControl={<ScreenRefreshControl onRefresh={handleRefresh} />}>
      <DeviceSection />
    </ScreenLayout>
  );
};

export default DeviceScreen;

export const DeviceHeaderActions = () => {
  return (
    <HeaderActions>
      <DeviceScreenNewAction />
      <DeviceScreenMenuAction />
    </HeaderActions>
  );
};

const DeviceScreenNewAction = () => {
  const navigation = useNavigation();
  const handlePress = () => {
    navigation.dispatch(
      CommonActions.navigate(QRScannerScreenName, {
        scope: [QRScannerScope.Device],
      }),
    );
  };
  return (
    <HeaderAction onPress={handlePress}>
      <Icon name="add-sharp" color="textPrimary" />
    </HeaderAction>
  );
};

const DeviceScreenSearchAction = () => {
  const handlePress = () => {
    console.log('search');
  };
  return (
    <HeaderAction onPress={handlePress}>
      <Icon name="search-sharp" color="textPrimary" />
    </HeaderAction>
  );
};

const DeviceScreenMenuAction = () => {
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
  const {data, isFetching} = useQuery<SearchResults<Device>>({
    queryKey: HomeDeviceRecentlyQueryKey,
    queryFn: async () => await searchDevices(condition),
  });

  const entities = !isFetching && data !== void 0 ? data[1] : [];
  const navigation = useNavigation();
  const handleItemPress = (entity: Device) => {
    navigation.dispatch(
      CommonActions.navigate(DeviceEditorScreenName, {
        id: entity.id,
      }),
    );
  };
  return (
    <Section>
      <SectionHeader style={{flexDirection: 'row', alignItems: 'center'}}>
        <Text font="bold" size="small" color="textSecondary" style={{flex: 1}}>
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
