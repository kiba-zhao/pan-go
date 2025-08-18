import {NativeModules} from 'react-native';

interface AgentNativeModuleSpec {
  execWithJSON(action: string, body: string): Promise<string>;
}

const AgentNativeModule =
  NativeModules.AgentNativeModule as AgentNativeModuleSpec;
export default AgentNativeModule;

export async function exec<Params extends unknown, Result extends unknown>(
  action: string,
  params?: Params,
) {
  const body = params !== void 0 ? JSON.stringify(params) : '';
  const results = await AgentNativeModule.execWithJSON(action, body);
  return JSON.parse(results) as Result;
}
