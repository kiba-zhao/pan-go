import {StyleSheet, View} from 'react-native';
import {withLayout} from './Layout';
import {OpacityPressable} from './Pressable';

const styles = StyleSheet.create({
  card: {
    boxShadow: [
      {
        color: 'rgba(0,0,0, .2)',
        offsetX: 0,
        offsetY: 2,
        blurRadius: 1,
        spreadDistance: -1,
      },
      {
        color: 'rgba(0,0,0, .14)',
        offsetX: 0,
        offsetY: 1,
        blurRadius: 1,
        spreadDistance: 0,
      },
      {
        color: 'rgba(0,0,0, .12)',
        offsetX: 0,
        offsetY: 1,
        blurRadius: 3,
        spreadDistance: 0,
      },
    ],
  },
});

export const Card = withLayout(View, {
  color: 'surface',
  style: styles.card,
  padding: 1,
  radius: 1,
});

export const PressableCard = OpacityPressable;
