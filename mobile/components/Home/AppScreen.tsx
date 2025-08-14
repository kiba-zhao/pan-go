import {StyleSheet, View} from 'react-native';

import {scale, verticalScale} from '../Common/SizeMatters';

import {CommonActions, useNavigation} from '../App/Navigation';
import {ScreenLayout} from '../App/ScreenBase';
import Box from '../Common/Box';
import {useTranslation} from '../Common/I18Next';
import Icon from '../Common/Icon';
import Paper from '../Common/Paper';
import Text from '../Common/Text';
import {useTheme} from '../Common/Theme';

import {ExplorerScreenName as DeviceStorageScreenName} from '../DeviceStorage/Screen';
import {ExplorerScreenName as StorageScreenName} from '../Storage/ScreenRoute';
import {HomeDeviceScreenName} from './ScreenRoute';

const AppScreen = () => {
  const navigation = useNavigation();
  const handleSwipeLeft = () => {
    navigation.dispatch(CommonActions.navigate(HomeDeviceScreenName));
  };

  return (
    <ScreenLayout swipeOpts={{onSwipeLeft: handleSwipeLeft}}>
      <StorageFragment />
    </ScreenLayout>
  );
};

export default AppScreen;

const StorageFragment = () => {
  return (
    <Paper bgColor="transparent" gap={1} style={styles.storageFragment}>
      <StorageBox />
      <DeviceStorageBox />
    </Paper>
  );
};

const StorageBox = () => {
  const {t} = useTranslation();
  const {sizes} = useTheme();

  const navigation = useNavigation();
  const handlePress = () => {
    navigation.dispatch(CommonActions.navigate(StorageScreenName));
  };
  return (
    <Box
      onPress={handlePress}
      style={{minWidth: scale(sizes.text) * 13}}
      borderColor="transparent">
      <Text font="bold" color="textPrimary">
        {t('screen.home.box.storage.title')}
      </Text>
      <Text color="textSecondary" size="small">
        {t('screen.home.box.storage.desc', {num: 14151})}
      </Text>
      <View style={{paddingVertical: verticalScale(sizes.base)}}>
        <Icon color="primary" name="folder" size="large" />
      </View>
    </Box>
  );
};

const DeviceStorageBox = () => {
  const {t} = useTranslation();
  const {sizes} = useTheme();

  const navigation = useNavigation();
  const handlePress = () => {
    navigation.dispatch(CommonActions.navigate(DeviceStorageScreenName));
  };
  return (
    <Box
      onPress={handlePress}
      style={[styles.deviceBoxPaper, {minWidth: scale(sizes.text) * 10}]}
      borderColor="transparent">
      <Text font="bold" color="textPrimary">
        {t('screen.home.box.device-storage.title')}
      </Text>
      <Text color="textSecondary" size="small">
        {t('screen.home.box.device-storage.desc', {enabled: 12, all: 32})}
      </Text>
      <View style={{paddingVertical: verticalScale(sizes.base)}}>
        <Icon color="primary" name="file-tray" size="large" />
      </View>
    </Box>
  );
};

const styles = StyleSheet.create({
  storageFragment: {
    flexWrap: 'wrap',
    flexDirection: 'row',
  },
  deviceBoxPaper: {
    alignItems: 'center',
    justifyContent: 'space-between',
  },
});
