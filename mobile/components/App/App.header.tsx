import {HeaderButton} from '@react-navigation/elements';
import {View} from 'react-native';

import {withLayout} from '../Common/Layout';

export const HeaderAction = withLayout(HeaderButton, {
  padding: 1,
});

export const HeaderActions = withLayout(View, {
  style: {alignItems: 'flex-end'},
  direction: 'row',
  gap: 0,
  margin: [0, -1, 0, 0],
});
