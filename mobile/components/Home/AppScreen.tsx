import {StyleSheet, View} from 'react-native';
import {scale, verticalScale} from 'react-native-size-matters';

import {CommonActions, useNavigation} from '../App/Navigation';
import Box from '../Common/Box';
import Icon from '../Common/Icon';
import Paper from '../Common/Paper';
import {ScreenLayout} from '../Common/ScreenBase';
import Text from '../Common/Text';
import {useTheme} from '../Common/Theme';

import {ExplorerScreenName as DeviceStorageScreenName} from '../DeviceStorage/Screen';
import {ExplorerScreenName as StorageScreenName} from '../Storage/Screen';

const AppScreen = () => {
  return (
    <ScreenLayout>
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
        Storage
      </Text>
      <Text color="textSecondary" size="small">
        14151 Files
      </Text>
      <View style={{paddingVertical: verticalScale(sizes.base)}}>
        <Icon color="primary" name="folder" size="large" />
      </View>
    </Box>
  );
};

const DeviceStorageBox = () => {
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
        Device
      </Text>
      <Text color="textSecondary" size="small">
        12 / 32
      </Text>
      <View style={{paddingVertical: verticalScale(sizes.base)}}>
        <Icon color="primary" name="radio" size="large" />
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
