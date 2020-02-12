# CHANGELOG-0.5

## [Unreleased](https://github.com/anz-bank/sysl/tree/HEAD)

[Full Changelog](https://github.com/anz-bank/sysl/compare/v0.5.2...HEAD)

**Closed issues:**

- Update all GitHub Actions workflow files [\#572](https://github.com/anz-bank/sysl/issues/572)
- Update release documentation [\#536](https://github.com/anz-bank/sysl/issues/536)
- Alias bug [\#473](https://github.com/anz-bank/sysl/issues/473)
- Better XSD import handling [\#369](https://github.com/anz-bank/sysl/issues/369)

**Merged pull requests:**

- Remove spurious logging [\#591](https://github.com/anz-bank/sysl/pull/591) ([camh-anz](https://github.com/camh-anz))
- Fix the importer so the import mode is set in the nested loader [\#589](https://github.com/anz-bank/sysl/pull/589) ([anz-gordonj7](https://github.com/anz-gordonj7))
- fix: remove misadded file [\#585](https://github.com/anz-bank/sysl/pull/585) ([ChloePlanet](https://github.com/ChloePlanet))
- Update GRPC proto generation tutorial [\#577](https://github.com/anz-bank/sysl/pull/577) ([joshcarp](https://github.com/joshcarp))
- Update CI workflows to use GoReleaser [\#573](https://github.com/anz-bank/sysl/pull/573) ([ChloePlanet](https://github.com/ChloePlanet))
- update releasing.md [\#571](https://github.com/anz-bank/sysl/pull/571) ([ChloePlanet](https://github.com/ChloePlanet))

## [v0.5.2](https://github.com/anz-bank/sysl/tree/v0.5.2) (2020-02-06)

[Full Changelog](https://github.com/anz-bank/sysl/compare/v0.5.1...v0.5.2)

**Implemented enhancements:**

- swagger/openAPI does not handle remote file references [\#509](https://github.com/anz-bank/sysl/issues/509)

**Fixed bugs:**

- Get latest git release info error [\#558](https://github.com/anz-bank/sysl/issues/558)

**Closed issues:**

- CI doesn't fail properly on Linux and Darwin [\#564](https://github.com/anz-bank/sysl/issues/564)
- New Command: sysl version [\#552](https://github.com/anz-bank/sysl/issues/552)
- New Command: sysl env [\#550](https://github.com/anz-bank/sysl/issues/550)
- Add sysl modules documentation [\#547](https://github.com/anz-bank/sysl/issues/547)
- codegen: invalid pointer type error using string type alias [\#546](https://github.com/anz-bank/sysl/issues/546)
- Build and upload sysl binaries while new version releases [\#534](https://github.com/anz-bank/sysl/issues/534)
- Docker image for go that sysl modules depends on [\#532](https://github.com/anz-bank/sysl/issues/532)
- sysl-modules: module's path is wrong when root marker exists [\#519](https://github.com/anz-bank/sysl/issues/519)
- Replace filepath functions with afero functions [\#518](https://github.com/anz-bank/sysl/issues/518)
- Adding a dependancy path to sysl transform [\#514](https://github.com/anz-bank/sysl/issues/514)
- No warnings on duplicated import [\#471](https://github.com/anz-bank/sysl/issues/471)
- Afero File Writing \(indirect\_1.png\) [\#359](https://github.com/anz-bank/sysl/issues/359)
- regenerate travis secure env variables for travis-ci.com [\#221](https://github.com/anz-bank/sysl/issues/221)
- .travis secure tokens: export opkM=\[secure\] is for what? [\#219](https://github.com/anz-bank/sysl/issues/219)
- Fix CI PATH management [\#167](https://github.com/anz-bank/sysl/issues/167)

**Merged pull requests:**

- go mod tidy [\#569](https://github.com/anz-bank/sysl/pull/569) ([ChloePlanet](https://github.com/ChloePlanet))
- Fixes \#564 [\#565](https://github.com/anz-bank/sysl/pull/565) ([nofun97](https://github.com/nofun97))
- Ensure the temlpate results are added to the results map [\#563](https://github.com/anz-bank/sysl/pull/563) ([anz-gordonj7](https://github.com/anz-gordonj7))
- Inject binary info through ldflags [\#562](https://github.com/anz-bank/sysl/pull/562) ([ChloePlanet](https://github.com/ChloePlanet))
- Vendor in aqwari.net/xml [\#561](https://github.com/anz-bank/sysl/pull/561) ([anz-gordonj7](https://github.com/anz-gordonj7))
- Bump coverage pass mark 80=\>85% [\#560](https://github.com/anz-bank/sysl/pull/560) ([anzdaddy](https://github.com/anzdaddy))
- Feature/really really mergeme [\#557](https://github.com/anz-bank/sysl/pull/557) ([anz-gordonj7](https://github.com/anz-gordonj7))
- Enrich issue templates [\#556](https://github.com/anz-bank/sysl/pull/556) ([ChloePlanet](https://github.com/ChloePlanet))
- feat: add release version to sysl info [\#555](https://github.com/anz-bank/sysl/pull/555) ([ChloePlanet](https://github.com/ChloePlanet))
- Add command sysl env [\#554](https://github.com/anz-bank/sysl/pull/554) ([ChloePlanet](https://github.com/ChloePlanet))
- Update PR template [\#549](https://github.com/anz-bank/sysl/pull/549) ([ChloePlanet](https://github.com/ChloePlanet))
- Add sysl modules docs [\#548](https://github.com/anz-bank/sysl/pull/548) ([ChloePlanet](https://github.com/ChloePlanet))
- dep-path should not be compulsory [\#544](https://github.com/anz-bank/sysl/pull/544) ([AriehSchneier](https://github.com/AriehSchneier))
- fix module path [\#542](https://github.com/anz-bank/sysl/pull/542) ([ChloePlanet](https://github.com/ChloePlanet))
- Enable sysl modules feature by default [\#540](https://github.com/anz-bank/sysl/pull/540) ([ChloePlanet](https://github.com/ChloePlanet))
- Remove .travisci and reljam mentions [\#539](https://github.com/anz-bank/sysl/pull/539) ([joshcarp](https://github.com/joshcarp))
- Build Docker image for sysl [\#535](https://github.com/anz-bank/sysl/pull/535) ([ChloePlanet](https://github.com/ChloePlanet))
- Release sysl binaries [\#533](https://github.com/anz-bank/sysl/pull/533) ([ChloePlanet](https://github.com/ChloePlanet))
- Added per file duplicate module import warning [\#530](https://github.com/anz-bank/sysl/pull/530) ([nydrani](https://github.com/nydrani))
- Add examples makefile target [\#527](https://github.com/anz-bank/sysl/pull/527) ([joshcarp](https://github.com/joshcarp))
- Fix typo in example docs [\#525](https://github.com/anz-bank/sysl/pull/525) ([cuminandpaprika](https://github.com/cuminandpaprika))
- Issue 473 [\#522](https://github.com/anz-bank/sysl/pull/522) ([ericzhang6222](https://github.com/ericzhang6222))
- Fix modules path when root marker exists [\#521](https://github.com/anz-bank/sysl/pull/521) ([ChloePlanet](https://github.com/ChloePlanet))
- Add new cmd param [\#517](https://github.com/anz-bank/sysl/pull/517) ([ashwinsajiv](https://github.com/ashwinsajiv))
- Put imported base paths in @basePath, not part of the actual SYSL path [\#508](https://github.com/anz-bank/sysl/pull/508) ([gkanz](https://github.com/gkanz))

## [v0.5.1](https://github.com/anz-bank/sysl/tree/v0.5.1) (2020-01-17)

[Full Changelog](https://github.com/anz-bank/sysl/compare/v0.5.0...v0.5.1)

**Implemented enhancements:**

- Make sysl playground [\#470](https://github.com/anz-bank/sysl/issues/470)

**Closed issues:**

- Missing newline at the end of the generated svg file [\#511](https://github.com/anz-bank/sysl/issues/511)
- output raw plantuml [\#503](https://github.com/anz-bank/sysl/issues/503)
- importing query params with hyphens generate bad sysl [\#499](https://github.com/anz-bank/sysl/issues/499)
- Upgrade CI to Go 1.13 [\#493](https://github.com/anz-bank/sysl/issues/493)
- Codegen generate wrong go code - Redundant semicolon [\#475](https://github.com/anz-bank/sysl/issues/475)
- Update go openapi importer to match Python importer [\#434](https://github.com/anz-bank/sysl/issues/434)
- Support Sysl imports across repos [\#415](https://github.com/anz-bank/sysl/issues/415)

**Merged pull requests:**

- fixes \#509, swagger now works with remote refs \(files only, no URLs\) [\#516](https://github.com/anz-bank/sysl/pull/516) ([nofun97](https://github.com/nofun97))
- Append newline to generated svg file [\#515](https://github.com/anz-bank/sysl/pull/515) ([ChloePlanet](https://github.com/ChloePlanet))
- Pass base path as a command line parameter for codegen [\#510](https://github.com/anz-bank/sysl/pull/510) ([gkanz](https://github.com/gkanz))
- Sysl-Playground Links [\#507](https://github.com/anz-bank/sysl/pull/507) ([joshcarp](https://github.com/joshcarp))
- Increase base path conversion between v2 and v3 alignment to spec [\#506](https://github.com/anz-bank/sysl/pull/506) ([gkanz](https://github.com/gkanz))
- Fix issue 475. [\#505](https://github.com/anz-bank/sysl/pull/505) ([ericzhang6222](https://github.com/ericzhang6222))
- added sysl-safe names checker and converter [\#502](https://github.com/anz-bank/sysl/pull/502) ([nofun97](https://github.com/nofun97))
- Import sysl files across repos [\#500](https://github.com/anz-bank/sysl/pull/500) ([ChloePlanet](https://github.com/ChloePlanet))
- upgrade to go 1.13 [\#496](https://github.com/anz-bank/sysl/pull/496) ([nofun97](https://github.com/nofun97))
- Feature/refactor listener [\#495](https://github.com/anz-bank/sysl/pull/495) ([anz-gordonj7](https://github.com/anz-gordonj7))
- Switch to displaying an actual diff when parse tests fail [\#494](https://github.com/anz-bank/sysl/pull/494) ([anz-gordonj7](https://github.com/anz-gordonj7))
- Postgres script generation using sysl [\#488](https://github.com/anz-bank/sysl/pull/488) ([sidhartha-priyadarshi](https://github.com/sidhartha-priyadarshi))
- Documentation [\#476](https://github.com/anz-bank/sysl/pull/476) ([joshcarp](https://github.com/joshcarp))

## [v0.5.0](https://github.com/anz-bank/sysl/tree/v0.5.0) (2019-12-20)

[Full Changelog](https://github.com/anz-bank/sysl/compare/v0.4.0...v0.5.0)

**Merged pull requests:**

- Replaces the swagger v2 importer with a converter to V3 and then import [\#492](https://github.com/anz-bank/sysl/pull/492) ([anz-gordonj7](https://github.com/anz-gordonj7))
- Fix installation docs [\#490](https://github.com/anz-bank/sysl/pull/490) ([juliaogris](https://github.com/juliaogris))
- Remove Python setup.cfg [\#489](https://github.com/anz-bank/sysl/pull/489) ([juliaogris](https://github.com/juliaogris))
- Move tests to root directory [\#486](https://github.com/anz-bank/sysl/pull/486) ([joshcarp](https://github.com/joshcarp))
- Initial work to import openapi 3 [\#483](https://github.com/anz-bank/sysl/pull/483) ([anz-gordonj7](https://github.com/anz-gordonj7))
- Documentation  [\#468](https://github.com/anz-bank/sysl/pull/468) ([joshcarp](https://github.com/joshcarp))



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
