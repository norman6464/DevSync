import { useTranslation } from 'react-i18next';
import { Monitor, Rocket, Target, FolderOpen, FileText, BookOpen, ChevronDown, ChevronUp, type LucideIcon } from 'lucide-react';
import { type RoadmapCategory, type Roadmap } from '../../api/roadmaps';

const CATEGORIES: { value: RoadmapCategory; label: string; Icon: LucideIcon }[] = [
  { value: 'language', label: 'roadmaps.categoryLanguage', Icon: Monitor },
  { value: 'framework', label: 'roadmaps.categoryFramework', Icon: Rocket },
  { value: 'skill', label: 'roadmaps.categorySkill', Icon: Target },
  { value: 'project', label: 'roadmaps.categoryProject', Icon: FolderOpen },
  { value: 'other', label: 'roadmaps.categoryOther', Icon: FileText },
];

const getCategoryInfo = (cat: RoadmapCategory) =>
  CATEGORIES.find(c => c.value === cat) || CATEGORIES[4];

interface RoadmapTemplatesSectionProps {
  templates: Roadmap[];
  showTemplates: boolean;
  expandedTemplate: number | null;
  creating: boolean;
  toggleTemplates: () => void;
  toggleExpandedTemplate: (id: number) => void;
  handleUseTemplate: (id: number) => void;
}

export default function RoadmapTemplatesSection({
  templates,
  showTemplates,
  expandedTemplate,
  creating,
  toggleTemplates,
  toggleExpandedTemplate,
  handleUseTemplate,
}: RoadmapTemplatesSectionProps) {
  const { t } = useTranslation();

  if (templates.length === 0) return null;

  return (
    <div className="mb-6">
      <button
        onClick={toggleTemplates}
        className="w-full flex items-center justify-between p-4 bg-gray-800 border border-gray-700 rounded-md hover:border-gray-600 transition-colors"
      >
        <div className="flex items-center gap-3">
          <BookOpen className="w-5 h-5 text-purple-400" />
          <div className="text-left">
            <h2 className="font-medium text-white">{t('roadmaps.templates')}</h2>
            <p className="text-sm text-gray-400">{t('roadmaps.templatesDescription')}</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-xs text-gray-500 bg-gray-700 px-2 py-1 rounded-full">
            {templates.length} {t('roadmaps.templatesAvailable')}
          </span>
          {showTemplates ? (
            <ChevronUp className="w-5 h-5 text-gray-400" />
          ) : (
            <ChevronDown className="w-5 h-5 text-gray-400" />
          )}
        </div>
      </button>

      {showTemplates && (
        <div className="mt-3 space-y-3">
          {templates.map(template => {
            const catInfo = getCategoryInfo(template.category);
            const CategoryIcon = catInfo.Icon;
            const isExpanded = expandedTemplate === template.id;

            return (
              <div
                key={template.id}
                className="bg-gray-800 border border-gray-700 rounded-md overflow-hidden"
              >
                <div className="p-4">
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex items-start gap-3 min-w-0 flex-1">
                      <CategoryIcon className="w-6 h-6 text-purple-400 flex-shrink-0 mt-0.5" />
                      <div className="min-w-0 flex-1">
                        <div className="flex items-center gap-2 flex-wrap">
                          <h3 className="font-medium text-white">{template.title}</h3>
                          <span className="px-2 py-0.5 text-xs rounded-full bg-purple-500/10 text-purple-400">
                            {t('roadmaps.template')}
                          </span>
                        </div>
                        {template.description && (
                          <p className="text-sm text-gray-400 mt-1">{template.description}</p>
                        )}
                        <div className="flex items-center gap-4 mt-2 text-xs text-gray-500">
                          <span>{t(catInfo.label)}</span>
                          <span>{template.step_count} {t('roadmaps.steps')}</span>
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <button
                        onClick={() => toggleExpandedTemplate(template.id)}
                        className="p-2 text-gray-400 hover:text-white transition-colors"
                        title={t('roadmaps.viewSteps')}
                      >
                        {isExpanded ? (
                          <ChevronUp className="w-4 h-4" />
                        ) : (
                          <ChevronDown className="w-4 h-4" />
                        )}
                      </button>
                      <button
                        onClick={() => handleUseTemplate(template.id)}
                        disabled={creating}
                        className="px-3 py-1.5 text-sm bg-purple-600 hover:bg-purple-500 disabled:bg-purple-800 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
                      >
                        {creating ? t('common.loading') : t('roadmaps.useTemplate')}
                      </button>
                    </div>
                  </div>
                </div>

                {/* Expanded Steps */}
                {isExpanded && template.steps && (
                  <div className="border-t border-gray-700 bg-gray-800/50 px-4 py-3">
                    <p className="text-xs font-medium text-gray-400 uppercase tracking-wide mb-2">
                      {t('roadmaps.templateSteps')}
                    </p>
                    <div className="space-y-1.5">
                      {template.steps.map((step, idx) => (
                        <div key={step.id} className="flex items-start gap-2">
                          <span className="text-xs text-gray-500 font-mono mt-0.5 w-5 text-right flex-shrink-0">
                            {idx + 1}.
                          </span>
                          <div className="min-w-0">
                            <p className="text-sm text-gray-300">{step.title}</p>
                            {step.description && (
                              <p className="text-xs text-gray-500 mt-0.5">{step.description}</p>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
