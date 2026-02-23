import { useState, useRef } from 'react';
import { UploadCloud } from 'lucide-react';

interface FileUploadProps {
  onUpload: (files: File[]) => void;
  accept?: string;
  multiple?: boolean;
  maxSizeMB?: number;
  disabled?: boolean;
  className?: string;
}

export default function FileUpload({
  onUpload,
  accept,
  multiple = false,
  maxSizeMB,
  disabled = false,
  className = '',
}: FileUploadProps) {
  const [isDragOver, setIsDragOver] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleFiles = (files: FileList | null) => {
    if (!files || disabled) return;
    onUpload(Array.from(files));
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
    if (!disabled) {
      handleFiles(e.dataTransfer.files);
    }
  };

  return (
    <div className={`${className}`.trim()}>
      <div
        data-testid="drop-zone"
        onClick={() => !disabled && inputRef.current?.click()}
        onDragOver={(e) => { e.preventDefault(); setIsDragOver(true); }}
        onDragLeave={() => setIsDragOver(false)}
        onDrop={handleDrop}
        className={`flex flex-col items-center justify-center p-8 border-2 border-dashed rounded-lg cursor-pointer transition-colors ${
          isDragOver ? 'border-blue-500 bg-blue-500/10' : 'border-gray-700 hover:border-gray-600'
        } ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
      >
        <UploadCloud className="w-10 h-10 text-gray-500 mb-3" />
        <p className="text-sm text-gray-300 mb-1">ファイルをドラッグ&ドロップ</p>
        <p className="text-xs text-gray-500">またはクリックして選択</p>
        {accept && (
          <p className="text-xs text-gray-500 mt-2">{accept}</p>
        )}
        {maxSizeMB && (
          <p className="text-xs text-gray-500 mt-1">最大 {maxSizeMB}MB</p>
        )}
      </div>
      <input
        ref={inputRef}
        type="file"
        accept={accept}
        multiple={multiple}
        disabled={disabled}
        onChange={(e) => handleFiles(e.target.files)}
        className="hidden"
      />
    </div>
  );
}
