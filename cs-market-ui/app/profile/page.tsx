'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  UserCircle2,
  Package,
  History,
  Settings,
  LogOut,
} from 'lucide-react';
import Image from 'next/image';

// Временные данные для примера
const mockUser = {
  id: '1',
  name: 'User123',
  avatar: 'https://images.unsplash.com/photo-1568602471122-7832951cc4c5?w=400',
  balance: 1250.50,
  totalTrades: 45,
  inventoryValue: 3780.25,
  steamLevel: 32,
};

const mockInventory = Array.from({ length: 6 }, (_, i) => ({
  id: `item-${i}`,
  name: `AK-47 | Азимов`,
  wear: 'Прямо с завода',
  price: 234.50 + i * 10,
  image: 'https://images.unsplash.com/photo-1608346128025-1896b97a6fa7?w=400',
}));

export default function ProfilePage() {
  const [activeTab, setActiveTab] = useState('inventory');

  return (
    <div className="space-y-6">
      {/* Профиль */}
      <Card className="p-6">
        <div className="flex flex-col md:flex-row gap-6 items-center md:items-start">
          <div className="relative w-32 h-32 rounded-full overflow-hidden">
            <Image
              src={mockUser.avatar}
              alt="Аватар"
              fill
              className="object-cover"
            />
          </div>
          
          <div className="flex-1 space-y-4 text-center md:text-left">
            <div>
              <h1 className="text-2xl font-bold">{mockUser.name}</h1>
              <p className="text-muted-foreground">Steam Level: {mockUser.steamLevel}</p>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              <div>
                <div className="text-lg font-bold">${mockUser.balance.toFixed(2)}</div>
                <div className="text-sm text-muted-foreground">Баланс</div>
              </div>
              <div>
                <div className="text-lg font-bold">{mockUser.totalTrades}</div>
                <div className="text-sm text-muted-foreground">Сделок</div>
              </div>
              <div>
                <div className="text-lg font-bold">${mockUser.inventoryValue.toFixed(2)}</div>
                <div className="text-sm text-muted-foreground">Стоимость инвентаря</div>
              </div>
            </div>
          </div>

          <div className="flex md:flex-col gap-2">
            <Button variant="outline" size="icon">
              <Settings className="h-4 w-4" />
            </Button>
            <Button variant="outline" size="icon" className="text-destructive">
              <LogOut className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </Card>

      {/* Вкладки */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid grid-cols-3 w-full md:w-[400px]">
          <TabsTrigger value="inventory">Инвентарь</TabsTrigger>
          <TabsTrigger value="trades">История</TabsTrigger>
          <TabsTrigger value="settings">Настройки</TabsTrigger>
        </TabsList>

        <TabsContent value="inventory" className="mt-6">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {mockInventory.map((item) => (
              <Card key={item.id} className="overflow-hidden">
                <div className="relative aspect-[16/9]">
                  <Image
                    src={item.image}
                    alt={item.name}
                    fill
                    className="object-cover"
                  />
                </div>
                <div className="p-4">
                  <h3 className="font-semibold">{item.name}</h3>
                  <p className="text-sm text-muted-foreground">{item.wear}</p>
                  <p className="mt-2 font-bold">${item.price.toFixed(2)}</p>
                </div>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="trades" className="mt-6">
          <Card className="p-6">
            <h3 className="text-xl font-semibold mb-4">История сделок</h3>
            <p className="text-muted-foreground">История сделок пока пуста</p>
          </Card>
        </TabsContent>

        <TabsContent value="settings" className="mt-6">
          <Card className="p-6">
            <h3 className="text-xl font-semibold mb-4">Настройки профиля</h3>
            <p className="text-muted-foreground">Настройки профиля будут доступны позже</p>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}