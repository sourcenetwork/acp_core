# ACP Core

ACP Core is a Go module which exposes an Access Control Engine.

The access control paradigm it implements can be classified as "descentralized administration discretionary access control" as defined in [SS94].
In short, "Discretionary Access Control" means that object owners are allowed to share their objects with other actors, "descentralized administration" refers to the fact that an object owner can delegate the ability to share access with other actors.

The engine is designed to be multi tenant.
It can accept multiple Policy definitions, each of which provides context isolation.

ACP Core was designed to be used as a library in some external context.
Currently ACP core provides the logic and functionality required by SourceHub's ACP module.

## Documentation

## Development Instructions

### Protobuff
### Testing
### Building
### WASM

## References:
SS94: https://ieeexplore.ieee.org/abstract/document/312842
