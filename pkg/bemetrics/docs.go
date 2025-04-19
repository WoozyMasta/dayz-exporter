// Package bemetrics provides a collector for DayZ server metrics, including:
// - Server information via A2S queries (players online, slots, queue, etc.)
// - Player statistics via BattlEye RCON (ping, online status, lobby/invalid players)
// - Ban information via BattlEye RCON (GUID and IP bans with durations)
//
// The package exposes metrics in Prometheus format and allows customization
// through additional labels. It handles both static server information and
// dynamic player/ban data with proper metric initialization and updates.
package bemetrics
