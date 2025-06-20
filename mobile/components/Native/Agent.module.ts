import {NativeModules} from 'react-native';

interface AgentNativeModuleSpec {
  execWithJSON(action: string, body: string): Promise<string>;
}

const AgentNativeModule =
  NativeModules.AgentNativeModule as AgentNativeModuleSpec;
export default AgentNativeModule;

export async function exec<T extends Record<string, unknown>>(
  action: string,
  params?: T,
) {
  const body = JSON.stringify(params);
  const results = await AgentNativeModule.execWithJSON(action, body);
  return JSON.parse(results);
}
