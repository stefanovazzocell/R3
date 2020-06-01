# Threat Model

Last Updated: May 24th, 2020.

Important definition in [this document](Definitions).

R3 is designed to protect users and users' content against a range of threats, but there is no such thing as perfect security. In this document we will describe the treat model that R3 aims to follow in its design.

Our main goal, in broad terms, is to allow users to share content securely and privately. A user should generally expect to create or view a Share without being identified and without having to share the private data in the share with the site owner.

Unfortunately is impossible to guarantee perfect security - and I don't claim to do so. The code might contain bugs and design issues that make it vulnerable to attacks; if you're aware of any please [contact me](https://stefanovazzoler.com/#contact).

## What R3 does NOT protect you from

Some things are out of scope of the protection that we attempt to offer due to the nature of how this program works or due to other causes outside of our control.

### Compromised Devices

If your device, the device of those who you share the Share ID to, or if the way you share the Share ID is compromised the Share ID could potentially be transferred to 3rd parties therefore undermining the encryption we offer. We have no control over the devices people use to access our service or how they transfer the Share ID.

### Weak Share ID

We try to generate a reasonably secure random Share ID but if the user modifies it we cannot guarantee that it is not compromised.

### Unauthorized backdoor

If this service of its servers get compromised we cannot guarantee the security of the encryption (or encryption at all). We do our best to protect our service against such compromise, but if we or any of our service providers gets hacked we cannot guarantee anything.

## What R3 does protect you from

After a security and encryption review of R3 is completed it should be safe for most users to use the service to share non-sensitive content. If you are a target of sophisticated (or targeted) attacks we advise you not to use our service.

In general we try to protect against the following threats.

### Future DB leaks

If our Database ever gets compromised and/or leaked at a future time you should have a reasonable expectation that your Share is secure against all but the most sophisticated targeted attacks.

### Unsophisticated Non-Targeted Attacks

We do our best to protect your data through the use of multiple level of security measures against unsophisticated and non-targeted attacks. We employ industry standard methods to secure and lock down our servers against both external threats and potential escalation of privileges if they get partially compromised.

### Service provider liability

Although this one is a little tricky as it depends on local laws, the server owner will not have insights over what is shared with the service and therefore makes it harder to make them liable for the content share. This is especially useful for private individuals like the creator of R3 that do not have the time to monitor users' shares and also is most certainly against any kind of monitoring (read: spying) of what users' share.

### Privacy

The users should have a reasonable expectation of privacy. If the service cannot get compromised and maliciously changed (including by the service provider) the users data and identities should remain private and not be collected.
