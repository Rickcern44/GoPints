## Purpose

Defines server-side import of a remotely-hosted image (beer image or brewery logo) into a keg's existing blob storage via an admin-supplied URL, gated behind a feature flag and validated against SSRF risks, giving admins a faster alternative to the existing upload-from-disk flow.

## Requirements

### Requirement: Feature flag gates URL import
The system SHALL expose URL-based image import (both endpoints and the admin UI affordance) only when the `remote_image_urls` feature flag is enabled, and SHALL default it to disabled.

#### Scenario: Flag disabled by default
- **WHEN** the `remote_image_urls` feature has never been explicitly enabled
- **THEN** `GET /api/features` reports `remote_image_urls: false`, the admin UI does not show a URL-input option, and the `from-url` endpoints return 404

#### Scenario: Flag enabled by admin
- **WHEN** an authenticated admin sets `PUT /api/features/remote_image_urls` to enabled
- **THEN** `GET /api/features` reports `remote_image_urls: true`, the admin UI shows the URL-input option, and the `from-url` endpoints become reachable

### Requirement: Server-side fetch of admin-supplied image URL
The system SHALL provide `POST /api/kegs/{id}/image/from-url` and `POST /api/kegs/{id}/brewery-image/from-url`, each accepting a JSON body `{ "url": string }`, requiring authentication, and requiring the `remote_image_urls` feature flag to be enabled.

#### Scenario: Successful import
- **WHEN** an authenticated admin submits a valid, reachable, direct (non-redirecting) image URL under `maxImageBytes` in size
- **THEN** the server downloads the image, stores it via the same blob storage used by direct upload, and the keg's `image_mime_type` (or `brewery_image_mime_type`) reflects the sniffed content type

#### Scenario: Unauthenticated request rejected
- **WHEN** a request to either `from-url` endpoint is made without valid admin authentication
- **THEN** the server returns 401/403 and does not attempt any outbound fetch

#### Scenario: Feature flag disabled
- **WHEN** `remote_image_urls` is disabled
- **THEN** both `from-url` endpoints return 404 regardless of authentication state

### Requirement: SSRF-safe fetch validation
The system SHALL validate every outbound image fetch against loopback, private, link-local, and unspecified address ranges, SHALL reject non-http(s) schemes, and SHALL NOT follow HTTP redirects.

#### Scenario: URL targets a private/loopback/link-local address
- **WHEN** the supplied URL resolves to an address in `127.0.0.0/8`, `::1`, `10/8`, `172.16/12`, `192.168/16`, `169.254.0.0/16` (including the cloud metadata address `169.254.169.254`), or `fc00::/7`/`fe80::/10`
- **THEN** the server rejects the request without attempting to connect, returning a client error

#### Scenario: URL uses a non-http(s) scheme
- **WHEN** the supplied URL's scheme is not `http` or `https` (e.g. `file://`, `ftp://`)
- **THEN** the server rejects the request without attempting to connect

#### Scenario: Remote server responds with a redirect
- **WHEN** the fetch receives any 3xx response
- **THEN** the server treats this as a failure and does not follow the redirect, returning an error indicating a direct image URL is required

#### Scenario: Fetch exceeds size limit
- **WHEN** the response body exceeds `maxImageBytes`
- **THEN** the server aborts the download and returns an error, storing nothing

#### Scenario: Fetch exceeds time limit
- **WHEN** the remote server does not respond within the configured timeout
- **THEN** the server aborts the request and returns an error, storing nothing

#### Scenario: Downloaded content is not a valid image
- **WHEN** the downloaded bytes do not sniff (via content detection, not the remote's declared header) as a supported image type
- **THEN** the server rejects the import and returns an error, storing nothing
