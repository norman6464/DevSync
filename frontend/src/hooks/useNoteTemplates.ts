import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import {
  getMyNoteTemplates,
  getDefaultNoteTemplate,
  createNoteTemplate,
  updateNoteTemplate,
  deleteNoteTemplate,
  useNoteTemplate as useTemplateAPI,
  type NoteTemplate,
  type CreateNoteTemplateRequest,
  type UpdateNoteTemplateRequest,
} from '../api/noteTemplates';
import { useAsyncData } from './useAsyncData';

export function useNoteTemplates() {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: templates, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getMyNoteTemplates();
      return data || [];
    },
    { initialData: [] }
  );

  const [localTemplates, setLocalTemplates] = useState<NoteTemplate[] | null>(null);
  const currentTemplates = localTemplates ?? templates;

  const setTemplates = useCallback((updater: NoteTemplate[] | ((prev: NoteTemplate[]) => NoteTemplate[])) => {
    setLocalTemplates(prev => {
      const current = prev ?? templates;
      return typeof updater === 'function' ? updater(current) : updater;
    });
  }, [templates]);

  const handleCreate = useCallback(async (data: CreateNoteTemplateRequest) => {
    setSaving(true);
    try {
      const { data: newTemplate } = await createNoteTemplate(data);
      setTemplates(prev => [newTemplate, ...prev]);
      toast.success(t('noteTemplates.created', { defaultValue: 'テンプレートを作成しました' }));
      return newTemplate;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, setTemplates]);

  const handleUpdate = useCallback(async (templateId: number, data: UpdateNoteTemplateRequest) => {
    setSaving(true);
    try {
      const { data: updated } = await updateNoteTemplate(templateId, data);
      setTemplates(prev => prev.map(t => t.id === updated.id ? updated : t));
      toast.success(t('noteTemplates.updated', { defaultValue: 'テンプレートを更新しました' }));
      return updated;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, setTemplates]);

  const handleDelete = useCallback(async (id: number) => {
    if (!confirm(t('noteTemplates.confirmDelete', { defaultValue: 'このテンプレートを削除しますか？' }))) return false;
    try {
      await deleteNoteTemplate(id);
      setTemplates(prev => prev.filter(t => t.id !== id));
      toast.success(t('noteTemplates.deleted', { defaultValue: 'テンプレートを削除しました' }));
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [t, setTemplates]);

  const handleUseTemplate = useCallback(async (id: number) => {
    try {
      const { data: note } = await useTemplateAPI(id);
      toast.success(t('noteTemplates.noteCreated', { defaultValue: 'テンプレートからノートを作成しました' }));
      return note;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    }
  }, [t]);

  return {
    templates: currentTemplates,
    loading,
    saving,
    createTemplate: handleCreate,
    updateTemplate: handleUpdate,
    deleteTemplate: handleDelete,
    useTemplate: handleUseTemplate,
    refetch,
  };
}

export function useDefaultNoteTemplate() {
  const { data: defaultTemplate, loading } = useAsyncData(
    async () => {
      try {
        const { data } = await getDefaultNoteTemplate();
        return data;
      } catch {
        return null;
      }
    },
    { initialData: null }
  );

  return {
    defaultTemplate,
    loading,
  };
}
