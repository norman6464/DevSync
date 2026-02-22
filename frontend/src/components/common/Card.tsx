import { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  onClick?: () => void;
  hoverable?: boolean;
  className?: string;
}

interface CardHeaderProps {
  children: ReactNode;
  className?: string;
}

interface CardBodyProps {
  children: ReactNode;
  className?: string;
}

interface CardFooterProps {
  children: ReactNode;
  className?: string;
}

interface CardTitleProps {
  children: ReactNode;
  icon?: ReactNode;
  className?: string;
}

function Card({ children, onClick, hoverable, className = '' }: CardProps) {
  const baseClasses = 'bg-gray-900 border border-gray-800 rounded-lg p-6';
  const hoverClasses = hoverable || onClick ? 'hover:border-gray-700 transition-colors' : '';
  const clickableClasses = onClick ? 'cursor-pointer' : '';

  return (
    <div
      className={`${baseClasses} ${hoverClasses} ${clickableClasses} ${className}`.trim()}
      onClick={onClick}
    >
      {children}
    </div>
  );
}

function CardHeader({ children, className = '' }: CardHeaderProps) {
  return (
    <div className={`mb-4 ${className}`.trim()}>
      {children}
    </div>
  );
}

function CardBody({ children, className = '' }: CardBodyProps) {
  return (
    <div className={className}>
      {children}
    </div>
  );
}

function CardFooter({ children, className = '' }: CardFooterProps) {
  return (
    <div className={`mt-4 pt-4 border-t border-gray-800 ${className}`.trim()}>
      {children}
    </div>
  );
}

function CardTitle({ children, icon, className = '' }: CardTitleProps) {
  return (
    <div className={`flex items-center gap-2 ${className}`.trim()}>
      {icon && <div className="text-blue-400">{icon}</div>}
      <h3 className="text-lg font-semibold text-white">{children}</h3>
    </div>
  );
}

Card.Header = CardHeader;
Card.Body = CardBody;
Card.Footer = CardFooter;
Card.Title = CardTitle;

export default Card;
