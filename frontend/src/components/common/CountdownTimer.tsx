import { useState, useEffect, useRef } from 'react';

interface CountdownTimerProps {
  targetDate: string;
  label?: string;
  expiredMessage?: string;
  compact?: boolean;
  onExpire?: () => void;
  className?: string;
}

function calculateTimeLeft(targetDate: string) {
  const difference = new Date(targetDate).getTime() - Date.now();

  if (difference <= 0) {
    return { days: 0, hours: 0, minutes: 0, seconds: 0, expired: true };
  }

  return {
    days: Math.floor(difference / (1000 * 60 * 60 * 24)),
    hours: Math.floor((difference / (1000 * 60 * 60)) % 24),
    minutes: Math.floor((difference / (1000 * 60)) % 60),
    seconds: Math.floor((difference / 1000) % 60),
    expired: false,
  };
}

export default function CountdownTimer({
  targetDate,
  label,
  expiredMessage = '期限終了',
  compact = false,
  onExpire,
  className = '',
}: CountdownTimerProps) {
  const [timeLeft, setTimeLeft] = useState(() => calculateTimeLeft(targetDate));
  const onExpireCalled = useRef(false);

  useEffect(() => {
    if (timeLeft.expired && !onExpireCalled.current) {
      onExpireCalled.current = true;
      onExpire?.();
    }
  }, [timeLeft.expired, onExpire]);

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft(calculateTimeLeft(targetDate));
    }, 1000);

    return () => clearInterval(timer);
  }, [targetDate]);

  if (timeLeft.expired) {
    return (
      <div className={`text-gray-400 ${className}`.trim()}>
        {expiredMessage}
      </div>
    );
  }

  const units = [
    { value: timeLeft.days, label: '日' },
    { value: timeLeft.hours, label: '時間' },
    { value: timeLeft.minutes, label: '分' },
    { value: timeLeft.seconds, label: '秒' },
  ];

  return (
    <div className={`${className}`.trim()}>
      {label && <div className="text-gray-400 text-sm mb-2">{label}</div>}
      <div className={`flex items-center gap-3 ${compact ? 'text-sm' : ''}`}>
        {units.map((unit) => (
          <div key={unit.label} className="flex items-baseline gap-1">
            <span className={`font-bold text-white ${compact ? 'text-sm' : 'text-2xl'}`}>
              {unit.value}
            </span>
            <span className="text-gray-400 text-xs">{unit.label}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
