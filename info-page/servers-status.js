// Time of day background images mapping
const TIME_BACKGROUNDS = {
  night: "night.jpg",
  dawn: "dawn.jpg",
  morning: "sunrise.jpg",
  afternoon: "noon.jpg",
  evening: "evening.jpg",
  sunset: "sunset.jpg"
};

// Global variables for update control
let updateInterval = null;
let isTabActive = true;
let firstStart = true;

// Event listeners for tab visibility changes
document.addEventListener('visibilitychange', handleVisibilityChange);
window.addEventListener('focus', () => { isTabActive = true; startAutoUpdate(); });
window.addEventListener('blur', () => { isTabActive = false; stopAutoUpdate(); });

function handleVisibilityChange() {
  isTabActive = !document.hidden;
  if (isTabActive) {
    loadServerStatus();
    startAutoUpdate();
  } else {
    stopAutoUpdate();
  }
}

function startAutoUpdate() {
  if (typeof UPDATE_FREQUENCY !== 'undefined' && UPDATE_FREQUENCY > 0) {
    const frequency = Math.max(UPDATE_FREQUENCY, 1000);
    const wasInactive = updateInterval === null;

    stopAutoUpdate();
    if (wasInactive && !firstStart) {
      loadServerStatus();
    }

    updateInterval = setInterval(loadServerStatus, frequency);
    firstStart = false;
  }
}

function stopAutoUpdate() {
  if (updateInterval) {
    clearInterval(updateInterval);
    updateInterval = null;
  }
}

// Main function to load and display server status
async function loadServerStatus() {
  const serverStatuses = [];

  for (const server of SERVERS) {
    try {
      let data;

      if (window.mockServersData && server.apiUrl.includes('mock')) {
        const mockKey = `mock-server-${SERVERS.indexOf(server)}`;
        data = window.mockServersData[mockKey] || window.mockServersData[Object.keys(window.mockServersData)[0]];
      } else {
        const response = await fetch(server.apiUrl);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        data = await response.json();
      }

      serverStatuses.push({ server, data });

    } catch (error) {
      console.error(`Error loading server ${server.name}:`, error);
      serverStatuses.push({ server, data: null });
    }
  }

  // Обновляем все карточки после загрузки данных
  updateAllServerCards(serverStatuses);

  // Обновляем год и версию
  document.getElementById('current-year').textContent = new Date().getFullYear();
}

// Create server card HTML
function createServerCard(serverConfig, serverData) {
  const card = document.createElement('div');
  card.className = 'server-card';
  card.dataset.server = serverConfig.name;

  // Set background based on time
  if (serverData.keywords?.time) {
    const time = parseTime(serverData.keywords.time);
    const timeOfDay = getTimeOfDay(time);
    card.style.setProperty('--bg-image', `url('${TIME_BACKGROUNDS[timeOfDay]}')`);
    card.dataset.time = timeOfDay;
  } else {
    // Default to 'afternoon' background if time is missing but server is up
    const fallbackTime = 'afternoon';
    card.style.setProperty('--bg-image', `url('${TIME_BACKGROUNDS[fallbackTime]}')`);
    card.dataset.time = fallbackTime;
  }

  // Server name
  const nameElement = document.createElement('h2');
  nameElement.className = 'server-name';
  nameElement.textContent = serverData.name || serverConfig.name;

  // Description
  const descElement = document.createElement('p');
  descElement.className = 'server-description';
  descElement.textContent = serverData.game || 'DayZ Server';

  // Map name
  const mapElement = document.createElement('p');
  mapElement.className = 'server-map';
  mapName = serverData.map || 'Unknown map';
  mapElement.textContent = mapName == 'enoch' ? 'Livonia' : mapName;

  // Players count
  const playersElement = document.createElement('p');
  playersElement.className = 'players-count';
  if (serverData.keywords?.lqs) {
    playersElement.innerHTML = `Players: ${serverData.players}/${serverData.max_players} ` +
      `(<span class="queue-count">+${serverData.keywords.lqs}</span>)`;
  } else {
    playersElement.textContent = `Players: ${serverData.players}/${serverData.max_players}`;
  }

  // Time of day
  const timeText = document.createElement('div');
  timeText.className = 'time-text';

  if (serverData.keywords?.time) {
    const time = parseTime(serverData.keywords.time);
    timeText.textContent = `Time: ${formatTime(time)}`;
  }

  // Server address
  const addressElement = document.createElement('div');
  addressElement.className = 'server-address';
  addressElement.textContent = `${serverConfig.ip}:${serverData.port}`;

  // Server extra info
  const serverDetails = document.createElement('div');
  serverDetails.className = 'server-details';
  serverDetails.innerHTML = `
    ${serverData.keywords?.dlc ? `
    <div class="detail-row">
      <span class="detail-value">DLC Required</span>
    </div>` : ''}

    ${serverData.keywords?.mod ? `
    <div class="detail-row">
      <span class="detail-value">Mods Required</span>
    </div>` : ''}

    ${serverData.keywords?.whitelist ? `
    <div class="detail-row">
      <span class="detail-value">Whitelist Enabled</span>
    </div>` : ''}

    ${serverData.public ? `
    <div class="detail-row">
      <span class="detail-label">Password:</span>
      <span class="detail-value">${serverData.public ? 'Yes' : 'No'}</span>
    </div>` : ''}

    <div class="detail-row">
      <span class="detail-label">Camera View:</span>
      <span class="detail-value">${serverData.keywords.no3rd ? '1PP' : '3PP'}</span>
    </div>

    ${serverData.keywords?.etm ? `
    <div class="detail-row">
      <span class="detail-label">Day Time Accel:</span>
      <span class="detail-value">x${serverData.keywords.etm.toFixed(1)}</span>
    </div>` : ''}

    ${serverData.keywords?.entm ? `
    <div class="detail-row">
      <span class="detail-label">Night Time Accel:</span>
      <span class="detail-value">x${serverData.keywords.entm.toFixed(1)}</span>
    </div>` : ''}

    <div class="detail-row">
      <span class="detail-label">Game Version:</span>
      <span class="detail-value">${serverData.version || 'unknown'}</span>
    </div>

    <div class="detail-row">
      <span class="detail-label">Server OS:</span>
      <span class="detail-value">${serverData.environment || 'unknown'}</span>
    </div>
  `;

  // Assemble card
  card.appendChild(nameElement);
  card.appendChild(descElement);
  card.appendChild(mapElement);
  card.appendChild(playersElement);
  card.appendChild(timeText);
  card.appendChild(serverDetails);
  card.appendChild(addressElement);

  return card;
}

