import Link from 'next/link';
import Image from 'next/image';

const PLATFORMS = [
  { name: 'X', src: '/icons/platforms/x.png' },
  { name: 'Instagram', src: '/icons/platforms/instagram.png' },
  { name: 'LinkedIn', src: '/icons/platforms/linkedin.png' },
  { name: 'Facebook', src: '/icons/platforms/facebook.png' },
  { name: 'TikTok', src: '/icons/platforms/tiktok.png' },
  { name: 'YouTube', src: '/icons/platforms/youtube.png' },
  { name: 'Discord', src: '/icons/platforms/discord.png' },
  { name: 'Reddit', src: '/icons/platforms/reddit.png' },
  { name: 'Threads', src: '/icons/platforms/threads.png' },
  { name: 'Pinterest', src: '/icons/platforms/pinterest.png' },
  { name: 'Bluesky', src: '/icons/platforms/bluesky.png' },
  { name: 'Mastodon', src: '/icons/platforms/mastodon.png' },
];

const FEATURES = [
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M12 2a10 10 0 1 0 10 10H12V2z"/><path d="M12 2a10 10 0 0 1 10 10"/><path d="M12 12l4-4"/>
      </svg>
    ),
    title: 'AI-Powered Content',
    desc: 'Generate engaging captions, hashtags, and visuals using built-in AI — optimized automatically for each platform\'s best practices.',
  },
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <rect x="3" y="4" width="18" height="18" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/>
      </svg>
    ),
    title: 'Smart Scheduling',
    desc: 'Queue posts across 28+ networks simultaneously. Set it once and publish everywhere with AI-optimal timing.',
  },
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M18 20V10M12 20V4M6 20v-6"/>
      </svg>
    ),
    title: 'Analytics & Insights',
    desc: 'Monitor engagement, track follower growth, and discover exactly what content resonates with your audience.',
  },
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/>
      </svg>
    ),
    title: 'Team Collaboration',
    desc: 'Invite teammates, assign roles, and set up approval workflows before posts go live across your channels.',
  },
];

const STEPS = [
  {
    num: '01',
    title: 'Connect Your Accounts',
    desc: 'Link all your social profiles in just a few clicks. Supports 28+ platforms including X, Instagram, LinkedIn, and more.',
  },
  {
    num: '02',
    title: 'Create & Schedule',
    desc: 'Write once, customize per platform. Use AI to generate content and schedule posts at the best time for engagement.',
  },
  {
    num: '03',
    title: 'Grow Your Audience',
    desc: 'Track performance with detailed analytics, iterate on what works, and watch your following grow consistently.',
  },
];

const TESTIMONIALS = [
  {
    name: 'Sarah K.',
    role: 'Content Creator',
    initials: 'SK',
    text: 'PostPanda cut my content workflow in half. The AI suggestions are spot-on and the calendar view makes planning a week of posts feel completely effortless.',
  },
  {
    name: 'Marcus T.',
    role: 'Marketing Lead',
    initials: 'MT',
    text: 'We manage 12 client accounts and PostPanda keeps everything organized. The analytics help us prove ROI to clients every single month.',
  },
  {
    name: 'Priya R.',
    role: 'Startup Founder',
    initials: 'PR',
    text: 'Being open-source means I can self-host and own my data completely. The team collaboration features are exactly what we needed to scale.',
  },
];

const PLANS = [
  {
    name: 'Starter',
    price: 'Free',
    period: 'forever',
    features: ['3 social channels', '10 scheduled posts / month', 'Basic analytics', 'Single user'],
    cta: 'Get Started',
    highlight: false,
  },
  {
    name: 'Pro',
    price: '$19',
    period: 'per month',
    features: ['Unlimited channels', 'Unlimited posts', 'Advanced analytics', 'AI content generation', 'Up to 5 team members'],
    cta: 'Start Free Trial',
    highlight: true,
  },
  {
    name: 'Business',
    price: '$49',
    period: 'per month',
    features: ['Everything in Pro', 'Unlimited team members', 'Priority support', 'Custom branding', 'Full API access'],
    cta: 'Contact Sales',
    highlight: false,
  },
];

