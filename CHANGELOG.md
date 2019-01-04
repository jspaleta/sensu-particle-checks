# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [0.0.2] - 2019-01-03

### Changes
-  Treat particle deviceID and access token as security sensitive inputs. Default is to pull them from env, and use cmdline arguments only as an override.
-  PARTICLE_DEVICEID => --device   
-  PARTICLE_TOKEN => --access_token  

## [0.0.1] - 2018-12-29

### Added
- Initial release
