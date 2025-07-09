declare module 'react-native-config' {
  export interface NativeConfig {
    MOCK_ENABLED?: string;
  }

  export const Config: NativeConfig;
  export default Config;
}
