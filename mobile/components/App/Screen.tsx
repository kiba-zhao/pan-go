import {useMemo} from 'react';
import {useColorScheme} from 'react-native';

import {DarkTheme, DefaultTheme, ThemeProvider} from '../Common/Theme';
import {
  createTheme as createNavigationTheme,
  Group,
  NavigationContainer,
  Navigator,
  ScreenLayout,
} from './Navigation';

import CameraScannerAppScreen from '../CameraScanner/Screen';
import DeviceAppScreen from '../Device/Screen';
import DeviceStorageAppScreen from '../DeviceStorage/Screen';
import {default as HomeAppScreen, HomeScreenName} from '../Home/Screen';
import StorageAppScreen from '../Storage/Screen';

const AppScreen = () => {
  const scheme = useColorScheme();

  const {theme, navigationTheme} = useMemo(
    () =>
      scheme === 'dark'
        ? {
            theme: DarkTheme,
            navigationTheme: createNavigationTheme(DarkTheme, true),
          }
        : {
            theme: DefaultTheme,
            navigationTheme: createNavigationTheme(DefaultTheme, false),
          },
    [scheme],
  );

  return (
    <ThemeProvider theme={theme}>
      <NavigationContainer theme={navigationTheme}>
        <Navigator
          initialRouteName={HomeScreenName}
          screenLayout={ScreenLayout}
          screenOptions={{
            headerShadowVisible: false,
            headerStyle: {backgroundColor: theme.colors.background},
            headerTitleStyle: {
              fontSize: theme.sizes.text + theme.sizes.base * 0.25,
            },
          }}>
          {HomeAppScreen()}
          <Group
            screenOptions={{
              animation: 'slide_from_right',
            }}>
            {StorageAppScreen()}
            {DeviceStorageAppScreen()}
            {DeviceAppScreen()}
            {CameraScannerAppScreen()}
          </Group>
        </Navigator>
      </NavigationContainer>
    </ThemeProvider>
  );
};

export default AppScreen;
