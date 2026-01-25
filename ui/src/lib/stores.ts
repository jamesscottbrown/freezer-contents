import { writable } from 'svelte/store';

export interface FreezerItem {
    Name: string;
    Date: string;
    Containers: string[];
}

export interface Freezer {
    Name: string;
    Contents: FreezerItem[];
}

export interface AppState {
    Freezers: Freezer[];
}

export const appState = writable<AppState | null>(null);
