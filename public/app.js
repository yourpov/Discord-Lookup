const qs = sel => document.querySelector(sel);

const elements = {
  input: qs('#uid'),
  btn: qs('#search'),
  toast: qs('#notif'),
  msg: qs('#msg'),
  empty: qs('#none'),
  card: qs('#card'),
  pfp: qs('#pfp'),
  banner: qs('#banner'),
  bannerWrap: qs('#banner-wrap'),
  name: qs('#name'),
  tag: qs('#tag'),
  uid: qs('#uid-text'),
  joined: qs('#joined'),
  type: qs('#type'),
  badges: qs('#badges'),
  badgeWrap: qs('#badge-wrap'),
  flags: qs('#flags'),
  age: qs('#age'),
  status: qs('#status'),
  avatarLink: qs('#avatar-link'),
  bannerLinkWrap: qs('#banner-link-wrap')
};

let loading = false;

const notify = (txt, kind = 'err') => {
  elements.msg.textContent = txt;
  const bg = kind === 'err' ? 'bg-red-600' : 'bg-emerald-600';
  elements.toast.className = `fixed top-4 right-4 px-5 py-3 rounded-lg ${bg} text-white font-medium shadow-lg z-50`;
  elements.toast.style.transform = 'translateX(0)';
  elements.toast.style.opacity = '1';
  
  setTimeout(() => {
    elements.toast.style.transform = 'translateX(400px)';
    elements.toast.style.opacity = '0';
  }, 3500);
};

const setLoad = state => {
  loading = state;
  elements.btn.disabled = state;
  elements.btn.innerHTML = state 
    ? '<div class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>' 
    : 'Search';
};

const formatDate = str => {
  if (!str) return 'unknown';
  const d = new Date(str);
  const diff = Math.floor((Date.now() - d) / 86400000);
  let ago;
  if (diff < 60) ago = `${diff}d ago`;
  else if (diff < 365) ago = `${Math.floor(diff / 30)}mo ago`;
  else ago = `${Math.floor(diff / 365)}y ago`;
  return `${d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })} (${ago})`;
};

const getAvatar = disc => {
  const n = parseInt(disc || 0) % 5;
  return `https://cdn.discordapp.com/embed/avatars/${n}.png`;
};

const renderProfile = user => {
  elements.empty.classList.add('hidden');
  elements.card.classList.remove('hidden');
  
  if (user.banner) {
    elements.banner.src = user.banner;
    elements.banner.classList.remove('hidden');
    elements.bannerWrap.style.background = '';
    elements.bannerLinkWrap.href = user.banner;
    elements.bannerLinkWrap.classList.remove('hidden');
  } else {
    elements.banner.classList.add('hidden');
    elements.bannerLinkWrap.classList.add('hidden');
    if (user.accent_color) {
      const hex = '#' + user.accent_color.toString(16).padStart(6, '0');
      elements.bannerWrap.style.background = `linear-gradient(135deg, ${hex}, ${hex}dd)`;
    }
  }
  
  const avatarUrl = user.avatar || getAvatar(user.discriminator);
  elements.pfp.src = avatarUrl;
  elements.avatarLink.href = avatarUrl;
  elements.name.textContent = user.display_name || user.username || 'unknown';
  elements.tag.textContent = '@' + (user.username || 'unknown');
  elements.uid.textContent = user.id || 'n/a';
  elements.joined.textContent = formatDate(user.created_at);
  
  elements.status.className = 'absolute bottom-0 right-0 w-8 h-8 rounded-full border-4 border-[#2b2d31] bg-green-500';
  
  elements.type.textContent = user.bot ? 'Bot Account' : user.system ? 'System Account' : 'Regular Account';
  
  const icons = {
    'House Bravery': 'img/bravery.png',
    'House Brilliance': 'img/brilliance.png',
    'House Balance': 'img/balance.png'
  };
  
  if (user.badges?.length) {
    elements.badgeWrap.classList.remove('hidden');
    elements.badges.innerHTML = user.badges.map(b => {
      const icon = icons[b] ? `<img src="${icons[b]}" class="w-5 h-5 inline-block mr-2" alt="">` : '✨ ';
      return `<div class="px-4 py-2 bg-[#1e1f22] border border-[#5865f2]/25 rounded-lg text-base text-gray-300 flex items-center">${icon}${b}</div>`;
    }).join('');
  } else {
    elements.badgeWrap.classList.add('hidden');
  }
  
  elements.flags.textContent = user.flags || '0';
  
  const days = user.created_at ? Math.floor((Date.now() - new Date(user.created_at)) / 86400000) : 0;
  
  let age;
  if (days < 60) {
    age = `${days} days`;
  } else if (days < 365) {
    const months = Math.floor(days / 30);
    age = `${months} ${months === 1 ? 'month' : 'months'}`;
  } else if (days < 3650) {
    const years = Math.floor(days / 365);
    age = `${years} ${years === 1 ? 'year' : 'years'}`;
  } else {
    const years = Math.floor(days / 365);
    age = `${years} years`;
  }
  
  elements.age.textContent = age;
  
  if (user.id === '1055337846657007648') {
    elements.type.innerHTML = '👨‍💻 <span class="gradient-text">Developer <3</span>';
  }
};

const copyId = (id, el) => {
  navigator.clipboard.writeText(id);
  const og = el.innerHTML;
  el.innerHTML = '<span class="text-green-500">copied</span>';
  setTimeout(() => el.innerHTML = og, 1500);
};

const search = async () => {
  const id = elements.input.value.trim();
  
  if (!id) {
    notify('missing discord id');
    elements.input.focus();
    return;
  }
  
  if (!/^\d+$/.test(id)) {
    notify('invalid id');
    return;
  }
  
  setLoad(true);
  
  try {
    const res = await fetch(`/lookup?id=${id}`);
    const data = await res.json();
    
    if (data.error) {
      notify(data.error);
      return;
    }
    
    if (data.id) {
      renderProfile(data);
    } else {
      notify('not found');
    }
  } catch (err) {
    notify('request failed');
    console.error(err);
  } finally {
    setLoad(false);
  }
};

window.copyId = copyId;

elements.input.addEventListener('keypress', e => {
  if (e.key === 'Enter' && !loading) search();
});

elements.btn.addEventListener('click', search);

const placeholders = ['paste user id', 'discord id here', '123456789012345678', 'enter id'];
elements.input.placeholder = placeholders[Math.floor(Math.random() * placeholders.length)];
elements.input.focus();
