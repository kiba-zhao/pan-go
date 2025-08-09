import type {
  GestureResponderEvent,
  PanResponderGestureState,
  PanResponderInstance,
} from 'react-native';
import {PanResponder} from 'react-native';

export type SwipePanResponderOpts = {
  swipeThreshold?: number;
  onSwipeLeft?: () => void;
  onSwipeRight?: () => void;
  onSwipeUp?: () => void;
  onSwipeDown?: () => void;
  onSwipe?: () => void;
};
export function createSwipePanResponder({
  swipeThreshold = 50,
  onSwipe,
  onSwipeLeft = onSwipe,
  onSwipeRight = onSwipe,
  onSwipeUp = onSwipe,
  onSwipeDown = onSwipe,
}: SwipePanResponderOpts): PanResponderInstance {
  return PanResponder.create({
    onMoveShouldSetPanResponder: () => true,
    onPanResponderMove: (
      event: GestureResponderEvent,
      gestureState: PanResponderGestureState,
    ) => {
      if (gestureState.dx > swipeThreshold) {
        onSwipeRight?.();
      } else if (gestureState.dx < swipeThreshold * -1) {
        onSwipeLeft?.();
      } else if (gestureState.dy > swipeThreshold) {
        onSwipeDown?.();
      } else if (gestureState.dy < swipeThreshold * -1) {
        onSwipeUp?.();
      }
    },
  });
}
