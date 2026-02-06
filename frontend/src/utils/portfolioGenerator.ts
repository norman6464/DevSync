import type { User } from '../types/user';
import type { GitHubLanguageStat, GitHubRepository } from '../types/github';
import type { LearningGoal } from '../api/goals';

export type PortfolioTheme = 'minimal' | 'modern' | 'gradient';

export interface PortfolioData {
  user: User;
  languages: GitHubLanguageStat[];
  repos: GitHubRepository[];
  goals: LearningGoal[];
  totalContributions: number;
  followerCount: number;
  followingCount: number;
}

const languageColors: Record<string, string> = {
  JavaScript: '#f1e05a',
  TypeScript: '#3178c6',
  Python: '#3572A5',
  Go: '#00ADD8',
  Rust: '#dea584',
  Java: '#b07219',
  'C++': '#f34b7d',
  C: '#555555',
  Ruby: '#701516',
  PHP: '#4F5D95',
  Swift: '#F05138',
  Kotlin: '#A97BFF',
  Dart: '#00B4AB',
  HTML: '#e34c26',
  CSS: '#563d7c',
  Vue: '#41b883',
  React: '#61dafb',
};

const getThemeStyles = (theme: PortfolioTheme): string => {
  switch (theme) {
    case 'minimal':
      return `
        :root { --bg: #ffffff; --text: #1a1a1a; --text-secondary: #666666; --card-bg: #f5f5f5; --border: #e0e0e0; --accent: #2563eb; }
        body { background: var(--bg); color: var(--text); }
      `;
    case 'modern':
      return `
        :root { --bg: #0f172a; --text: #f8fafc; --text-secondary: #94a3b8; --card-bg: #1e293b; --border: #334155; --accent: #22c55e; }
        body { background: var(--bg); color: var(--text); }
      `;
    case 'gradient':
      return `
        :root { --bg: linear-gradient(135deg, #667eea 0%, #764ba2 100%); --text: #ffffff; --text-secondary: rgba(255,255,255,0.8); --card-bg: rgba(255,255,255,0.1); --border: rgba(255,255,255,0.2); --accent: #fbbf24; }
        body { background: var(--bg); color: var(--text); min-height: 100vh; }
        .card { backdrop-filter: blur(10px); }
      `;
  }
};

export const generatePortfolioHTML = (data: PortfolioData, theme: PortfolioTheme): string => {
  const { user, languages, repos, goals, totalContributions, followerCount, followingCount } = data;
  const topLanguages = languages.slice(0, 6);
  const topRepos = repos.slice(0, 6);
  const activeGoals = goals.filter(g => g.status === 'active').slice(0, 4);

  const skills = user.skills_languages ? user.skills_languages.split(',') : [];
  const frameworks = user.skills_frameworks ? user.skills_frameworks.split(',') : [];

  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>${user.name} - Developer Portfolio</title>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    ${getThemeStyles(theme)}
    body { font-family: 'Inter', sans-serif; line-height: 1.6; }
    .container { max-width: 1000px; margin: 0 auto; padding: 3rem 1.5rem; }
    .header { text-align: center; margin-bottom: 3rem; }
    .avatar { width: 120px; height: 120px; border-radius: 50%; margin-bottom: 1rem; object-fit: cover; border: 4px solid var(--accent); }
    .name { font-size: 2.5rem; font-weight: 700; margin-bottom: 0.5rem; }
    .bio { color: var(--text-secondary); font-size: 1.1rem; max-width: 600px; margin: 0 auto; }
    .stats { display: flex; justify-content: center; gap: 2rem; margin-top: 1.5rem; flex-wrap: wrap; }
    .stat { text-align: center; }
    .stat-value { font-size: 1.5rem; font-weight: 700; color: var(--accent); }
    .stat-label { font-size: 0.875rem; color: var(--text-secondary); }
    .section { margin-bottom: 3rem; }
    .section-title { font-size: 1.5rem; font-weight: 600; margin-bottom: 1.5rem; padding-bottom: 0.5rem; border-bottom: 2px solid var(--border); }
    .card { background: var(--card-bg); border: 1px solid var(--border); border-radius: 12px; padding: 1.25rem; }
    .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1rem; }
    .lang-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 0.75rem; }
    .lang-item { display: flex; align-items: center; gap: 0.5rem; }
    .lang-dot { width: 12px; height: 12px; border-radius: 50%; }
    .lang-name { font-weight: 500; }
    .lang-count { color: var(--text-secondary); font-size: 0.875rem; }
    .repo-name { font-weight: 600; font-size: 1.1rem; margin-bottom: 0.5rem; color: var(--accent); }
    .repo-desc { color: var(--text-secondary); font-size: 0.9rem; margin-bottom: 0.75rem; }
    .repo-meta { display: flex; gap: 1rem; font-size: 0.875rem; color: var(--text-secondary); }
    .skill-tags { display: flex; flex-wrap: wrap; gap: 0.5rem; }
    .skill-tag { background: var(--card-bg); border: 1px solid var(--border); padding: 0.375rem 0.75rem; border-radius: 9999px; font-size: 0.875rem; }
    .goal-item { margin-bottom: 1rem; }
    .goal-title { font-weight: 500; margin-bottom: 0.5rem; }
    .progress-bar { height: 8px; background: var(--border); border-radius: 4px; overflow: hidden; }
    .progress-fill { height: 100%; background: var(--accent); border-radius: 4px; transition: width 0.3s; }
    .footer { text-align: center; color: var(--text-secondary); font-size: 0.875rem; padding-top: 2rem; border-top: 1px solid var(--border); }
    a { color: var(--accent); text-decoration: none; }
    a:hover { text-decoration: underline; }
  </style>
