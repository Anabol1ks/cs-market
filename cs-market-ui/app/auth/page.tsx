'use client';

import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Stamp as Steam, UserCircle2 } from 'lucide-react';
import Image from 'next/image';

export default function AuthPage() {
  const handleSteamLogin = () => {
		window.location.href = 'http://localhost:8080/auth/steam'
	}

  return (
		<div className='min-h-[80vh] flex items-center justify-center'>
			<Card className='w-full max-w-md p-6 space-y-6'>
				<div className='text-center space-y-2'>
					<h1 className='text-2xl font-bold'>Авторизация</h1>
					<p className='text-muted-foreground'>
						Войдите через Steam или VK для доступа к торговой площадке
					</p>
				</div>

				<div className='space-y-4'>
					<Button
						className='w-full h-12 text-lg'
						variant='outline'
						onClick={handleSteamLogin}
					>
						<Steam className='mr-2 h-5 w-5' />
						Войти через Steam
					</Button>

					<Button className='w-full h-12 text-lg bg-[#0077FF] hover:bg-[#0066CC]'>
						<Image
							src='https://upload.wikimedia.org/wikipedia/commons/2/21/VK.com-logo.svg'
							alt='VK Logo'
							width={24}
							height={24}
							className='mr-2'
						/>
						Войти через VK
					</Button>
				</div>

				<div className='text-center text-sm text-muted-foreground'>
					Авторизуясь, вы соглашаетесь с правилами использования сервиса
				</div>
			</Card>
		</div>
	)
}