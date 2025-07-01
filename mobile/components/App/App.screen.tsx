import {useMemo} from 'react';
import {useColorScheme} from 'react-native';

import {DarkTheme, DefaultTheme, ThemeProvider} from '../Common/Theme';
import {
  createTheme as createNavigationTheme,
  Group,
  NavigationContainer,
  Navigator,
} from './App.navigation';

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
        <Navigator initialRouteName={HomeScreenName}>
          <Group>{HomeAppScreen()}</Group>
          <Group>
            {StorageAppScreen()}
            {DeviceStorageAppScreen()}
          </Group>
        </Navigator>
      </NavigationContainer>
    </ThemeProvider>
  );
};

export default AppScreen;
