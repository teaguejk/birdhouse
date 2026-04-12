import { Link } from "react-router-dom";
import { GoogleLogin } from "@react-oauth/google";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

export function Navbar() {
  const { user, login, logout, isLoading } = useAuth();

  return (
    <nav className="flex items-center border-b px-6 py-4 gap-6">
      <Link to="/" className="text-lg font-semibold">
        Birdhouse
      </Link>
      <Link to="/" className="text-sm text-muted-foreground hover:text-foreground">
        Home
      </Link>
      {user?.isAdmin && (
        <Link to="/admin" className="text-sm text-muted-foreground hover:text-foreground">
          Admin
        </Link>
      )}
      <div className="ml-auto flex items-center gap-3">
        {isLoading ? null : user ? (
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2">
              <span className="text-sm leading-none text-muted-foreground">{user.name || user.email}</span>
              {user.isAdmin && <Badge variant="secondary" className="leading-none">Admin</Badge>}
            </div>
            <Button variant="outline" size="sm" onClick={logout}>
              Logout
            </Button>
          </div>
        ) : (
          <GoogleLogin
            onSuccess={(response) => {
              if (response.credential) {
                login(response.credential);
              }
            }}
            onError={() => console.error("Google login failed")}
            size="medium"
            theme="filled_black"
          />
        )}
      </div>
    </nav>
  );
}
