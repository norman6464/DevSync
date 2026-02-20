import { useTranslation } from 'react-i18next';

interface Repo {
  id: number;
  name: string;
  full_name: string;
  description: string;
  language: string;
  stars: number;
  forks: number;
}

interface ProfileRepositoriesSectionProps {
  repos: Repo[];
  githubUsername: string | null;
}

export default function ProfileRepositoriesSection({ repos, githubUsername }: ProfileRepositoriesSectionProps) {
  const { t } = useTranslation();

  if (repos.length === 0) return null;

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide">{t('profile.repositories')}</h2>
        {githubUsername && (
          <a href={`https://github.com/${githubUsername}?tab=repositories`} target="_blank" rel="noopener noreferrer" className="text-xs text-gray-500 hover:text-blue-400 transition-colors flex items-center gap-1">
            {t('profile.viewAllOnGitHub')}
            <svg className="w-3 h-3" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" /></svg>
          </a>
        )}
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {repos.slice(0, 6).map((repo) => (
          <a key={repo.id} href={`https://github.com/${repo.full_name}`} target="_blank" rel="noopener noreferrer" className="bg-gray-900 border border-gray-800 rounded-md p-4 hover:border-gray-600 transition-colors group">
            <div className="flex items-start gap-2">
              <svg className="w-4 h-4 text-gray-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z" /></svg>
              <div className="min-w-0 flex-1">
                <div className="font-medium text-sm text-blue-400 group-hover:text-blue-300 truncate">{repo.name}</div>
                {repo.description && <p className="text-xs text-gray-500 mt-1 line-clamp-2">{repo.description}</p>}
                <div className="flex items-center gap-3 mt-2">
                  {repo.language && <span className="flex items-center gap-1 text-xs text-gray-400"><span className="w-2.5 h-2.5 rounded-full bg-yellow-400" />{repo.language}</span>}
                  {repo.stars > 0 && <span className="flex items-center gap-1 text-xs text-gray-400"><svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M11.48 3.499a.562.562 0 0 1 1.04 0l2.125 5.111a.563.563 0 0 0 .475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 0 0-.182.557l1.285 5.385a.562.562 0 0 1-.84.61l-4.725-2.885a.562.562 0 0 0-.586 0L6.982 20.54a.562.562 0 0 1-.84-.61l1.285-5.386a.562.562 0 0 0-.182-.557l-4.204-3.602a.562.562 0 0 1 .321-.988l5.518-.442a.563.563 0 0 0 .475-.345L11.48 3.5Z" /></svg>{repo.stars}</span>}
                  {repo.forks > 0 && <span className="flex items-center gap-1 text-xs text-gray-400"><svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M7.217 10.907a2.25 2.25 0 1 0 0 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186 9.566-5.314m-9.566 7.5 9.566 5.314m0 0a2.25 2.25 0 1 0 3.935 2.186 2.25 2.25 0 0 0-3.935-2.186Zm0-12.814a2.25 2.25 0 1 0 3.933-2.185 2.25 2.25 0 0 0-3.933 2.185Z" /></svg>{repo.forks}</span>}
                </div>
              </div>
            </div>
          </a>
        ))}
      </div>
    </div>
  );
}
