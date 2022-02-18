# What?

sapfun - Utility that takes control over your video card coolers to keep it cool and steady. Works with `amdgpu' kernel module only.

# Why?

I have a Sapphire RX 6800 Pulse video card. A few days after buying it I discovered that it suffers from strong overheat. The reason is in the cooling system. Coolers don't start (or keeps rotates slowly) even when the GPU is going to halt due to overheat. I googled WTF is going on and realized that it's "absolutely fine"! Cause there's a bloatware for m$ windows that "solves" this problem. It's called "AMD Radeon Software". You can install it and create a custom cooling profile to make your fans works at high speed when a video card is gonna melt your PC. Awesome!

But what if I am a linuxer? There are a few programs which can do what I need:
- [radeon-profile](https://github.com/marazmista/radeon-profile)
- [horrible bash script](https://gist.github.com/danger89/d13b92e00ad1a5139d58c74ba95a6bc8)
- [and thousands of them](https://unix.stackexchange.com/questions/627182/how-to-lock-fan-speed-for-amd-gpu-in-ubuntu-20-04)

The first one looks promising but I don't need a GUI. The others look not so good as I want.

So there's a great point to create a new pet project! I've been learning Go for some time and felt that I'm ready to make first steps outside the "hello world".

# How?

The program works as a daemon. When the GPU temperature rises it runs the coolers faster and vice versa.

# Build and install

You will need:
- Go >= 1.16. The only standard library is used and no extra modules required.
- `amdgpu` kernel module;

Run `make` to build the binary `sapfun`.
Run `make DESTDIR=... install` to install it wherever you want.

`PKGBUILD` for Arch Linux and systemd unit are also provided in this repository.

# Usage

Run `sapfun` directly unit under root privileges or start a systemd unit. No command line options or configs are available.

# License

GPL.
