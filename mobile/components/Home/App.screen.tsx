import {StyleSheet, View} from 'react-native';
import {scale, verticalScale} from 'react-native-size-matters';

import {CommonActions, useNavigation} from '../App/App.navigation';
import {
  CardPanel,
  CardText,
  CardTitleText,
  PressableCard,
} from '../Common/Card';
import {default as Icon} from '../Common/Icon';
import {Layout, RowLayout} from '../Common/Layout';
import {useTheme} from '../Common/Theme';
import {Screen} from './Home.navigation';
import {newTabBarIcon} from './Home.tab';

import {DeviceStorageScreenName} from '../DeviceStorage/DeviceStorage.screen';
import {StorageScreenName} from '../Storage/Storage.screen';

export const AppName = `app`;

const AppScreen = () => {
  return (
    <Layout>
      <ExtFSFragment />
    </Layout>
  );
};

export default AppScreen;

const AppIcon = newTabBarIcon({
  focusedName: 'apps-sharp',
  defaultName: 'apps-outline',
});

export const AppHomeScreen = () => (
  <Screen
    name={AppName}
    component={AppScreen}
    options={{
      title: AppName.toUpperCase(),
      tabBarIcon: AppIcon,
    }}
  />
);

const ExtFSFragment = () => {
  const {sizes} = useTheme();
  return (
    <RowLayout style={{gap: scale(sizes.base)}}>
      <StorageBox />
      <DeviceStorageBox />
    </RowLayout>
  );
};

const StorageBox = () => {
  const {sizes, colors} = useTheme();

  const navigation = useNavigation();
  const handlePress = () => {
    navigation.dispatch(CommonActions.navigate(StorageScreenName));
  };
  return (
    <PressableCard onPress={handlePress}>
      <CardPanel style={{minWidth: scale(sizes.text) * 13}}>
        <CardTitleText>Storage</CardTitleText>
        <CardText>14151 Files</CardText>
        <View style={{paddingVertical: verticalScale(sizes.base)}}>
          <Icon color={colors.primary} name="folder" size={scale(sizes.h2)} />
        </View>
      </CardPanel>
    </PressableCard>
  );
};

const DeviceStorageBox = () => {
  const {sizes, colors} = useTheme();

  const navigation = useNavigation();
  const handlePress = () => {
    navigation.dispatch(CommonActions.navigate(DeviceStorageScreenName));
  };
  return (
    <PressableCard onPress={handlePress}>
      <CardPanel
        style={{...styles.deviceBoxPanel, minWidth: scale(sizes.text) * 10}}>
        <CardTitleText>Device</CardTitleText>
        <CardText>11 / 32</CardText>
        <View style={{paddingVertical: verticalScale(sizes.base)}}>
          <Icon color={colors.primary} name="radio" size={scale(sizes.h2)} />
        </View>
      </CardPanel>
    </PressableCard>
  );
};

const styles = StyleSheet.create({
  deviceBoxPanel: {
    alignItems: 'center',
    justifyContent: 'space-between',
  },
});
