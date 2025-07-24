import {Platform} from 'react-native';

const sizes = Platform.select({
  default: {
    base: 8,
    text: 14,
    textRatio: 14 / 8,
    border: 1,
    borderRatio: 1 / 8,
    // padding: 20,
    // // font sizes
    // h1: 44,
    // h2: 40,
    // h3: 32,
    // h4: 24,
    // h5: 18,
    // p: 16,
  },
});

export default sizes;
