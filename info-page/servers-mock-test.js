const SERVERS = [
  {    name: "Server 1",    apiUrl: "mock",    ip: "10.123.1.1"  },
  {    name: "Server 2",    apiUrl: "mock",    ip: "10.123.1.12"  },
  {    name: "Server 3",    apiUrl: "mock",    ip: "10.123.1.33"  },
  {    name: "Server 4",    apiUrl: "mock",    ip: "10.123.2.2"  },
  {    name: "Server 5",    apiUrl: "mock",    ip: "10.123.1.15"  },
  {    name: "Server 6",    apiUrl: "mock",    ip: "10.123.7.212"  }
];

function generateRealisticTime(serverIndex) {
  const times = [
    {hours: 1, minutes: 10},   // night.jpg
    {hours: 4, minutes: 12},   // dawn.jpg
    {hours: 7, minutes: 23},   // sunrise.jpg
    {hours: 12, minutes: 2},   // noon.jpg
    {hours: 18, minutes: 37},  // sunset.jpg
    {hours: 21, minutes: 3}    // evening.jpg
  ];

  const { hours, minutes } = times[serverIndex % times.length];
  return (hours * 3600 + minutes * 60) * 1e9;
}

function generateMockServers(count = 6) {
  const maps = ["chernarusplus", "enoch", "sakhal", "namalsk", "deerisle", "banov"];
  const servers = {};

  for (let i = 0; i < count; i++) {
    const hasQueue = Math.random() < 0.25;

    servers[`mock-server-${i}`] = {
      name: `${maps[i].charAt(0).toUpperCase() + maps[i].slice(1)} Server ${i % 2 ? "1" : "3"}PP`,
      game: `Some cool ${maps[i]} server description`,
      map: maps[i],
      players: hasQueue ? 60 : Math.floor(Math.random() * 50) + 10,
      max_players: 60,
      port: 2302 + i * 10,
      version: "1.27.159674",
      environment: hasQueue ? "Windows" : "Linux",
      public: Math.random() < 0.3,
      keywords: {
        time: generateRealisticTime(i),
        lqs: hasQueue ? Math.floor(Math.random() * 5) + 1 : 0,
        dlc: i == 2,
        mod: Math.random() == 1,
        whitelist: Math.random() < 0.2,
        no3rd: i % 2,
        etm: Math.floor(Math.random() * 6) + 2,
        entm: Math.floor(Math.random() * 4) + 1
      }
    };
  }

  return servers;
}

window.mockServersData = generateMockServers();
