import {HeaderButton} from '@react-navigation/elements';
import type {ViewStyle} from 'react-native';
import {View} from 'react-native';

import {scale} from '../Common/SizeMatters';
import {withTheme} from '../Common/StyleSheet';

export const HeaderAction = HeaderButton;

export const HeaderActions = withTheme(
  View,
  ({sizes}) =>
    ({
      alignItems: 'flex-end',
      flexDirection: 'row',
      paddingRight: scale(sizes.base),
    } as ViewStyle),
);
