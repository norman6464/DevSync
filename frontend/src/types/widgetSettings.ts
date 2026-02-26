export interface WidgetConfig {
  key: string;
  visible: boolean;
  order: number;
}

export interface WidgetSettings {
  id: number;
  user_id: number;
  settings: string;
  created_at: string;
  updated_at: string;
}

export interface UpdateWidgetSettingsRequest {
  settings: WidgetConfig[];
}
