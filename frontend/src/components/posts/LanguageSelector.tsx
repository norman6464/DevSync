import { useTranslation } from 'react-i18next';
import { selectClass } from '../../constants/styles';

const LANGUAGES = [
  'go', 'typescript', 'javascript', 'python', 'java', 'rust', 'ruby', 'php',
  'c', 'cpp', 'csharp', 'swift', 'kotlin', 'dart', 'sql', 'html', 'css',
  'bash', 'yaml', 'json', 'toml', 'dockerfile', 'graphql', 'markdown',
  'xml', 'scala', 'elixir', 'haskell', 'lua', 'plaintext',
];

interface LanguageSelectorProps {
  value: string;
  onChange: (value: string) => void;
}

export default function LanguageSelector({ value, onChange }: LanguageSelectorProps) {
  const { t } = useTranslation();

  return (
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className={selectClass}
    >
      <option value="">{t('post.selectLanguage')}</option>
      {LANGUAGES.map((lang) => (
        <option key={lang} value={lang}>
          {lang}
        </option>
      ))}
    </select>
  );
}
