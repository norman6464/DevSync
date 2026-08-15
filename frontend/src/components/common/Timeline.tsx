import type { ReactNode } from 'react';

interface TimelineEvent {
  id: string;
  title: string;
  date: string;
  description?: string;
  content?: ReactNode;
  active?: boolean;
}

interface TimelineProps {
  events: TimelineEvent[];
  className?: string;
}

export default function Timeline({
  events,
  className = '',
}: TimelineProps) {
  return (
    <div className={`relative ${className}`.trim()}>
      {events.map((event, index) => {
        const isLast = index === events.length - 1;

        return (
          <div key={event.id} className="flex gap-4">
            <div className="flex flex-col items-center">
              <div
                className={`w-3 h-3 rounded-full flex-shrink-0 ${
                  event.active ? 'bg-blue-500' : 'bg-gray-600'
                }`}
              />
              {!isLast && (
                <div className="w-0.5 flex-1 bg-gray-700 min-h-[2rem]" />
              )}
            </div>
            <div className="pb-6">
              <span className="text-xs text-gray-500">{event.date}</span>
              <h4 className={`text-sm font-medium ${event.active ? 'text-white' : 'text-gray-300'}`}>
                {event.title}
              </h4>
              {event.description && (
                <p className="text-xs text-gray-500 mt-1">{event.description}</p>
              )}
              {event.content}
            </div>
          </div>
        );
      })}
    </div>
  );
}
