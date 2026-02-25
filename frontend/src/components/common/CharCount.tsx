/** フィールドの文字数カウンター表示。フォームのインプット下に配置する。 */
export default function CharCount({ value, max }: { value: string; max: number }) {
  return (
    <p className="text-xs text-gray-500 text-right mt-1">
      {value.length}/{max}
    </p>
  );
}
