import client from './client';

export interface EmailPreferences {
  email_weekly_report: boolean;
  email_language: string;
}

export const getEmailPreferences = () =>
  client.get<EmailPreferences>('/email-preferences');

export const updateEmailPreferences = (data: Partial<EmailPreferences>) =>
  client.put<EmailPreferences>('/email-preferences', data);
