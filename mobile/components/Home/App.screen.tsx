import {StyleSheet, View} from 'react-native';
import {scale, verticalScale} from 'react-native-size-matters';

import {CommonActions, useNavigation} from '../App/App.navigation';
import {Card, PressableCard} from '../Common/Card';
import Icon from '../Common/Icon';
import {Layout, RowLayout} from '../Common/Layout';
import Text from '../Common/Text';
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
      <Card style={{minWidth: scale(sizes.text) * 13}}>
        <Text font="bold" color="textPrimary">
          Storage
        </Text>
        <Text color="textSecondary" size="small">
          14151 Files
        </Text>
        <View style={{paddingVertical: verticalScale(sizes.base)}}>
          <Icon color="primary" name="folder" size="large" />
        </View>
      </Card>
    </PressableCard>
  );
};

const DeviceStorageBox = () => {
  const {sizes} = useTheme();

  const navigation = useNavigation();
  const handlePress = () => {
    navigation.dispatch(CommonActions.navigate(DeviceStorageScreenName));
  };
  return (
    <PressableCard onPress={handlePress}>
      <Card
        style={{...styles.deviceBoxPanel, minWidth: scale(sizes.text) * 10}}>
        <Text font="bold" color="textPrimary">
          Device
        </Text>
        <Text color="textSecondary" size="small">
          12 / 32
        </Text>
        <View style={{paddingVertical: verticalScale(sizes.base)}}>
          <Icon color="primary" name="radio" size="large" />
        </View>
      </Card>
    </PressableCard>
  );
};

const styles = StyleSheet.create({
  deviceBoxPanel: {
    alignItems: 'center',
    justifyContent: 'space-between',
  },
});
