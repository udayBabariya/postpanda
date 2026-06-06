import type { ReactNode } from 'react';
import { Plus_Jakarta_Sans } from 'next/font/google';
import clsx from 'clsx';
import '../global.scss';

const jakartaSans = Plus_Jakarta_Sans({
  weight: ['400', '500', '600', '700'],
  subsets: ['latin'],
});

export const metadata = {
  title: 'PostPanda — AI Social Media Scheduling',
  description:
    'Schedule posts, grow your audience, and manage all your social channels from one powerful open-source platform.',
};

export default function LandingLayout({ children }: { children: ReactNode }) {
  return (
    <html className="dark" lang="en">
      <head>
        <link rel="icon" type="image/svg+xml" href="/favicon.svg" />
        <link rel="icon" href="/favicon.ico" sizes="any" />
      </head>
      <body className={clsx(jakartaSans.className, 'bg-[#0e0e0e] text-white antialiased')}>
        {children}
      </body>
    </html>
  );
}