function updateAllServerCards(serverStatuses) {
  const container = document.getElementById('servers-container');

  let totalPlayers = 0;
  let totalMaxPlayers = 0;
  let latestVersion = null;

  for (const { server, data } of serverStatuses) {
    const existingCard = document.querySelector(`.server-card[data-server="${server.name}"]`);

    let newCard = data
      ? createServerCard(server, data)
      : createErrorCard(server);

    newCard.dataset.server = server.name;

    if (existingCard) {
      container.replaceChild(newCard, existingCard);
    } else {
      container.appendChild(newCard);
    }

    // Collect data for footer
    if (data?.version) {
      latestVersion = data.version;
    }

    if (typeof data?.players === 'number' && typeof data?.max_players === 'number') {
      totalPlayers += data.players;
      totalMaxPlayers += data.max_players;
    }
  }

  // Update footer info
  document.getElementById('current-year').textContent = new Date().getFullYear();

  if (latestVersion) {
    document.getElementById('game-version').textContent = latestVersion;
  }

  document.getElementById('total-online').textContent = `${totalPlayers}/${totalMaxPlayers}`;
}

// Create error card when server is unavailable
function createErrorCard(serverConfig) {
  const card = document.createElement('div');
  card.className = 'server-card';
  card.style.setProperty('--bg-image', "url('offline.jpg')");
  card.dataset.server = serverConfig.name;

  const nameElement = document.createElement('h2');
  nameElement.className = 'server-name';
  nameElement.textContent = serverConfig.name;

  const errorElement = document.createElement('p');
  errorElement.className = 'server-description';
  errorElement.textContent = 'Server is currently unavailable';
  errorElement.style.color = '#ff5555';

  card.appendChild(nameElement);
  card.appendChild(errorElement);

  return card;
}

// Helper functions for time handling
function parseTime(duration) {
  // Convert nanoseconds to hours and minutes
  const hours = Math.floor(duration / 3600000000000) % 24;
  const minutes = Math.floor((duration % 3600000000000) / 60000000000);
  return { hours, minutes };
}

function formatTime(time) {
  return `${String(time.hours).padStart(2, '0')}:${String(time.minutes).padStart(2, '0')}`;
}

function getTimeOfDay(time) {
  const hour = time.hours + time.minutes / 60; // Combine hours + minutes

  // Time periods for Czech Republic in late summer:
  // - Astronomical dawn: ~4:00
  // - Sunrise: ~6:00
  // - Sunset: ~19:30
  // - Astronomical dusk: ~21:00
  if (hour >= 23 || hour < 4) return 'night';
  if (hour >= 4 && hour < 6) return 'dawn';
  if (hour >= 6 && hour < 9) return 'morning';
  if (hour >= 9 && hour < 18) return 'afternoon';
  if (hour >= 18 && hour < 21) return 'sunset';
  return 'evening';
}

// Load server status when page loads and start auto-update if enabled
document.addEventListener('DOMContentLoaded', () => {
  loadServerStatus();
  startAutoUpdate();
});
