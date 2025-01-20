'use client'

import { useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Stamp as Steam, UserCircle2 } from 'lucide-react'
import Image from 'next/image'
import { useRouter } from 'next/navigation' // Правильный импорт
import Coockies from 'js-cookie'


export default function AuthPage() {
	const router = useRouter()
	
	const steamAuthUrl = `${process.env.NEXT_PUBLIC_API_URL}/auth/steam`
	const handleSteamLogin = () => {
		if (steamAuthUrl) {
      window.location.href = steamAuthUrl; // Использование переменной из .env
    } else {
      console.error('Steam authentication URL is not defined in .env');
    }
	}
	useEffect(() => {
		const params = new URLSearchParams(window.location.search)
		const accessToken = params.get('access_token')
		const refreshToken = params.get('refresh_token')

		if (accessToken && refreshToken) {
			Coockies.set('access_token', accessToken)
			Coockies.set('refresh_token', refreshToken)
			router.replace('/') // Перенаправление на защищённую страницу
		}
	}, [router])

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
