import React from 'react';

export const LogoTextComponent = () => {
  return (
    <div className="flex items-center gap-2">
      <svg
        width="32"
        height="32"
        viewBox="0 0 32 32"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        {/* Green rounded background */}
        <rect width="32" height="32" rx="8" fill="#22c55e" />
        {/* Left ear */}
        <circle cx="9" cy="9" r="5.5" fill="#0f172a" />
        {/* Right ear */}
        <circle cx="23" cy="9" r="5.5" fill="#0f172a" />
        {/* Face */}
        <ellipse cx="16" cy="19" rx="11" ry="10" fill="white" />
        {/* Left eye patch */}
        <ellipse cx="11.5" cy="17" rx="3.5" ry="4" transform="rotate(-10 11.5 17)" fill="#0f172a" />
        {/* Right eye patch */}
        <ellipse cx="20.5" cy="17" rx="3.5" ry="4" transform="rotate(10 20.5 17)" fill="#0f172a" />
        {/* Left eye */}
        <circle cx="11.5" cy="16.8" r="1.6" fill="white" />
        {/* Right eye */}
        <circle cx="20.5" cy="16.8" r="1.6" fill="white" />
        {/* Left pupil */}
        <circle cx="11.8" cy="17" r="0.8" fill="#0f172a" />
        {/* Right pupil */}
        <circle cx="20.8" cy="17" r="0.8" fill="#0f172a" />
        {/* Nose */}
        <ellipse cx="16" cy="21.5" rx="2.2" ry="1.4" fill="#0f172a" />
        {/* Mouth */}
        <path d="M13.5 23.5 Q16 25.5 18.5 23.5" stroke="#0f172a" strokeWidth="0.9" fill="none" strokeLinecap="round" />
      </svg>
      <span className="text-[20px] font-bold tracking-tight">PostPanda</span>
    </div>
  );
};
