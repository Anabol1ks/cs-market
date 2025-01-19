'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Slider } from '@/components/ui/slider';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { WeaponType, WearCondition, Rarity } from '@/app/types';

const weaponTypes: WeaponType[] = [
  'Rifle',
  'Pistol',
  'SMG',
  'Sniper Rifle',
  'Shotgun',
  'Machine Gun',
  'Knife',
];

const wearConditions: WearCondition[] = [
  'Factory New',
  'Minimal Wear',
  'Field-Tested',
  'Well-Worn',
  'Battle-Scarred',
];

const rarities: Rarity[] = [
  'Consumer Grade',
  'Industrial Grade',
  'Mil-Spec',
  'Restricted',
  'Classified',
  'Covert',
  '★ Covert',
];

export function FilterSidebar() {
  const [priceRange, setPriceRange] = useState([0, 1000]);
  const [floatRange, setFloatRange] = useState([0, 1]);

  return (
    <div className="space-y-6">
      <div className="space-y-4">
        <h3 className="text-lg font-semibold">Filters</h3>
        
        {/* Price Range */}
        <div className="space-y-2">
          <Label>Price Range ($)</Label>
          <div className="flex items-center space-x-2">
            <Input
              type="number"
              value={priceRange[0]}
              onChange={(e) => setPriceRange([+e.target.value, priceRange[1]])}
              className="w-24"
            />
            <span>-</span>
            <Input
              type="number"
              value={priceRange[1]}
              onChange={(e) => setPriceRange([priceRange[0], +e.target.value])}
              className="w-24"
            />
          </div>
          <Slider
            value={priceRange}
            min={0}
            max={1000}
            step={1}
            onValueChange={setPriceRange}
          />
        </div>

        <Separator />

        {/* Float Range */}
        <div className="space-y-2">
          <Label>Float Value</Label>
          <div className="flex items-center space-x-2">
            <Input
              type="number"
              value={floatRange[0]}
              onChange={(e) => setFloatRange([+e.target.value, floatRange[1]])}
              className="w-24"
              step="0.0001"
            />
            <span>-</span>
            <Input
              type="number"
              value={floatRange[1]}
              onChange={(e) => setFloatRange([floatRange[0], +e.target.value])}
              className="w-24"
              step="0.0001"
            />
          </div>
          <Slider
            value={floatRange}
            min={0}
            max={1}
            step={0.0001}
            onValueChange={setFloatRange}
          />
        </div>

        <Separator />

        {/* Weapon Type */}
        <div className="space-y-2">
          <Label>Weapon Type</Label>
          <Select>
            <SelectTrigger>
              <SelectValue placeholder="Select weapon type" />
            </SelectTrigger>
            <SelectContent>
              {weaponTypes.map((type) => (
                <SelectItem key={type} value={type}>
                  {type}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Wear Condition */}
        <div className="space-y-2">
          <Label>Wear Condition</Label>
          <Select>
            <SelectTrigger>
              <SelectValue placeholder="Select condition" />
            </SelectTrigger>
            <SelectContent>
              {wearConditions.map((condition) => (
                <SelectItem key={condition} value={condition}>
                  {condition}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Rarity */}
        <div className="space-y-2">
          <Label>Rarity</Label>
          <Select>
            <SelectTrigger>
              <SelectValue placeholder="Select rarity" />
            </SelectTrigger>
            <SelectContent>
              {rarities.map((rarity) => (
                <SelectItem key={rarity} value={rarity}>
                  {rarity}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <Separator />

        {/* Apply Filters Button */}
        <Button className="w-full">Apply Filters</Button>
      </div>
    </div>
  );
}