import type {Key, PropsWithChildren} from 'react';
import {Fragment, useMemo} from 'react';
import {StyleSheet, View, type ViewStyle} from 'react-native';
import type {BoxPressableProps, BoxProps} from './Box';
import Box from './Box';
import type {PaperProps} from './Paper';
import Paper from './Paper';
import {scale} from './SizeMatters';
import {withTheme as withThemeStyle} from './StyleSheet';
import Text from './Text';
import {useTheme} from './Theme';

const styles = StyleSheet.create({
  sectionList: {
    borderWidth: 1,
  },
  sectionItem: {
    flexDirection: 'row',
    alignItems: 'center',
  },
});

const Section = withThemeStyle(
  View,
  ({sizes}) =>
    ({
      paddingBottom: scale(sizes.base),
    } as ViewStyle),
);

export default Section;

export const SectionHeader = withThemeStyle(
  View,
  ({sizes}) =>
    ({
      paddingBottom: scale(sizes.base),
      paddingHorizontal: scale(sizes.base * 1.5),
      borderWidth: scale(sizes.border),
      borderColor: 'transparent',
    } as ViewStyle),
);

type SectionItemVariant = 'default' | 'row' | 'row-start' | 'row-end';
type SectionItemProps = BoxProps & {
  variant?: SectionItemVariant;
};
export const SectionItem = ({
  variant = 'default',
  radius,
  bgColor,
  borderWidth,
  padding,
  children,
  ...props
}: SectionItemProps) => {
  const {sizes} = useTheme();

  const radius_ = useMemo<BoxProps['radius']>(() => {
    if (radius === void 0) {
      if (variant === 'row-start') {
        return [1, 1, 0, 0];
      }
      if (variant === 'row-end') {
        return [0, 0, 1, 1];
      }
      if (variant === 'row') {
        return [0, 0, 0, 0];
      }
    }
    return radius;
  }, [radius, variant]);

  const padding_ = useMemo<BoxProps['padding']>(() => {
    if (padding === void 0) {
      if (variant === 'row-start') {
        return [1.5 + sizes.borderRatio / 2, 1.5];
      }
      if (variant === 'row-end') {
        return [1.5 + sizes.borderRatio / 2, 1.5];
      }
      if (variant === 'row') {
        return [1.5 + sizes.borderRatio, 1.5];
      }
      return 1.5;
    }
    return padding;
  }, [padding, variant]);

  const borderWidth_ = useMemo<BoxProps['borderWidth']>(() => {
    if (borderWidth === void 0) {
      if (variant === 'row-start') {
        return [1, 1, 0, 1];
      }
      if (variant === 'row-end') {
        return [0, 1, 1, 1];
      }
      if (variant === 'row') {
        return [0, 1, 0, 1];
      }
    }
    return borderWidth;
  }, [borderWidth, variant]);

  return (
    <Box
      radius={radius_}
      padding={padding_}
      borderWidth={borderWidth_}
      style={styles.sectionItem}
      {...props}>
      {children}
    </Box>
  );
};

type SectionListItemProps<Entity extends unknown> = Omit<
  BoxProps,
  keyof BoxPressableProps
> & {
  onPress?: (entity: Entity, index: number, entities: Entity[]) => void;
  onLongPress?: (entity: Entity, index: number, entities: Entity[]) => void;
  onPressIn?: (entity: Entity, index: number, entities: Entity[]) => void;
  onPressOut?: (entity: Entity, index: number, entities: Entity[]) => void;
};

type SectionListProps<Entity extends unknown> = PropsWithChildren<{
  entities?: Entity[];
  renderItem: (
    entity: Entity,
    index: number,
    entities: Entity[],
  ) => SectionItemProps['children'];
  extractKey: (entity: Entity, index: number, entities: Entity[]) => Key;
  itemProps?: SectionListItemProps<Entity>;
  variantItemProps?: Record<
    SectionItemVariant,
    SectionListItemProps<Entity> | undefined
  >;
}>;
export const SectionList = <Entity extends unknown>({
  entities = [],
  renderItem,
  extractKey,
  children,
  itemProps = {},
  variantItemProps = {} as Record<
    SectionItemVariant,
    SectionListItemProps<Entity> | undefined
  >,
}: SectionListProps<Entity>) => {
  if (entities === void 0 || entities.length <= 0) {
    return children ?? <SectionEmpty />;
  }

  return (
    <Fragment>
      {entities.map((entity, index, entities) => {
        const key = extractKey(entity, index, entities);
        const variant = generateSectionItemVariant(index, entities.length);
        const {onPress, onLongPress, onPressIn, onPressOut, ...props} =
          variantItemProps[variant] || itemProps;

        if (onPress) {
          (props as BoxProps).onPress = event =>
            onPress(entity, index, entities);
        }
        if (onLongPress) {
          (props as BoxProps).onLongPress = event =>
            onLongPress(entity, index, entities);
        }
        if (onPressIn) {
          (props as BoxProps).onPressIn = event =>
            onPressIn(entity, index, entities);
        }
        if (onPressOut) {
          (props as BoxProps).onPressOut = event =>
            onPressOut(entity, index, entities);
        }

        return (
          <SectionItem {...props} key={key} variant={variant}>
            {renderItem(entity, index, entities)}
          </SectionItem>
        );
      })}
    </Fragment>
  );
};

function generateSectionItemVariant(
  offset: number,
  total: number,
): SectionItemVariant {
  if (total === 1) return 'default';
  if (total > 1 && offset === 0) return 'row-start';
  if (total > 1 && offset === total - 1) return 'row-end';
  return 'row';
}

export const SectionEmpty = ({children, ...props}: PaperProps) => {
  return (
    <Paper padding={[6, 0]} style={{alignItems: 'center'}} {...props}>
      {children ?? (
        <Text font="bold" color="textDisabled" size="title">
          No Content
        </Text>
      )}
    </Paper>
  );
};
