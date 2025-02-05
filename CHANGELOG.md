# Changelog

Significant features added between versions will be contained in the changelog, as well as any breaking changes.

## [0.4.0] - 2025-02-05
Features:
- Add RevealRegistrationCmd which allows registering an object at a specific timestamp

Improvements and Fixes:
- Improved Object Owner query lookup performance
- Adds benchmark tests for Object and Relationship creation

## [0.3.0] - 2024-11-06

Features:
- TransferObject command
- AmendRegistration command
- UnarchiveObject command

Breaking changes:
- Registering an archived object no longer unarchives it, use UnarchiveObject instead
- UnregisterObject command removed, use ArchiveObject instead

## [0.2.0] - 2024-08-30

Features:
- Adds a WASM module for JS runtime which exports PlaygroundService

## [0.1.0] - 2024-08-20

ACP Core initial release.
Contains implementation of Source's Decentralized Administration Discretionary Access Control engine.
