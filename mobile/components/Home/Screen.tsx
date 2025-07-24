import {Screen as AppNavigationScreen} from '../App/Navigation';

import {Navigator, Screen} from './Navigation';

import {useTheme} from '../Common/Theme';
import {newTabBarIcon} from './ScreenBase';
import {
  HomeAppScreenName,
  HomeDeviceScreenName,
  HomeScreenName,
  HomeSettingsScreenName,
} from './ScreenRoute';

import {default as AppScreen} from './AppScreen';
import {DeviceHeaderActions, default as DeviceScreen} from './DeviceScreen';
import {default as SettingsScreen} from './SettingsScreen';

const HomeAppScreen = () => (
  <AppNavigationScreen
    name={HomeScreenName}
    component={HomeScreen}
    options={{title: HomeScreenName.toUpperCase(), headerShown: false}}
  />
);
export default HomeAppScreen;

export {HomeScreenName};

const HomeScreen = () => {
  const {colors, sizes} = useTheme();
  return (
    <Navigator
      initialRouteName={HomeAppScreenName}
      screenOptions={{
        headerStyle: {backgroundColor: colors.background},
        headerShadowVisible: false,
        headerTitleStyle: {
          fontSize: sizes.text + sizes.base * 0.5,
        },
      }}>
      {AppHomeScreen()}
      {DeviceHomeScreen()}
      {SettingsHomeScreen()}
    </Navigator>
  );
};

// const AppScreen = lazy(() => import('./AppScreen'));

const AppHomeScreen = () => (
  <Screen
    name={HomeAppScreenName}
    component={AppScreen}
    options={{
      title: 'Home',
      tabBarIcon: AppIcon,
    }}
  />
);

const AppIcon = newTabBarIcon({
  focusedName: 'apps-sharp',
  defaultName: 'apps-outline',
});

// const DeviceScreen = lazy(() => import('./DeviceScreen'));
// const DeviceHeaderActions = lazy(() =>
//   import('./DeviceScreen').then(m => ({default: m.DeviceHeaderActions})),
// );

const DeviceHomeScreen = () => (
  <Screen
    name={HomeDeviceScreenName}
    component={DeviceScreen}
    options={{
      title: 'Device',
      tabBarIcon: DeviceIcon,
      headerRight: DeviceHeaderActions,
    }}
  />
);

const DeviceIcon = newTabBarIcon({
  focusedName: 'radio-sharp',
  defaultName: 'radio-outline',
});

// const DeviceHeaderRight = withLazy(DeviceHeaderActions);

// const SettingsScreen = lazy(() => import('./SettingsScreen'));

const SettingsHomeScreen = () => (
  <Screen
    name={HomeSettingsScreenName}
    component={SettingsScreen}
    options={{
      title: 'Settings',
      tabBarIcon: SettingsIcon,
    }}
  />
);

const SettingsIcon = newTabBarIcon({
  focusedName: 'settings-sharp',
  defaultName: 'settings-outline',
});
