# Privacy Policy

## TL/DR
1. We collect as little data as possible (and only if you view/create a link).
2. We delete the data we collect as soon as possible (when a share expires or when it is no longer useful).
3. We use [CloudFlare](https://www.cloudflare.com/privacypolicy/).

## What we collect
If you just browse our website we don't collect any data.
If you view a share we only modify the share record to decrease the view count.
If you create or edit a share only collect the following:
1. An hash derived by the ID.
2. The encrypted data that you're sharing.
3. Number of views and time before expiration.
4. [optional] A hash derived by your edit password.
5. A non-identifiable version of your hashed IP* protected using [K-anonymity](https://en.wikipedia.org/wiki/K-anonymity).
*Deleted within 2h, everything else deleted after the share expire.

## Why we collect
We collect the afformentioned data to provide the service and to protect our service from abuse.

## Data Exchange
We do not (and will never) sell, license or sub-license any of the data submitted directly or indirectly by our users with any person or entity.
We will only share the encrypted payload of a link IF a user provides a valid ID-derived hash matching the share and such share has not expired.

## Data Deletion
We do NOT collect PII (Personal Identifiable Information); any IP-derived hash with k-anonymity will be automatically deleted within 2h from creation.
Unless there are extraordinary circumstances we won't delete data manually as we cannot verify link ownership (unless you specified a edit password hash and a ID-derived hash).

## Data Request
We do NOT collect PII (Personal Identifiable Information) therefore we cannot provide any personal data that you upload.
You can download data for a given share using our service. We will NOT provide any expiration data (remaining views / TTL) to users.

## Others
We rely on [CloudFlare](https://www.cloudflare.com/privacypolicy/) to provide the service, please review their Privacy Policy.
Note that CloudFlare is only used for the regular domain. The ".onion" hidden service is served directly by us, therefore you can skip the note above.