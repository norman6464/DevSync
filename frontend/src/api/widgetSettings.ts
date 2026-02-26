import client from './client';
import type { WidgetSettings, WidgetConfig } from '../types/widgetSettings';

export const getWidgetSettings = () =>
  client.get<WidgetSettings>('/widget-settings');

export const updateWidgetSettings = (settings: WidgetConfig[]) =>
  client.put('/widget-settings', { settings });
