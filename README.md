I want to share my location with my friends, and many of them want to share theirs with me.

Why do location sharing apps all suck?
- Life360: Sells all your data to advertisers and insurers. Ever wonder why their car crash detection is so advanced compared to the rest of the program?
- Apple Find My: now he sings (iOS), now he sobs (Android)
- Google Maps Sharing: NUH UH
- OwnTracks/Traccar/Dawarich/Overland/Home Assistant Companion/GPSLogger: anyone actually understand how to set these up for social use?
- Grid, WhereFam: E2EE makes this *too* private, trusted servers do infact provide convenience benefits

<img src="assets/venn.svg" width=500 />

# SIYA: Single(s) In Your Area
Self-hosted location sharing app for you and your friends :)

In terms of architecture: the backend is in Go exposing a Chi API; the frontend is in SvelteKit (TS). It connects to an existing Traccar server and exposes a website that acts as the viewing portal for clients, which can be added as a web app on mobile.

Advantages
- You alone control the data
- No ads or subscriptions
- iOS and Android compatible with web app

Disadvantages
- Requires a server
- Requires both an app and a website (which can be added as a web app) to function properly. One for sending, one for viewing. There is only hate in my heart for app store licensing costs.

# Roadmap

## Backend

- [x] Fetch devices and latest positions from Traccar (REST).
- [x] Store latest positions in memory.
- [ ] Refresh locations continuously (REST).
- [ ] Refresh locations continuously (websockets).
- [ ] Setup Chi API.
- [ ] Stream live locations (REST).
- [ ] Stream live locations (websockets).
- [ ] Get history in range from Traccar.
- [ ] Expose history in range.

## Auth

- [ ] Siya accounts. 
- [ ] Sign-in.
- [ ] Harden security.
- [ ] Link accounts to devices.
- [ ] Perms mechanism e.g. circles - TBD.
- [ ] Fine-grained perms (live, history, etc).
- [ ] Temporary sharing rules.

## Frontend

- [ ] Make contact with backend.
- [ ] Map.
- [ ] Show current locations on map.
- [ ] Live update functionality (last update time, connection status, stale marker).
- [ ] Controls for sign-in, access.
- [ ] Location history.
