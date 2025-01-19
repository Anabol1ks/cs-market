'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Grid2X2, List } from 'lucide-react';
import Image from 'next/image';
import { cn } from '@/lib/utils';

// Temporary mock data
const mockSkins = Array.from({ length: 12 }, (_, i) => ({
  id: `skin-${i}`,
  name: `AK-47 | Asiimov ${i}`,
  wear: 'Factory New',
  floatValue: 0.01234,
  price: 199.99 + i,
  rarity: 'Covert',
  imageUrl: 'https://images.unsplash.com/photo-1608346128025-1896b97a6fa7?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3',
}));

export function SkinGrid() {
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Featured Skins</h2>
        <div className="flex items-center space-x-2">
          <Button
            variant={viewMode === 'grid' ? 'default' : 'outline'}
            size="icon"
            onClick={() => setViewMode('grid')}
          >
            <Grid2X2 className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === 'list' ? 'default' : 'outline'}
            size="icon"
            onClick={() => setViewMode('list')}
          >
            <List className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <div className={cn(
        'grid gap-4',
        viewMode === 'grid' 
          ? 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4' 
          : 'grid-cols-1'
      )}>
        {mockSkins.map((skin) => (
          <Card 
            key={skin.id}
            className={cn(
              'overflow-hidden transition-all hover:border-primary',
              viewMode === 'list' && 'flex'
            )}
          >
            <div className={cn(
              'relative',
              viewMode === 'grid' ? 'aspect-[16/9]' : 'w-48 shrink-0'
            )}>
              <Image
                src={skin.imageUrl}
                alt={skin.name}
                fill
                className="object-cover"
              />
            </div>
            <div className="p-4">
              <h3 className="font-semibold truncate">{skin.name}</h3>
              <div className="mt-1 text-sm text-muted-foreground">
                {skin.wear} ({skin.floatValue})
              </div>
              <div className="mt-2 font-bold">${skin.price.toFixed(2)}</div>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}