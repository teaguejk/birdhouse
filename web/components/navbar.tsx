"use client";

import Link from "next/link";
import { GoogleLogin } from "@react-oauth/google";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/button";

export function Navbar() {
  const { user, login, logout, isLoading } = useAuth();

  return (
    <nav className="flex items-center border-b px-6 py-4 gap-6">
      <Link href="/" className="text-lg font-semibold">
        Birdhouse
      </Link>
      <Link href="/" className="text-sm text-muted-foreground hover:text-foreground">
        Feed
      </Link>
      {user?.isAdmin && (
        <Link href="/admin" className="text-sm text-muted-foreground hover:text-foreground">
          Admin
        </Link>
      )}
      <div className="ml-auto flex items-center gap-4">
        {isLoading ? null : user ? (
          <>
            <span className="text-sm text-muted-foreground">{user.name || user.email}</span>
            <Button variant="outline" size="sm" onClick={logout}>
              Logout
            </Button>
          </>
        ) : (
          <GoogleLogin
            onSuccess={(response) => {
              if (response.credential) {
                login(response.credential);
              }
            }}
            onError={() => console.error("Google login failed")}
            size="medium"
            theme="outline"
          />
        )}
      </div>
    </nav>
  );
}
