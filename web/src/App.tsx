import { BrowserRouter, Routes, Route } from "react-router-dom";
import { GoogleOAuthProvider } from "@react-oauth/google";
import { AuthProvider } from "@/context/AuthContext";
import { Navbar } from "@/components/navbar";
import Home from "@/routes/Home";
import AdminLayout from "@/routes/admin/AdminLayout";
import DeviceList from "@/routes/admin/DeviceList";
import CreateDevice from "@/routes/admin/CreateDevice";
import DeviceDetail from "@/routes/admin/DeviceDetail";

const GOOGLE_CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID || "";

export default function App() {
  return (
    <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
      <AuthProvider>
        <BrowserRouter>
          <div className="min-h-screen bg-background text-foreground antialiased">
            <Navbar />
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/admin" element={<AdminLayout />}>
                <Route index element={<DeviceList />} />
                <Route path="devices/new" element={<CreateDevice />} />
                <Route path="devices/:id" element={<DeviceDetail />} />
              </Route>
            </Routes>
          </div>
        </BrowserRouter>
      </AuthProvider>
    </GoogleOAuthProvider>
  );
}
