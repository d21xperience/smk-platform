# @sdp/storage

Shared Storage abstraction for School Digital Platform (SDP).

## Overview

Package ini menyediakan unified storage interface untuk seluruh aplikasi SDP:

- `siakad-tu`
- `siakad-guru`
- `psb`
- `website`
- Aplikasi masa depan

## Fitur

- LocalStorageAdapter — Persistent storage (browser localStorage)
- SessionStorageAdapter — Session-scoped storage (browser sessionStorage)
- MemoryStorageAdapter — In-memory storage (testing/SSR/fallback)
- Context-aware key building — Mencegah kebocoran data antar context
- Auto-detect — Fallback otomatis jika storage tidak tersedia
- JSON serialization — Automatic serialize/deserialize
- Namespace isolation — Prefix support