</head>
<body>
  <div class="container">
    <header class="header">
      ${user.avatar_url ? `<img src="${user.avatar_url}" alt="${user.name}" class="avatar">` : ''}
      <h1 class="name">${user.name}</h1>
      ${user.bio ? `<p class="bio">${user.bio}</p>` : ''}
      <div class="stats">
        <div class="stat">
          <div class="stat-value">${totalContributions.toLocaleString()}</div>
          <div class="stat-label">Contributions</div>
        </div>
        <div class="stat">
          <div class="stat-value">${followerCount}</div>
          <div class="stat-label">Followers</div>
        </div>
        <div class="stat">
          <div class="stat-value">${followingCount}</div>
          <div class="stat-label">Following</div>
        </div>
      </div>
    </header>

    ${topLanguages.length > 0 ? `
    <section class="section">
      <h2 class="section-title">Languages</h2>
      <div class="lang-grid">
        ${topLanguages.map(lang => `
          <div class="lang-item">
            <span class="lang-dot" style="background: ${languageColors[lang.language] || '#888'}"></span>
            <span class="lang-name">${lang.language}</span>
            <span class="lang-count">(${lang.repo_count})</span>
          </div>
        `).join('')}
      </div>
    </section>
    ` : ''}

    ${skills.length > 0 || frameworks.length > 0 ? `
    <section class="section">
      <h2 class="section-title">Skills</h2>
      <div class="skill-tags">
        ${skills.map((s: string) => `<span class="skill-tag">${s}</span>`).join('')}
        ${frameworks.map((f: string) => `<span class="skill-tag">${f}</span>`).join('')}
      </div>
    </section>
    ` : ''}

    ${topRepos.length > 0 ? `
    <section class="section">
      <h2 class="section-title">Projects</h2>
      <div class="grid">
        ${topRepos.map(repo => `
          <div class="card">
            <div class="repo-name">${repo.name}</div>
            ${repo.description ? `<p class="repo-desc">${repo.description}</p>` : ''}
            <div class="repo-meta">
              ${repo.language ? `<span>${repo.language}</span>` : ''}
              <span>⭐ ${repo.stars}</span>
              <span>🍴 ${repo.forks}</span>
            </div>
          </div>
        `).join('')}
      </div>
    </section>
    ` : ''}

    ${activeGoals.length > 0 ? `
    <section class="section">
      <h2 class="section-title">Learning Goals</h2>
      <div class="card">
        ${activeGoals.map(goal => `
          <div class="goal-item">
            <div class="goal-title">${goal.title}</div>
            <div class="progress-bar">
              <div class="progress-fill" style="width: ${goal.progress}%"></div>
            </div>
          </div>
        `).join('')}
      </div>
    </section>
    ` : ''}

    <footer class="footer">
      <p>Generated with DevSync • ${new Date().toLocaleDateString()}</p>
    </footer>
  </div>
</body>
</html>`;
};

export const downloadPortfolio = (html: string, filename: string): void => {
  const blob = new Blob([html], { type: 'text/html' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
};
