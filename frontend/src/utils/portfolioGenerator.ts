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

const getLanguageColor = (lang: string): string => {
  const colors: Record<string, string> = {
    JavaScript: '#f7df1e',
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
    Vue: '#41b883',
    HTML: '#e34c26',
    CSS: '#563d7c',
    Shell: '#89e051',
  };
  return colors[lang] || '#6b7280';
};

const generateMinimalTemplate = (data: PortfolioData): string => {
  const { user, languages, repos, goals, totalContributions } = data;

  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>${user.name} - Developer Portfolio</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
      background: #ffffff;
      color: #1a1a2e;
      line-height: 1.6;
      min-height: 100vh;
    }
    .container { max-width: 800px; margin: 0 auto; padding: 60px 24px; }
    .header { text-align: center; margin-bottom: 60px; }
    .avatar {
      width: 120px; height: 120px; border-radius: 50%;
      object-fit: cover; margin-bottom: 20px;
      border: 3px solid #e5e7eb;
    }
    .name { font-size: 2.5rem; font-weight: 700; margin-bottom: 8px; }
    .bio { color: #6b7280; font-size: 1.1rem; max-width: 500px; margin: 0 auto; }
    .links { display: flex; gap: 16px; justify-content: center; margin-top: 20px; }
    .links a {
      color: #6b7280; text-decoration: none; font-size: 0.9rem;
      transition: color 0.2s;
    }
    .links a:hover { color: #3b82f6; }
    .section { margin-bottom: 50px; }
    .section-title {
      font-size: 0.75rem; font-weight: 600; text-transform: uppercase;
      letter-spacing: 0.1em; color: #9ca3af; margin-bottom: 20px;
    }
    .stats-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; text-align: center; }
    .stat-value { font-size: 2rem; font-weight: 700; color: #1a1a2e; }
    .stat-label { font-size: 0.85rem; color: #6b7280; }
    .skills-row { display: flex; flex-wrap: wrap; gap: 8px; }
    .skill-tag {
      padding: 6px 14px; background: #f3f4f6; border-radius: 20px;
      font-size: 0.85rem; color: #374151;
    }
    .language-bar { height: 8px; border-radius: 4px; background: #f3f4f6; overflow: hidden; display: flex; }
    .language-segment { height: 100%; }
    .language-legend { display: flex; flex-wrap: wrap; gap: 12px; margin-top: 12px; }
    .language-item { display: flex; align-items: center; gap: 6px; font-size: 0.85rem; color: #6b7280; }
    .language-dot { width: 10px; height: 10px; border-radius: 50%; }
    .repo-grid { display: grid; gap: 16px; }
    .repo-card {
      padding: 20px; border: 1px solid #e5e7eb; border-radius: 12px;
      transition: border-color 0.2s;
    }
    .repo-card:hover { border-color: #3b82f6; }
    .repo-name { font-weight: 600; color: #1a1a2e; margin-bottom: 6px; }
    .repo-desc { font-size: 0.9rem; color: #6b7280; margin-bottom: 10px; }
    .repo-meta { display: flex; gap: 16px; font-size: 0.8rem; color: #9ca3af; }
    .footer { text-align: center; padding-top: 40px; border-top: 1px solid #e5e7eb; color: #9ca3af; font-size: 0.85rem; }
    .footer a { color: #3b82f6; text-decoration: none; }
  </style>
</head>
<body>
  <div class="container">
    <header class="header">
      ${user.avatar_url ? `<img src="${user.avatar_url}" alt="${user.name}" class="avatar">` : ''}
      <h1 class="name">${user.name}</h1>
      ${user.bio ? `<p class="bio">${user.bio}</p>` : ''}
      <div class="links">
        ${user.github_username ? `<a href="https://github.com/${user.github_username}" target="_blank">GitHub</a>` : ''}
        ${user.zenn_username ? `<a href="https://zenn.dev/${user.zenn_username}" target="_blank">Zenn</a>` : ''}
        ${user.qiita_username ? `<a href="https://qiita.com/${user.qiita_username}" target="_blank">Qiita</a>` : ''}
      </div>
    </header>

    ${totalContributions > 0 ? `
    <section class="section">
      <h2 class="section-title">Statistics</h2>
      <div class="stats-grid">
        <div>
          <div class="stat-value">${totalContributions.toLocaleString()}</div>
          <div class="stat-label">Contributions</div>
        </div>
        <div>
          <div class="stat-value">${repos.length}</div>
          <div class="stat-label">Repositories</div>
        </div>
        <div>
          <div class="stat-value">${languages.length}</div>
          <div class="stat-label">Languages</div>
        </div>
      </div>
    </section>
    ` : ''}

    ${(user.skills_languages || user.skills_frameworks) ? `
    <section class="section">
      <h2 class="section-title">Skills</h2>
      <div class="skills-row">
        ${(user.skills_languages || '').split(',').filter(Boolean).map(s => `<span class="skill-tag">${s.trim()}</span>`).join('')}
        ${(user.skills_frameworks || '').split(',').filter(Boolean).map(s => `<span class="skill-tag">${s.trim()}</span>`).join('')}
      </div>
    </section>
    ` : ''}

    ${languages.length > 0 ? `
    <section class="section">
      <h2 class="section-title">Languages</h2>
      <div class="language-bar">
        ${languages.slice(0, 6).map(l => {
          const total = languages.reduce((sum, lang) => sum + lang.bytes, 0);
          const percent = (l.bytes / total) * 100;
          return `<div class="language-segment" style="width: ${percent}%; background: ${getLanguageColor(l.language)}"></div>`;
        }).join('')}
      </div>
      <div class="language-legend">
        ${languages.slice(0, 6).map(l => {
          const total = languages.reduce((sum, lang) => sum + lang.bytes, 0);
          const percent = ((l.bytes / total) * 100).toFixed(1);
          return `<div class="language-item"><span class="language-dot" style="background: ${getLanguageColor(l.language)}"></span>${l.language} ${percent}%</div>`;
        }).join('')}
      </div>
    </section>
    ` : ''}

    ${repos.length > 0 ? `
    <section class="section">
      <h2 class="section-title">Featured Repositories</h2>
      <div class="repo-grid">
        ${repos.slice(0, 4).map(r => `
          <div class="repo-card">
            <div class="repo-name">${r.name}</div>
            ${r.description ? `<div class="repo-desc">${r.description}</div>` : ''}
            <div class="repo-meta">
              ${r.language ? `<span>${r.language}</span>` : ''}
              ${r.stars > 0 ? `<span>⭐ ${r.stars}</span>` : ''}
              ${r.forks > 0 ? `<span>🔱 ${r.forks}</span>` : ''}
            </div>
          </div>
        `).join('')}
      </div>
    </section>
    ` : ''}

    ${goals.filter(g => g.status === 'active').length > 0 ? `
    <section class="section">
      <h2 class="section-title">Currently Learning</h2>
      <div class="skills-row">
        ${goals.filter(g => g.status === 'active').slice(0, 5).map(g => `<span class="skill-tag">${g.title}</span>`).join('')}
      </div>
    </section>
    ` : ''}

    <footer class="footer">
      Generated with <a href="https://devsync.app" target="_blank">DevSync</a>
    </footer>
  </div>
</body>
</html>`;
};

const generateModernTemplate = (data: PortfolioData): string => {
  const { user, languages, repos, goals, totalContributions, followerCount } = data;

  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>${user.name} - Developer Portfolio</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;600;700&display=swap" rel="stylesheet">
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: 'Space Grotesk', -apple-system, BlinkMacSystemFont, sans-serif;
      background: #0a0a0f;
      color: #e5e7eb;
      line-height: 1.6;
      min-height: 100vh;
    }
    .container { max-width: 1000px; margin: 0 auto; padding: 40px 24px; }
    .hero {
      display: grid; grid-template-columns: auto 1fr; gap: 40px;
      align-items: center; padding: 40px; background: #111118;
      border-radius: 24px; margin-bottom: 40px; border: 1px solid #1f1f2e;
    }
    @media (max-width: 640px) {
      .hero { grid-template-columns: 1fr; text-align: center; }
    }
    .avatar {
      width: 160px; height: 160px; border-radius: 20px;
      object-fit: cover; border: 4px solid #3b82f6;
    }
    .name { font-size: 2.5rem; font-weight: 700; margin-bottom: 8px; }
    .bio { color: #9ca3af; font-size: 1.1rem; margin-bottom: 16px; }
    .links { display: flex; gap: 12px; flex-wrap: wrap; }
    @media (max-width: 640px) { .links { justify-content: center; } }
    .link-btn {
      display: inline-flex; align-items: center; gap: 8px;
      padding: 10px 20px; background: #1f1f2e; border-radius: 10px;
      color: #e5e7eb; text-decoration: none; font-size: 0.9rem;
      transition: all 0.2s;
    }
    .link-btn:hover { background: #3b82f6; }
    .stats-row {
      display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px;
      margin-bottom: 40px;
    }
    @media (max-width: 640px) { .stats-row { grid-template-columns: repeat(2, 1fr); } }
    .stat-card {
      background: #111118; padding: 24px; border-radius: 16px;
      text-align: center; border: 1px solid #1f1f2e;
    }
    .stat-value { font-size: 2.5rem; font-weight: 700; color: #3b82f6; }
    .stat-label { font-size: 0.85rem; color: #6b7280; margin-top: 4px; }
    .section { margin-bottom: 40px; }
    .section-header {
      display: flex; align-items: center; gap: 12px;
      margin-bottom: 20px;
    }
    .section-icon {
      width: 36px; height: 36px; background: #3b82f6;
      border-radius: 10px; display: flex; align-items: center;
      justify-content: center; font-size: 1.2rem;
    }
    .section-title { font-size: 1.2rem; font-weight: 600; }
    .skills-grid { display: flex; flex-wrap: wrap; gap: 10px; }
    .skill-chip {
      padding: 10px 18px; background: linear-gradient(135deg, #1f1f2e, #111118);
      border-radius: 12px; font-size: 0.9rem; border: 1px solid #2d2d3a;
    }
    .languages-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
    @media (max-width: 640px) { .languages-grid { grid-template-columns: repeat(2, 1fr); } }
    .lang-card {
      background: #111118; padding: 16px; border-radius: 12px;
      border: 1px solid #1f1f2e; display: flex; align-items: center; gap: 12px;
    }
    .lang-dot { width: 12px; height: 12px; border-radius: 50%; }
    .lang-name { font-weight: 500; }
    .lang-percent { color: #6b7280; font-size: 0.85rem; margin-left: auto; }
    .repo-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px; }
    @media (max-width: 640px) { .repo-grid { grid-template-columns: 1fr; } }
    .repo-card {
      background: #111118; padding: 24px; border-radius: 16px;
      border: 1px solid #1f1f2e; transition: all 0.2s;
    }
    .repo-card:hover { border-color: #3b82f6; transform: translateY(-2px); }
    .repo-name { font-weight: 600; font-size: 1.1rem; margin-bottom: 8px; color: #3b82f6; }
    .repo-desc { font-size: 0.9rem; color: #9ca3af; margin-bottom: 12px; line-height: 1.5; }
    .repo-meta { display: flex; gap: 16px; font-size: 0.85rem; color: #6b7280; }
    .footer {
      text-align: center; padding-top: 40px; border-top: 1px solid #1f1f2e;
      color: #6b7280; font-size: 0.9rem;
    }
    .footer a { color: #3b82f6; text-decoration: none; }
  </style>
</head>
<body>
  <div class="container">
    <div class="hero">
      ${user.avatar_url ? `<img src="${user.avatar_url}" alt="${user.name}" class="avatar">` : '<div class="avatar" style="background: #1f1f2e;"></div>'}
      <div>
        <h1 class="name">${user.name}</h1>
        ${user.bio ? `<p class="bio">${user.bio}</p>` : ''}
        <div class="links">
          ${user.github_username ? `<a href="https://github.com/${user.github_username}" target="_blank" class="link-btn">📦 GitHub</a>` : ''}
          ${user.zenn_username ? `<a href="https://zenn.dev/${user.zenn_username}" target="_blank" class="link-btn">📝 Zenn</a>` : ''}
          ${user.qiita_username ? `<a href="https://qiita.com/${user.qiita_username}" target="_blank" class="link-btn">📘 Qiita</a>` : ''}
        </div>
      </div>
    </div>

    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-value">${totalContributions.toLocaleString()}</div>
        <div class="stat-label">Contributions</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">${repos.length}</div>
        <div class="stat-label">Repositories</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">${followerCount}</div>
        <div class="stat-label">Followers</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">${languages.length}</div>
        <div class="stat-label">Languages</div>
      </div>
    </div>

    ${(user.skills_languages || user.skills_frameworks) ? `
    <section class="section">
      <div class="section-header">
        <div class="section-icon">💡</div>
        <h2 class="section-title">Skills & Technologies</h2>
      </div>
      <div class="skills-grid">
        ${(user.skills_languages || '').split(',').filter(Boolean).map(s => `<span class="skill-chip">${s.trim()}</span>`).join('')}
        ${(user.skills_frameworks || '').split(',').filter(Boolean).map(s => `<span class="skill-chip">${s.trim()}</span>`).join('')}
      </div>
    </section>
    ` : ''}

    ${languages.length > 0 ? `
    <section class="section">
      <div class="section-header">
        <div class="section-icon">📊</div>
        <h2 class="section-title">Top Languages</h2>
      </div>
      <div class="languages-grid">
        ${languages.slice(0, 6).map(l => {
          const total = languages.reduce((sum, lang) => sum + lang.bytes, 0);
          const percent = ((l.bytes / total) * 100).toFixed(1);
          return `
            <div class="lang-card">
              <span class="lang-dot" style="background: ${getLanguageColor(l.language)}"></span>
              <span class="lang-name">${l.language}</span>
              <span class="lang-percent">${percent}%</span>
            </div>
          `;
        }).join('')}
      </div>
    </section>
    ` : ''}

    ${repos.length > 0 ? `
    <section class="section">
      <div class="section-header">
        <div class="section-icon">🚀</div>
        <h2 class="section-title">Featured Projects</h2>
      </div>
      <div class="repo-grid">
        ${repos.slice(0, 4).map(r => `
          <a href="https://github.com/${r.full_name}" target="_blank" class="repo-card" style="text-decoration: none; color: inherit;">
            <div class="repo-name">${r.name}</div>
            ${r.description ? `<div class="repo-desc">${r.description}</div>` : ''}
            <div class="repo-meta">
              ${r.language ? `<span>🔹 ${r.language}</span>` : ''}
              ${r.stars > 0 ? `<span>⭐ ${r.stars}</span>` : ''}
              ${r.forks > 0 ? `<span>🔱 ${r.forks}</span>` : ''}
            </div>
          </a>
        `).join('')}
      </div>
    </section>
    ` : ''}

    ${goals.filter(g => g.status === 'active').length > 0 ? `
    <section class="section">
      <div class="section-header">
        <div class="section-icon">🎯</div>
        <h2 class="section-title">Learning Goals</h2>
      </div>
      <div class="skills-grid">
        ${goals.filter(g => g.status === 'active').slice(0, 5).map(g => `<span class="skill-chip">📚 ${g.title}</span>`).join('')}
      </div>
    </section>
    ` : ''}

    <footer class="footer">
      Built with <a href="https://devsync.app" target="_blank">DevSync</a>
    </footer>
  </div>
</body>
</html>`;
};

const generateGradientTemplate = (data: PortfolioData): string => {
  const { user, languages, repos, goals, totalContributions, followerCount } = data;

  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>${user.name} - Developer Portfolio</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link href="https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600;700&display=swap" rel="stylesheet">
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: 'Poppins', -apple-system, BlinkMacSystemFont, sans-serif;
      background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
      color: #e5e7eb;
      line-height: 1.6;
      min-height: 100vh;
    }
    .bg-blur {
      position: fixed; top: 0; left: 0; right: 0; bottom: 0;
      background: radial-gradient(circle at 20% 80%, rgba(120, 119, 198, 0.3), transparent 40%),
                  radial-gradient(circle at 80% 20%, rgba(59, 130, 246, 0.3), transparent 40%);
      pointer-events: none; z-index: 0;
    }
    .container { max-width: 900px; margin: 0 auto; padding: 60px 24px; position: relative; z-index: 1; }
    .hero { text-align: center; margin-bottom: 60px; }
    .avatar {
      width: 150px; height: 150px; border-radius: 50%;
      object-fit: cover; margin-bottom: 24px;
      border: 4px solid rgba(255,255,255,0.2);
      box-shadow: 0 20px 40px rgba(0,0,0,0.3);
    }
    .name {
      font-size: 3rem; font-weight: 700;
      background: linear-gradient(135deg, #fff, #a5b4fc);
      -webkit-background-clip: text; -webkit-text-fill-color: transparent;
      margin-bottom: 12px;
    }
    .bio { color: rgba(255,255,255,0.7); font-size: 1.1rem; max-width: 500px; margin: 0 auto 24px; }
    .links { display: flex; gap: 12px; justify-content: center; flex-wrap: wrap; }
    .link-pill {
      display: inline-flex; align-items: center; gap: 8px;
      padding: 12px 24px; background: rgba(255,255,255,0.1);
      backdrop-filter: blur(10px); border-radius: 50px;
      color: #fff; text-decoration: none; font-size: 0.9rem;
      border: 1px solid rgba(255,255,255,0.15);
      transition: all 0.3s;
    }
    .link-pill:hover { background: rgba(255,255,255,0.2); transform: translateY(-2px); }
    .stats-row {
      display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px;
      margin-bottom: 50px;
    }
    @media (max-width: 640px) { .stats-row { grid-template-columns: repeat(2, 1fr); } }
    .stat-card {
      background: rgba(255,255,255,0.05); backdrop-filter: blur(10px);
      padding: 28px; border-radius: 20px; text-align: center;
      border: 1px solid rgba(255,255,255,0.1);
    }
    .stat-value {
      font-size: 2.5rem; font-weight: 700;
      background: linear-gradient(135deg, #a5b4fc, #3b82f6);
      -webkit-background-clip: text; -webkit-text-fill-color: transparent;
    }
    .stat-label { font-size: 0.85rem; color: rgba(255,255,255,0.6); margin-top: 4px; }
    .section { margin-bottom: 50px; }
    .section-title {
      font-size: 1.3rem; font-weight: 600; margin-bottom: 20px;
      display: flex; align-items: center; gap: 12px;
    }
    .title-line { flex: 1; height: 1px; background: linear-gradient(90deg, rgba(255,255,255,0.2), transparent); }
    .glass-card {
      background: rgba(255,255,255,0.05); backdrop-filter: blur(10px);
      padding: 24px; border-radius: 20px;
      border: 1px solid rgba(255,255,255,0.1);
    }
    .skills-cloud { display: flex; flex-wrap: wrap; gap: 10px; }
    .skill-bubble {
      padding: 10px 20px; background: linear-gradient(135deg, rgba(165, 180, 252, 0.2), rgba(59, 130, 246, 0.2));
      border-radius: 50px; font-size: 0.9rem;
      border: 1px solid rgba(255,255,255,0.1);
    }
    .lang-list { display: flex; flex-direction: column; gap: 12px; }
    .lang-item { display: flex; align-items: center; gap: 12px; }
    .lang-bar-bg { flex: 1; height: 8px; background: rgba(255,255,255,0.1); border-radius: 4px; overflow: hidden; }
    .lang-bar-fill { height: 100%; border-radius: 4px; }
    .lang-name { min-width: 100px; font-size: 0.9rem; }
    .lang-percent { min-width: 50px; text-align: right; font-size: 0.85rem; color: rgba(255,255,255,0.6); }
    .projects-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 20px; }
    @media (max-width: 640px) { .projects-grid { grid-template-columns: 1fr; } }
    .project-card {
      background: rgba(255,255,255,0.05); backdrop-filter: blur(10px);
      padding: 24px; border-radius: 20px;
      border: 1px solid rgba(255,255,255,0.1);
      transition: all 0.3s; text-decoration: none; color: inherit;
    }
    .project-card:hover { transform: translateY(-4px); border-color: rgba(165, 180, 252, 0.5); }
    .project-name {
      font-weight: 600; font-size: 1.1rem; margin-bottom: 8px;
      color: #a5b4fc;
    }
    .project-desc { font-size: 0.9rem; color: rgba(255,255,255,0.6); margin-bottom: 12px; }
    .project-meta { display: flex; gap: 16px; font-size: 0.8rem; color: rgba(255,255,255,0.5); }
    .footer { text-align: center; padding-top: 40px; color: rgba(255,255,255,0.5); font-size: 0.9rem; }
    .footer a { color: #a5b4fc; text-decoration: none; }
  </style>
</head>
<body>
  <div class="bg-blur"></div>
  <div class="container">
    <header class="hero">
      ${user.avatar_url ? `<img src="${user.avatar_url}" alt="${user.name}" class="avatar">` : ''}
      <h1 class="name">${user.name}</h1>
      ${user.bio ? `<p class="bio">${user.bio}</p>` : ''}
      <div class="links">
        ${user.github_username ? `<a href="https://github.com/${user.github_username}" target="_blank" class="link-pill">🐙 GitHub</a>` : ''}
        ${user.zenn_username ? `<a href="https://zenn.dev/${user.zenn_username}" target="_blank" class="link-pill">📝 Zenn</a>` : ''}
        ${user.qiita_username ? `<a href="https://qiita.com/${user.qiita_username}" target="_blank" class="link-pill">📘 Qiita</a>` : ''}
      </div>
    </header>

    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-value">${totalContributions.toLocaleString()}</div>
        <div class="stat-label">Contributions</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">${repos.length}</div>
        <div class="stat-label">Repositories</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">${followerCount}</div>
        <div class="stat-label">Followers</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">${languages.length}</div>
        <div class="stat-label">Languages</div>
      </div>
    </div>

    ${(user.skills_languages || user.skills_frameworks) ? `
    <section class="section">
      <h2 class="section-title">✨ Skills <span class="title-line"></span></h2>
      <div class="glass-card">
        <div class="skills-cloud">
          ${(user.skills_languages || '').split(',').filter(Boolean).map(s => `<span class="skill-bubble">${s.trim()}</span>`).join('')}
          ${(user.skills_frameworks || '').split(',').filter(Boolean).map(s => `<span class="skill-bubble">${s.trim()}</span>`).join('')}
        </div>
      </div>
    </section>
    ` : ''}

    ${languages.length > 0 ? `
    <section class="section">
      <h2 class="section-title">📊 Languages <span class="title-line"></span></h2>
      <div class="glass-card">
        <div class="lang-list">
          ${languages.slice(0, 5).map(l => {
            const total = languages.reduce((sum, lang) => sum + lang.bytes, 0);
            const percent = ((l.bytes / total) * 100).toFixed(1);
            return `
              <div class="lang-item">
                <span class="lang-name">${l.language}</span>
                <div class="lang-bar-bg">
                  <div class="lang-bar-fill" style="width: ${percent}%; background: ${getLanguageColor(l.language)}"></div>
                </div>
                <span class="lang-percent">${percent}%</span>
              </div>
            `;
          }).join('')}
        </div>
      </div>
    </section>
    ` : ''}

    ${repos.length > 0 ? `
    <section class="section">
      <h2 class="section-title">🚀 Projects <span class="title-line"></span></h2>
      <div class="projects-grid">
        ${repos.slice(0, 4).map(r => `
          <a href="https://github.com/${r.full_name}" target="_blank" class="project-card">
            <div class="project-name">${r.name}</div>
            ${r.description ? `<div class="project-desc">${r.description}</div>` : ''}
            <div class="project-meta">
              ${r.language ? `<span>${r.language}</span>` : ''}
              ${r.stars > 0 ? `<span>⭐ ${r.stars}</span>` : ''}
            </div>
          </a>
        `).join('')}
      </div>
    </section>
    ` : ''}

    ${goals.filter(g => g.status === 'active').length > 0 ? `
    <section class="section">
      <h2 class="section-title">🎯 Learning <span class="title-line"></span></h2>
      <div class="glass-card">
        <div class="skills-cloud">
          ${goals.filter(g => g.status === 'active').slice(0, 5).map(g => `<span class="skill-bubble">${g.title}</span>`).join('')}
        </div>
      </div>
    </section>
    ` : ''}

    <footer class="footer">
      Crafted with <a href="https://devsync.app" target="_blank">DevSync</a>
    </footer>
  </div>
</body>
</html>`;
};

export const generatePortfolioHTML = (
  data: PortfolioData,
  theme: PortfolioTheme
): string => {
  switch (theme) {
    case 'minimal':
      return generateMinimalTemplate(data);
    case 'modern':
      return generateModernTemplate(data);
    case 'gradient':
      return generateGradientTemplate(data);
    default:
      return generateModernTemplate(data);
  }
};

export const downloadPortfolio = (html: string, filename: string): void => {
  const blob = new Blob([html], { type: 'text/html' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
};
