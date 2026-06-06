'use client';

export const Logo = () => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="60"
      height="60"
      viewBox="0 0 60 60"
      fill="none"
      className="mt-[8px] min-w-[60px] min-h-[60px]"
    >
      {/* Green rounded background */}
      <rect width="60" height="60" rx="14" fill="#22c55e" />
      {/* Left ear */}
      <circle cx="16" cy="16" r="11" fill="#0f172a" />
      {/* Right ear */}
      <circle cx="44" cy="16" r="11" fill="#0f172a" />
      {/* Face */}
      <ellipse cx="30" cy="36" rx="21" ry="20" fill="white" />
      {/* Left eye patch */}
      <ellipse cx="21" cy="32" rx="7" ry="8" transform="rotate(-10 21 32)" fill="#0f172a" />
      {/* Right eye patch */}
      <ellipse cx="39" cy="32" rx="7" ry="8" transform="rotate(10 39 32)" fill="#0f172a" />
      {/* Left eye */}
      <circle cx="21" cy="31.5" r="3.2" fill="white" />
      {/* Right eye */}
      <circle cx="39" cy="31.5" r="3.2" fill="white" />
      {/* Left pupil */}
      <circle cx="21.6" cy="32" r="1.6" fill="#0f172a" />
      {/* Right pupil */}
      <circle cx="39.6" cy="32" r="1.6" fill="#0f172a" />
      {/* Nose */}
      <ellipse cx="30" cy="40" rx="4.2" ry="2.8" fill="#0f172a" />
      {/* Mouth */}
      <path d="M25.5 44 Q30 48 34.5 44" stroke="#0f172a" strokeWidth="1.6" fill="none" strokeLinecap="round" />
    </svg>
  );
};
