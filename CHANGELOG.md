## [1.3.2](https://github.com/cactus-platform/cmaestro-core/compare/v1.3.1...v1.3.2) (2026-08-30)


### Bug Fixes

* changes Repository::CreateRevision behaviour if ID is used ([902afb4](https://github.com/cactus-platform/cmaestro-core/commit/902afb41f2325b0dd2b5c9d86f35af3aad28c044))

## [1.3.1](https://github.com/cactus-platform/cmaestro-core/compare/v1.3.0...v1.3.1) (2026-08-30)


### Bug Fixes

* updates repository to a one-to-many relationship with artifact ([63140f5](https://github.com/cactus-platform/cmaestro-core/commit/63140f51dcd83cfd1ec34ce3f40dd97317dc87e2))

# [1.3.0](https://github.com/cactus-platform/cmaestro-core/compare/v1.2.1...v1.3.0) (2026-08-30)


### Features

* updates repository to a one-to-many relationship with artifact ([2db7236](https://github.com/cactus-platform/cmaestro-core/commit/2db72369ea5922a6433e1928dd70660e5cb2da71))

## [1.2.1](https://github.com/cactus-platform/cmaestro-core/compare/v1.2.0...v1.2.1) (2026-08-30)


### Bug Fixes

* Repository management takes Repository model as input ([bb8122b](https://github.com/cactus-platform/cmaestro-core/commit/bb8122b720b50815560ce013f0e3d22f9d9a670e))

# [1.2.0](https://github.com/cactus-platform/cmaestro-core/compare/v1.1.0...v1.2.0) (2026-08-30)


### Features

* adds repository management foundation (Model, Repository, Service) ([bfc6ac8](https://github.com/cactus-platform/cmaestro-core/commit/bfc6ac805401cbcd6d4da666f7ff2f5b34b03a51))

# [1.1.0](https://github.com/cactus-platform/cmaestro-core/compare/v1.0.1...v1.1.0) (2026-08-30)


### Features

* moves artifact+ingest controls from cactus to cactus-core ([09d55ea](https://github.com/cactus-platform/cmaestro-core/commit/09d55eab5dd0dae4d81d547fa4a38c044d04fc38))

## [1.0.1](https://github.com/cactus-platform/cmaestro-core/compare/v1.0.0...v1.0.1) (2026-08-28)


### Bug Fixes

* update workflow name in release.yaml ([cae2ea0](https://github.com/cactus-platform/cmaestro-core/commit/cae2ea0525faa18e98a50300937de065c18404c0))

# 1.0.0 (2026-08-28)


### Bug Fixes

* changes module name ([8368825](https://github.com/cactus-platform/cmaestro-core/commit/8368825d430db2c478f395235bd156f3fdf11961))


### Features

* adds ArtifactRepository implementation ([694c846](https://github.com/cactus-platform/cmaestro-core/commit/694c8463c85447d0e6747d6e4be3f8e85986cf9e))
* adds bucket (seaweedfs) connector ([bc69cbb](https://github.com/cactus-platform/cmaestro-core/commit/bc69cbbdd8f0887631e96dbba944e1c430039c84))
* adds dbutil utils ([16c0307](https://github.com/cactus-platform/cmaestro-core/commit/16c03079d62731ef9e622353a8aa2964ab4938f3))
* adds models foundation ([cab031a](https://github.com/cactus-platform/cmaestro-core/commit/cab031aeddf2b955d562d3b43105fd0bbf368349))
* adds postgres conn helper ([0ee9389](https://github.com/cactus-platform/cmaestro-core/commit/0ee93894b00ce258e4c27cb59f3793964bb639d3))
* adds redis helper (as keyval) + install redis client dependencies ([0a32b8a](https://github.com/cactus-platform/cmaestro-core/commit/0a32b8affd7e37a565b81b1ae515f308c3e24d44))
* adds semantic release workflow ([0872cef](https://github.com/cactus-platform/cmaestro-core/commit/0872cef8c8c06a57f34d155a06e67aea2b2f8197))
* implements changelog auto writer using semantic release workflow ([97d29da](https://github.com/cactus-platform/cmaestro-core/commit/97d29da672d3871d877b024a10f99d8710a01137))
