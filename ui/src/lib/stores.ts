import { writable } from "svelte/store";

export interface FreezerItem {
  Name: string;
  Date: string;
  Containers: string[];
}

export interface Freezer {
  Name: string;
  Contents: FreezerItem[];
}

export interface ContainerType {
  Value: string;
  Label: string;
  Prefix: string;
}

export interface AppState {
  ContainerTypes: ContainerType[];
  Freezers: Freezer[];
}

export const appState = writable<AppState | null>(null);

// Error store for displaying error messages to the user
export const errorMessage = writable<string | null>(null);

// Loading states for different operations
export interface LoadingState {
  initialLoad: boolean;
  addItem: boolean;
  removeContainer: boolean;
  moveContainer: boolean;
  editItem: boolean;
}

export const loadingState = writable<LoadingState>({
  initialLoad: true,
  addItem: false,
  removeContainer: false,
  moveContainer: false,
  editItem: false,
});

// Helper to clear error after a delay
export function clearErrorAfterDelay(delayMs: number = 5000) {
  setTimeout(() => {
    errorMessage.set(null);
  }, delayMs);
}
