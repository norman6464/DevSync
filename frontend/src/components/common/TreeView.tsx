import { useState } from 'react';
import { ChevronRight } from 'lucide-react';

interface TreeNode {
  id: string;
  label: string;
  children?: TreeNode[];
}

interface TreeViewProps {
  nodes: TreeNode[];
  onSelect?: (id: string) => void;
  defaultExpandedIds?: string[];
  className?: string;
}

function TreeItem({
  node,
  depth,
  expandedIds,
  toggleExpand,
  onSelect,
}: {
  node: TreeNode;
  depth: number;
  expandedIds: Set<string>;
  toggleExpand: (id: string) => void;
  onSelect?: (id: string) => void;
}) {
  const hasChildren = node.children && node.children.length > 0;
  const isExpanded = expandedIds.has(node.id);

  const handleClick = () => {
    if (hasChildren) {
      toggleExpand(node.id);
    }
    onSelect?.(node.id);
  };

  return (
    <div>
      <button
        type="button"
        onClick={handleClick}
        className="flex items-center gap-1 w-full text-left py-1 px-2 text-sm text-gray-300 hover:bg-gray-800 rounded transition-colors"
        style={{ paddingLeft: `${depth * 16 + 8}px` }}
      >
        {hasChildren && (
          <ChevronRight
            className={`w-4 h-4 text-gray-500 transition-transform flex-shrink-0 ${isExpanded ? 'rotate-90' : ''}`}
          />
        )}
        <span>{node.label}</span>
      </button>
      {hasChildren && isExpanded && (
        <div>
          {node.children!.map((child) => (
            <TreeItem
              key={child.id}
              node={child}
              depth={depth + 1}
              expandedIds={expandedIds}
              toggleExpand={toggleExpand}
              onSelect={onSelect}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export default function TreeView({
  nodes,
  onSelect,
  defaultExpandedIds = [],
  className = '',
}: TreeViewProps) {
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set(defaultExpandedIds));

  const toggleExpand = (id: string) => {
    setExpandedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  return (
    <div className={`${className}`.trim()}>
      {nodes.map((node) => (
        <TreeItem
          key={node.id}
          node={node}
          depth={0}
          expandedIds={expandedIds}
          toggleExpand={toggleExpand}
          onSelect={onSelect}
        />
      ))}
    </div>
  );
}
