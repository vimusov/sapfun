# What?

sapfun - Utility that takes control over your AMD video card coolers to keep it cool and steady.

# Why?

My Sapphire RX 6800 Pulse can not control its own coolers! So that is a kludge to do it programmatically.

# How?

The program works as a daemon. When the GPU temperature rises it runs the coolers faster and vice versa.

# Requirements

- Go >= 1.16;
- `amdgpu` kernel module;

# Usage

Run `sapfun` directly unit under root privileges or start a systemd unit. No command line options or configs are available.

# License

GPL.
