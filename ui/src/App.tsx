
import { useEffect } from 'react';
import RootLayout from './components/RootLayout';
import { usePlaygroundStore } from './lib/playgroundStore';

function App() {
    const [initPlayground] = usePlaygroundStore((state) => [state.initPlayground]);

    useEffect(() => {
        void initPlayground();
    }, [initPlayground])

    return <>
        <RootLayout />
    </>
}

export default App
