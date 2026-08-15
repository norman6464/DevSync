import '@testing-library/jest-dom';
import i18n from '../i18n';

// Node 22+ はグローバルの localStorage を持つが、--localstorage-file 無しでは
// undefined を返すため、jsdom の実装より優先されてテストが落ちる。
// 無いときだけメモリ実装を差し込む（ブラウザ相当の挙動があれば触らない）。
if (globalThis.localStorage == null) {
  const store = new Map<string, string>();
  const memoryStorage: Storage = {
    get length() {
      return store.size;
    },
    key: (index: number) => [...store.keys()][index] ?? null,
    getItem: (key: string) => store.get(key) ?? null,
    setItem: (key: string, value: string) => {
      store.set(key, String(value));
    },
    removeItem: (key: string) => {
      store.delete(key);
    },
    clear: () => {
      store.clear();
    },
  };
  Object.defineProperty(globalThis, 'localStorage', {
    value: memoryStorage,
    writable: true,
    configurable: true,
  });
}

i18n.changeLanguage('ja');
