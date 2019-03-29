# Ferret Server
<p align="center">
	<a href="https://goreportcard.com/report/github.com/MontFerret/ferret-server">
		<img alt="Go Report Status" src="https://goreportcard.com/badge/github.com/MontFerret/ferret-server">
	</a>
	<a href="https://travis-ci.com/MontFerret/ferret-server">
		<img alt="Build Status" src="https://travis-ci.com/MontFerret/ferret-server.svg?branch=master">
	</a>
	<a href="https://codecov.io/gh/MontFerret/ferret-server">
		<img src="https://codecov.io/gh/MontFerret/ferret-server/branch/master/graph/badge.svg" />
	</a>
	<a href="https://discord.gg/kzet32U">
		<img alt="Discord Chat" src="https://img.shields.io/discord/501533080880676864.svg">
	</a>
	<a href="https://github.com/MontFerret/ferret-server/releases">
		<img alt="Ferret release" src="https://img.shields.io/github/release/MontFerret/ferret-server.svg">
	</a>
	<a href="http://opensource.org/licenses/MIT">
		<img alt="MIT License" src="http://img.shields.io/badge/license-MIT-brightgreen.svg">
	</a>
</p>

Server for advanced web scraping.    
[Open API defintion.](https://next.stoplight.io/montferret/ferret-server/)

# Features
- Scripts persistence
- Scraped data persistence
- Script execution scheduling
- Integration with 3rd party systems
- Web Hooks
- Security

# WIP
Be aware, that the project is under heavy development.    
There is no documentation and some things may change in the final release.

# Installation
### Binary
You can download latest binaries from [here](https://github.com/MontFerret/ferret-server/releases).

### Source code
#### Production
* Go >=1.11
* Chrome or Docker
* ArangoDB

#### Development
* GNU Make
* ANTLR4 >=4.7.1

## Quick start

```sh
ferret-server --db=http://0.0.0.0:8529
```
