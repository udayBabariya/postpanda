import Link from 'next/link';
import Image from 'next/image';

// ── Brand SVGs ────────────────────────────────────────────────────────────────
const PandaLogo = ({ size = 36 }: { size?: number }) => (
  <svg width={size} height={size} viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle cx="20" cy="23" r="15" fill="white" />
    <circle cx="8"  cy="9"  r="5.5" fill="#111" />
    <circle cx="32" cy="9"  r="5.5" fill="#111" />
    <circle cx="8"  cy="9"  r="2.5" fill="#333" />
    <circle cx="32" cy="9"  r="2.5" fill="#333" />
    <ellipse cx="14" cy="21" rx="4.5" ry="5"   fill="#111" transform="rotate(-12 14 21)" />
    <ellipse cx="26" cy="21" rx="4.5" ry="5"   fill="#111" transform="rotate(12 26 21)" />
    <circle cx="14" cy="21" r="2.2" fill="white" />
    <circle cx="26" cy="21" r="2.2" fill="white" />
    <circle cx="14.8" cy="21.8" r="1.2" fill="#111" />
    <circle cx="26.8" cy="21.8" r="1.2" fill="#111" />
    <circle cx="15.3" cy="20.7" r="0.4" fill="white" />
    <circle cx="27.3" cy="20.7" r="0.4" fill="white" />
    <ellipse cx="20" cy="27" rx="2.2" ry="1.5" fill="#111" />
    <path d="M17 29.5 Q20 31.5 23 29.5" stroke="#111" strokeWidth="1.3" fill="none" strokeLinecap="round" />
    <ellipse cx="11" cy="27" rx="2.5" ry="1.5" fill="#fca5a5" opacity="0.5" />
    <ellipse cx="29" cy="27" rx="2.5" ry="1.5" fill="#fca5a5" opacity="0.5" />
  </svg>
);

const BambooLeaf = ({ className }: { className?: string }) => (
  <svg className={className} width="80" height="50" viewBox="0 0 80 50" fill="none">
    <path d="M8 45 C22 32 45 14 74 6 C60 20 38 34 8 45Z" fill="#22c55e" opacity="0.25" />
    <path d="M8 45 C22 32 45 14 74 6" stroke="#22c55e" strokeWidth="1.2" opacity="0.3" fill="none" />
  </svg>
);

const BambooStalk = ({ className }: { className?: string }) => (
  <svg className={className} width="16" height="160" viewBox="0 0 16 160" fill="none">
    <rect x="5" y="0"   width="6" height="48" rx="3" fill="#22c55e" opacity="0.18" />
    <rect x="5" y="52"  width="6" height="48" rx="3" fill="#22c55e" opacity="0.18" />
    <rect x="5" y="104" width="6" height="56" rx="3" fill="#22c55e" opacity="0.18" />
    <rect x="3" y="48"  width="10" height="4" rx="2" fill="#22c55e" opacity="0.25" />
    <rect x="3" y="100" width="10" height="4" rx="2" fill="#22c55e" opacity="0.25" />
    <path d="M11 22 C18 16 28 20 22 30 C19 26 14 23 11 22Z" fill="#22c55e" opacity="0.3" />
    <path d="M5 78 C-2 72 -12 76 -6 86 C-3 82 2 79 5 78Z" fill="#22c55e" opacity="0.3" />
  </svg>
);

