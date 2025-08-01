import type {Device, SearchResults} from '@pango/data';
import {createContext, Fragment, useContext, useMemo, useReducer} from 'react';
import {CommonActions, useNavigation} from '../App/Navigation';
import {QRScannerScreenName} from '../CameraScanner/Screen';
import {QRScannerScope} from '../CameraScanner/ScreenRoute';

import type {Dispatch, PropsWithChildren} from 'react';
import {ActivityIndicator} from 'react-native';
import Icon from '../Common/Icon';
import {FlexRowLayout, RowLayout} from '../Common/Layout';
import {useIsFetching, useQuery, useQueryClient} from '../Common/ReactQuery';
import {ScreenLayout, ScreenRefreshControl} from '../Common/ScreenBase';
import {
  default as Section,
  SectionHeader,
  SectionList,
} from '../Common/Section';
import {scale} from '../Common/SizeMatters';
import Text from '../Common/Text';
import {useTheme} from '../Common/Theme';
import {
  DisabledQueryKey as HomeDeviceDisabledQueryKey,
  OnlineQueryKey as HomeDeviceOnlineQueryKey,
  RecentlyQueryKey as HomeDeviceRecentlyQueryKey,
} from '../Device/ReactQuery';
import {
  DeviceEditorScreenName,
  DeviceSearchScreenName,
} from '../Device/ScreenRoute';
import {SearchCondition, searchDevices} from '../Spec/Device';
import {HeaderAction, HeaderActions} from './ScreenBase';

type DeviceScreenState = {
  queryKey: any[];
};
type DeviceScreenStateAction = DeviceScreenState;

const reducer = (
  state: DeviceScreenState,
  action: DeviceScreenStateAction,
) => ({
  ...state,
  ...action,
});

const Context = createContext<DeviceScreenState>({
  queryKey: HomeDeviceRecentlyQueryKey,
});
const DispatchContext = createContext<Dispatch<DeviceScreenStateAction> | null>(
  null,
);

const DeviceScreenStateProvider = ({children}: PropsWithChildren<{}>) => {
  const [state, dispatch] = useReducer(reducer, {
    queryKey: HomeDeviceRecentlyQueryKey,
  });
  return (
    <Context.Provider value={state}>
      <DispatchContext.Provider value={dispatch}>
        {children}
      </DispatchContext.Provider>
    </Context.Provider>
  );
};

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
      <DeviceScreenStateProvider>
        <DeviceSection />
      </DeviceScreenStateProvider>
    </ScreenLayout>
  );
};

export default DeviceScreen;

export const DeviceHeaderActions = () => {
  return (
    <HeaderActions>
      <DeviceScreenSearchAction />
      <DeviceScreenMenuAction />
    </HeaderActions>
  );
};

const DeviceScreenSearchAction = () => {
  const navigation = useNavigation();

  const handlePress = () => {
    navigation.dispatch(CommonActions.navigate(DeviceSearchScreenName));
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
  const {queryKey} = useContext(Context);
  const condition = queryKey.at(-1) as SearchCondition;
  const {data} = useQuery<SearchResults<Device>>({
    queryKey,
    queryFn: async () => await searchDevices(condition),
    placeholderData: _ => _,
  });

  const entities = data !== void 0 ? data[1] : [];
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
      <DeviceSectionHeader />
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

const DeviceSectionHeader = () => {
  const {sizes} = useTheme();

  const {queryKey} = useContext(Context);

  const isFetching = useIsFetching({
    queryKey: queryKey,
  });

  return (
    <SectionHeader
      style={{flexDirection: 'row', alignItems: 'center', paddingLeft: 0}}>
      <RowLayout style={{gap: scale(sizes.base), flex: 1}}>
        <DeviceTab text="Recently" queryKey={HomeDeviceRecentlyQueryKey} />
        <DeviceTab text="Online" queryKey={HomeDeviceOnlineQueryKey} />
        <DeviceTab text="Disabled" queryKey={HomeDeviceDisabledQueryKey} />
      </RowLayout>
      {isFetching > 0 ? (
        <ActivityIndicator size={sizes.text} />
      ) : (
        <DeviceNewAction />
      )}
    </SectionHeader>
  );
};

type DeviceTabProps = {text: string; queryKey: DeviceScreenState['queryKey']};
const DeviceTab = ({text, queryKey: tabQueryKey}: DeviceTabProps) => {
  const {queryKey} = useContext(Context);
  const dispatch = useContext(DispatchContext);

  const isFetching = useIsFetching({
    queryKey: queryKey,
  });

  const disabled = useMemo(
    () => isFetching > 0 || queryKey === tabQueryKey,
    [isFetching, queryKey, tabQueryKey],
  );

  const {color, bgColor} = useMemo(() => {
    if (queryKey === tabQueryKey) {
      return {color: 'textPrimary', bgColor: 'surface'};
    }
    return {color: 'textSecondary', bgColor: 'transparent'};
  }, [queryKey, tabQueryKey]);

  const handlePress = () => {
    dispatch?.({queryKey: tabQueryKey});
  };

  return (
    <Text
      disabled={disabled}
      font="bold"
      color={color}
      size="small"
      bgColor={bgColor}
      radius={0.5}
      onPress={handlePress}
      padding={1}>
      {text}
    </Text>
  );
};

const DeviceNewAction = () => {
  const navigation = useNavigation();
  const handlePressNew = () => {
    navigation.dispatch(
      CommonActions.navigate(QRScannerScreenName, {
        scope: [QRScannerScope.Device],
      }),
    );
  };
  return (
    <Text size="small" font="bold" color="primary" onPress={handlePressNew}>
      New
    </Text>
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
      <FlexRowLayout style={{alignItems: 'center'}}>
        <Text color="textPrimary" style={{flex: 1}}>
          {entity.name}
        </Text>
        <Icon name="chevron-forward-outline" size="small" color="textPrimary" />
      </FlexRowLayout>
    </Fragment>
  );
};
