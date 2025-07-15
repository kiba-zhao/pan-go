import {useMemo} from 'react';
import {useColorScheme} from 'react-native';

import {DarkTheme, DefaultTheme, ThemeProvider} from '../Common/Theme';
import {
  createTheme as createNavigationTheme,
  Group,
  NavigationContainer,
  Navigator,
} from './App.navigation';

import {DeviceAppScreen} from '../Device/Device.screen';
import {DeviceStorageAppScreen} from '../DeviceStorage/DeviceStorage.screen';
import {HomeAppScreen, HomeScreenName} from '../Home/Home.screen';
import {StorageAppScreen} from '../Storage/Storage.screen';

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
          </Group>
        </Navigator>
      </NavigationContainer>
    </ThemeProvider>
  );
};

export default AppScreen;
