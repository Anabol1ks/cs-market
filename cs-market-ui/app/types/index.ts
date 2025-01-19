export type WeaponType = 
  | 'Rifle'
  | 'Pistol' 
  | 'SMG'
  | 'Sniper Rifle'
  | 'Shotgun'
  | 'Machine Gun'
  | 'Knife';

export type WearCondition = 
  | 'Factory New'
  | 'Minimal Wear'
  | 'Field-Tested'
  | 'Well-Worn'
  | 'Battle-Scarred';

export type Rarity =
  | 'Consumer Grade'
  | 'Industrial Grade'
  | 'Mil-Spec'
  | 'Restricted'
  | 'Classified'
  | 'Covert'
  | '★ Covert';

export interface Skin {
  id: string;
  name: string;
  weaponType: WeaponType;
  wear: WearCondition;
  floatValue: number;
  price: number;
  rarity: Rarity;
  collection: string;
  imageUrl: string;
  pattern: number;
  stickers: string[];
  createdAt: string;
  popularity: number;
}