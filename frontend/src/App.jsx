import { AuthProvider } from './context/AuthContext';
import { ThemeProvider } from './context/ThemeContext';
import { RouterProvider } from 'react-router-dom';
import router from './routes/AppRoutes';

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <RouterProvider router={router} />
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