const CheckIcon = () => (
  <svg width="16" height="16" viewBox="0 0 16 16" fill="none" className="shrink-0">
    <circle cx="8" cy="8" r="7" stroke="#612bd3" strokeWidth="1.5" />
    <path d="M5 8l2 2 4-4" stroke="#612bd3" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

const ArrowRight = () => (
  <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
    <path d="M3 8h10M9 4l4 4-4 4" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-[#0e0e0e] text-white overflow-x-hidden">

      {/* ── Navbar ── */}
      <nav className="fixed top-0 left-0 right-0 z-50 bg-[#0e0e0e]/80 backdrop-blur-md border-b border-[#252525]">
        <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-2.5">
            <Image src="/postiz.svg" alt="PostPanda" width={30} height={30} />
            <span className="font-bold text-lg tracking-tight">PostPanda</span>
          </Link>

          <div className="hidden md:flex items-center gap-8 text-sm text-[#9c9c9c]">
            <a href="#features" className="hover:text-white transition-colors">Features</a>
            <a href="#how-it-works" className="hover:text-white transition-colors">How it Works</a>
            <a href="#pricing" className="hover:text-white transition-colors">Pricing</a>
          </div>

          <div className="flex items-center gap-3">
            <Link
              href="/auth/login"
              className="text-sm text-[#9c9c9c] hover:text-white transition-colors px-4 py-2 hidden sm:block"
            >
              Sign In
            </Link>
            <Link
              href="/auth"
              className="text-sm bg-[#612bd3] hover:bg-[#7236f1] transition-colors text-white px-5 py-2 rounded-lg font-medium"
            >
              Get Started
            </Link>
          </div>
        </div>
      </nav>

      {/* ── Hero ── */}
      <section className="relative pt-36 pb-28 px-6 flex flex-col items-center text-center overflow-hidden">
        {/* purple glow */}
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[900px] h-[600px] bg-[#612bd3]/15 blur-[140px] rounded-full pointer-events-none" />
        {/* grid dots */}
        <div
          className="absolute inset-0 pointer-events-none opacity-[0.06]"
          style={{
            backgroundImage:
              'radial-gradient(circle, #ffffff 1px, transparent 1px)',
            backgroundSize: '40px 40px',
          }}
        />

        <div className="relative z-10 max-w-4xl mx-auto">
          <div className="inline-flex items-center gap-2 bg-[#612bd3]/10 border border-[#612bd3]/30 text-[#b49de8] text-sm px-4 py-1.5 rounded-full mb-8">
            <span className="w-1.5 h-1.5 bg-[#612bd3] rounded-full animate-pulse" />
            Open-source · Self-hostable · Free to start
          </div>

          <h1 className="text-5xl md:text-7xl font-bold leading-[1.1] mb-6 tracking-tight">
            Schedule.{' '}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#7236f1] to-[#d82d7e]">
              Publish.
            </span>{' '}
            Grow.
          </h1>

          <p className="text-xl text-[#9c9c9c] max-w-2xl mx-auto mb-10 leading-relaxed">
            The ultimate AI-powered social media scheduling tool. Manage all your channels,
            create content with AI, and grow your audience — from one powerful dashboard.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              href="/auth"
              className="flex items-center gap-2 bg-[#612bd3] hover:bg-[#7236f1] transition-all text-white px-8 py-3.5 rounded-xl font-semibold text-base shadow-lg shadow-[#612bd3]/25 w-full sm:w-auto justify-center"
            >
              Get Started Free <ArrowRight />
            </Link>
            <a
              href="#how-it-works"
              className="flex items-center gap-2 bg-[#1a1919] hover:bg-[#252525] transition-all border border-[#2b2a2a] text-white px-8 py-3.5 rounded-xl font-semibold text-base w-full sm:w-auto justify-center"
            >
              See how it works
            </a>
          </div>

          <p className="text-[#555] text-sm mt-6">
            No credit card required &nbsp;·&nbsp; Free forever plan &nbsp;·&nbsp; Self-hostable
          </p>
        </div>
      </section>

      {/* ── Platforms ── */}
      <section className="py-16 px-6 border-y border-[#1e1e1e]">
        <div className="max-w-6xl mx-auto">
          <p className="text-center text-[#555] text-xs font-semibold uppercase tracking-widest mb-10">
            Works with 28+ platforms
          </p>
          <div className="flex flex-wrap items-center justify-center gap-5">
            {PLATFORMS.map((p) => (
              <div key={p.name} className="flex flex-col items-center gap-2 group cursor-default">
                <div className="w-12 h-12 rounded-xl bg-[#1a1919] border border-[#252525] flex items-center justify-center group-hover:border-[#612bd3]/40 transition-colors">
                  <Image src={p.src} alt={p.name} width={28} height={28} className="rounded object-contain" />
                </div>
                <span className="text-[11px] text-[#555] group-hover:text-[#9c9c9c] transition-colors">{p.name}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Features ── */}
      <section id="features" className="py-28 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-5 tracking-tight">
              Everything you need to<br />dominate social media
            </h2>
            <p className="text-[#9c9c9c] text-lg max-w-2xl mx-auto">
              From AI content creation to in-depth analytics, PostPanda has every tool you need
              to build a thriving social presence.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            {FEATURES.map((f) => (
              <div
                key={f.title}
                className="bg-[#111] border border-[#1e1e1e] rounded-2xl p-8 hover:border-[#612bd3]/40 hover:bg-[#1a1919] transition-all group"
              >
                <div className="text-[#612bd3] mb-5 group-hover:text-[#7236f1] transition-colors">
                  {f.icon}
                </div>
                <h3 className="text-xl font-semibold mb-2">{f.title}</h3>
                <p className="text-[#9c9c9c] leading-relaxed text-[15px]">{f.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Stats ── */}
      <section className="py-16 px-6 bg-[#0a0a0a] border-y border-[#1e1e1e]">
        <div className="max-w-4xl mx-auto grid grid-cols-1 sm:grid-cols-3 gap-10 text-center">
          {[
            { value: '28+', label: 'Platforms supported' },
            { value: '10M+', label: 'Posts scheduled' },
            { value: '50k+', label: 'Active users' },
          ].map((s) => (
            <div key={s.label}>
              <div className="text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-[#7236f1] to-[#d82d7e] mb-2">
                {s.value}
              </div>
              <div className="text-[#9c9c9c] text-sm">{s.label}</div>
            </div>
          ))}
        </div>
      </section>

      {/* ── How it works ── */}
      <section id="how-it-works" className="py-28 px-6">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-5 tracking-tight">
              Get up and running in minutes
            </h2>
            <p className="text-[#9c9c9c] text-lg">
              Three simple steps to transform your social media strategy.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-10">
            {STEPS.map((step, i) => (
              <div key={step.num} className="relative flex flex-col">
                {i < STEPS.length - 1 && (
                  <div className="hidden md:block absolute top-9 left-[calc(100%+20px)] right-0 h-px bg-gradient-to-r from-[#612bd3]/40 to-transparent w-[calc(100%-40px)]" />
                )}
                <div className="text-6xl font-black text-[#612bd3]/20 mb-4 leading-none">
                  {step.num}
                </div>
                <h3 className="text-xl font-semibold mb-3">{step.title}</h3>
                <p className="text-[#9c9c9c] text-[15px] leading-relaxed">{step.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Testimonials ── */}
      <section className="py-28 px-6 bg-[#0a0a0a]">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-5 tracking-tight">
              Loved by creators & marketers
            </h2>
            <p className="text-[#9c9c9c] text-lg">
              Join thousands of teams already using PostPanda to grow their audience.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {TESTIMONIALS.map((t) => (
              <div
                key={t.name}
                className="bg-[#111] border border-[#1e1e1e] rounded-2xl p-7 flex flex-col gap-5"
              >
                <div className="flex items-center gap-0.5 text-[#f59e0b]">
                  {Array.from({ length: 5 }).map((_, i) => (
                    <svg key={i} width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                    </svg>
                  ))}
                </div>
                <p className="text-[#c8c8c8] text-sm leading-relaxed flex-1">
                  &ldquo;{t.text}&rdquo;
                </p>
                <div className="flex items-center gap-3 pt-2 border-t border-[#1e1e1e]">
                  <div className="w-9 h-9 rounded-full bg-gradient-to-br from-[#612bd3] to-[#d82d7e] flex items-center justify-center text-xs font-bold shrink-0">
                    {t.initials}
                  </div>
                  <div>
                    <div className="font-semibold text-sm">{t.name}</div>
                    <div className="text-[#9c9c9c] text-xs">{t.role}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Pricing ── */}
      <section id="pricing" className="py-28 px-6">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-5 tracking-tight">
              Simple, transparent pricing
            </h2>
            <p className="text-[#9c9c9c] text-lg">Start free. Scale as you grow. No hidden fees.</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {PLANS.map((plan) => (
              <div
                key={plan.name}
                className={`relative rounded-2xl p-8 flex flex-col gap-6 ${
                  plan.highlight
                    ? 'bg-[#1a1919] border-2 border-[#612bd3] shadow-lg shadow-[#612bd3]/10'
                    : 'bg-[#111] border border-[#1e1e1e]'
                }`}
              >
                {plan.highlight && (
                  <div className="absolute -top-3.5 left-1/2 -translate-x-1/2 bg-[#612bd3] text-white text-xs font-bold px-3 py-1 rounded-full whitespace-nowrap">
                    Most Popular
                  </div>
                )}

                <div>
                  <div className="text-[#9c9c9c] text-sm font-medium mb-2">{plan.name}</div>
                  <div className="flex items-baseline gap-1">
                    <span className="text-4xl font-bold">{plan.price}</span>
                    <span className="text-[#9c9c9c] text-sm">/ {plan.period}</span>
                  </div>
                </div>

                <ul className="space-y-3 flex-1">
                  {plan.features.map((f) => (
                    <li key={f} className="flex items-center gap-2.5 text-sm text-[#c8c8c8]">
                      <CheckIcon />
                      {f}
                    </li>
                  ))}
                </ul>

                <Link
                  href="/auth"
                  className={`block text-center py-3 rounded-xl font-semibold text-sm transition-all ${
                    plan.highlight
                      ? 'bg-[#612bd3] hover:bg-[#7236f1] text-white'
                      : 'bg-[#1e1e1e] hover:bg-[#252525] text-white border border-[#2b2a2a]'
                  }`}
                >
                  {plan.cta}
                </Link>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── CTA Banner ── */}
      <section className="py-24 px-6 bg-[#0a0a0a]">
        <div className="max-w-4xl mx-auto">
          <div className="relative bg-[#111] border border-[#1e1e1e] rounded-3xl p-16 overflow-hidden text-center">
            <div className="absolute inset-0 bg-gradient-to-br from-[#612bd3]/10 via-transparent to-[#d82d7e]/10 pointer-events-none" />
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[300px] bg-[#612bd3]/15 blur-[100px] rounded-full pointer-events-none" />

            <div className="relative z-10">
              <h2 className="text-4xl md:text-5xl font-bold mb-4 tracking-tight">
                Ready to grow your{' '}
                <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#7236f1] to-[#d82d7e]">
                  social presence?
                </span>
              </h2>
              <p className="text-[#9c9c9c] text-lg mb-8 max-w-xl mx-auto">
                Join thousands of creators and marketers who trust PostPanda
                to manage their social media strategy.
              </p>
              <Link
                href="/auth"
                className="inline-flex items-center gap-2 bg-[#612bd3] hover:bg-[#7236f1] transition-all text-white px-10 py-4 rounded-xl font-semibold text-lg shadow-xl shadow-[#612bd3]/25"
              >
                Start for Free <ArrowRight />
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* ── Footer ── */}
      <footer className="border-t border-[#1e1e1e] py-12 px-6">
        <div className="max-w-6xl mx-auto flex flex-col md:flex-row items-center justify-between gap-6">
          <Link href="/" className="flex items-center gap-2">
            <Image src="/postiz.svg" alt="PostPanda" width={24} height={24} />
            <span className="font-semibold">PostPanda</span>
          </Link>

          <div className="flex flex-wrap items-center justify-center gap-6 text-sm text-[#9c9c9c]">
            <a href="#features" className="hover:text-white transition-colors">Features</a>
            <a href="#how-it-works" className="hover:text-white transition-colors">How it Works</a>
            <a href="#pricing" className="hover:text-white transition-colors">Pricing</a>
            <Link href="/auth/login" className="hover:text-white transition-colors">Sign In</Link>
            <Link href="/auth" className="hover:text-white transition-colors">Register</Link>
          </div>

          <p className="text-sm text-[#444]">© 2024 PostPanda · AGPL-3.0</p>
        </div>
      </footer>

    </div>
  );
}
