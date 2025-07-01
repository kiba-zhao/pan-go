import {SafeAreaProvider} from 'react-native-safe-area-context';
import AppScreen from './components/App/App.screen';

const App = () => (
  <SafeAreaProvider>
    <AppScreen />
  </SafeAreaProvider>
);

export default App;
