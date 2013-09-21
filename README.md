Description
===========

Fspoolr is a dynamic artifact generation and serving framework. Artifacts are packages hat can be generated using instructions given in artifact description files. Artifact descriptions contain dependencies and build instructions. While fspoolr is running, any time a dependency changes, the artifact's build procedure will be rerun.

Todo
====

- Add support for caching successful artifact builds
- Add multiple build targets for artifacts
- Refactor into distributed architecture by creating socket stub to intermediate between filemonitor and statemanager