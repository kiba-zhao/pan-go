import type {Device} from '@pango/data';
import {
  createContext,
  Fragment,
  useContext,
  useEffect,
  useMemo,
  useReducer,
  useState,
} from 'react';
import {CommonActions, useNavigation} from '../App/Navigation';
import {QRScannerScreenName} from '../CameraScanner/Screen';
import {QRScannerScope} from '../CameraScanner/ScreenRoute';

import type {Dispatch, PropsWithChildren} from 'react';
import {
  ActivityIndicator,
  RefreshControl,
  RefreshControlProps,
  StyleSheet,
  View,
} from 'react-native';
import {ScreenLayout} from '../App/ScreenBase';
import {useTranslation} from '../Common/I18Next';
import Icon from '../Common/Icon';
import {FlexRowLayout, RowLayout} from '../Common/Layout';
import {
  QueryObserver,
  useIsFetching,
  useQuery,
  useQueryClient,
} from '../Common/ReactQuery';
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
import type {FetchDeviceCondition, FetchDevicesResults} from '../Spec/Device';
import {fetchDevices} from '../Spec/Device';
import {HeaderAction, HeaderActions} from './ScreenBase';
import {HomeAppScreenName, HomeSettingsScreenName} from './ScreenRoute';

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
  const navigation = useNavigation();
  const handleSwipeLeft = () => {
    navigation.dispatch(CommonActions.navigate(HomeSettingsScreenName));
  };
  const handleSwipeRight = () => {
    navigation.dispatch(CommonActions.navigate(HomeAppScreenName));
  };
  return (
    <DeviceScreenStateProvider>
      <ScreenLayout
        refreshControl={<DeviceScreenRefreshControl />}
        swipeOpts={{
          onSwipeLeft: handleSwipeLeft,
          onSwipeRight: handleSwipeRight,
        }}>
        <DeviceSection />
        <DeviceActionSection />
      </ScreenLayout>
    </DeviceScreenStateProvider>
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

type DeviceScreenRefreshControlProps = Partial<
  Pick<RefreshControlProps, 'refreshing'>
> &
  Omit<RefreshControlProps, 'refreshing'>;
const DeviceScreenRefreshControl = ({
  refreshing,
  ...props
}: DeviceScreenRefreshControlProps) => {
  const {queryKey} = useContext(Context);
  const queryClient = useQueryClient();
  const handleRefresh = () => {
    queryClient.refetchQueries({
      queryKey: queryKey,
      type: 'active',
    });
  };
  return (
    <RefreshControl
      {...props}
      refreshing={refreshing || false}
      onRefresh={handleRefresh}
    />
  );
};

const DeviceSection = () => {
  const {queryKey} = useContext(Context);
  const condition = queryKey.at(-1) as FetchDeviceCondition;
  const {data} = useQuery<FetchDevicesResults>({
    queryKey,
    queryFn: async () => await fetchDevices(condition),
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

  const {t} = useTranslation();
  const {queryKey} = useContext(Context);

  const isFetching = useIsFetching({
    queryKey: queryKey,
  });

  return (
    <SectionHeader
      style={{flexDirection: 'row', alignItems: 'center', paddingLeft: 0}}>
      <RowLayout style={{gap: scale(sizes.base), flex: 1}}>
        <DeviceTab
          text={t('screen.device.tabs.recently')}
          queryKey={HomeDeviceRecentlyQueryKey}
        />
        <DeviceTab
          text={t('screen.device.tabs.online')}
          queryKey={HomeDeviceOnlineQueryKey}
        />
        <DeviceTab
          text={t('screen.device.tabs.disabled')}
          queryKey={HomeDeviceDisabledQueryKey}
        />
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
  const {t} = useTranslation();
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
      {t('action.new')}
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

const DeviceActionSection = () => {
  const [hasMore, setHasMore] = useState(false);

  const queryClient = useQueryClient();
  const {queryKey} = useContext(Context);

  useEffect(() => {
    const observer = new QueryObserver<FetchDevicesResults>(queryClient, {
      queryKey: queryKey,
    });

    return observer.subscribe(({data, isFetching}) => {
      if (isFetching) {
        return;
      }
      const meta = data !== void 0 ? data[0] : void 0;
      setHasMore(meta?.next !== void 0);
    });
  }, [queryKey, queryClient]);

  return <Fragment>{hasMore ? <DeviceMoreAction /> : void 0}</Fragment>;
};

const DeviceMoreAction = () => {
  const {t} = useTranslation();
  const navigation = useNavigation();
  const handlePress = () => {
    navigation.dispatch(CommonActions.navigate(DeviceSearchScreenName));
  };
  return (
    <View style={[styles.moreContainer]}>
      <Text size="small" font="bold" color="primary" onPress={handlePress}>
        {t('action.find-more')}
      </Text>
    </View>
  );
};

const styles = StyleSheet.create({
  moreContainer: {
    alignItems: 'center',
    justifyContent: 'center',
  },
});
