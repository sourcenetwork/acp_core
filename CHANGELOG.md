# Changelog

Significant features added between versions will be contained in the changelog, as well as any breaking changes.

## [0.5.3] - 2025-08-12
Features:
- Add parsed Policy object to ValidatePolicy method in Engine

## [0.5.2] - 2025-04-07
Features:
- Add EventPolicyEdited type

## [0.5.1] - 2025-03-28
Fixes:
- Fixed an issue where the ValidatePolicy method did not return validation errors as part of the response body

## [0.5.0] - 2025-03-27
Features:
- Added EditPolicy engine method
- Added EditPolicy metadata engine method
- Added the `spec` field to the policy dsl, which sets a the policy specification type
- Add CreatePolicyWithSpecification method which constrains a Policy to only be created if the expected spec is satisfied 

Breaking changes:
- acp_core was updated to always include the `owner` relation as one of the authorized relations in a permission expressions

## [0.4.1] - 2025-02-26
Fixes:
- Fixed a bug where sometimes List Policy would fail due to Zanzi returning a malformed record

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
