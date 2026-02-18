import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { type RoadmapCategory, type Roadmap } from '../api/roadmaps';
import { useRoadmaps, useRoadmapTemplates } from './useRoadmaps';

export function useRoadmapForm() {
  const navigate = useNavigate();
  const {
    roadmaps, loading, saving, activeRoadmaps, completedRoadmaps,
    createRoadmap, updateRoadmap, deleteRoadmap, refetch,
  } = useRoadmaps();
  const { templates, loading: templatesLoading, creating, createFromTemplate } = useRoadmapTemplates();

  const [showForm, setShowForm] = useState(false);
  const [editingRoadmap, setEditingRoadmap] = useState<Roadmap | null>(null);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState<RoadmapCategory>('other');
  const [isPublic, setIsPublic] = useState(false);
  const [showTemplates, setShowTemplates] = useState(false);
  const [expandedTemplate, setExpandedTemplate] = useState<number | null>(null);

  const resetForm = useCallback(() => {
    setTitle('');
    setDescription('');
    setCategory('other');
    setIsPublic(false);
    setEditingRoadmap(null);
    setShowForm(false);
  }, []);

  const handleSubmit = useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    if (editingRoadmap) {
      const result = await updateRoadmap(editingRoadmap.id, { title, description, category, is_public: isPublic });
      if (result) resetForm();
    } else {
      const result = await createRoadmap({ title, description, category, is_public: isPublic });
      if (result) {
        resetForm();
        navigate(`/roadmaps/${result.id}`);
      }
    }
  }, [title, description, category, isPublic, editingRoadmap, updateRoadmap, createRoadmap, resetForm, navigate]);

  const handleEdit = useCallback((roadmap: Roadmap) => {
    setEditingRoadmap(roadmap);
    setTitle(roadmap.title);
    setDescription(roadmap.description);
    setCategory(roadmap.category);
    setIsPublic(roadmap.is_public);
    setShowForm(true);
  }, []);

  const handleUseTemplate = useCallback(async (templateId: number) => {
    const result = await createFromTemplate(templateId);
    if (result) {
      await refetch();
      navigate(`/roadmaps/${result.id}`);
    }
  }, [createFromTemplate, refetch, navigate]);

  const toggleTemplates = useCallback(() => {
    setShowTemplates(prev => !prev);
  }, []);

  const toggleExpandedTemplate = useCallback((templateId: number) => {
    setExpandedTemplate(prev => prev === templateId ? null : templateId);
  }, []);

  return {
    // データ
    roadmaps,
    activeRoadmaps,
    completedRoadmaps,
    templates,
    loading,
    saving,
    templatesLoading,
    creating,
    // フォーム状態
    showForm,
    setShowForm,
    editingRoadmap,
    title,
    setTitle,
    description,
    setDescription,
    category,
    setCategory,
    isPublic,
    setIsPublic,
    // テンプレート状態
    showTemplates,
    expandedTemplate,
    // アクション
    resetForm,
    handleSubmit,
    handleEdit,
    handleUseTemplate,
    deleteRoadmap,
    toggleTemplates,
    toggleExpandedTemplate,
    navigate,
  };
}