const CheckIcon = () => (
  <svg width="16" height="16" viewBox="0 0 16 16" fill="none" className="shrink-0">
    <circle cx="8" cy="8" r="7" stroke="#22c55e" strokeWidth="1.5" />
    <path d="M5 8l2 2 4-4" stroke="#22c55e" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

const ArrowRight = () => (
  <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
    <path d="M3 8h10M9 4l4 4-4 4" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

// ── Static data ───────────────────────────────────────────────────────────────
const PLATFORMS = [
  { name: 'X / Twitter',  src: '/icons/platforms/x.png' },
  { name: 'Instagram',    src: '/icons/platforms/instagram.png' },
  { name: 'LinkedIn',     src: '/icons/platforms/linkedin.png' },
  { name: 'Facebook',     src: '/icons/platforms/facebook.png' },
  { name: 'TikTok',       src: '/icons/platforms/tiktok.png' },
  { name: 'YouTube',      src: '/icons/platforms/youtube.png' },
  { name: 'Discord',      src: '/icons/platforms/discord.png' },
  { name: 'Reddit',       src: '/icons/platforms/reddit.png' },
  { name: 'Threads',      src: '/icons/platforms/threads.png' },
  { name: 'Pinterest',    src: '/icons/platforms/pinterest.png' },
  { name: 'Bluesky',      src: '/icons/platforms/bluesky.png' },
  { name: 'Mastodon',     src: '/icons/platforms/mastodon.png' },
  { name: 'Telegram',     src: '/icons/platforms/telegram.png' },
  { name: 'Slack',        src: '/icons/platforms/slack.png' },
  { name: 'Dribbble',     src: '/icons/platforms/dribbble.png' },
  { name: 'Twitch',       src: '/icons/platforms/twitch.png' },
];

const FEATURES = [
  {
    emoji: '🗓️',
    title: 'Smart Scheduling',
    desc: 'Queue posts across 30+ networks simultaneously. Set it once and publish everywhere with AI-optimal timing for peak engagement.',
  },
  {
    emoji: '🤖',
    title: 'AI Content Agent',
    desc: 'Your personal AI agent drafts captions, generates hashtags, and creates images — trained on what performs best for every platform.',
  },
  {
    emoji: '🎨',
    title: 'Design Studio',
    desc: 'Create stunning visuals with a built-in Canva-like editor. Generate AI images, resize for each platform, and publish in one click.',
  },
  {
    emoji: '⚡',
    title: 'Auto Actions',
    desc: 'Set milestone triggers: auto-post, auto-like, or auto-comment when you hit a goal. Your social media works even when you sleep.',
  },
  {
    emoji: '📊',
    title: 'Analytics & Insights',
    desc: 'Per-channel and per-post metrics — impressions, reach, engagement rates, follower growth. Know exactly what content wins.',
  },
  {
    emoji: '👥',
    title: 'Team Collaboration',
    desc: 'Invite teammates, assign Admin or Member roles, set approval workflows, and manage multiple brands from one workspace.',
  },
];

const WHO_FOR = [
  {
    emoji: '✍️',
    title: 'Content Creators',
    subtitle: 'Plan a week of content in under an hour',
    points: [
      'AI-generated captions tailored per platform',
      'Visual calendar to plan your content mix',
      'Recycle evergreen posts automatically',
      'Media library for all your assets',
    ],
  },
  {
    emoji: '🏢',
    title: 'Marketing Teams',
    subtitle: 'Collaborate, approve, and scale',
    points: [
      'Multi-brand workspaces with role-based access',
      'Approval workflows before anything goes live',
      'Manage 100+ channels across all clients',
      'White-label and custom branding options',
    ],
  },
  {
    emoji: '🧑‍💻',
    title: 'Developers & AI',
    subtitle: 'Build on top of PostPanda or automate it',
    points: [
      'REST API + OAuth2 SDK for custom integrations',
      'MCP server — prompt PostPanda from Claude or GPT',
      'n8n and Make.com nodes out of the box',
      'Webhooks for real-time event triggers',
    ],
  },
];

const STEPS = [
  {
    num: '01',
    title: 'Connect Your Accounts',
    desc: 'Link all your social profiles in a few clicks. Supports 30+ platforms including X, Instagram, LinkedIn, TikTok, and more.',
  },
  {
    num: '02',
    title: 'Create & Schedule',
    desc: 'Write once, customize per platform. Let AI draft your content and schedule posts at the optimal time for maximum reach.',
  },
  {
    num: '03',
    title: 'Automate & Grow',
    desc: 'Set auto-actions on milestones, recycle top-performing posts, and watch your following grow with zero manual effort.',
  },
];

const INTEGRATIONS = [
  { name: 'Claude',      emoji: '🤖', desc: 'Prompt PostPanda via MCP' },
  { name: 'ChatGPT',     emoji: '💬', desc: 'AI content generation' },
  { name: 'n8n',         emoji: '⚙️', desc: 'Workflow automation' },
  { name: 'Make.com',    emoji: '🔄', desc: 'No-code automation' },
  { name: 'Zapier',      emoji: '⚡', desc: 'Connect 5,000+ apps' },
  { name: 'REST API',    emoji: '🔌', desc: 'Build your own tools' },
  { name: 'Webhooks',    emoji: '🪝', desc: 'Real-time event triggers' },
  { name: 'OAuth2 SDK',  emoji: '🔐', desc: 'Multi-user integrations' },
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
    features: ['Unlimited channels', 'Unlimited posts', 'Advanced analytics', 'AI content generation', 'Up to 5 team members', 'REST API + webhooks'],
    cta: 'Start Free Trial',
    highlight: true,
  },
  {
    name: 'Business',
    price: '$49',
    period: 'per month',
    features: ['Everything in Pro', 'Unlimited team members', 'Priority support', 'Custom branding', 'Full API access', 'Multi-brand workspaces'],
    cta: 'Contact Sales',
    highlight: false,
  },
];

const FAQS = [
  {
    q: 'Can I recycle and re-use posts?',
    a: 'Yes. Mark any post as "evergreen" and PostPanda will automatically re-schedule it on a custom cadence — great for top-performing content you want to keep circulating.',
  },
  {
    q: 'What platforms are supported?',
    a: 'PostPanda supports 30+ networks including X, Instagram, LinkedIn, Facebook, TikTok, YouTube, Threads, Pinterest, Bluesky, Mastodon, Telegram, Discord, Slack, Reddit, Dribbble, Twitch, Google My Business, WordPress, Medium, Hashnode, and more.',
  },
  {
    q: 'Is there an API I can build on?',
    a: 'Absolutely. Every paid plan includes REST API access with OAuth2 and webhooks. You can build your own posting app on top of PostPanda, or connect it to n8n, Make.com, or Zapier.',
  },
  {
    q: 'Can AI agents control PostPanda?',
    a: 'Yes — PostPanda ships a CLI and MCP server. Connect it to Claude, ChatGPT, Cursor, or Codex and let your AI agents draft, schedule, and publish posts autonomously.',
  },
  {
    q: 'How does team collaboration work?',
    a: 'Admins manage channels, billing, and team settings. Members can draft, schedule, and review posts. You can set up approval workflows so nothing publishes without sign-off.',
  },
  {
    q: 'Can I self-host PostPanda?',
    a: 'Yes! PostPanda is fully open-source (AGPL-3.0). You can self-host on your own infrastructure, own all your data, and customize it however you like.',
  },
  {
    q: 'How is cross-posting handled?',
    a: 'Write your post once, then customize the copy, media, and hashtags per platform before scheduling. One draft — tailored content for every network.',
  },
  {
    q: 'Can I manage multiple brands?',
    a: 'Pro and Business plans support multiple workspaces, each with their own channel sets, team members, and analytics — perfect for agencies managing many clients.',
  },
];

// ── Page ──────────────────────────────────────────────────────────────────────
export default function LandingPage() {
  return (
    <div className="min-h-screen bg-[#0a0a0a] text-white overflow-x-hidden">

      {/* ── Navbar ── */}
      <nav className="fixed top-0 left-0 right-0 z-50 bg-[#0a0a0a]/85 backdrop-blur-md border-b border-[#1a2e1a]">
        <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-2.5">
            <PandaLogo size={34} />
            <span className="font-bold text-lg tracking-tight">PostPanda</span>
          </Link>
          <div className="hidden md:flex items-center gap-8 text-sm text-[#9ca3af]">
            <a href="#features"     className="hover:text-white transition-colors">Features</a>
            <a href="#who-its-for"  className="hover:text-white transition-colors">Who it's for</a>
            <a href="#how-it-works" className="hover:text-white transition-colors">How it Works</a>
            <a href="#pricing"      className="hover:text-white transition-colors">Pricing</a>
            <a href="#faq"          className="hover:text-white transition-colors">FAQ</a>
          </div>
          <div className="flex items-center gap-3">
            <Link href="/auth/login" className="text-sm text-[#9ca3af] hover:text-white transition-colors px-4 py-2 hidden sm:block">
              Sign In
            </Link>
            <Link href="/auth" className="text-sm bg-[#22c55e] hover:bg-[#16a34a] transition-colors text-black px-5 py-2 rounded-lg font-semibold">
              Get Started
            </Link>
          </div>
        </div>
      </nav>

      {/* ── Hero ── */}
      <section className="relative pt-36 pb-28 px-6 flex flex-col items-center text-center overflow-hidden">
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[900px] h-[600px] bg-[#22c55e]/10 blur-[160px] rounded-full pointer-events-none" />
        <div className="absolute inset-0 pointer-events-none opacity-[0.04]"
          style={{ backgroundImage: 'radial-gradient(circle, #ffffff 1px, transparent 1px)', backgroundSize: '40px 40px' }} />
        <BambooStalk className="absolute left-[4%]  top-20 hidden xl:block" />
        <BambooStalk className="absolute left-[8%]  top-36 hidden xl:block opacity-60" />
        <BambooStalk className="absolute right-[4%] top-16 hidden xl:block" />
        <BambooStalk className="absolute right-[9%] top-40 hidden xl:block opacity-60" />
        <BambooLeaf className="absolute top-24 left-[6%]   rotate-[-20deg] hidden lg:block" />
        <BambooLeaf className="absolute bottom-16 right-[5%] rotate-[160deg] hidden lg:block" />
        <BambooLeaf className="absolute top-40 right-[12%] rotate-[30deg]  hidden xl:block opacity-70" />
        <div className="absolute right-0 top-1/2 -translate-y-1/2 opacity-[0.035] pointer-events-none select-none hidden xl:block">
          <svg width="520" height="520" viewBox="0 0 100 100">
            <circle cx="50" cy="57" r="37" fill="#4ade80" />
            <circle cx="21" cy="21" r="13" fill="#4ade80" />
            <circle cx="79" cy="21" r="13" fill="#4ade80" />
          </svg>
        </div>

        <div className="relative z-10 max-w-4xl mx-auto">
          <div className="inline-flex items-center gap-2 bg-[#22c55e]/10 border border-[#22c55e]/30 text-[#86efac] text-sm px-4 py-1.5 rounded-full mb-8">
            <span className="text-base">🐼</span>
            Open-source · Self-hostable · Free to start
          </div>

          <h1 className="text-5xl md:text-7xl font-bold leading-[1.1] mb-6 tracking-tight">
            Run social media{' '}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#4ade80] to-[#2dd4bf]">
              on autopilot
            </span>
            <br />with AI agents.
          </h1>

          <p className="text-xl text-[#9ca3af] max-w-2xl mx-auto mb-10 leading-relaxed">
            Plan, generate, and schedule posts automatically to <strong className="text-white">30+ social networks</strong>.
            Let AI agents draft your content, design visuals, and publish — while you focus on growing.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link href="/auth"
              className="flex items-center gap-2 bg-[#22c55e] hover:bg-[#16a34a] transition-all text-black px-8 py-3.5 rounded-xl font-semibold text-base shadow-lg shadow-[#22c55e]/20 w-full sm:w-auto justify-center">
              Start for Free <ArrowRight />
            </Link>
            <a href="#how-it-works"
              className="flex items-center gap-2 bg-[#111] hover:bg-[#1a1a1a] transition-all border border-[#1e2e1e] text-white px-8 py-3.5 rounded-xl font-semibold text-base w-full sm:w-auto justify-center">
              See how it works
            </a>
          </div>

          <p className="text-[#444] text-sm mt-6">
            No credit card required &nbsp;·&nbsp; Free forever plan &nbsp;·&nbsp; Self-hostable
          </p>

          {/* AI agent logos strip */}
          <div className="mt-12 flex flex-col items-center gap-3">
            <p className="text-[#555] text-xs uppercase tracking-widest">Works with your AI stack</p>
            <div className="flex items-center gap-4 flex-wrap justify-center">
              {['Claude', 'ChatGPT', 'Cursor', 'Codex'].map((ai) => (
                <span key={ai} className="text-xs text-[#9ca3af] bg-[#111] border border-[#1e2e1e] px-3 py-1.5 rounded-full">
                  {ai}
                </span>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* ── Platforms ── */}
      <section className="py-16 px-6 border-y border-[#1a2e1a]">
        <div className="max-w-6xl mx-auto">
          <p className="text-center text-[#4ade80]/60 text-xs font-semibold uppercase tracking-widest mb-10">
            Works with 30+ platforms
          </p>
          <div className="flex flex-wrap items-center justify-center gap-5">
            {PLATFORMS.map((p) => (
              <div key={p.name} className="flex flex-col items-center gap-2 group cursor-default">
                <div className="w-12 h-12 rounded-xl bg-[#111] border border-[#1e1e1e] flex items-center justify-center group-hover:border-[#22c55e]/40 transition-colors">
                  <Image src={p.src} alt={p.name} width={28} height={28} className="rounded object-contain" />
                </div>
                <span className="text-[11px] text-[#555] group-hover:text-[#9ca3af] transition-colors">{p.name}</span>
              </div>
            ))}
            <div className="flex flex-col items-center gap-2 cursor-default">
              <div className="w-12 h-12 rounded-xl bg-[#111] border border-[#1e1e1e] border-dashed flex items-center justify-center">
                <span className="text-[#555] text-xs font-bold">+14</span>
              </div>
              <span className="text-[11px] text-[#555]">and more</span>
            </div>
          </div>
        </div>
      </section>

      {/* ── Features ── */}
      <section id="features" className="py-28 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 bg-[#22c55e]/10 border border-[#22c55e]/20 text-[#86efac] text-xs font-semibold uppercase tracking-widest px-3 py-1 rounded-full mb-5">
              🎋 Features
            </div>
            <h2 className="text-4xl md:text-5xl font-bold mb-5 tracking-tight">
              Everything you need to<br />dominate social media
            </h2>
            <p className="text-[#9ca3af] text-lg max-w-2xl mx-auto">
              From AI content creation to automation and analytics, PostPanda covers your entire
              social media workflow in one dashboard.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
            {FEATURES.map((f) => (
              <div key={f.title}
                className="bg-[#0f0f0f] border border-[#1a1a1a] rounded-2xl p-7 hover:border-[#22c55e]/30 hover:bg-[#111] transition-all group flex flex-col gap-4">
                <div className="text-3xl">{f.emoji}</div>
                <div>
                  <h3 className="text-lg font-semibold mb-2">{f.title}</h3>
                  <p className="text-[#9ca3af] leading-relaxed text-sm">{f.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Stats ── */}
      <section className="py-16 px-6 bg-[#060606] border-y border-[#1a2e1a]">
        <div className="max-w-4xl mx-auto grid grid-cols-1 sm:grid-cols-3 gap-10 text-center">
          {[
            { value: '30+',  label: 'Platforms supported' },
            { value: '10M+', label: 'Posts scheduled' },
            { value: '50k+', label: 'Active users' },
          ].map((s) => (
            <div key={s.label}>
              <div className="text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-[#4ade80] to-[#2dd4bf] mb-2">
                {s.value}
              </div>
              <div className="text-[#9ca3af] text-sm">{s.label}</div>
            </div>
          ))}
        </div>
      </section>

      {/* ── Who it's for ── */}
      <section id="who-its-for" className="py-28 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 bg-[#22c55e]/10 border border-[#22c55e]/20 text-[#86efac] text-xs font-semibold uppercase tracking-widest px-3 py-1 rounded-full mb-5">
              🐼 Who it's for
            </div>
            <h2 className="text-4xl md:text-5xl font-bold mb-5 tracking-tight">
              Built for every type of creator
            </h2>
            <p className="text-[#9ca3af] text-lg max-w-2xl mx-auto">
              Whether you're a solo creator, a marketing team, or an AI developer — PostPanda
              has a workflow that fits perfectly.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {WHO_FOR.map((w) => (
              <div key={w.title}
                className="bg-[#0f0f0f] border border-[#1a1a1a] rounded-2xl p-8 flex flex-col gap-5 hover:border-[#22c55e]/30 transition-all">
                <div className="text-4xl">{w.emoji}</div>
                <div>
                  <h3 className="text-xl font-bold mb-1">{w.title}</h3>
                  <p className="text-[#4ade80] text-sm font-medium">{w.subtitle}</p>
                </div>
                <ul className="space-y-2.5 flex-1">
                  {w.points.map((pt) => (
                    <li key={pt} className="flex items-start gap-2.5 text-sm text-[#c8c8c8]">
                      <CheckIcon />
                      <span>{pt}</span>
                    </li>
                  ))}
                </ul>
                <Link href="/auth"
                  className="text-center py-2.5 rounded-xl border border-[#22c55e]/30 text-[#4ade80] text-sm font-medium hover:bg-[#22c55e]/10 transition-all">
                  Get started →
                </Link>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── How it works ── */}
      <section id="how-it-works" className="py-28 px-6 bg-[#060606]">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 bg-[#22c55e]/10 border border-[#22c55e]/20 text-[#86efac] text-xs font-semibold uppercase tracking-widest px-3 py-1 rounded-full mb-5">
              🎋 How it Works
            </div>
            <h2 className="text-4xl md:text-5xl font-bold mb-5 tracking-tight">
              Up and running in minutes
            </h2>
            <p className="text-[#9ca3af] text-lg">
              Three steps from signup to your first scheduled post.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-10">
            {STEPS.map((step, i) => (
              <div key={step.num} className="relative flex flex-col">
                {i < STEPS.length - 1 && (
                  <div className="hidden md:block absolute top-9 left-[calc(100%+20px)] h-px bg-gradient-to-r from-[#22c55e]/40 to-transparent w-[calc(100%-40px)]" />
                )}
                <div className="text-6xl font-black text-[#22c55e]/15 mb-4 leading-none select-none">
                  {step.num}
                </div>
                <h3 className="text-xl font-semibold mb-3">{step.title}</h3>
                <p className="text-[#9ca3af] text-[15px] leading-relaxed">{step.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Integrations ── */}
      <section className="py-28 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 bg-[#22c55e]/10 border border-[#22c55e]/20 text-[#86efac] text-xs font-semibold uppercase tracking-widest px-3 py-1 rounded-full mb-5">
              🔌 Integrations
            </div>
            <h2 className="text-4xl md:text-5xl font-bold mb-5 tracking-tight">
              Works with your entire stack
            </h2>
            <p className="text-[#9ca3af] text-lg max-w-2xl mx-auto">
              PostPanda plugs into the AI tools and automation platforms you already use.
              Build workflows, let AI agents post autonomously, or wire it all into your CI/CD.
            </p>
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
            {INTEGRATIONS.map((item) => (
              <div key={item.name}
                className="bg-[#0f0f0f] border border-[#1a1a1a] rounded-xl p-5 flex flex-col gap-2 hover:border-[#22c55e]/30 hover:bg-[#111] transition-all text-center">
                <div className="text-3xl">{item.emoji}</div>
                <div className="font-semibold text-sm">{item.name}</div>
                <div className="text-[#9ca3af] text-xs">{item.desc}</div>
              </div>
            ))}
          </div>

          {/* MCP highlight */}
          <div className="mt-8 bg-[#0f0f0f] border border-[#22c55e]/20 rounded-2xl p-8 flex flex-col md:flex-row items-center gap-6">
            <div className="text-5xl shrink-0">🤖</div>
            <div className="flex-1 text-center md:text-left">
              <h3 className="text-xl font-bold mb-2">Let AI agents post for you</h3>
              <p className="text-[#9ca3af] text-sm leading-relaxed">
                PostPanda ships a CLI and MCP server. Connect it to Claude, ChatGPT, Cursor, or Codex
                and let your AI agents draft, schedule, and publish posts autonomously — no human in the loop needed.
              </p>
            </div>
            <Link href="/auth"
              className="shrink-0 bg-[#22c55e] hover:bg-[#16a34a] text-black px-6 py-3 rounded-xl font-semibold text-sm transition-all whitespace-nowrap">
              Connect your AI →
            </Link>
          </div>
        </div>
      </section>

      {/* ── Testimonials ── */}
      <section className="py-28 px-6 bg-[#060606]">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-5 tracking-tight">
              Loved by creators &amp; marketers
            </h2>
            <p className="text-[#9ca3af] text-lg">
              Join thousands of teams already using PostPanda to grow their audience.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {TESTIMONIALS.map((t) => (
              <div key={t.name}
                className="bg-[#0f0f0f] border border-[#1a1a1a] rounded-2xl p-7 flex flex-col gap-5">
                <div className="flex items-center gap-0.5 text-[#f59e0b]">
                  {Array.from({ length: 5 }).map((_, i) => (
                    <svg key={i} width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                    </svg>
                  ))}
                </div>
                <p className="text-[#c8c8c8] text-sm leading-relaxed flex-1">&ldquo;{t.text}&rdquo;</p>
                <div className="flex items-center gap-3 pt-2 border-t border-[#1a1a1a]">
                  <div className="w-9 h-9 rounded-full bg-gradient-to-br from-[#22c55e] to-[#2dd4bf] flex items-center justify-center text-xs font-bold text-black shrink-0">
                    {t.initials}
                  </div>
                  <div>
                    <div className="font-semibold text-sm">{t.name}</div>
                    <div className="text-[#9ca3af] text-xs">{t.role}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── FAQ ── */}
      <section id="faq" className="py-28 px-6">
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 bg-[#22c55e]/10 border border-[#22c55e]/20 text-[#86efac] text-xs font-semibold uppercase tracking-widest px-3 py-1 rounded-full mb-5">
              🐼 FAQ
            </div>
            <h2 className="text-4xl md:text-5xl font-bold mb-5 tracking-tight">
              Common questions
            </h2>
            <p className="text-[#9ca3af] text-lg">Everything you need to know before getting started.</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            {FAQS.map((faq) => (
              <div key={faq.q}
                className="bg-[#0f0f0f] border border-[#1a1a1a] rounded-2xl p-6 hover:border-[#22c55e]/20 transition-all">
                <h3 className="font-semibold text-white mb-2 text-[15px]">{faq.q}</h3>
                <p className="text-[#9ca3af] text-sm leading-relaxed">{faq.a}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Pricing ── */}
      <section id="pricing" className="py-28 px-6 bg-[#060606]">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 bg-[#22c55e]/10 border border-[#22c55e]/20 text-[#86efac] text-xs font-semibold uppercase tracking-widest px-3 py-1 rounded-full mb-5">
              🎋 Pricing
            </div>
            <h2 className="text-4xl md:text-5xl font-bold mb-5 tracking-tight">
              Simple, transparent pricing
            </h2>
            <p className="text-[#9ca3af] text-lg">Start free. Scale as you grow. No hidden fees.</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {PLANS.map((plan) => (
              <div key={plan.name}
                className={`relative rounded-2xl p-8 flex flex-col gap-6 ${
                  plan.highlight
                    ? 'bg-[#111] border-2 border-[#22c55e] shadow-lg shadow-[#22c55e]/10'
                    : 'bg-[#0f0f0f] border border-[#1a1a1a]'
                }`}>
                {plan.highlight && (
                  <div className="absolute -top-3.5 left-1/2 -translate-x-1/2 bg-[#22c55e] text-black text-xs font-bold px-3 py-1 rounded-full whitespace-nowrap">
                    Most Popular 🐼
                  </div>
                )}
                <div>
                  <div className="text-[#9ca3af] text-sm font-medium mb-2">{plan.name}</div>
                  <div className="flex items-baseline gap-1">
                    <span className="text-4xl font-bold">{plan.price}</span>
                    <span className="text-[#9ca3af] text-sm">/ {plan.period}</span>
                  </div>
                </div>
                <ul className="space-y-3 flex-1">
                  {plan.features.map((f) => (
                    <li key={f} className="flex items-center gap-2.5 text-sm text-[#c8c8c8]">
                      <CheckIcon />{f}
                    </li>
                  ))}
                </ul>
                <Link href="/auth"
                  className={`block text-center py-3 rounded-xl font-semibold text-sm transition-all ${
                    plan.highlight
                      ? 'bg-[#22c55e] hover:bg-[#16a34a] text-black'
                      : 'bg-[#151515] hover:bg-[#1e1e1e] text-white border border-[#222]'
                  }`}>
                  {plan.cta}
                </Link>
              </div>
            ))}
          </div>

          <p className="text-center text-[#555] text-sm mt-8">
            All plans include a <span className="text-[#9ca3af]">7-day free trial</span>. No credit card required.
          </p>
        </div>
      </section>

      {/* ── CTA Banner ── */}
      <section className="py-24 px-6">
        <div className="max-w-4xl mx-auto">
          <div className="relative bg-[#0f0f0f] border border-[#1a2e1a] rounded-3xl p-16 overflow-hidden text-center">
            <div className="absolute inset-0 bg-gradient-to-br from-[#22c55e]/8 via-transparent to-[#2dd4bf]/8 pointer-events-none" />
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[300px] bg-[#22c55e]/10 blur-[120px] rounded-full pointer-events-none" />
            <BambooLeaf className="absolute -bottom-4 -left-4  rotate-[20deg]  opacity-40 hidden md:block" />
            <BambooLeaf className="absolute -top-4    -right-4 rotate-[200deg] opacity-40 hidden md:block" />
            <div className="relative z-10">
              <div className="text-5xl mb-4">🐼</div>
              <h2 className="text-4xl md:text-5xl font-bold mb-4 tracking-tight">
                Ready to grow your{' '}
                <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#4ade80] to-[#2dd4bf]">
                  social presence?
                </span>
              </h2>
              <p className="text-[#9ca3af] text-lg mb-8 max-w-xl mx-auto">
                Join thousands of creators and marketers who trust PostPanda
                to manage their social media strategy — on autopilot.
              </p>
              <Link href="/auth"
                className="inline-flex items-center gap-2 bg-[#22c55e] hover:bg-[#16a34a] transition-all text-black px-10 py-4 rounded-xl font-semibold text-lg shadow-xl shadow-[#22c55e]/20">
                Start for Free <ArrowRight />
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* ── Footer ── */}
      <footer className="border-t border-[#1a2e1a] py-12 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="flex flex-col md:flex-row items-start justify-between gap-10 mb-10">
            <div className="flex flex-col gap-3 max-w-xs">
              <Link href="/" className="flex items-center gap-2">
                <PandaLogo size={28} />
                <span className="font-bold text-base">PostPanda</span>
              </Link>
              <p className="text-[#555] text-sm leading-relaxed">
                The open-source AI-powered social media scheduling platform. Schedule, publish, and grow — on autopilot.
              </p>
            </div>

            <div className="flex flex-wrap gap-12 text-sm">
              <div className="flex flex-col gap-3">
                <div className="text-[#4ade80] font-semibold text-xs uppercase tracking-widest">Product</div>
                <a href="#features"     className="text-[#9ca3af] hover:text-white transition-colors">Features</a>
                <a href="#pricing"      className="text-[#9ca3af] hover:text-white transition-colors">Pricing</a>
                <a href="#how-it-works" className="text-[#9ca3af] hover:text-white transition-colors">How it Works</a>
                <a href="#faq"          className="text-[#9ca3af] hover:text-white transition-colors">FAQ</a>
              </div>
              <div className="flex flex-col gap-3">
                <div className="text-[#4ade80] font-semibold text-xs uppercase tracking-widest">Platforms</div>
                <span className="text-[#9ca3af]">X / Twitter</span>
                <span className="text-[#9ca3af]">Instagram</span>
                <span className="text-[#9ca3af]">LinkedIn</span>
                <span className="text-[#9ca3af]">TikTok + 26 more</span>
              </div>
              <div className="flex flex-col gap-3">
                <div className="text-[#4ade80] font-semibold text-xs uppercase tracking-widest">Account</div>
                <Link href="/auth/login" className="text-[#9ca3af] hover:text-white transition-colors">Sign In</Link>
                <Link href="/auth"       className="text-[#9ca3af] hover:text-white transition-colors">Register</Link>
              </div>
            </div>
          </div>

          <div className="border-t border-[#1a1a1a] pt-6 flex flex-col sm:flex-row items-center justify-between gap-4">
            <p className="text-sm text-[#444]">© 2024 PostPanda · Open-source under AGPL-3.0</p>
            <div className="flex items-center gap-1 text-xs text-[#555]">
              <span>🐼</span>
              <span>Made with bamboo & ☕</span>
            </div>
          </div>
        </div>
      </footer>

    </div>
  );
}
