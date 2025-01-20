'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import Cookies from 'js-cookie';
import { 
  Search, 
  ShoppingCart, 
  Menu,
  X,
  LogIn,
  UserCircle2,
} from 'lucide-react';

export function Navbar() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false); // TODO: Replace with actual auth state

  useEffect(() => {
		const checkAuth = async () => {
			const accessToken = Cookies.get('access_token')

			if (!accessToken) {
				setIsLoggedIn(false)
				return
			}

			try {
				const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/verify`, {
					headers: {
						Authorization: `Bearer ${accessToken}`,
					},
				})

				if (response.ok) {
					setIsLoggedIn(true)
				} else {
					setIsLoggedIn(false)
				}
			} catch (error) {
				setIsLoggedIn(false)
        Cookies.remove('access_token')
			}
		}

		checkAuth()
	}, [])

  return (
    <nav className="bg-background border-b border-border">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-2">
            <span className="text-xl font-bold text-primary">CS2 Market</span>
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-4 flex-1 justify-center max-w-2xl">
            <div className="relative w-full max-w-xl">
              <Input
                type="search"
                placeholder="Поиск скинов..."
                className="w-full pl-10"
              />
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            </div>
          </div>

          {/* Desktop Actions */}
          <div className="hidden md:flex items-center space-x-4">
            <Button variant="ghost" size="icon">
              <ShoppingCart className="h-5 w-5" />
            </Button>
            {isLoggedIn ? (
              <Link href="/profile">
                <Button variant="ghost" size="icon">
                  <UserCircle2 className="h-5 w-5" />
                </Button>
              </Link>
            ) : (
              <Link href="/auth">
                <Button>
                  <LogIn className="mr-2 h-4 w-4" />
                  Войти
                </Button>
              </Link>
            )}
          </div>

          {/* Mobile Menu Button */}
          <div className="md:hidden">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setIsMenuOpen(!isMenuOpen)}
            >
              {isMenuOpen ? (
                <X className="h-5 w-5" />
              ) : (
                <Menu className="h-5 w-5" />
              )}
            </Button>
          </div>
        </div>

        {/* Mobile Menu */}
        {isMenuOpen && (
          <div className="md:hidden py-4 space-y-4">
            <div className="relative">
              <Input
                type="search"
                placeholder="Поиск скинов..."
                className="w-full pl-10"
              />
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            </div>
            <div className="flex flex-col space-y-2">
              <Button variant="ghost" className="justify-start">
                <ShoppingCart className="mr-2 h-5 w-5" />
                Корзина
              </Button>
              {isLoggedIn ? (
                <Link href="/profile">
                  <Button variant="ghost" className="justify-start w-full">
                    <UserCircle2 className="mr-2 h-5 w-5" />
                    Профиль
                  </Button>
                </Link>
              ) : (
                <Link href="/auth">
                  <Button className="justify-start w-full">
                    <LogIn className="mr-2 h-4 w-4" />
                    Войти
                  </Button>
                </Link>
              )}
            </div>
          </div>
        )}
      </div>
    </nav>
  );
}